package core

import (
	"log"
	"strconv"
	"sync"
	"time"
)

//缓存表，缓存服务器中包含多个缓存表，name作为唯一标识
type CacheTable struct {
	sync.RWMutex                                                                     //锁
	name                      string                                                 //表名
	items                     map[interface{}]*CacheItem                             //缓存项
	createTime                time.Time                                              //创建时间
	lastAccessTime            time.Time                                              //最后访问时间
	lastModifyTime            time.Time                                              //最后修改时间
	accessCount               int64                                                  //访问次数
	cleanupTimer              *time.Timer                                            //清除定时器
	cleanupInterval           time.Duration                                          //清除时间间隔
	itemEventCallBack         func(table *CacheTable, item *CacheItem, event string) //缓存项事件回调
	tableEventCallBack        func(table *CacheTable, event string)                  //缓存表事件回调
	startLoadFromDiskCallBack func(table *CacheTable)                                //开始从硬盘加载回调
	endLoadFromDiskCallBack   func(table *CacheTable)                                //完成从硬盘加载回调
	startDumpToDiskCallBack   func(table *CacheTable)                                //开始写入硬盘回调
	endDumpToDiskCallBack     func(table *CacheTable)                                //完成写入硬盘回调
}

//获取缓存项数量
func (table *CacheTable) ItemCount() int {
	table.RLock()
	defer table.RUnlock()
	return len(table.items)
}

//访问表时更新表信息
func (table *CacheTable) Access() {
	table.Lock()
	now := time.Now()
	table.accessCount++
	table.lastAccessTime = now
	table.lastModifyTime = now
	table.Unlock()
}

//缓存项是否存在
func (table *CacheTable) IsExist(key interface{}) bool {
	//使用读锁，判断是否存在
	table.RLock()
	_, ok := table.items[key]
	table.RUnlock()
	log.Println("查询键是否存在，表名：", table.name, "键为：", key, "，结果：", ok)
	//更新表信息
	table.Access()
	return ok
}

//添加缓存项，如果已经存在，会返回原缓存项，如果不存在，则返回的是nil
func (table *CacheTable) Set(key interface{}, value interface{}, liveTime time.Duration) *CacheItem {
	item := NewCacheItem(key, value, liveTime)
	return table.AddItem(key, item)
}

//添加缓存项，如果已经存在，会返回原缓存项，如果不存在，则返回的是nil
func (table *CacheTable) AddItem(key interface{}, item *CacheItem) *CacheItem {
	//修改表属性
	table.Lock()
	old, ok := table.items[key]
	now := time.Now()
	interval := table.cleanupInterval
	table.items[key] = item
	table.accessCount++
	table.lastAccessTime = now
	table.lastModifyTime = now
	//启动定时检查，条件为：该项设定了过期时间，并且过期时间小于间隔时间，或没有设置过间隔时间
	//如果大于间隔时间，那么说明会有更早的定时会启动，此时无需再设定
	liveTime := item.liveTime
	if liveTime > 0 && (interval == 0 || liveTime < interval) {
		table.cleanupTimer = time.AfterFunc(liveTime, func() {
			go table.expireCheck()
		})
	}
	table.Unlock()
	log.Println("新增缓存项，键：", key, "，值：", item.value, "，原来是否存在：", ok)
	//调用回调函数
	if ok {
		table.itemEventCallBack(table, item, EVENT_ITEM_MODIFY)
	} else {
		table.itemEventCallBack(table, item, EVENT_ITEM_ADD)
	}
	return old
}

//删除缓存项，删除成功返回true，删除失败返回false
func (table *CacheTable) Delete(key interface{}) bool {
	table.RLock()
	item, ok := table.items[key]
	table.RUnlock()
	if ok {
		//更新表信息
		table.Lock()
		now := time.Now()
		delete(table.items, key)
		table.accessCount++
		table.lastAccessTime = now
		table.lastModifyTime = now
		table.Unlock()
		log.Println("删除缓存项，删除成功，表名：", table.name, "，键名：", key)
		//回调通知
		table.itemEventCallBack(table, item, EVENT_ITEM_DELETE)
		return true
	} else {
		log.Println("删除缓存项，缓存项不存在，删除失败，表名：", table.name, "，键名：", key)
		return false
	}
}

