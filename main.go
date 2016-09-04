package main

import (
	"log"
	"os"

	"github.com/louch2010/gocache/conf"
	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/server"
)

func main() {
	args := os.Args
	log.Println("启动参数：", args)
	//初始化配置文件
	configPath := ""
	if len(args) > 1 {
		configPath = args[1]
	}
	err := conf.InitConfig(configPath)
	if err != nil {
		log.Println("初始化配置文件失败！", err)
		return
	}
	//初始化缓存
	err = core.InitCache()
	if err != nil {
		log.Println("初始化缓存失败！", err)
		return
	}
	//启动服务
	port := conf.GetSystemConfig().MustInt("server", "port", 1334)
	aliveTime := conf.GetSystemConfig().MustInt("server", "aliveTime", 30)
	err = server.Start(port, aliveTime)
	if err != nil {
		log.Println("启动服务失败！", err)
		return
	}
}
