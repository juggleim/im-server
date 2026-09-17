package services

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestPurgeSchedulerMergesLatestCutoffWithoutConcurrentRun(t *testing.T) {
	scheduler := newPurgeScheduler(1, 10, time.Minute)
	started := make(chan int64, 2)
	release := make(chan struct{}, 2)
	var running atomic.Int32
	var maxRunning atomic.Int32

	f := func(cutoff int64) error {
		current := running.Add(1)
		for {
			old := maxRunning.Load()
			if current <= old || maxRunning.CompareAndSwap(old, current) {
				break
			}
		}
		started <- cutoff
		<-release
		running.Add(-1)
		return nil
	}

	if scheduler.Submit("app:inbox", 1, 10, f) {
		t.Fatal("first submission should only establish the throttle timestamp")
	}
	if !scheduler.Submit("app:inbox", 60_002, 20, f) {
		t.Fatal("second submission should be queued")
	}
	if cutoff := waitCutoff(t, started); cutoff != 20 {
		t.Fatalf("first cutoff = %d, want 20", cutoff)
	}

	if scheduler.Submit("app:inbox", 120_003, 30, f) {
		t.Fatal("running key must be merged instead of queued concurrently")
	}
	release <- struct{}{}
	if cutoff := waitCutoff(t, started); cutoff != 30 {
		t.Fatalf("merged cutoff = %d, want 30", cutoff)
	}
	release <- struct{}{}

	if maxRunning.Load() != 1 {
		t.Fatalf("same key ran concurrently: max=%d", maxRunning.Load())
	}
}

func TestPurgeSchedulerBoundsWorkerConcurrency(t *testing.T) {
	const workers = 2
	scheduler := newPurgeScheduler(workers, 16, 0)
	release := make(chan struct{})
	done := make(chan struct{}, 4)
	var running atomic.Int32
	var maxRunning atomic.Int32

	f := func(int64) error {
		current := running.Add(1)
		for {
			old := maxRunning.Load()
			if current <= old || maxRunning.CompareAndSwap(old, current) {
				break
			}
		}
		<-release
		running.Add(-1)
		done <- struct{}{}
		return nil
	}

	for i := 0; i < 4; i++ {
		key := string(rune('a' + i))
		scheduler.Submit(key, 0, 1, f)
		if !scheduler.Submit(key, 1, 2, f) {
			t.Fatalf("key %q was not queued", key)
		}
	}

	deadline := time.Now().Add(time.Second)
	for running.Load() != workers && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if running.Load() != workers {
		t.Fatalf("running workers = %d, want %d", running.Load(), workers)
	}
	if maxRunning.Load() > workers {
		t.Fatalf("worker concurrency = %d, want <= %d", maxRunning.Load(), workers)
	}

	close(release)
	for i := 0; i < 4; i++ {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for queued tasks")
		}
	}
}

func TestPurgeSchedulerQueueIsBounded(t *testing.T) {
	scheduler := &purgeScheduler{
		interval:       0,
		requeueReserve: 1,
		tasks:          make(map[string]*purgeTaskState),
		queue:          make(chan string, 2),
	}
	f := func(int64) error { return nil }

	scheduler.Submit("a", 0, 1, f)
	if !scheduler.Submit("a", 1, 2, f) {
		t.Fatal("first key should fit in the externally usable queue slot")
	}
	scheduler.Submit("b", 0, 1, f)
	if scheduler.Submit("b", 1, 2, f) {
		t.Fatal("second key should be rejected while the bounded queue is full")
	}
}

func waitCutoff(t *testing.T, values <-chan int64) int64 {
	t.Helper()
	select {
	case value := <-values:
		return value
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for purge task")
		return 0
	}
}