//获取缓存项
func (table *CacheTable) Get(key interface{}) *CacheItem {
	//获取
	table.RLock()
	item, ok := table.items[key]
	table.RUnlock()
	if !ok {
		log.Println("获取缓存项失败，缓存项不存在，表名：", table.name, "，键名：", key)
		return nil
	}
	log.Println("获取缓存项，表名：", table.name, "，键名：", key, "，值为：", item.value)
	//修改表信息
	table.Access()
	//对过期的缓存项进行处理
	now := time.Now()
	if item.liveTime > 0 && now.Sub(item.createTime) >= item.liveTime {
		log.Println("获取缓存项失败，绑在已过期，表名：", table.name, "，键名：", key)
		return nil
	}
	return item
}

//设置缓存项事件回调
func (table *CacheTable) SetItemEventCallBack(f func(table *CacheTable, item *CacheItem, event string)) {
	table.Lock()
	defer table.Unlock()
	table.itemEventCallBack = f
}

//设置缓存表事件回调
func (table *CacheTable) SetTableEventCallBack(f func(table *CacheTable, event string)) {
	table.Lock()
	defer table.Unlock()
	table.tableEventCallBack = f
}

//设置开始从硬盘加载回调
func (table *CacheTable) SetStartLoadFromDiskCallBack(f func(table *CacheTable)) {
	table.Lock()
	defer table.Unlock()
	table.startLoadFromDiskCallBack = f
}

//设置完成从硬盘加载回调
func (table *CacheTable) SetEndLoadFromDiskCallBack(f func(table *CacheTable)) {
	table.Lock()
	defer table.Unlock()
	table.endLoadFromDiskCallBack = f
}

//设置开始写入硬盘回调
func (table *CacheTable) SetStartDumpToDiskCallBack(f func(table *CacheTable)) {
	table.Lock()
	defer table.Unlock()
	table.startDumpToDiskCallBack = f
}

//设置完成写入硬盘回调
func (table *CacheTable) SetEndDumpToDiskCallBack(f func(table *CacheTable)) {
	table.Lock()
	defer table.Unlock()
	table.endDumpToDiskCallBack = f
}

//过期检查
func (table *CacheTable) expireCheck() {
	log.Println("准备启动过期检查，表名：", table.name)
	table.Lock()
	//已经在执行，那么就先停止定时
	if table.cleanupTimer != nil {
		table.cleanupTimer.Stop()
	}
	items := table.items
	table.Unlock()
	log.Println("----开始启动过期检查，表名：", table.name, "，表容量：", strconv.Itoa(len(items)))
	//遍历
	now := time.Now()
	interval := 0 * time.Second
	for key, item := range items {
		item.RLock()
		liveTime := item.liveTime
		//lastAccessTime := item.lastAccessTime
		createTime := item.createTime
		item.RUnlock()
		//没有设置存活时间，表示永久缓存，不做处理
		if liveTime == 0 {
			continue
		}
		//当前时间 - 创建时间 > 存活时间，表示已经过期，需要删除
		if now.Sub(createTime) >= liveTime {
			table.Delete(key)
			continue
		}
		//清除时间间隔没有设置（即为0）或存活时间小于间隔时间，则将间隔时间设置为 存活时间 + 创建时间 - 当前时间
		if interval == 0 || liveTime < interval {
			interval = liveTime - now.Sub(createTime)
		}
	}
	table.Lock()
	//设置间隔时间，并在 间隔时间后，再次启动过期检查
	table.cleanupInterval = interval
	if interval > 0 {
		table.cleanupTimer = time.AfterFunc(interval, func() {
			go table.expireCheck()
		})
	}
	table.Unlock()
	log.Println("----过期检查结束，表名：", table.name, "，耗时：", strconv.FormatInt(time.Now().Unix()-now.Unix(), 10))
}

//新建缓存表
func NewCacheTable(name string) *CacheTable {
	now := time.Now()
	table := CacheTable{
		name:            name,
		items:           make(map[interface{}]*CacheItem),
		createTime:      now,
		lastAccessTime:  now,
		lastModifyTime:  now,
		accessCount:     0,
		cleanupTimer:    nil,
		cleanupInterval: 0,
		itemEventCallBack: func(table *CacheTable, item *CacheItem, event string) {
			log.Println("缓存项事件回调...事件类型：", event)
		},
		tableEventCallBack: func(table *CacheTable, event string) {
			log.Println("缓存表事件回调...事件类型：", event)
		},
		startLoadFromDiskCallBack: func(table *CacheTable) {
			log.Println("开始从硬盘加载回调...")
		},
		endLoadFromDiskCallBack: func(table *CacheTable) {
			log.Println("完成从硬盘加载回调...")
		},
		startDumpToDiskCallBack: func(table *CacheTable) {
			log.Println("开始写入硬盘回调...")
		},
		endDumpToDiskCallBack: func(table *CacheTable) {
			log.Println("完成写入硬盘回调...")
		},
	}
	return &table
}
