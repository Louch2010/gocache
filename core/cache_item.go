package core

import (
	"sync"
	"time"
)

//缓存项
type CacheItem struct {
	sync.RWMutex                 //锁
	key            interface{}   //键
	value          interface{}   //值
	liveTime       time.Duration //存活时间
	createTime     time.Time     //创建时间
	lastAccessTime time.Time     //最后访问时间
	lastModifyTime time.Time     //最后修改时间
	accessCount    int64         //访问次数
}

//键
func (item *CacheItem) Key() interface{} {
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

//最后修改时间
func (item *CacheItem) LastModifyTime() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.lastModifyTime
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

//修改
func (item *CacheItem) Modify() {
	item.Lock()
	defer item.Unlock()
	item.lastModifyTime = time.Now()
}

//新建Item
func NewCacheItem(key interface{}, value interface{}, liveTime time.Duration) *CacheItem {
	now := time.Now()
	item := CacheItem{
		key:            key,
		value:          value,
		liveTime:       liveTime,
		createTime:     now,
		lastAccessTime: now,
		lastModifyTime: now,
		accessCount:    0,
	}
	return &item
}
