package tools

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSinglePoolContinuesAfterTaskPanic(t *testing.T) {
	pools := NewSinglePools(1, 100, true)
	pool := pools.GetPool("test")
	executed := make(chan struct{})

	if !pool.Submit(func() {
		panic("test panic")
	}) {
		t.Fatal("failed to submit panic task")
	}
	if !pool.Submit(func() {
		close(executed)
	}) {
		t.Fatal("failed to submit task after panic")
	}

	select {
	case <-executed:
	case <-time.After(time.Second):
		t.Fatal("single pool worker stopped after task panic")
	}
}

func TestSinglePoolsQueueStatsAndConcurrentCreation(t *testing.T) {
	pools := NewSinglePools(16, 2, true)
	started := make(chan struct{})
	release := make(chan struct{})
	pool := pools.GetPool("hot")
	pool.Submit(func() {
		close(started)
		<-release
	})
	<-started
	pool.Submit(func() {})
	pool.Submit(func() {})
	stats := pools.QueueStats()
	if stats.QueuedTasks != 2 || stats.FullPools != 1 || stats.MaxQueueLen != 2 {
		t.Fatalf("unexpected full queue stats: %+v", stats)
	}
	close(release)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pools.GetPool(fmt.Sprint(i))
			pools.QueueStats()
		}(i)
	}
	wg.Wait()
}

func TestSinglePoolsPoolIndexMatchesGetPool(t *testing.T) {
	pools := NewSinglePools(16, 2, true)
	key := "app_user"
	index := pools.PoolIndex(key)
	if index < 0 || index >= 16 {
		t.Fatalf("pool index out of range: %d", index)
	}
	if pools.GetPool(key) != pools.pools[index] {
		t.Fatalf("GetPool did not use reported pool index %d", index)
	}
}
