package server

import (
	"strings"

	"github.com/louch2010/gocache/conf"
	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/log"
	"github.com/louch2010/goutil"
)

//解析请求
func ParserRequest(request string, token string, client Client) ServerRespMsg {
	//去除换行、空格
	request = goutil.StringUtil().TrimToEmpty(request)
	log.Debug("开始处理请求，token：", token, "，请求内容为：", request)
	//请求内容为空时，不处理
	if goutil.StringUtil().IsEmpty(request) {
		return GetServerRespMsg(MESSAGE_SUCCESS, "", nil, false, nil)
	}
	arr := strings.SplitN(request, " ", 2)
	head := strings.ToUpper(goutil.StringUtil().TrimToEmpty(arr[0])) //请求头
	body := ""                                                       //请求体
	if len(arr) == 2 {
		body = goutil.StringUtil().TrimToEmpty(arr[1])
	}
	log.Debug("请求头为：", head, "，请求体为：", body)
	var response ServerRespMsg //响应内容
	//会话信息校验
	openSession := conf.GetSystemConfig().MustBool("client", "openSession", true)
	isLogin := false
	if len(client.token) > 0 {
		isLogin = true
	} else {
		client, isLogin = GetSession(token)
	}
	//没有登录
	if !isLogin {
		//需要登录，而且也不是免登录的命令
		if openSession && !IsAnonymCommnd(head) {
			return GetServerRespMsg(MESSAGE_COMMND_NO_LOGIN, "", ERROR_COMMND_NO_LOGIN, false, nil)
		}
		//模拟登录
		if !openSession {
			table := conf.GetSystemConfig().MustValue("table", "default", core.DEFAULT_TABLE_NAME)
			client = Client{
				table: table,
				token: token,
			}
		}
	}
	log.Debug("会话信息：", client)
	//解析
	switch head {
	//心跳检测
	case REQUEST_TYPE_PING:
		response = GetServerRespMsg(MESSAGE_SUCCESS, MESSAGE_PONG, nil, false, &client)
		break
	//查看帮助
	case REQUEST_TYPE_HELP:
		response = HandleHelpCommnd(body, client)
		break
	//退出
	case REQUEST_TYPE_EXIT:
		response = GetServerRespMsg(MESSAGE_SUCCESS, "", nil, true, &client)
		DestroySession(token)
		log.Debug("客户端主动退出，请求处理完毕")
		break
	//打开连接
	case REQUEST_TYPE_CONNECT:
		response = HandleConnectCommnd(body, token)
		break
	//新增
	case REQUEST_TYPE_SET:
		response = HandleSetCommnd(body, client)
		break
	//获取
	case REQUEST_TYPE_GET:
		response = HandleGetCommnd(body, client)
		break
	//删除
	case REQUEST_TYPE_DELETE:
		response = HandleDeleteCommnd(body, client)
		break
	//是否存在
	case REQUEST_TYPE_EXIST:
		response = HandleExistCommnd(body, client)
		break
	//切换表
	case REQUEST_TYPE_USE:
		response = HandleUseCommnd(body, client)
		break
	//显示表信息
	case REQUEST_TYPE_SHOWT:
		response = HandleShowtCommnd(body, client)
		break
	//显示项信息
	case REQUEST_TYPE_SHOWI:
		response = HandleShowiCommnd(body, client)
		break
	//服务器信息
	case REQUEST_TYPE_INFO:
		response = HandleInfoCommnd(body, client)
		break
	//命令不正确
	default:
		response = GetServerRespMsg(MESSAGE_COMMND_NOT_FOUND, "", ERROR_COMMND_NOT_FOUND, false, &client)
	}
	return response
}

//创建会话
func CreateSession(token string, c Client) bool {
	//缓存登录信息
	table, _ := core.GetSysTable()
	table.Set(token, c, 0)
	//创建表信息
	core.Cache(c.table)
	return true
}

//获取会话
func GetSession(token string) (Client, bool) {
	table, _ := core.GetSysTable()
	item := table.Get(token)
	if item == nil {
		return Client{}, false
	}
	value, falg := item.Value().(Client)
	return value, falg
}

//销毁会话
func DestroySession(token string) bool {
	table, _ := core.GetSysTable()
	return table.Delete(token)
}

//判断是否为免登录命令
func IsAnonymCommnd(commnd string) bool {
	anonymCommnd := conf.SystemConfigFile.MustValue("server", "anonymCommnd", "ping,connect,exit,help")
	list := strings.Split(strings.ToUpper(anonymCommnd), ",")
	for _, c := range list {
		if commnd == c {
			return true
		}
	}
	return false
}
