package tools

import (
	"sync"
	"sync/atomic"
)

type SinglePools struct {
	pools     []*SinglePool
	lock      *sync.RWMutex
	isBlocked bool
}

type SinglePool struct {
	taskChan  chan func()
	isActived atomic.Bool
	isBlocked bool
	done      chan struct{}
}

func NewSinglePools(count int, isBlocked bool) *SinglePools {
	return &SinglePools{
		pools:     make([]*SinglePool, count),
		lock:      &sync.RWMutex{},
		isBlocked: isBlocked,
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
				isBlocked: pools.isBlocked,
			}
			pool.isActived.Store(true)
			pools.pools[index] = pool
			pool.Start()
		}
	}
	return pool
}

func (pool *SinglePool) Start() {
	pool.done = make(chan struct{})
	go func() {
		for {
			select {
			case task := <-pool.taskChan:
				task()
			case <-pool.done:
				// 把队列里剩余任务消费完再退出
				for task := range pool.taskChan {
					task()
				}
				return
			}
		}
	}()
}

func (pool *SinglePool) Stop() {
	if pool.isActived.CompareAndSwap(true, false) {
		close(pool.done)
	}
}

func (pool *SinglePool) Submit(task func()) bool {
	if pool.isActived.Load() {
		if pool.isBlocked {
			pool.taskChan <- task
			return true
		}
		select {
		case pool.taskChan <- task:
			return true
		default:
			return false
		}
	}
	return false
}
