package conf

import (
	"log"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
)

//配置文件
var HelpConfigFile, SystemConfigFile *goconfig.ConfigFile

//配置文件路径
const CONFIG_SYSTEM_FILE = "sys.ini"

//初始化配置文件
func InitConfig(path string) error {
	//加载帮助文件
	helpConfigFile, err := goconfig.LoadFromReader(strings.NewReader(CONFIG_HELP_CONTENT_EN))
	if err != nil {
		log.Println("加载帮助文件失败！", err)
		return err
	}
	//加载系统配置文件
	systemConfigFile, err := goconfig.LoadFromReader(strings.NewReader(CONFIG_SYSTEM_DEFAULT))
	if err != nil {
		log.Println("加载默认配置文件失败！", err)
		return err
	}
	//如果没有指定配置文件，则加载默认配置文件
	if len(strings.TrimSpace(path)) == 0 {
		file, err := os.Open(CONFIG_SYSTEM_FILE)
		defer file.Close()
		//默认配置文件不存在，则创建
		if err != nil && os.IsNotExist(err) {
			log.Println("配置文件文件不存在，创建：", CONFIG_SYSTEM_FILE)
			goconfig.SaveConfigFile(systemConfigFile, CONFIG_SYSTEM_FILE)
		} else {
			path = CONFIG_SYSTEM_FILE
		}
	}
	//如果用户指定了配置文件，或存在默认配置，则覆盖系统预设配置文件
	if len(strings.TrimSpace(path)) > 0 {
		log.Println("加载系统配置文件，文件路径：", path)
		userConfigFile, err := goconfig.LoadConfigFile(path)
		if err != nil {
			log.Println("加配置文件失败！", err)
			return err
		}
		for _, sec := range systemConfigFile.GetSectionList() {
			m, _ := systemConfigFile.GetSection(sec)
			for k, _ := range m {
				if nv, err := userConfigFile.GetValue(sec, k); err == nil {
					systemConfigFile.SetValue(sec, k, nv)
				}
			}
		}
		log.Println("加载系统配置文件完成")
	}
	HelpConfigFile = helpConfigFile
	SystemConfigFile = systemConfigFile
	return nil
}

//获取帮助配置
func GetHelpConfig() *goconfig.ConfigFile {
	return HelpConfigFile
}

//获取系统配置
func GetSystemConfig() *goconfig.ConfigFile {
	return SystemConfigFile
}

//获取日志配置
func GetLogConfig(level string, format string, path string, roll string, consoleOn bool) string {
	config := `
	<seelog>
	    <outputs formatid="main">`
	if consoleOn {
		config += `
			<filter levels="` + level + `">
	        	<console />
	    	</filter>
		`
	}
	config += `
			<filter levels="` + level + `">
				<rollingfile type="date" filename="` + path + `" datepattern="` + roll + `" maxrolls="7" />
	        </filter>
	    </outputs>
	    <formats>
	        <format id="main" format="` + format + `"/>
	    </formats>
	</seelog>
	`
	return config
}

//帮助
const CONFIG_HELP_CONTENT_EN = `
[connect]
Desc=connect to the server
Format=connect [-t'table'] [-a'pwd'] [-i'ip'] [-p'port'] [-e'e1,e2...']
[exit]
Desc=close connect and exit
Format=exit
[ping]
Desc=use for check the server is running
Format=ping
[help]
Desc=show commnd manual
Format=help
[set]
Desc=set key-value
Format=set key value [time]
[get]
Desc=get key-value
Format=get key
[delete]
Desc=delete key-value
Format=delete key
[exist]
Desc=check the key-value is exist
Format=exist key
[info]
Desc=show system info
Format=info
[use]
Desc=change cache table
Format=use [table name]
[showt]
Desc=show table info
Format=showt [table name]
[showi]
Desc=show item info
Format=showi 'item key'
`

//系统默认配置
const CONFIG_SYSTEM_DEFAULT = `
appname=gocache
version=1.0
author=luocihang@126.com

[server]
#server port
port=1334
#connect password
password=
#max size for connect pool
maxPoolSize=10
#core size for connect pool
corePoolSize=5
#connect alive time, unit is second
aliveTime=3000
#system table use to cache connection info
sysTable=sys
#the commnds which could be execute without login
anonymCommnd=ping,connect,exit,help,info

[table]
#default table name. if client connect server without assign table name or 'openSession' is false, 
#use this name as default table name
default=default

[client]
#client connect type: long|short
connectType=long
#open session: true|false
openSession=true

[dump]
#dump file directory
dir=./data
#dump file name extension
ext=.dmp

[log]
#log level: trace|debug|info|critical|warn|error
level = trace
#use console log: true|false
console.on=true
#log format
format=%Date %Time [%LEV] %Msg%n
#log file path
path=./log/gocache.log
#
roll=02.01.2006
`
