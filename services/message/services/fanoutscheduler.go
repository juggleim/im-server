package services

import (
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type fanoutTask func(queueWait time.Duration)

type queuedFanoutTask struct {
	enqueuedAt time.Time
	run        fanoutTask
}

type fanoutSchedulerStats struct {
	Pending    int `json:"pending"`
	Running    int `json:"running"`
	ReadyKeys  int `json:"ready_keys"`
	ActiveKeys int `json:"active_keys"`
	Capacity   int `json:"capacity"`
	Workers    int `json:"workers"`
}

// keyedFanoutScheduler executes tasks for the same key in FIFO order while
// allowing different keys to run concurrently. Pending includes running tasks.
type keyedFanoutScheduler struct {
	mu        sync.Mutex
	notFull   *sync.Cond
	queues    map[string][]queuedFanoutTask
	scheduled map[string]bool
	ready     chan string
	pending   int
	running   int
	capacity  int
	workers   int
}

func newKeyedFanoutScheduler(workers, capacity int) *keyedFanoutScheduler {
	if workers <= 0 {
		workers = 1
	}
	if capacity < workers {
		capacity = workers
	}
	scheduler := &keyedFanoutScheduler{
		queues:    make(map[string][]queuedFanoutTask),
		scheduled: make(map[string]bool),
		ready:     make(chan string, capacity),
		capacity:  capacity,
		workers:   workers,
	}
	scheduler.notFull = sync.NewCond(&scheduler.mu)
	for i := 0; i < workers; i++ {
		go scheduler.work()
	}
	return scheduler
}

func (scheduler *keyedFanoutScheduler) Submit(key string, task fanoutTask) (time.Duration, bool) {
	if task == nil {
		return 0, false
	}
	started := time.Now()
	blocked := false
	scheduler.mu.Lock()
	for scheduler.pending >= scheduler.capacity {
		blocked = true
		scheduler.notFull.Wait()
	}
	scheduler.queues[key] = append(scheduler.queues[key], queuedFanoutTask{
		enqueuedAt: time.Now(),
		run:        task,
	})
	scheduler.pending++
	if !scheduler.scheduled[key] {
		scheduler.scheduled[key] = true
		scheduler.ready <- key
	}
	scheduler.mu.Unlock()
	return time.Since(started), blocked
}

func (scheduler *keyedFanoutScheduler) Stats() fanoutSchedulerStats {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()
	return fanoutSchedulerStats{
		Pending:    scheduler.pending,
		Running:    scheduler.running,
		ReadyKeys:  len(scheduler.ready),
		ActiveKeys: len(scheduler.scheduled),
		Capacity:   scheduler.capacity,
		Workers:    scheduler.workers,
	}
}

func (scheduler *keyedFanoutScheduler) work() {
	for key := range scheduler.ready {
		task := scheduler.take(key)
		if task.run != nil {
			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Printf("msg dispatch fanout panic: %v\n%s", err, debug.Stack())
					}
				}()
				task.run(time.Since(task.enqueuedAt))
			}()
		}
		scheduler.complete(key)
	}
}

func (scheduler *keyedFanoutScheduler) take(key string) queuedFanoutTask {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()
	tasks := scheduler.queues[key]
	if len(tasks) == 0 {
		return queuedFanoutTask{}
	}
	task := tasks[0]
	if len(tasks) == 1 {
		scheduler.queues[key] = nil
	} else {
		tasks[0] = queuedFanoutTask{}
		scheduler.queues[key] = tasks[1:]
	}
	scheduler.running++
	return task
}

func (scheduler *keyedFanoutScheduler) complete(key string) {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()
	scheduler.running--
	scheduler.pending--
	if len(scheduler.queues[key]) > 0 {
		scheduler.ready <- key
	} else {
		delete(scheduler.queues, key)
		delete(scheduler.scheduled, key)
	}
	scheduler.notFull.Signal()
}
