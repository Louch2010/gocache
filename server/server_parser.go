package server

import (
	"log"
	"strings"

	"github.com/louch2010/goutil"
)

//解析请求
func ParserRequest(request string) string {
	//去除换行、空格
	request = goutil.StringUtil().TrimToEmpty(request)
	log.Println("开始处理请求，请求内容为：", request)
	//请求内容为空时，不处理
	if goutil.StringUtil().IsEmpty(request) {
		return ""
	}
	//
	arr := strings.SplitN(request, " ", 2)
	head := arr[0] //请求头
	body := ""     //请求体
	if len(arr) == 2 {
		body = arr[1]
	}
	response := "" //响应内容

	log.Println("请求头为：", head, "，请求体为：", body)

	//解析
	switch strings.ToUpper(head) {
	case REQUEST_TYPE_PING:
		response = "PONG"
		break
	case REQUEST_TYPE_CONNECT:
		response = "connect suc"
		break
	default:
		response = "commond not found"
	}

	return response
}
