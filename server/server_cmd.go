package server

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/louch2010/gocache/conf"
	"github.com/louch2010/gocache/core"
	"github.com/louch2010/goutil"
)

//帮助命令处理
func HandleHelpCommnd(body string, client Client) string {
	response := ""
	help := conf.GetHelpConfig()
	if len(body) == 0 { //没有请求体，则显示所有命令名称
		for index, sec := range help.GetSectionList() {
			response += "[" + strconv.Itoa(index+1) + "] " + sec + "\r\n"
		}
		response += "use 'help commnd' to see detail info"
	} else {
		body = strings.ToLower(body)
		sec, err := help.GetSection(body)
		if err != nil {
			response = "no help for the commnd"
		} else {
			response += "[" + body + "]\r\n"
			for k, v := range sec {
				response += k + ": " + v + "\r\n"
			}
		}
	}
	return response
}

//连接命令处理connect [-t'table'] [-a'pwd'] [-i'ip'] [-p'port'] [-e'e1,e2...']
func HandleConnectCommnd(body, token string) string {
	table := conf.GetSystemConfig().MustValue("table", "default", core.DEFAULT_TABLE_NAME)
	var pwd, ip, port, event, protocol string
	args := strings.Split(body, " ")
	for _, arg := range args {
		//参数长度小于3或不是以-开头，说明参数不对，直接跳过
		if len(arg) < 3 || !strings.HasPrefix(arg, "-") {
			continue
		}
		paramType := arg[1]
		param := arg[2:len(arg)]
		switch paramType {
		case 't':
			table = param
			break
		case 'a':
			pwd = param
			break
			break
		case 'i':
			ip = param
			break
		case 'p':
			port = param
			break
		case 'e':
			event = param
			break
		default:
		}
	}
	//密码校验
	syspwd := conf.GetSystemConfig().MustValue("server", "password", "")
	if len(syspwd) > 0 {
		if len(pwd) == 0 {
			return ERROR_AUTHORITY_NO_PWD.Error()
		}
		if syspwd != pwd {
			return ERROR_AUTHORITY_PWD_ERROR.Error()
		}
	}
	//端口校验
	portInt := 0
	if len(port) > 0 {
		p, err := strconv.Atoi(port)
		if err != nil {
			log.Println("端口转换错误，", err)
			return ERROR_PORT_ERROR.Error()
			portInt = p
		}
	}
	//存储连接信息
	client := Client{
		host:        ip,
		port:        portInt,
		table:       table,
		listenEvent: strings.Split(event, ","),
		protocol:    protocol,
		token:       token,
	}
	CreateSession(token, client)
	return MESSAGE_SUCCESS
}

//Set命令处理
func HandleSetCommnd(body string, client Client) string {
	if len(body) == 0 {
		return ERROR_COMMND_PARAM_ERROR.Error()
	}
	//参数处理
	args := strings.Split(body, " ")
	if len(args) < 2 || len(args) > 3 {
		return ERROR_COMMND_PARAM_ERROR.Error()
	}
	var liveTime int = 0
	if len(args) == 3 {
		t, err := strconv.Atoi(args[2])
		if err != nil {
			log.Println("参数转换错误，liveTime：", args[2], err)
			return ERROR_COMMND_PARAM_ERROR.Error()
		}
		liveTime = t
	}
	//增加缓存项
	table, err := core.Cache(client.table)
	if err != nil {
		log.Println("执行Set命令出错，获取表失败，表名：", client.table)
		return MESSAGE_ERROR
	}
	table.Set(args[0], args[1], time.Duration(liveTime)*time.Second)
	return MESSAGE_SUCCESS
}

//Get命令处理
func HandleGetCommnd(body string, client Client) string {
	if len(body) == 0 {
		return ERROR_COMMND_PARAM_ERROR.Error()
	}
	//参数处理
	args := strings.Split(body, " ")
	if len(args) != 1 {
		return ERROR_COMMND_PARAM_ERROR.Error()
	}
	table, err := core.Cache(client.table)
	if err != nil {
		log.Println("执行Get命令出错，获取表失败，表名：", client.table)
		return MESSAGE_ERROR
	}
	item := table.Get(body)
	if item == nil {
		return ERROR_ITEM_NOT_EXIST.Error()
	}
	return item.Value().(string)
}

