package tasks

import (
	"sync"
	"time"

	"github.com/Jeffail/tunny"
)

var pool *tunny.Pool

// var taskCache *sync.Map
var taskCache map[string]*TaskItem
var taskLock *sync.Mutex
var taskChan chan *TaskItem
var stopChan chan struct{}
var stopped bool
var stopOnce sync.Once

type TaskItem struct {
	key               string
	f                 func()
	latestExecuteTime int64
	pending           bool
}

func init() {
	pool = tunny.NewCallback(64)
	taskCache = make(map[string]*TaskItem)
	taskLock = &sync.Mutex{}
	taskChan = make(chan *TaskItem, 1000)
	stopChan = make(chan struct{})
}

func StartTaskExecute() {
	go func() {
		for {
			select {
			case task, ok := <-taskChan:
				if !ok {
					return
				}
				if task != nil && task.f != nil {
					go pool.Process(task.f)
				}
			case <-stopChan:
				return
			}
		}
	}()
}

func StopTaskExecute() {
	stopOnce.Do(func() {
		taskLock.Lock()
		stopped = true
		close(stopChan)
		close(taskChan)
		taskLock.Unlock()
		pool.Close()
	})
}

func TaskExecute(key string, interval int64, f func()) {
	taskLock.Lock()
	defer taskLock.Unlock()

	if stopped {
		return
	}

	curr := time.Now().UnixMilli()
	if val, exist := taskCache[key]; exist && val != nil {
		if val.pending || val.latestExecuteTime+interval >= curr {
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
		select {
		case taskChan <- &TaskItem{key: key, f: wrapped}:
		default:
			val.pending = false
		}
	} else {
		taskCache[key] = &TaskItem{
			key:               key,
			latestExecuteTime: curr,
			pending:           false,
		}
	}
}
