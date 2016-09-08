package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/louch2010/gocache/conf"
	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/log"
	"github.com/louch2010/goutil"
)

//帮助命令处理
func HandleHelpCommnd(body string, client Client) ServerRespMsg {
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
	return GetServerRespMsg(MESSAGE_SUCCESS, response, nil, false, &client)
}

//连接命令处理connect [-t'table'] [-a'pwd'] [-i'ip'] [-p'port'] [-e'e1,e2...']
func HandleConnectCommnd(body, token string) ServerRespMsg {
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
			return GetServerRespMsg(MESSAGE_NO_PWD, "", ERROR_AUTHORITY_NO_PWD, false, nil)
		}
		if syspwd != pwd {
			return GetServerRespMsg(MESSAGE_PWD_ERROR, "", ERROR_AUTHORITY_PWD_ERROR, false, nil)
		}
	}
	//端口校验
	portInt := 0
	if len(port) > 0 {
		p, err := strconv.Atoi(port)
		if err != nil {
			log.Info("端口转换错误，", err)
			return GetServerRespMsg(MESSAGE_PORT_ERROR, "", ERROR_PORT_ERROR, false, nil)
		}
		portInt = p
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
	return GetServerRespMsg(MESSAGE_SUCCESS, token, nil, false, &client)
}

//Set命令处理
func HandleSetCommnd(body string, client Client) ServerRespMsg {
	//请求体校验
	resp, check := checkBody(body, 2, 3)
	if !check {
		return resp
	}
	args := strings.Split(body, " ")
	var liveTime int = 0
	if len(args) == 3 {
		t, err := strconv.Atoi(args[2])
		if err != nil {
			log.Info("参数转换错误，liveTime：", args[2], err)
			return GetServerRespMsg(MESSAGE_COMMND_PARAM_ERROR, "", ERROR_COMMND_PARAM_ERROR, false, &client)
		}
		liveTime = t
	}
	//增加缓存项
	table, err := core.Cache(client.table)
	if err != nil {
		log.Error("执行Set命令出错，获取表失败，表名：", client.table)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, false, &client)
	}
	item := table.Set(args[0], args[1], time.Duration(liveTime)*time.Second)
	return GetServerRespMsg(MESSAGE_SUCCESS, item, nil, false, &client)
}

//Get命令处理
func HandleGetCommnd(body string, client Client) ServerRespMsg {
	//请求体校验
	resp, check := checkBody(body, 1, 1)
	if !check {
		return resp
	}
	table, err := core.Cache(client.table)
	if err != nil {
		log.Error("执行Get命令出错，获取表失败，表名：", client.table)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, false, &client)
	}
	item := table.Get(body)
	if item == nil {
		return GetServerRespMsg(MESSAGE_ITEM_NOT_EXIST, "", ERROR_ITEM_NOT_EXIST, false, &client)
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, item.Value(), nil, false, &client)
}

//Delete命令处理
func HandleDeleteCommnd(body string, client Client) ServerRespMsg {
	//请求体校验
	resp, check := checkBody(body, 1, 1)
	if !check {
		return resp
	}
	table, err := core.Cache(client.table)
	if err != nil {
		log.Error("执行Delete命令出错，获取表失败，表名：", client.table)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, false, &client)
	}
	if table.Delete(body) {
		return GetServerRespMsg(MESSAGE_SUCCESS, "", nil, false, &client)
	}
	return GetServerRespMsg(MESSAGE_ITEM_NOT_EXIST, "", ERROR_ITEM_NOT_EXIST, false, &client)
}

//Exist命令处理
func HandleExistCommnd(body string, client Client) ServerRespMsg {
	//请求体校验
	resp, check := checkBody(body, 1, 1)
	if !check {
		return resp
	}
	table, err := core.Cache(client.table)
	if err != nil {
		log.Error("执行Exist命令出错，获取表失败，表名：", client.table)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, false, &client)
	}
	if table.IsExist(body) {
		return GetServerRespMsg(MESSAGE_ITEM_NOT_EXIST, "", ERROR_ITEM_NOT_EXIST, false, &client)
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, "", nil, false, &client)
}

//切换表
func HandleUseCommnd(body string, client Client) ServerRespMsg {
	//请求体校验
	resp, check := checkBody(body, 1, 1)
	if !check {
		return resp
	}
	client.table = body
	if CreateSession(client.token, client) {
		return GetServerRespMsg(MESSAGE_SUCCESS, "", nil, false, &client)
	}
	return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, false, &client)
}

//显示表信息
func HandleShowtCommnd(body string, client Client) ServerRespMsg {
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
			return GetServerRespMsg(MESSAGE_TABLE_NOT_EXIST, response, ERROR_TABLE_NOT_EXIST, false, &client)
		}
		response += "name:" + table.Name() + "\r\n"
		response += "itemCount: " + strconv.Itoa(table.ItemCount()) + "\r\n"
		response += "createTime: " + goutil.DateUtil().TimeFullFormat(table.CreateTime()) + "\r\n"
		response += "lastAccessTime: " + goutil.DateUtil().TimeFullFormat(table.LastAccessTime()) + "\r\n"
		response += "lastModifyTime: " + goutil.DateUtil().TimeFullFormat(table.LastModifyTime()) + "\r\n"
		response += "accessCount: " + strconv.FormatInt(table.AccessCount(), 10)
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, response, nil, false, &client)
}

//显示项信息
func HandleShowiCommnd(body string, client Client) ServerRespMsg {
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
			return GetServerRespMsg(MESSAGE_ITEM_NOT_EXIST, "", ERROR_ITEM_NOT_EXIST, false, &client)
		}
		response += "key: " + item.Key().(string) + "\r\n"
		response += "value: " + item.Value().(string) + "\r\n"
		response += "liveTime: " + item.LiveTime().String() + "\r\n"
		response += "createTime: " + goutil.DateUtil().TimeFullFormat(item.CreateTime()) + "\r\n"
		response += "lastAccessTime: " + goutil.DateUtil().TimeFullFormat(item.LastAccessTime()) + "\r\n"
		response += "lastModifyTime: " + goutil.DateUtil().TimeFullFormat(item.LastModifyTime()) + "\r\n"
		response += "accessCount: " + strconv.FormatInt(item.AccessCount(), 10) + "\r\n"
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, response, nil, false, &client)
}

//服务器信息
func HandleInfoCommnd(body string, client Client) ServerRespMsg {
	info, _ := conf.GetSystemConfig().GetSection("")
	response := ""
	for k, v := range info {
		response += k + ": " + v + "\r\n"
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, response, nil, false, &client)
}

//请求体检查
func checkBody(body string, minBodyLen int, maxBodyLen int) (ServerRespMsg, bool) {
	if len(body) == 0 && minBodyLen != 0 {
		return GetServerRespMsg(MESSAGE_COMMND_PARAM_ERROR, "", ERROR_COMMND_PARAM_ERROR, false, nil), false
	}
	//参数处理
	args := strings.Split(body, " ")
	if len(args) < minBodyLen || len(args) > maxBodyLen {
		return GetServerRespMsg(MESSAGE_COMMND_PARAM_ERROR, "", ERROR_COMMND_PARAM_ERROR, false, nil), false
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, "", nil, false, nil), true
}