//Delete命令处理
func HandleDeleteCommnd(body string, client Client) string {
	if len(body) == 0 {
		return ERROR_COMMND_PARAM_ERROR.Error()
	}
	//参数处理
	args := strings.Split(body, " ")
	if len(args) != 1 {
		return ERROR_COMMND_PARAM_ERROR.Error()
	}
	table, err := core.Cache(client.table)
	if err != nil {
		log.Println("执行Delete命令出错，获取表失败，表名：", client.table)
		return MESSAGE_ERROR
	}
	if table.Delete(body) {
		return MESSAGE_SUCCESS
	}
	return MESSAGE_ERROR
}

//Exist命令处理
func HandleExistCommnd(body string, client Client) string {
	if len(body) == 0 {
		return ERROR_COMMND_PARAM_ERROR.Error()
	}
	//参数处理
	args := strings.Split(body, " ")
	if len(args) != 1 {
		return ERROR_COMMND_PARAM_ERROR.Error()
	}
	table, err := core.Cache(client.table)
	if err != nil {
		log.Println("执行Exist命令出错，获取表失败，表名：", client.table)
		return MESSAGE_ERROR
	}
	if table.IsExist(body) {
		return MESSAGE_SUCCESS
	}
	return MESSAGE_ERROR
}

//切换表
func HandleUseCommnd(body string, client Client) string {
	if len(body) == 0 {
		return ERROR_COMMND_PARAM_ERROR.Error()
	}
	client.table = body
	if CreateSession(client.token, client) {
		return MESSAGE_SUCCESS
	}
	return MESSAGE_ERROR
}

//显示表信息
func HandleShowtCommnd(body string, client Client) string {
	response := ""
	if len(body) == 0 { //没有请求体，则显示所有表名
		list := core.GetCacheTables()
		index := 1
		for k, _ := range list {
			if k == client.table {
				response += "[* " + strconv.Itoa(index) + "] "
			} else {
				response += "[" + strconv.Itoa(index) + "] "
			}
			response += k + "\r\n"
			index += 1
		}
		response += "use 'showt tableName' to see detail info"
	} else {
		table, ok := core.GetCacheTable(body)
		if !ok {
			return ERROR_TABLE_NOT_EXIST.Error()
		}
		response += "name:" + table.Name() + "\r\n"
		response += "itemCount: " + strconv.Itoa(table.ItemCount()) + "\r\n"
		response += "createTime: " + goutil.DateUtil().TimeFullFormat(table.CreateTime()) + "\r\n"
		response += "lastAccessTime: " + goutil.DateUtil().TimeFullFormat(table.LastAccessTime()) + "\r\n"
		response += "lastModifyTime: " + goutil.DateUtil().TimeFullFormat(table.LastModifyTime()) + "\r\n"
		response += "accessCount: " + strconv.FormatInt(table.AccessCount(), 10)
	}
	return response
}

//显示项信息
func HandleShowiCommnd(body string, client Client) string {
	response := ""
	table, _ := core.Cache(client.table)
	if len(body) == 0 { //没有请求体，则显示所有项
		index := 1
		for k, _ := range table.GetItems() {
			response += "[" + strconv.Itoa(index) + "] " + k.(string) + "\r\n"
			index += 1
		}
		response += "use 'showi key' to see detail info"
	} else {
		item := table.Get(body)
		if item == nil {
			return ERROR_ITEM_NOT_EXIST.Error()
		}

		//fmt.Fprintln(response, "key: ", item.Key())
		//fmt.Fprintln(response, "value: ", item.Value())
		response += "key: " + item.Key().(string) + "\r\n"
		response += "value: " + item.Value().(string) + "\r\n"
		response += "liveTime: " + item.LiveTime().String() + "\r\n"
		response += "createTime: " + goutil.DateUtil().TimeFullFormat(item.CreateTime()) + "\r\n"
		response += "lastAccessTime: " + goutil.DateUtil().TimeFullFormat(item.LastAccessTime()) + "\r\n"
		response += "lastModifyTime: " + goutil.DateUtil().TimeFullFormat(item.LastModifyTime()) + "\r\n"
		response += "accessCount: " + strconv.FormatInt(item.AccessCount(), 10) + "\r\n"
	}
	return response
}

//服务器信息
func HandleInfoCommnd(body string, client Client) string {
	info, _ := conf.GetSystemConfig().GetSection("")
	response := ""
	for k, v := range info {
		response += k + ": " + v + "\r\n"
	}
	return response
}
