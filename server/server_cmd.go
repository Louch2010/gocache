package server

import (
	"container/list"
	"strconv"
	"strings"

	"github.com/louch2010/gocache/conf"
	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/log"
	"github.com/louch2010/goutil"
)

//帮助命令处理
func HandleHelpCommnd(body string, client *Client) ServerRespMsg {
	response := ""
	help := conf.GetHelpConfig()
	args, resp, check := initParam(body, 0, 1)
	if !check {
		return resp
	}
	if len(args) == 0 { //没有请求体，则显示所有命令名称
		for index, sec := range help.GetSectionList() {
			response += "[" + strconv.Itoa(index+1) + "] " + sec + "\r\n"
		}
		response += "use 'help commnd' to see detail info"
	} else {
		cmd := strings.ToLower(args[0])
		sec, err := help.GetSection(cmd)
		if err != nil {
			response = "no help for the commnd"
		} else {
			response += "[" + cmd + "]\r\n"
			for k, v := range sec {
				response += k + ": " + v + "\r\n"
			}
		}
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, response, nil, client)
}

//连接命令处理connect [-t'table'] [-a'pwd'] [-i'ip'] [-p'port'] [-e'e1,e2...']
func HandleConnectCommnd(body, token string) ServerRespMsg {
	table := conf.GetSystemConfig().MustValue("table", "default", core.DEFAULT_TABLE_NAME)
	protocol := PROTOCOL_RESPONSE_DEFAULT
	var pwd, ip, port, event string
	args, resp, check := initParam(body, 0, 5)
	if !check {
		return resp
	}
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
		case 'm':
			protocol = param
			break
		default:
		}
	}
	//密码校验
	syspwd := conf.GetSystemConfig().MustValue("server", "password", "")
	if len(syspwd) > 0 {
		if len(pwd) == 0 {
			return GetServerRespMsg(MESSAGE_NO_PWD, "", ERROR_AUTHORITY_NO_PWD, nil)
		}
		if syspwd != pwd {
			return GetServerRespMsg(MESSAGE_PWD_ERROR, "", ERROR_AUTHORITY_PWD_ERROR, nil)
		}
	}
	//端口校验
	portInt := 0
	if len(port) > 0 {
		p, err := strconv.Atoi(port)
		if err != nil {
			log.Info("端口转换错误，", err)
			return GetServerRespMsg(MESSAGE_PORT_ERROR, "", ERROR_PORT_ERROR, nil)
		}
		portInt = p
	}
	//协议校验
	if len(protocol) > 0 && protocol != PROTOCOL_RESPONSE_JSON && protocol != PROTOCOL_RESPONSE_TERMINAL {
		log.Info("协议错误：", protocol)
		return GetServerRespMsg(MESSAGE_PROTOCOL_ERROR, "", ERROR_PROTOCOL_ERROR, nil)
	}
	//获取表
	cacheTable, err := core.Cache(table)
	if err != nil {
		log.Error("连接时获取表失败！", err)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, nil)
	}
	//存储连接信息
	client := &Client{
		host:        ip,
		port:        portInt,
		table:       table,
		cacheTable:  cacheTable,
		listenEvent: strings.Split(event, ","),
		protocol:    protocol,
		token:       token,
	}
	CreateSession(token, client)
	return GetServerRespMsg(MESSAGE_SUCCESS, token, nil, client)
}

//Delete命令处理
func HandleDeleteCommnd(body string, client *Client) ServerRespMsg {
	//请求体校验
	args, resp, check := initParam(body, 1, 1)
	if !check {
		return resp
	}
	if client.cacheTable.Delete(args[0]) {
		return GetServerRespMsg(MESSAGE_SUCCESS, "", nil, client)
	}
	return GetServerRespMsg(MESSAGE_ITEM_NOT_EXIST, "", ERROR_ITEM_NOT_EXIST, client)
}

//Exist命令处理
func HandleExistCommnd(body string, client *Client) ServerRespMsg {
	//请求体校验
	args, resp, check := initParam(body, 1, 1)
	if !check {
		return resp
	}
	response := GetServerRespMsg(MESSAGE_SUCCESS, client.cacheTable.IsExist(args[0]), nil, client)
	response.DataType = DATA_TYPE_BOOL
	return response
}

//切换表
func HandleUseCommnd(body string, client *Client) ServerRespMsg {
	//请求体校验
	args, resp, check := initParam(body, 1, 1)
	if !check {
		return resp
	}
	cacheTable, err := core.Cache(args[0])
	if err != nil {
		log.Error("切换表时获取表失败！", err)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, nil)
	}
	client.table = args[0]
	client.cacheTable = cacheTable
	if CreateSession(client.token, client) {
		return GetServerRespMsg(MESSAGE_SUCCESS, "", nil, client)
	}
	return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, client)
}

