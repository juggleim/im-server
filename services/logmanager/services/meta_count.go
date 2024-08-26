package services

import (
	"sync"
	"sync/atomic"
)

type MetaCount struct {
	countMap map[string]*int64
	lock     sync.RWMutex
}

func NewMetaCount() *MetaCount {
	return &MetaCount{
		countMap: make(map[string]*int64),
	}
}

func (mc *MetaCount) Increment(key string, initialValueFn func(key string) int64) (newValue int64) {
	mc.lock.RLock()
	c, ok := mc.countMap[key]
	mc.lock.RUnlock()
	if ok {
		newValue = atomic.AddInt64(c, 1)
	} else {
		mc.lock.Lock()
		var v = initialValueFn(key)
		mc.countMap[key] = &v
		newValue = v
		mc.lock.Unlock()
	}

	return
}
