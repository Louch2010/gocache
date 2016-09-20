package gdb

import (
	"errors"
)

const (
	GOCACHE         = "GOCACHE"
	VERSION         = "0001"
	DATABASE        = "DATABASE"
	TABLE           = "TABLE"
	EOF             = "EOF"
	LIVETIME_ALWAYS = "A"
)

const (
	//魔数长度
	LEN_GOCACHE = 7
	//版本号长度
	LEN_VERSION = 4
	//库标识长度
	LEN_DATABASE = 8
	//表标识长度
	LEN_TABLE = 5
	//结束符长度
	LEN_EOF = 3
	//校验码长度
	LEN_CHECK_SUM = 32
	//库大小长度
	LEN_DATABASE_SIZE = 4
	//表大小长度
	LEN_TABLE_SIZE = 8
	//键长度，包括表名、键名
	LEN_KEY = 3
	//值长度
	LEN_VALUE = 6
	//数据类型长度
	LEN_DATATYPE = 1
	//存活时间长度
	LEN_LIVETIME = 14
	//永久存活时间长度
	LEN_LIVETIME_ALWAYS = 1
)

//数据类型
const (
	TYPE_STRING = "1" //字符
	TYPE_BOOL   = "2" //布尔
	TYPE_NUMBER = "3" //数字
	TYPE_MAP    = "4" //Map
	TYPE_SET    = "5" //集合
	TYPE_LIST   = "6" //列表
	TYPE_ZSET   = "7" //有序集合
	TYPE_OBJECT = "8" //对象
)

//数据类型
const (
	DATA_TYPE_STRING = "string" //字符
	DATA_TYPE_BOOL   = "bool"   //布尔
	DATA_TYPE_NUMBER = "number" //数字
	DATA_TYPE_MAP    = "map"    //Map
	DATA_TYPE_SET    = "set"    //集合
	DATA_TYPE_LIST   = "list"   //列表
	DATA_TYPE_ZSET   = "zset"   //有序集合
	DATA_TYPE_OBJECT = "object" //对象
)

var (
	GDB_FILE_INVALID       = errors.New("invalid gdb file")
	GDB_FILE_VERSION_ERROR = errors.New("gdb file version is not support")
	GDB_FILE_CHECK_ERROR   = errors.New("gdb file is broken")
	GDB_FILE_FORMAT_ERROR  = errors.New("error gdb file format")
)
