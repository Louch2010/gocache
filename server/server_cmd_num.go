package server

import (
	"strconv"
	"time"

	"github.com/louch2010/gocache/log"
)

//NSet命令处理
func HandleNSetCommnd(body string, client *Client) ServerRespMsg {
	//请求体校验
	args, resp, check := initParam(body, 2, 3)
	if !check {
		return resp
	}
	var liveTime int = 0
	if len(args) == 3 {
		t, err := strconv.Atoi(args[2])
		if err != nil {
			log.Info("参数转换错误，liveTime：", args[2], err)
			return GetServerRespMsg(MESSAGE_COMMAND_PARAM_ERROR, "", ERROR_COMMAND_PARAM_ERROR, client)
		}
		liveTime = t
	}
	f, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		log.Info("参数类型错误，nset value：", args[0], err)
		return GetServerRespMsg(MESSAGE_COMMAND_NOT_SUPPORT_DATA, "", ERROR_COMMAND_NOT_SUPPORT_DATA, client)
	}
	//增加缓存项
	item := client.cacheTable.Set(args[0], f, time.Duration(liveTime)*time.Second, DATA_TYPE_NUMBER)
	return GetServerRespMsg(MESSAGE_SUCCESS, item, nil, client)
}

//NGet命令处理
func HandleNGetCommnd(body string, client *Client) ServerRespMsg {
	//请求体校验
	args, resp, check := initParam(body, 1, 1)
	if !check {
		return resp
	}
	item := client.cacheTable.Get(args[0])
	if item == nil {
		return GetServerRespMsg(MESSAGE_ITEM_NOT_EXIST, "", ERROR_ITEM_NOT_EXIST, client)
	}
	//数据类型校验
	if item.DataType() != DATA_TYPE_NUMBER {
		return GetServerRespMsg(MESSAGE_COMMAND_NOT_SUPPORT_DATA, "", ERROR_COMMAND_NOT_SUPPORT_DATA, client)
	}
	response := GetServerRespMsg(MESSAGE_SUCCESS, item.Value(), nil, client)
	response.DataType = DATA_TYPE_NUMBER
	return response
}

//Incr命令处理
func HandleIncrCommnd(body string, client *Client) ServerRespMsg {
	//请求体校验
	args, resp, check := initParam(body, 1, 1)
	if !check {
		return resp
	}
	item := client.cacheTable.Get(args[0])
	var v float64 = 1
	//不存在，则设置为0，存在增加1
	if item == nil {
		client.cacheTable.Set(args[0], v, 0, DATA_TYPE_NUMBER)
	} else {
		//数据类型校验
		if item.DataType() != DATA_TYPE_NUMBER {
			return GetServerRespMsg(MESSAGE_COMMAND_NOT_SUPPORT_DATA, "", ERROR_COMMAND_NOT_SUPPORT_DATA, client)
		}
		o, _ := item.Value().(float64)
		v = o + v
		item.SetValue(v)
	}
	response := GetServerRespMsg(MESSAGE_SUCCESS, v, nil, client)
	response.DataType = DATA_TYPE_NUMBER
	return response
}

//IncrBy命令处理
func HandleIncrByCommnd(body string, client *Client) ServerRespMsg {
	//请求体校验
	args, resp, check := initParam(body, 2, 2)
	if !check {
		return resp
	}
	item := client.cacheTable.Get(args[0])
	v, _ := strconv.ParseFloat(args[1], 10)
	//不存在，则设置为0，存在增加1
	if item == nil {
		client.cacheTable.Set(args[0], v, 0, DATA_TYPE_NUMBER)
	} else {
		//数据类型校验
		if item.DataType() != DATA_TYPE_NUMBER {
			return GetServerRespMsg(MESSAGE_COMMAND_NOT_SUPPORT_DATA, "", ERROR_COMMAND_NOT_SUPPORT_DATA, client)
		}
		o, _ := item.Value().(float64)
		v = o + v
		item.SetValue(v)
	}
	response := GetServerRespMsg(MESSAGE_SUCCESS, v, nil, client)
	response.DataType = DATA_TYPE_NUMBER
	return response
}
