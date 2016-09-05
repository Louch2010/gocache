package server

import (
	"log"
	"strings"

	"github.com/louch2010/gocache/conf"
	"github.com/louch2010/gocache/core"
	"github.com/louch2010/goutil"
)

//解析请求
func ParserRequest(request string, token string) (string, bool) {
	//去除换行、空格
	request = goutil.StringUtil().TrimToEmpty(request)
	log.Println("开始处理请求，token：", token, "，请求内容为：", request)
	//请求内容为空时，不处理
	if goutil.StringUtil().IsEmpty(request) {
		return "", false
	}
	arr := strings.SplitN(request, " ", 2)
	head := strings.ToUpper(goutil.StringUtil().TrimToEmpty(arr[0])) //请求头
	body := ""                                                       //请求体
	if len(arr) == 2 {
		body = goutil.StringUtil().TrimToEmpty(arr[1])
	}
	log.Println("请求头为：", head, "，请求体为：", body)
	response := "" //响应内容
	clo := false   //是否关闭
	//获取会话信息
	client, isLogin := GetSession(token)
	//未登录则强制先进行登录
	if !isLogin && !IsAnonymCommnd(head) {
		return ERROR_COMMND_NO_LOGIN.Error(), clo
	}
	log.Println("会话信息：", client)
	//解析
	switch head {
	//心跳检测
	case REQUEST_TYPE_PING:
		response = MESSAGE_PONG
		break
	//查看帮助
	case REQUEST_TYPE_HELP:
		response = HandleHelpCommnd(body, client)
		break
	//退出
	case REQUEST_TYPE_EXIT:
		response = MESSAGE_EXIT
		clo = true
		DestroySession(token)
		log.Println("客户端主动退出，请求处理完毕")
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
	//命令不正确
	default:
		response = ERROR_COMMND_NOT_FOUND.Error()
	}
	return response, clo
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
