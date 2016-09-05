package server

//请求类型
const (
	REQUEST_TYPE_PING    = "PING"    //心跳检测
	REQUEST_TYPE_CONNECT = "CONNECT" //连接
	REQUEST_TYPE_EXIT    = "EXIT"    //断开连接
	REQUEST_TYPE_SET     = "SET"     //添加
	REQUEST_TYPE_GET     = "GET"     //获取
	REQUEST_TYPE_DELETE  = "DELETE"  //删除
	REQUEST_TYPE_EXIST   = "EXIST"   //存在
	REQUEST_TYPE_EVENT   = "EVENT"   //事件
	REQUEST_TYPE_USE     = "USE"     //切换表
	REQUEST_TYPE_SHOWT   = "SHOWT"   //显示表信息
	REQUEST_TYPE_SHOWI   = "SHOWI"   //显示项信息
	REQUEST_TYPE_INFO    = "INFO"    //显示系统信息
	REQUEST_TYPE_HELP    = "HELP"    //帮助
)

//协议类型
const (
	PROTOCOL_DEFAULT  = "TERMINAL"
	PROTOCOL_JSON     = "JSON"
	PROTOCOL_TERMINAL = "TERMINAL"
)

//连接
type Connect struct {
	host        string   //地址
	port        int      //端口
	table       string   //表名
	listenEvent []string //侦听事件
	protocol    string   //通讯协议
}

//客户端
type Client struct {
	host        string   //地址
	port        int      //端口
	table       string   //表名
	listenEvent []string //侦听事件
	protocol    string   //通讯协议
	token       string   //令牌
}

//新增
type Add struct {
	key      string //键
	value    string //值
	liveTime int64  //存活时间
}

//获取
type Get struct {
	key string //键
}

//删除
type Delete struct {
	key string //键
}

//存在
type Exist struct {
	key string //键
}

//表
type Table struct {
	name           string //表名
	itemsCount     int64  //缓存项数量
	createTime     string //创建时间
	lastAccessTime string //最后访问时间
	lastModifyTime string //最后修改时间
	accessCount    int64  //访问次数
}

//项
type Item struct {
	key            string //键
	value          string //值
	liveTime       int64  //存活时间
	createTime     string //创建时间
	lastAccessTime string //最后访问时间
	lastModifyTime string //最后修改时间
	accessCount    int64  //访问次数
}

//事件
type Event struct {
	eventType string //事件类型
	table     Table  //表
	item      Item   //项
}
