package server

import (
	"encoding/json"
	"strconv"

	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/log"
)

//请求命令
const (
	REQUEST_TYPE_PING    = "PING"    //心跳检测
	REQUEST_TYPE_CONNECT = "CONNECT" //连接
	REQUEST_TYPE_EXIT    = "EXIT"    //断开连接
	REQUEST_TYPE_DELETE  = "DELETE"  //删除
	REQUEST_TYPE_EXIST   = "EXIST"   //存在
	REQUEST_TYPE_EVENT   = "EVENT"   //事件
	REQUEST_TYPE_USE     = "USE"     //切换表
	REQUEST_TYPE_SHOWT   = "SHOWT"   //显示表信息
	REQUEST_TYPE_SHOWI   = "SHOWI"   //显示项信息
	REQUEST_TYPE_INFO    = "INFO"    //显示系统信息
	REQUEST_TYPE_HELP    = "HELP"    //帮助

	REQUEST_TYPE_SET    = "SET"    //添加string
	REQUEST_TYPE_GET    = "GET"    //获取string
	REQUEST_TYPE_APPEND = "APPEND" //追加string
	REQUEST_TYPE_STRLEN = "STRLEN" //值的长度string
	REQUEST_TYPE_SETNX  = "SETNX"  //不存在则设置string

	REQUEST_TYPE_NSET   = "NSET"   //添加number
	REQUEST_TYPE_NGET   = "NGET"   //获取number
	REQUEST_TYPE_INCR   = "INCR"   //增加1 number
	REQUEST_TYPE_INCRBY = "INCRBY" //增加指定值 number

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

//协议类型
const (
	//默认使用终端
	PROTOCOL_RESPONSE_DEFAULT = "TERMINAL"
	//JSON，一般用于各语言客户端
	PROTOCOL_RESPONSE_JSON = "JSON"
	//终端，用于telnet等方式
	PROTOCOL_RESPONSE_TERMINAL = "TERMINAL"
)

//标识字符
const (
	FLAG_CHAR_SOCKET_COMMND_END            = '\n'
	FLAG_CHAR_SOCKET_TERMINAL_RESPONSE_END = "\r\n ->"
	FLAG_CHAR_SOCKET_JSON_RESPONSE_END     = "\r\n!--!>"
)

//客户端
type Client struct {
	host        string           //地址
	port        int              //端口
	table       string           //表名
	cacheTable  *core.CacheTable //表指针
	listenEvent []string         //侦听事件
	protocol    string           //通讯协议
	token       string           //令牌
}

//响应信息
type ServerRespMsg struct {
	Code     string      //响应码
	Data     interface{} //响应数据
	DataType string      //数据类型
	Clo      bool        //是否关闭连接
	Err      error       //错误信息
	Client   *Client     //客户端对象
}

//JSON响应
type JsonRespMsg struct {
	Code     string
	Msg      string
	Data     interface{}
	DataType string //数据类型
}

func GetServerRespMsg(code string, data interface{}, err error, client *Client) ServerRespMsg {
	resp := ServerRespMsg{
		Code:     code,
		Data:     data,
		DataType: DATA_TYPE_STRING,
		Err:      err,
		Clo:      false,
		Client:   client,
	}
	return resp
}

//根据连接协议，将响应内容进行封装
func TransferResponse(response ServerRespMsg) string {
	protocol := ""
	if response.Client != nil {
		protocol = response.Client.protocol
	}
	//终端方式：有错误，则输出错误信息，没有错误，则直接输出响应信息
	if protocol == "" || protocol == PROTOCOL_RESPONSE_TERMINAL {
		if response.Err != nil {
			return response.Err.Error() + FLAG_CHAR_SOCKET_TERMINAL_RESPONSE_END
		}
		return toString(response.Data) + FLAG_CHAR_SOCKET_TERMINAL_RESPONSE_END
	}
	//JSON方式：对响应信息进行json封装
	if protocol == PROTOCOL_RESPONSE_JSON {
		msg := MESSAGE_SUCCESS
		if response.Err != nil {
			msg = response.Err.Error()
		}
		obj := JsonRespMsg{
			Code:     response.Code,
			Msg:      msg,
			Data:     response.Data,
			DataType: response.DataType,
		}
		j, _ := json.Marshal(obj)
		return string(j) + FLAG_CHAR_SOCKET_JSON_RESPONSE_END
	}
	return ""
}

//转string
func toString(v interface{}) string {
	response := ""
	switch conv := v.(type) {
	case string:
		response = conv
		break
	case int:
		response = strconv.Itoa(conv)
		break
	case bool:
		response = strconv.FormatBool(conv)
		break
	case float64:
		response = strconv.FormatFloat(conv, 'E', -1, 64)
		break
	case *core.CacheItem:
		if conv != nil {
			tmp, _ := json.Marshal(conv.Value())
			response = string(tmp)
		}
		break
	default:
		log.Error("类型转换异常")
	}
	return response
}

//连接
//type Connect struct {
//	host        string   //地址
//	port        int      //端口
//	table       string   //表名
//	listenEvent []string //侦听事件
//	protocol    string   //通讯协议
//}

//新增
//type Add struct {
//	key      string //键
//	value    string //值
//	liveTime int64  //存活时间
//}

////获取
//type Get struct {
//	key string //键
//}

////删除
//type Delete struct {
//	key string //键
//}

////存在
//type Exist struct {
//	key string //键
//}

////表
//type Table struct {
//	name           string //表名
//	itemsCount     int64  //缓存项数量
//	createTime     string //创建时间
//	lastAccessTime string //最后访问时间
//	lastModifyTime string //最后修改时间
//	accessCount    int64  //访问次数
//}

////项
//type Item struct {
//	key            string //键
//	value          string //值
//	liveTime       int64  //存活时间
//	createTime     string //创建时间
//	lastAccessTime string //最后访问时间
//	lastModifyTime string //最后修改时间
//	accessCount    int64  //访问次数
//}

//事件
//type Event struct {
//	eventType string //事件类型
//	table     Table  //表
//	item      Item   //项
//}
