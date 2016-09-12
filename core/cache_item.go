package core

import (
	"sync"
	"time"
)

//缓存项
type CacheItem struct {
	sync.RWMutex                 //锁
	key            string        //键
	value          interface{}   //值
	liveTime       time.Duration //存活时间
	createTime     time.Time     //创建时间
	lastAccessTime time.Time     //最后访问时间
	accessCount    int64         //访问次数
	dataType       string        //数据类型
}

//键
func (item *CacheItem) Key() string {
	item.RLock()
	defer item.RUnlock()
	return item.key
}

//获取值
func (item *CacheItem) Value() interface{} {
	item.RLock()
	v := item.value
	defer item.RUnlock()
	return v
}

//获取值
func (item *CacheItem) SetValue(value interface{}) {
	item.Lock()
	item.value = value
	defer item.Unlock()
}

//存活时间
func (item *CacheItem) LiveTime() time.Duration {
	item.RLock()
	defer item.RUnlock()
	return item.liveTime
}

//创建时间
func (item *CacheItem) CreateTime() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.createTime
}

//最后访问时间
func (item *CacheItem) LastAccessTime() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.lastAccessTime
}

//访问次数
func (item *CacheItem) AccessCount() int64 {
	item.RLock()
	defer item.RUnlock()
	return item.accessCount
}

//访问
func (item *CacheItem) Access() {
	item.Lock()
	defer item.Unlock()
	item.accessCount++
	item.lastAccessTime = time.Now()
}

//数据类型
func (item *CacheItem) SetDataType(dataType string) {
	item.Lock()
	defer item.Unlock()
	item.dataType = dataType
}

func (item *CacheItem) DataType() string {
	item.RLock()
	defer item.RUnlock()
	return item.dataType
}

//新建Item
func NewCacheItem(key string, value interface{}, liveTime time.Duration, dataType string) *CacheItem {
	now := time.Now()
	item := CacheItem{
		key:            key,
		value:          value,
		liveTime:       liveTime,
		createTime:     now,
		lastAccessTime: now,
		accessCount:    0,
		dataType:       dataType,
	}
	return &item
}
