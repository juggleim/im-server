package services

import (
	"im-server/services/commonservices/logs"
	"sync"
	"time"
)

type purgeFunc func(cutoff int64) error

type purgeTaskState struct {
	latestCutoff   int64
	lastSubmitTime int64
	queued         bool
	running        bool
	f              purgeFunc
}

// purgeScheduler bounds both execution concurrency and the number of queued
// keys. Repeated submissions for the same key are merged by retaining the
// greatest cutoff; a key can never run concurrently with itself.
type purgeScheduler struct {
	mu             sync.Mutex
	interval       int64
	requeueReserve int
	tasks          map[string]*purgeTaskState
	queue          chan string
}

func newPurgeScheduler(workerCount, queueCapacity int, interval time.Duration) *purgeScheduler {
	s := &purgeScheduler{
		interval:       interval.Milliseconds(),
		requeueReserve: workerCount,
		tasks:          make(map[string]*purgeTaskState),
		queue:          make(chan string, queueCapacity),
	}
	for i := 0; i < workerCount; i++ {
		go s.worker()
	}
	return s
}

func (s *purgeScheduler) Submit(key string, currentTime, cutoff int64, f purgeFunc) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.tasks[key]
	if state == nil {
		// Preserve the existing behavior: the first message establishes the
		// throttle timestamp and a later message triggers the first purge.
		s.tasks[key] = &purgeTaskState{
			latestCutoff:   cutoff,
			lastSubmitTime: currentTime,
			f:              f,
		}
		return false
	}

	if cutoff > state.latestCutoff {
		state.latestCutoff = cutoff
		state.f = f
	}
	if state.queued || state.running || currentTime-state.lastSubmitTime <= s.interval {
		return false
	}

	previousSubmitTime := state.lastSubmitTime
	state.lastSubmitTime = currentTime
	state.queued = true
	// Keep one queue slot per worker available for a key whose cutoff advances
	// while it is running. This makes requeueing lossless without an auxiliary
	// goroutine or an unbounded overflow queue.
	if len(s.queue) >= cap(s.queue)-s.requeueReserve {
		state.queued = false
		state.lastSubmitTime = previousSubmitTime
		return false
	}
	select {
	case s.queue <- key:
		return true
	default:
		// Do not block the message save path. Roll back the throttle time so
		// the next message can retry submission immediately.
		state.queued = false
		state.lastSubmitTime = previousSubmitTime
		return false
	}
}

func (s *purgeScheduler) worker() {
	for key := range s.queue {
		s.run(key)
	}
}

func (s *purgeScheduler) run(key string) {
	s.mu.Lock()
	state := s.tasks[key]
	if state == nil || !state.queued {
		s.mu.Unlock()
		return
	}
	state.queued = false
	state.running = true
	cutoff := state.latestCutoff
	f := state.f
	s.mu.Unlock()

	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				logs.NewLogEntity().Errorf("panic while purging offline messages. key:%s err:%v", key, recovered)
			}
		}()
		if f != nil {
			_ = f(cutoff)
		}
	}()

	s.mu.Lock()
	state = s.tasks[key]
	state.running = false
	if state.latestCutoff > cutoff {
		// A newer cutoff arrived while this key was running. Requeue exactly
		// once so the latest request is not lost and other keys get a turn.
		state.queued = true
		s.queue <- key
	}
	s.mu.Unlock()
}
