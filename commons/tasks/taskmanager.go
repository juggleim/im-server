package tasks

import (
	"im-server/commons/logs"
	"im-server/commons/taskpools"
	"sync"
	"time"
)

var pool *taskpools.AsyncTaskPool

var taskCache map[string]*TaskItem
var taskLock *sync.Mutex
var stopped bool
var stopOnce sync.Once

type TaskItem struct {
	key               string
	f                 func()
	latestExecuteTime int64
	pending           bool
}

func init() {
	var err error
	pool, err = taskpools.NewAsyncTaskPool(64, 10000)
	if err != nil {
		panic("tasks: NewAsyncTaskPool: " + err.Error())
	}
	taskCache = make(map[string]*TaskItem)
	taskLock = &sync.Mutex{}
}

func StopTaskExecute() {
	stopOnce.Do(func() {
		taskLock.Lock()
		stopped = true
		taskLock.Unlock()
		pool.Release()
	})
}

func resetPending(key string) {
	taskLock.Lock()
	if v, ok := taskCache[key]; ok && v != nil {
		v.pending = false
	}
	taskLock.Unlock()
}

func TaskExecute(key string, interval int64, f func()) {
	taskLock.Lock()
	if stopped {
		taskLock.Unlock()
		return
	}

	curr := time.Now().UnixMilli()
	val := taskCache[key]
	if val == nil {
		val = &TaskItem{key: key, latestExecuteTime: 0}
		taskCache[key] = val
	}
	if val.pending || val.latestExecuteTime+interval >= curr {
		taskLock.Unlock()
		return
	}

	val.pending = true
	val.f = f
	wrapped := func() {
		defer func() {
			taskLock.Lock()
			defer taskLock.Unlock()
			if v, ok := taskCache[key]; ok && v != nil {
				v.latestExecuteTime = time.Now().UnixMilli()
				v.pending = false
			}
		}()
		f()
	}

	stoppedNow := stopped
	taskLock.Unlock()

	if stoppedNow {
		resetPending(key)
		return
	}

	if err := pool.Submit(wrapped); err != nil {
		logs.Error(err)
		resetPending(key)
	}
}
