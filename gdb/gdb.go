package gdb

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/louch2010/gocache/log"
	"github.com/louch2010/goutil"
)

var open = true
var triggers [][]int
var filePath = ""

//初始化GDB
func InitGDB(dumpOn bool, dumpTrigger string, dumpFilePath string) error {
	open = dumpOn
	filePath = dumpFilePath
	if !dumpOn || goutil.StringUtil().IsEmpty(dumpFilePath) || goutil.StringUtil().IsEmpty(dumpTrigger) {
		open = false
		log.Info("GDB没有开启或定时为空，系统无需执行持久化文件处理!")
		return nil
	}
	unit := strings.Split(dumpTrigger, ",")
	triggers = make([][]int, len(unit))
	for i := 0; i < len(triggers); i++ {
		t := strings.SplitN(unit[i], " ", 2)
		time, err := strconv.Atoi(t[0])
		if err != nil {
			log.Error("初始化GDB时数据格式转换异常！", err)
			return err
		}
		modify, err := strconv.Atoi(t[1])
		if err != nil {
			log.Error("初始化GDB时数据格式转换异常！", err)
			return err
		}
		triggers[i] = []int{time, modify}
	}
	return LoadDB(filePath)
}

//加载本地gdb文件
func LoadDB(path string) error {
	log.Info("加载本地持久化gdb文件，文件路径：", path)
	start := time.Now()
	//文件校验
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("gdb文件不存在，无需加载：", path)
		} else {
			log.Error("加载gdb文件失败！", err)
			return err
		}
	}
	//文件解析
	err = parseGDB(file)
	if err != nil {
		log.Error("解析gdb文件失败！", err)
		return err
	}
	end := time.Now()
	cost := end.Sub(start)
	log.Info("加载本地持久化文件完成，耗时：", cost)
	return nil
}

//检查是否需要持久化
func CheckSave(dirty int, lastSave time.Time) {
	if !open {
		return
	}
}

//执行持久化操作
func SaveDB(dirty int, lastSave time.Time) error {
	if !open {
		return nil
	}
	return nil
}