//显示表信息
func HandleShowtCommnd(body string, client *Client) ServerRespMsg {
	response := ""
	args, resp, check := initParam(body, 0, 1)
	if !check {
		return resp
	}
	if len(args) == 0 { //没有请求体，则显示所有表名
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
		table, ok := core.GetCacheTable(args[0])
		if !ok {
			return GetServerRespMsg(MESSAGE_TABLE_NOT_EXIST, response, ERROR_TABLE_NOT_EXIST, client)
		}
		response += "name:" + table.Name() + "\r\n"
		response += "itemCount: " + strconv.Itoa(table.ItemCount()) + "\r\n"
		response += "createTime: " + goutil.DateUtil().TimeFullFormat(table.CreateTime()) + "\r\n"
		response += "lastAccessTime: " + goutil.DateUtil().TimeFullFormat(table.LastAccessTime()) + "\r\n"
		response += "lastModifyTime: " + goutil.DateUtil().TimeFullFormat(table.LastModifyTime()) + "\r\n"
		response += "accessCount: " + strconv.FormatInt(table.AccessCount(), 10)
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, response, nil, client)
}

//显示项信息
func HandleShowiCommnd(body string, client *Client) ServerRespMsg {
	response := ""
	table, _ := core.Cache(client.table)
	args, resp, check := initParam(body, 0, 1)
	if !check {
		return resp
	}
	if len(args) == 0 { //没有请求体，则显示所有项
		index := 1
		for k, _ := range table.GetItems() {
			response += "[" + strconv.Itoa(index) + "] " + k + "\r\n"
			index += 1
		}
		response += "use 'showi key' to see detail info"
	} else {
		item := table.Get(args[0])
		if item == nil {
			return GetServerRespMsg(MESSAGE_ITEM_NOT_EXIST, "", ERROR_ITEM_NOT_EXIST, client)
		}
		response += "key: " + item.Key() + "\r\n"
		response += "value: " + toString(item) + "\r\n"
		response += "liveTime: " + item.LiveTime().String() + "\r\n"
		response += "createTime: " + goutil.DateUtil().TimeFullFormat(item.CreateTime()) + "\r\n"
		response += "lastAccessTime: " + goutil.DateUtil().TimeFullFormat(item.LastAccessTime()) + "\r\n"
		response += "accessCount: " + strconv.FormatInt(item.AccessCount(), 10) + "\r\n"
		response += "dataType: " + item.DataType() + "\r\n"
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, response, nil, client)
}

//服务器信息
func HandleInfoCommnd(body string, client *Client) ServerRespMsg {
	_, resp, check := initParam(body, 0, 0)
	if !check {
		return resp
	}
	info, _ := conf.GetSystemConfig().GetSection("")
	response := ""
	for k, v := range info {
		response += k + ": " + v + "\r\n"
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, response, nil, client)
}

//初始化参数
func initParam(body string, minBodyLen int, maxBodyLen int) ([]string, ServerRespMsg, bool) {
	result := make([]string, 0)
	if len(body) == 0 {
		if minBodyLen != 0 {
			return nil, GetServerRespMsg(MESSAGE_COMMAND_PARAM_ERROR, "", ERROR_COMMAND_PARAM_ERROR, nil), false
		}
		return result, GetServerRespMsg(MESSAGE_SUCCESS, "", nil, nil), true
	}
	//如果包含引号，则需要特殊处理
	if strings.Contains(body, "\"") {
		l := list.New()
		open := false
		buffer := ""
		for _, c := range body {
			if '"' == c {
				if open {
					l.PushBack(buffer)
					buffer = ""
				}
				open = !open
				continue
			}
			if ' ' == c && !open {
				if len(buffer) > 0 {
					l.PushBack(buffer)
					buffer = ""
				}
				continue
			}
			buffer += string(c)
		}
		result = make([]string, l.Len())
		var i = 0
		for e := l.Front(); e != nil; e = e.Next() {
			result[i] = e.Value.(string)
			i = i + 1
		}
	} else {
		body = strings.Replace(body, "  ", " ", 99)
		result = strings.Split(body, " ")
	}
	log.Debug("初始化请求参数完成，请求参数为：", result, "，长度为：", len(result))
	if len(result) < minBodyLen || len(result) > maxBodyLen {
		return nil, GetServerRespMsg(MESSAGE_COMMAND_PARAM_ERROR, "", ERROR_COMMAND_PARAM_ERROR, nil), false
	}
	return result, GetServerRespMsg(MESSAGE_SUCCESS, "", nil, nil), true
}
