package tools

import (
	"sync"
	"time"
)

type BatchExecutorPool struct {
	executors []*BatchExecutor
	lock      *sync.RWMutex

	batchSize     int
	checkDuration time.Duration
	execFun       func(tasks []interface{})
}

func NewBatchExecutorPool(count int, batchSize int, checkDuration time.Duration, exec func(tasks []interface{})) *BatchExecutorPool {
	return &BatchExecutorPool{
		executors: make([]*BatchExecutor, count),
		lock:      &sync.RWMutex{},

		batchSize:     batchSize,
		checkDuration: checkDuration,
		execFun:       exec,
	}
}

func (pool *BatchExecutorPool) GetBatchExecutor(key string) *BatchExecutor {
	hash := HashStr(key)
	index := int(hash % uint32(len(pool.executors)))
	executor := pool.executors[index]
	if executor == nil {
		pool.lock.Lock()
		defer pool.lock.Unlock()
		executor = pool.executors[index]
		if executor == nil {
			executor = NewBatchExecutor(pool.batchSize, pool.checkDuration, pool.execFun)
			pool.executors[index] = executor
		}
	}
	return executor
}

func (pool *BatchExecutorPool) Stop() {
	for _, executor := range pool.executors {
		if executor != nil {
			executor.Stop()
		}
	}
}

func NewBatchExecutor(batchSize int, duration time.Duration, exec func([]interface{})) *BatchExecutor {
	executor := &BatchExecutor{
		lock:          &sync.RWMutex{},
		batchSize:     batchSize,
		checkDuration: duration,
		executeFun:    exec,

		taskBuffer: make(map[string]interface{}),
	}
	executor.start()
	return executor
}

type BatchExecutor struct {
	// tasks         []interface{}
	taskBuffer    map[string]interface{}
	lock          *sync.RWMutex
	batchSize     int
	checkTimer    *time.Ticker
	checkDuration time.Duration

	executeFun func(tasks []interface{})
}

func (executor *BatchExecutor) start() {
	if executor.checkTimer != nil {
		executor.checkTimer.Stop()
	}
	executor.checkTimer = time.NewTicker(executor.checkDuration)
	go func() {
		for range executor.checkTimer.C {
			tasks := executor.featchTasks()
			if len(tasks) > 0 && executor.executeFun != nil {
				executor.executeFun(tasks)
			}
		}
	}()
}

func (executor *BatchExecutor) innerAppend(uniqId string, task interface{}) []interface{} {
	executor.lock.Lock()
	defer executor.lock.Unlock()
	if len(executor.taskBuffer) < executor.batchSize {
		executor.taskBuffer[uniqId] = task
	} else {
		tasks := []interface{}{}
		for _, v := range executor.taskBuffer {
			tasks = append(tasks, v)
		}
		executor.taskBuffer = make(map[string]interface{})
		return tasks
	}
	return nil
}

func (executor *BatchExecutor) Append(uniqId string, task interface{}) {
	tasks := executor.innerAppend(uniqId, task)
	if len(tasks) > 0 && executor.executeFun != nil {
		executor.executeFun(tasks)
	}
}

func (executor *BatchExecutor) featchTasks() []interface{} {
	executor.lock.Lock()
	defer executor.lock.Unlock()
	tasks := []interface{}{}
	for _, v := range executor.taskBuffer {
		tasks = append(tasks, v)
	}
	executor.taskBuffer = make(map[string]interface{})
	return tasks
}

func (executor *BatchExecutor) Stop() {
	if executor.checkTimer != nil {
		executor.checkTimer.Stop()
	}
}
