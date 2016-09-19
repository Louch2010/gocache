package gdb

const (
	GOCACHE         = "GOCACHE"
	VERSION         = "0001"
	DATABASE        = "DATABASE"
	EOF             = "EOF"
	CHECK_SUM       = "CHECK_SUM"
	LIVETIME_ALWAYS = "A"
)

const (
	LEN_GOCACHE         = 7
	LEN_VERSION         = 4
	LEN_DATABASE        = 8
	LEN_EOF             = 3
	LEN_CHECK_SUM       = 32
	LEN_KEY             = 3
	LEN_VALUE           = 6
	LEN_LIVETIME        = 14
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
