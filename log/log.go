package log

import (
	default_log "log"
	"strings"

	"github.com/cihub/seelog"
	"github.com/louch2010/gocache/conf"
	"github.com/louch2010/goutil"
)

var logFactory seelog.LoggerInterface

//初始化日志
func InitLog(level string, format string, path string, roll string, consoleOn bool) error {
	//初始化日志级别
	LEVEL_LIST := []string{"error", "warn", "critical", "info", "debug", "trace"}
	level = strings.ToLower(goutil.StringUtil().TrimToEmpty(level))
	consv := 0
	for index, v := range LEVEL_LIST {
		if v == level {
			consv = index
			break
		}
	}
	levels := ""
	for i := 0; i <= consv; i++ {
		levels += LEVEL_LIST[i]
		if i != consv {
			levels += ","
		}
	}
	//初始化日志
	logConf := conf.GetLogConfig(levels, format, path, roll, consoleOn)
	log, err := seelog.LoggerFromConfigAsString(logConf)
	if err != nil {
		default_log.Println("初始化日志失败！", err)
		return err
	}
	logFactory = log
	return nil
}

func Debug(v ...interface{}) {
	logFactory.Debug(v)
}

func Info(v ...interface{}) {
	logFactory.Info(v)
}

func Error(v ...interface{}) {
	logFactory.Error(v)
}

func Warn(v ...interface{}) {
	logFactory.Warn(v)
}

func Trace(v ...interface{}) {
	logFactory.Trace(v)
}
