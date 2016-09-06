package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	InitLog("debug", "%Date %Time [%LEV] %Msg%n", "gocache.log", "02.01.2006", true)
	Debug("这是debug日志")
	Info("这是info日志")
	Error("这是error日志")
	Warn("这是warm日志")
	Trace("这是trace日志")
}
