package tools

import (
	"sync"
)

type SinglePools struct {
	pools []*SinglePool
	lock  *sync.RWMutex
}
type SinglePool struct {
	taskChan  chan func()
	isActived bool
}

func NewSinglePools(count int) *SinglePools {
	return &SinglePools{
		pools: make([]*SinglePool, count),
		lock:  &sync.RWMutex{},
	}
}

func (pools *SinglePools) GetPool(key string) *SinglePool {
	hash := HashStr(key)
	index := int(hash % uint32(len(pools.pools)))
	pool := pools.pools[index]
	if pool == nil {
		pools.lock.Lock()
		defer pools.lock.Unlock()
		pool = pools.pools[index]
		if pool == nil {
			pool = &SinglePool{
				taskChan:  make(chan func(), 100),
				isActived: true,
			}
			pools.pools[index] = pool
			pool.Start()
		}
	}
	return pool
}

func (pool *SinglePool) Start() {
	go func() {
		for pool.isActived {
			task := <-pool.taskChan
			task()
		}
		close(pool.taskChan)
	}()
}

func (pool *SinglePool) Submit(task func()) bool {
	if pool.isActived {
		pool.taskChan <- task
		return true
	}
	return false
}
