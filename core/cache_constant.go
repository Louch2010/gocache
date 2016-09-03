package core

import (
	"errors"
)

//事件
const (
	//缓存项新增
	EVENT_ITEM_ADD = "EVENT_ITEM_ADD"
	//缓存项修改
	EVENT_ITEM_MODIFY = "EVENT_ITEM_MODIFY"
	//缓存项删除
	EVENT_ITEM_DELETE = "EVENT_ITEM_DELETE"
	//缓存表新增
	EVENT_TABLE_ADD = "EVENT_TABLE_ADD"
	//缓存表删除
	EVENT_TABLE_DELETE = "EVENT_TABLE_DELETE"
)

//默认常量
const (
	//默认表名
	DEFAULT_TABLE_NAME = "default"
)

//异常
var (
	//缓存表名错误
	ERROR_TABLE_NAME = errors.New("ERROR_TABLE_NAME")
)
