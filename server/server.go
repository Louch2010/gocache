package server

import (
	"bufio"
	"encoding/json"
	"io"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/log"
	"github.com/louch2010/goutil"
)

//服务器运行状态标识
var ServerStatusFlag bool = false

//启动服务
func Start(port int, timeout int) error {
	if ServerStatusFlag {
		log.Error("服务已经在运行，无需再次启动")
		return ERROR_SERVER_ALREADY_START
	}
	ServerStatusFlag = true
	log.Info("启动服务，端口号：", port, "，连接超时时间：", timeout)
	//定义端口地址
	addr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Error("TCP地址初始化失败！", err)
		return err
	}
	//启动端口侦听
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error("启动端口侦听失败！", err)
		return err
	}
	for ServerStatusFlag {
		//当接收到了请求，则返回一个conn
		conn, err := listener.Accept()
		//生成唯一token
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		token := goutil.DateUtil().Time14Now() + strconv.Itoa(r.Int())
		log.Info("接收到请求，token：", token)
		if err != nil {
			log.Error("接收请求时出错！", err)
			continue
		}
		//使用长连接方式处理
		go handleLongConn(conn, timeout, token)
	}
	log.Info("服务已停止！")
	return nil
}

//停止服务
func Stop() {
	log.Info("停止服务...")
	ServerStatusFlag = false
}

//短连接处理
func handleShortConn(conn net.Conn, timeout int, token string) {
	log.Debug("开始处理短连接请求...")
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	//将请求读入缓存，并读取其中的一行
	buff := bufio.NewReader(conn)
	line, _ := buff.ReadString('\n')
	//解析请求并响应
	response := ParserRequest(line, token, Client{})
	conn.Write([]byte(toString(response.Data)))
	log.Debug("请求处理完成，响应状态为：", response.Code, "响应内容为：", response.Data)
	conn.Close()
}

//长连接处理
func handleLongConn(conn net.Conn, timeout int, token string) {
	log.Debug("开始处理长连接请求...")
	//客户端信息
	client := Client{}
	for {
		//将请求内容写入buff
		buff := bufio.NewReader(conn)
		//只读取一行内容
		line, err := buff.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Info("连接已关闭！")
			} else {
				log.Error("读取连接内容失败！", err)
			}
			conn.Close()
			return
		}
		manage := make(chan string)
		//心跳计时
		go heartBeating(conn, manage, timeout)
		//检测每次Client是否有数据传来
		go gravelChannel(line, manage)
		//解析请求并响应
		response := ParserRequest(line, token, client)
		data := ""
		if response.Err != nil {
			data = response.Err.Error()
		} else {
			data = toString(response.Data)
			if len(data) == 0 {
				data = response.Code
			}
			//将client进行缓存
			if response.Client != nil {
				client = *response.Client
			}
		}
		if response.Clo {
			io.WriteString(conn, data)
			conn.Close()
		} else {
			io.WriteString(conn, data+"\r\n -> ")
		}
		log.Debug("请求处理完成，响应状态为：", response.Code, "响应内容为：", data)
	}
}

//心跳计时，根据gravelChannel判断Client是否在设定时间内发来信息
func heartBeating(conn net.Conn, manage chan string, timeout int) {
	select {
	case <-manage:
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	case <-time.After(time.Second * time.Duration(timeout)):
		log.Info("客户端连接超时，自动关闭连接")
		conn.Close()
	}
}

func gravelChannel(content string, manage chan string) {
	manage <- content
	close(manage)
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
