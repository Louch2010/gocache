package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/log"
)

//NSet命令处理
func HandleNSetCommnd(body string, client Client) ServerRespMsg {
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

//NGet命令处理
func HandleNGetCommnd(body string, client Client) ServerRespMsg {
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
	response := GetServerRespMsg(MESSAGE_SUCCESS, item.Value(), nil, &client)
	response.DataType = DATA_TYPE_NUMBER
	return response
}

//Incr命令处理
func HandleIncrCommnd(body string, client Client) ServerRespMsg {
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
	var v float64 = 1
	//不存在，则设置为0，存在增加1
	if item == nil {
		table.Set(body, v, 0)
	} else {
		o, _ := item.Value().(float64)
		v = o + v
		item.SetValue(v)
	}
	response := GetServerRespMsg(MESSAGE_SUCCESS, v, nil, &client)
	response.DataType = DATA_TYPE_NUMBER
	return response
}

//IncrBy命令处理
func HandleIncrByCommnd(body string, client Client) ServerRespMsg {
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
	v, _ := strconv.ParseFloat(args[1], 10)
	//不存在，则设置为0，存在增加1
	if item == nil {
		table.Set(args[0], v, 0)
	} else {
		o, _ := item.Value().(float64)
		v = o + v
		item.SetValue(v)
	}
	response := GetServerRespMsg(MESSAGE_SUCCESS, v, nil, &client)
	response.DataType = DATA_TYPE_NUMBER
	return response
}
