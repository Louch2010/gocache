package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

//服务器运行状态标识
var ServerStatusFlag bool = false

//启动服务
func Start(port int, timeout int) error {
	if ServerStatusFlag {
		log.Println("服务已经在运行，无需再次启动")
		return ERROR_SERVER_ALREADY_START
	}
	ServerStatusFlag = true
	log.Println("启动服务，端口号：", port, "，连接超时时间：", timeout)
	//定义端口地址
	host := ":" + strconv.Itoa(port)
	addr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		log.Println("TCP地址初始化失败！", err)
		return err
	}
	//启动端口侦听
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println("启动端口侦听失败！", err)
		return err
	}
	for ServerStatusFlag {
		//当接收到了请求，则返回一个conn
		conn, err := listener.Accept()
		log.Println("接收到请求...")
		if err != nil {
			log.Println("接收请求时出错！", err)
		}
		//使用长连接方式处理
		go handleLongConn(conn, timeout)
	}
	log.Println("服务停止完成！")
	return nil
}

//停止服务
func Stop() {
	log.Println("停止服务...")
	ServerStatusFlag = false
}

//短连接处理
func handleShortConn(conn net.Conn, timeout int) {
	log.Println("开始处理请求...")
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	//将请求读入缓存，并读取其中的一行
	buff := bufio.NewReader(conn)
	line, _ := buff.ReadString('\n')
	//解析请求并响应
	response, _ := ParserRequest(line)
	conn.Write([]byte(response))
	log.Println("请求处理完成，响应内容为：", response)
	conn.Close()
}

//长连接处理
func handleLongConn(conn net.Conn, timeout int) {
	for {
		//将请求内容写入buff
		buff := bufio.NewReader(conn)
		//只读取一行内容
		line, err := buff.ReadString('\n')
		if err != nil {
			log.Println("读取连接内容失败！", err)
			conn.Close()
			return
		}
		messnager := make(chan string)
		//心跳计时
		go heartBeating(conn, messnager, timeout)
		//检测每次Client是否有数据传来
		go gravelChannel(line, messnager)
		//解析请求并响应
		response, clo := ParserRequest(line)
		if clo {
			io.WriteString(conn, response)
			conn.Close()
		} else {
			io.WriteString(conn, response+"\r\n -> ")
		}
		log.Println("请求处理完成，响应内容为：", response)
	}
}

//心跳计时，根据gravelChannel判断Client是否在设定时间内发来信息
func heartBeating(conn net.Conn, readerChannel chan string, timeout int) {
	select {
	case <-readerChannel:
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	case <-time.After(time.Second * time.Duration(timeout)):
		log.Println("客户端连接超时，自动关闭连接")
		conn.Close()
	}
}

func gravelChannel(content string, mess chan string) {
	mess <- content
	close(mess)
}
