package caches

import (
	"container/list"
	"sync"
	"time"
)

type EphemeralCacheItem struct {
	key       interface{}
	value     interface{}
	addedTime int64
}

type EphemeralCache struct {
	timer     *time.Ticker
	evictList *list.List
	items     map[interface{}]*list.Element
	onEvict   EphemeralCacheEvict
	lock      *sync.RWMutex
}
type EphemeralCacheEvict func(key interface{}, value interface{})

func NewEphemeralCache(checkInterval, maxLife time.Duration, onEvict EphemeralCacheEvict) *EphemeralCache {
	cache := &EphemeralCache{
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   onEvict,
		lock:      &sync.RWMutex{},
	}
	cache.AddTimeoutAfterCreate(checkInterval, maxLife)
	return cache
}

func (c *EphemeralCache) Len() (int, int) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	l := len(c.items)
	x := c.evictList.Len()
	return l, x
}

func (c *EphemeralCache) AddTimeoutAfterCreate(checkInterval, maxLife time.Duration) *EphemeralCache {
	if c.timer != nil {
		c.timer.Stop()
	}
	c.timer = time.NewTicker(checkInterval)
	go func() {
		for task := range c.timer.C {
			current := time.Now().UnixMilli()
			if current-task.UnixMilli() > 500 {
				continue
			}
			timeLine := current - int64(maxLife)/(1000*1000)
			c.cleanOldestByCreatedTime(timeLine)
		}
	}()
	return c
}
func (c *EphemeralCache) cleanOldestByCreatedTime(timeLine int64) {
	for {
		ele := c.evictList.Back()
		if ele != nil {
			kv := ele.Value.(*EphemeralCacheItem)
			if kv.addedTime < timeLine {
				c.Remove(kv.key)
			} else {
				break
			}
		} else {
			break
		}
	}
}

func (c *EphemeralCache) Add(key, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	//check for existing item
	if ent, ok := c.items[key]; ok {
		ent.Value.(*EphemeralCacheItem).value = val
		return
	}
	//add new item
	item := &EphemeralCacheItem{
		key:       key,
		value:     val,
		addedTime: time.Now().UnixMilli(),
	}
	entry := c.evictList.PushFront(item)
	c.items[key] = entry
}

func (c *EphemeralCache) Upsert(key interface{}, f func(oldVal interface{}) interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	//check for existing item
	if ent, ok := c.items[key]; ok {
		item := ent.Value.(*EphemeralCacheItem)
		newVal := f(item.value)
		if newVal != nil {
			item.value = newVal
		}
		return
	}
	//add new item
	newVal := f(nil)
	if newVal != nil {
		item := &EphemeralCacheItem{
			key:       key,
			value:     newVal,
			addedTime: time.Now().UnixMilli(),
		}
		entry := c.evictList.PushFront(item)
		c.items[key] = entry
	}
}

func (c *EphemeralCache) Remove(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if ent, ok := c.items[key]; ok {
		c.evictList.Remove(ent)
		kv := ent.Value.(*EphemeralCacheItem)
		delete(c.items, kv.key)
		if c.onEvict != nil {
			c.onEvict(kv.key, kv.value)
		}
	}
}
