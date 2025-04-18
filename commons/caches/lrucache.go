package caches

import (
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/simplelru"
)

type lruCacheItem struct {
	value        interface{}
	lastReadTime int64
	addedTime    int64
}

type LruCache struct {
	name               string
	lastRecord         int64
	lru                simplelru.LRUCache
	lock               sync.RWMutex
	readTimeoutChecker *time.Ticker
	MaxLifeCycle       time.Duration
	valueCreator       func(key interface{}) interface{}

	batchEvict func(items []CacheItem)
	batchSize  int
}

type CacheItem struct {
	Key   interface{}
	Value interface{}
}

func NewLruCacheWithAddReadTimeout(name string, size int, onEvict simplelru.EvictCallback, timeoutAfterRead time.Duration, timeoutAfterCreate time.Duration) *LruCache {
	cache := NewLruCache(name, size, onEvict)
	cache.AddTimeoutAfterRead(timeoutAfterRead)
	cache.AddTimeoutAfterCreate(timeoutAfterCreate)
	return cache
}

func NewLruCacheWithReadTimeout(name string, size int, onEvict simplelru.EvictCallback, timeoutAfterRead time.Duration) *LruCache {
	cache := NewLruCache(name, size, onEvict)
	cache.AddTimeoutAfterRead(timeoutAfterRead)
	return cache
}

func NewLruCache(name string, size int, onEvict simplelru.EvictCallback) *LruCache {
	myLru, _ := simplelru.NewLRU(size, func(key, value interface{}) {
		if onEvict != nil && value != nil {
			cacheItem, ok := value.(*lruCacheItem)
			if ok {
				onEvict(key, cacheItem.value)
			}
		}
	})
	cache := &LruCache{
		name: name,
		lru:  myLru,
	}
	return cache
}

func (c *LruCache) SetBatchEvict(batchSize int, f func(items []CacheItem)) *LruCache {
	if batchSize <= 0 {
		return c
	}
	c.batchSize = batchSize
	c.batchEvict = f
	return c
}

func (c *LruCache) SetValueCreator(creator func(interface{}) interface{}) *LruCache {
	c.valueCreator = creator
	return c
}

func (c *LruCache) AddTimeoutAfterCreate(timeout time.Duration) *LruCache {
	c.MaxLifeCycle = timeout
	return c
}

func (c *LruCache) AddTimeoutAfterRead(timeout time.Duration) *LruCache {
	if c.readTimeoutChecker != nil {
		c.readTimeoutChecker.Stop()
	}
	c.readTimeoutChecker = time.NewTicker(time.Second)
	go func() {
		for task := range c.readTimeoutChecker.C {
			current := time.Now().UnixMilli()
			if current-task.UnixMilli() > 500 {
				continue
			}
			timeLine := current - int64(timeout)/(1000*1000)
			c.cleanOldestByReadTime(timeLine)
		}
	}()
	return c
}

func (c *LruCache) cleanOldestByReadTime(timeLine int64) {
	cacheItems := []CacheItem{}
	for {
		itemKey, itemValue, ok := c.lru.GetOldest()
		if ok {
			valObj := itemValue.(*lruCacheItem)
			lastReadTime := valObj.lastReadTime
			if lastReadTime < timeLine {
				c.Remove(itemKey)
				if c.batchEvict != nil {
					cacheItems = append(cacheItems, CacheItem{
						Key:   itemKey,
						Value: valObj.value,
					})
					if len(cacheItems) >= c.batchSize {
						c.batchEvict(cacheItems)
						cacheItems = []CacheItem{}
					}
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	if c.batchEvict != nil && len(cacheItems) > 0 {
		c.batchEvict(cacheItems)
	}
}

func (c *LruCache) Add(key, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.innerAdd(key, value)
}

func (c *LruCache) AddIfAbsent(key, value interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.innerContains(key) {
		old, ok := c.innerGet(key)
		if ok {
			return old, false
		}
	}
	c.innerAdd(key, value)
	return value, true
}

func (c *LruCache) AddIfAbsendNoGetOldVal(key, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.innerContains(key) {
		return false
	}
	c.innerAdd(key, value)
	return true
}

func (c *LruCache) innerAdd(key, value interface{}) bool {
	nowTime := time.Now().UnixMilli()
	return c.lru.Add(key, &lruCacheItem{
		value:        value,
		lastReadTime: nowTime,
		addedTime:    nowTime,
	})
}

func (c *LruCache) Get(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.innerGet(key)
}

func (c *LruCache) innerGet(key interface{}) (interface{}, bool) {
	item, ok := c.lru.Get(key)
	if ok {
		cacheItem := item.(*lruCacheItem)
		if c.MaxLifeCycle > 0 {
			timeLine := time.Now().UnixMilli() - int64(c.MaxLifeCycle)/1000/1000
			if cacheItem.addedTime < timeLine { //remove
				c.lru.Remove(key)
				return nil, false
			}
		}
		cacheItem.lastReadTime = time.Now().UnixMilli()
		return cacheItem.value, ok
	} else {
		return nil, ok
	}
}

func (c *LruCache) GetByDefault(key interface{}, defaultValue interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	val, ok := c.innerGet(key)
	if ok {
		return val, ok
	} else {
		return defaultValue, ok
	}
}

func (c *LruCache) GetByCreator(key interface{}, creator func() interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	val, ok := c.innerGet(key)
	if ok {
		return val, ok
	} else {
		if creator != nil {
			newVal := creator()
			if newVal != nil {
				c.innerAdd(key, newVal)
				return newVal, true
			}
		} else {
			if c.valueCreator != nil {
				newVal := c.valueCreator(key)
				if newVal != nil {
					c.innerAdd(key, newVal)
					return newVal, true
				}
			}
		}
	}
	return nil, ok
}

func (c *LruCache) Contains(key interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Contains(key)
}

func (c *LruCache) innerContains(key interface{}) bool {
	return c.lru.Contains(key)
}

func (c *LruCache) Peek(key interface{}) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	item, ok := c.lru.Peek(key)
	if ok {
		cacheItem := item.(*lruCacheItem)
		if c.MaxLifeCycle > 0 {
			timeLine := time.Now().UnixMilli() - int64(c.MaxLifeCycle)/1000/1000
			if cacheItem.addedTime < timeLine {
				c.lru.Remove(key)
				return nil, false
			}
		}
		return cacheItem.value, ok
	} else {
		return nil, ok
	}
}

func (c *LruCache) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.lru.Purge()
}

func (c *LruCache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Len()
}

func (c *LruCache) ReSize(size int) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Resize(size)
}

func (c *LruCache) Remove(key interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Remove(key)
}

func (c *LruCache) Keys() []interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.lru.Keys()
}
