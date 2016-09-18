package main

import (
	default_log "log"
	"os"
	//"runtime"

	"github.com/louch2010/gocache/conf"
	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/gdb"
	"github.com/louch2010/gocache/log"
	"github.com/louch2010/gocache/server"
)

func main() {
	//设置cpu数量
	//runtime.GOMAXPROCS(runtime.NumCPU())
	args := os.Args
	default_log.Println("启动参数：", args)
	//初始化配置文件
	configPath := ""
	if len(args) > 1 {
		configPath = args[1]
	}
	err := conf.InitConfig(configPath)
	if err != nil {
		default_log.Println("初始化配置文件失败！", err)
		return
	}
	//初始化日志
	level := conf.GetSystemConfig().MustValue("log", "level", "info")
	format := conf.GetSystemConfig().MustValue("log", "format", "%Date/%Time [%LEV] %Msg%n")
	path := conf.GetSystemConfig().MustValue("log", "path", "./log/gocache.log")
	roll := conf.GetSystemConfig().MustValue("log", "roll", "02.01.2006")
	consoleOn := conf.GetSystemConfig().MustBool("log", "console.on", true)
	err = log.InitLog(level, format, path, roll, consoleOn)
	if err != nil {
		default_log.Println("初始化日志失败！", err)
		return
	}
	//初始化缓存
	log.Info("初始化缓存...")
	err = core.InitCache()
	if err != nil {
		log.Error("初始化缓存失败！", err)
		return
	}
	log.Info("初始化缓存完成")
	//初始化、加载持久化文件
	log.Info("初始化、加载持久化文件...")
	dumpOn := conf.GetSystemConfig().MustBool("dump", "dump.on", true)
	dumpTrigger := conf.GetSystemConfig().MustValue("dump", "trigger", "")
	dumpFilePath := conf.GetSystemConfig().MustValue("dump", "filePath", "./data/dump.gdb")
	err = gdb.InitGDB(dumpOn, dumpTrigger, dumpFilePath)
	if err != nil {
		log.Error("初始化、加载持久化文件失败！", err)
		return
	}
	//启动服务
	port := conf.GetSystemConfig().MustInt("server", "port", 1334)
	aliveTime := conf.GetSystemConfig().MustInt("server", "aliveTime", 30)
	connectType := conf.GetSystemConfig().MustValue("server", "connectType", "long")
	err = server.Start(port, aliveTime, connectType)
	if err != nil {
		log.Error("启动服务失败！", err)
		return
	}
}
