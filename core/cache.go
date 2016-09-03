package core

import (
	"log"
	"sync"

	"github.com/louch2010/goutil"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

//获取缓存表，如果不存在，则新建缓存表；存在则直接返回
func Cache(table string) (*CacheTable, error) {
	if goutil.StringUtil().IsEmpty(table) {
		return nil, ERROR_TABLE_NAME
	}
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()
	if !ok {
		log.Println("缓存表不存在，新建缓存表，表名为：", table)
		t = NewCacheTable(table)
		mutex.Lock()
		cache[table] = t
		mutex.Unlock()
		//事件回调通知
		t.tableEventCallBack(t, EVENT_TABLE_ADD)
	}
	return t, nil
}

//获取默认表，表名为
func DefaultCache() *CacheTable {
	t, _ := Cache(DEFAULT_TABLE_NAME)
	return t
}
