package server

import (
	"log"
	"strconv"
	"strings"

	"github.com/louch2010/gocache/conf"
	"github.com/louch2010/goutil"
)

//解析请求
func ParserRequest(request string) (string, bool) {
	//去除换行、空格
	request = goutil.StringUtil().TrimToEmpty(request)
	log.Println("开始处理请求，请求内容为：", request)
	//请求内容为空时，不处理
	if goutil.StringUtil().IsEmpty(request) {
		return "", false
	}
	arr := strings.SplitN(request, " ", 2)
	head := goutil.StringUtil().TrimToEmpty(arr[0]) //请求头
	body := ""                                      //请求体
	if len(arr) == 2 {
		body = goutil.StringUtil().TrimToEmpty(arr[1])
	}
	log.Println("请求头为：", head, "，请求体为：", body)
	//解析
	response := "" //响应内容
	clo := false   //是否关闭
	switch strings.ToUpper(head) {
	//心跳检测
	case REQUEST_TYPE_PING:
		response = "PONG"
		break
	//查看帮助
	case REQUEST_TYPE_HELP:
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
					response += k + ":" + v + "\r\n"
				}
			}
		}
		break
	//退出
	case REQUEST_TYPE_EXIT:
		response = "Bye"
		clo = true
		log.Println("客户端主动退出，请求处理完毕")
		break
	//打开连接
	case REQUEST_TYPE_CONNECT:
		response = "connect suc"
		break
	//新增
	case REQUEST_TYPE_SET:
		response = "set"
		break
	//获取
	case REQUEST_TYPE_GET:
		response = "get"
		break
	//删除
	case REQUEST_TYPE_DELETE:
		response = "delete"
		break
	//是否存在
	case REQUEST_TYPE_EXIST:
		response = "exist"
		break
	//命令不正确
	default:
		response = "commnd not found"
	}

	return response, clo
}
