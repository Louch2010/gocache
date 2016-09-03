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

//访问
func (item *CacheItem) Access() {
	item.Lock()
	defer item.Unlock()
	now := time.Now()
	item.accessCount++
	item.lastAccessTime = now
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
