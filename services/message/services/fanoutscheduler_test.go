package services

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestKeyedFanoutSchedulerPreservesKeyOrderAndRunsDifferentKeysConcurrently(t *testing.T) {
	scheduler := newKeyedFanoutScheduler(2, 10)
	releaseFirst := make(chan struct{})
	firstStarted := make(chan struct{})
	otherKeyRan := make(chan struct{})
	ordered := make(chan int, 2)

	scheduler.Submit("group-a", func(time.Duration) {
		close(firstStarted)
		<-releaseFirst
		ordered <- 1
	})
	<-firstStarted
	scheduler.Submit("group-a", func(time.Duration) { ordered <- 2 })
	scheduler.Submit("group-b", func(time.Duration) { close(otherKeyRan) })

	select {
	case <-otherKeyRan:
	case <-time.After(time.Second):
		t.Fatal("different key did not run concurrently")
	}
	close(releaseFirst)
	for want := 1; want <= 2; want++ {
		select {
		case got := <-ordered:
			if got != want {
				t.Fatalf("execution order = %d, want %d", got, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for ordered task %d", want)
		}
	}
}

func TestKeyedFanoutSchedulerAppliesBackpressure(t *testing.T) {
	scheduler := newKeyedFanoutScheduler(1, 1)
	release := make(chan struct{})
	started := make(chan struct{})
	scheduler.Submit("group-a", func(time.Duration) {
		close(started)
		<-release
	})
	<-started

	submitReturned := make(chan bool)
	secondRan := make(chan struct{})
	go func() {
		_, blocked := scheduler.Submit("group-b", func(time.Duration) { close(secondRan) })
		submitReturned <- blocked
	}()
	select {
	case <-submitReturned:
		t.Fatal("submit returned while scheduler was full")
	case <-time.After(30 * time.Millisecond):
	}

	close(release)
	select {
	case blocked := <-submitReturned:
		if !blocked {
			t.Fatal("submit did not report backpressure")
		}
	case <-time.After(time.Second):
		t.Fatal("submit did not resume after capacity became available")
	}
	select {
	case <-secondRan:
	case <-time.After(time.Second):
		t.Fatal("second task did not run")
	}
}

func TestKeyedFanoutSchedulerContinuesAfterPanic(t *testing.T) {
	scheduler := newKeyedFanoutScheduler(1, 2)
	done := make(chan struct{})
	var once sync.Once
	scheduler.Submit("group-a", func(time.Duration) { panic("test") })
	scheduler.Submit("group-a", func(time.Duration) { once.Do(func() { close(done) }) })
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("scheduler stopped after task panic")
	}
}

func TestProcessDispatchMemberBatchesProcessesEachMemberWithinLimit(t *testing.T) {
	members := make([]string, 2501)
	for i := range members {
		members[i] = string(rune(i + 1))
	}
	seen := make(map[string]int, len(members))
	var seenMu sync.Mutex
	var active atomic.Int32
	var maxActive atomic.Int32

	processDispatchMemberBatches(context.Background(), members, func(receiverId string) {
		current := active.Add(1)
		for old := maxActive.Load(); current > old; old = maxActive.Load() {
			if maxActive.CompareAndSwap(old, current) {
				break
			}
		}
		time.Sleep(100 * time.Microsecond)
		seenMu.Lock()
		seen[receiverId]++
		seenMu.Unlock()
		active.Add(-1)
	})

	if len(seen) != len(members) {
		t.Fatalf("processed members = %d, want %d", len(seen), len(members))
	}
	for _, receiverId := range members {
		if seen[receiverId] != 1 {
			t.Fatalf("member %q processed %d times, want 1", receiverId, seen[receiverId])
		}
	}
	if got := maxActive.Load(); got < 2 || got > msgDispatchLocalBatchConcurrency {
		t.Fatalf("max concurrency = %d, want between 2 and %d", got, msgDispatchLocalBatchConcurrency)
	}
}
