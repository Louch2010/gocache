package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/log"
)

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
			return GetServerRespMsg(MESSAGE_COMMND_PARAM_ERROR, "", ERROR_COMMND_PARAM_ERROR, &client)
		}
		liveTime = t
	}
	//增加缓存项
	table, err := core.Cache(client.table)
	if err != nil {
		log.Error("执行Set命令出错，获取表失败，表名：", client.table)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, &client)
	}
	item := table.Set(args[0], args[1], time.Duration(liveTime)*time.Second)
	return GetServerRespMsg(MESSAGE_SUCCESS, item, nil, &client)
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
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, &client)
	}
	item := table.Get(body)
	//不存在，则设置
	if item == nil {
		return GetServerRespMsg(MESSAGE_ITEM_NOT_EXIST, "", ERROR_ITEM_NOT_EXIST, &client)
	}
	return GetServerRespMsg(MESSAGE_SUCCESS, item.Value(), nil, &client)
}

//Append命令处理
func HandleAppendCommnd(body string, client Client) ServerRespMsg {
	//请求体校验
	resp, check := checkBody(body, 2, 2)
	if !check {
		return resp
	}
	table, err := core.Cache(client.table)
	if err != nil {
		log.Error("执行Get命令出错，获取表失败，表名：", client.table)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, &client)
	}
	args := strings.Split(body, " ")
	item := table.Get(args[0])
	//不存在，则设置
	if item == nil {
		table.Set(args[0], args[1], 0)
		return GetServerRespMsg(MESSAGE_ITEM_NOT_EXIST, "", ERROR_ITEM_NOT_EXIST, &client)
	}
	v := item.Value().(string) + args[1]
	item.SetValue(v)
	return GetServerRespMsg(MESSAGE_SUCCESS, v, nil, &client)
}

//StrLen命令处理
func HandleStrLenCommnd(body string, client Client) ServerRespMsg {
	//请求体校验
	resp, check := checkBody(body, 1, 1)
	if !check {
		return resp
	}
	table, err := core.Cache(client.table)
	if err != nil {
		log.Error("执行Get命令出错，获取表失败，表名：", client.table)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, &client)
	}
	item := table.Get(body)
	length := 0
	if item != nil {
		length = len(item.Value().(string))
	}
	response := GetServerRespMsg(MESSAGE_SUCCESS, length, nil, &client)
	response.DataType = DATA_TYPE_NUMBER
	return response
}

//SetNx命令处理
func HandleSetNxCommnd(body string, client Client) ServerRespMsg {
	//请求体校验
	resp, check := checkBody(body, 2, 2)
	if !check {
		return resp
	}
	table, err := core.Cache(client.table)
	if err != nil {
		log.Error("执行Get命令出错，获取表失败，表名：", client.table)
		return GetServerRespMsg(MESSAGE_ERROR, "", ERROR_SYSTEM, &client)
	}
	item := table.Get(body)
	flag := false
	//不存在，则设置
	if item == nil {
		flag = true
		args := strings.Split(body, " ")
		table.Set(args[0], args[1], 0)
	}
	response := GetServerRespMsg(MESSAGE_SUCCESS, flag, nil, &client)
	response.DataType = DATA_TYPE_BOOL
	return response
}
