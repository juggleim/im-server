package tools

import (
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

type SinglePools struct {
	pools     []*SinglePool
	lock      *sync.RWMutex
	isBlocked bool
	queueSize int
}
type SinglePool struct {
	taskChan  chan func()
	isActived atomic.Bool
	isBlocked atomic.Bool
}

func (pool *SinglePool) QueueLen() int {
	return len(pool.taskChan)
}

func (pool *SinglePool) QueueCapacity() int {
	return cap(pool.taskChan)
}

type SinglePoolsQueueStats struct {
	ActivePools int `json:"active_pools"`
	QueuedTasks int `json:"queued_tasks"`
	FullPools   int `json:"full_pools"`
	MaxQueueLen int `json:"max_queue_len_current"`
}

func (pools *SinglePools) QueueStats() SinglePoolsQueueStats {
	pools.lock.RLock()
	defer pools.lock.RUnlock()
	var stats SinglePoolsQueueStats
	for _, pool := range pools.pools {
		if pool == nil {
			continue
		}
		stats.ActivePools++
		length := pool.QueueLen()
		stats.QueuedTasks += length
		if length > stats.MaxQueueLen {
			stats.MaxQueueLen = length
		}
		if length == pool.QueueCapacity() {
			stats.FullPools++
		}
	}
	return stats
}

func NewSinglePools(count, queueSize int, isBlocked bool) *SinglePools {
	if queueSize <= 0 {
		queueSize = 100
	}
	return &SinglePools{
		pools:     make([]*SinglePool, count),
		lock:      &sync.RWMutex{},
		isBlocked: isBlocked,
		queueSize: queueSize,
	}
}

func (pools *SinglePools) PoolIndex(key string) int {
	return int(HashStr(key) % uint32(len(pools.pools)))
}

func (pools *SinglePools) GetPool(key string) *SinglePool {
	index := pools.PoolIndex(key)
	pools.lock.RLock()
	pool := pools.pools[index]
	pools.lock.RUnlock()
	if pool == nil {
		pools.lock.Lock()
		defer pools.lock.Unlock()
		pool = pools.pools[index]
		if pool == nil {
			pool = &SinglePool{
				taskChan: make(chan func(), pools.queueSize),
			}
			pool.isActived.Store(true)
			pool.isBlocked.Store(pools.isBlocked)
			pools.pools[index] = pool
			pool.Start()
		}
	}
	return pool
}

func (pool *SinglePool) Start() {
	go func() {
		for pool.isActived.Load() {
			task := <-pool.taskChan
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Printf("SinglePool task panic: %v\n%s", err, debug.Stack())
					}
				}()
				task()
			}()
		}
		close(pool.taskChan)
	}()
}

func (pool *SinglePool) NoBlockSubmit(task func()) bool {
	if pool.isActived.Load() {
		select {
		case pool.taskChan <- task:
			return true
		default:
			return false
		}
	}
	return false
}

func (pool *SinglePool) Submit(task func()) bool {
	if pool.isActived.Load() {
		if pool.isBlocked.Load() {
			pool.taskChan <- task
			return true
		} else {
			select {
			case pool.taskChan <- task:
				return true
			default:
				return false
			}
		}
	}
	return false
}
