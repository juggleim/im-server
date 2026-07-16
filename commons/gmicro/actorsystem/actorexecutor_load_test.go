package actorsystem

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

type blockingLoadActor struct {
	release   <-chan struct{}
	completed *sync.WaitGroup
}

func (a *blockingLoadActor) OnReceive(context.Context, proto.Message) {
	<-a.release
	a.completed.Done()
}

type concurrentLoadActor struct {
	current   *atomic.Int32
	maximum   *atomic.Int32
	completed *sync.WaitGroup
}

func (a *concurrentLoadActor) OnReceive(context.Context, proto.Message) {
	current := a.current.Add(1)
	for {
		maximum := a.maximum.Load()
		if current <= maximum || a.maximum.CompareAndSwap(maximum, current) {
			break
		}
	}
	time.Sleep(5 * time.Millisecond)
	a.current.Add(-1)
	a.completed.Done()
}

func TestActorExecutorRetainsConfiguredConcurrency(t *testing.T) {
	const (
		workers = 8
		jobs    = 80
	)
	var current atomic.Int32
	var maximum atomic.Int32
	var completed sync.WaitGroup
	completed.Add(jobs)
	started := time.Now()
	executor := NewActorExecutor(workers, func() IUntypedActor {
		return &concurrentLoadActor{current: &current, maximum: &maximum, completed: &completed}
	})

	for i := 0; i < jobs; i++ {
		executor.wraperChan <- wraper{}
	}
	done := make(chan struct{})
	go func() {
		completed.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for actor jobs")
	}

	if got := maximum.Load(); got != workers {
		t.Fatalf("configured concurrency was not retained: want=%d got=%d", workers, got)
	}
	t.Logf("concurrency load: jobs=%d workers=%d elapsed=%s", jobs, workers, time.Since(started))
}

func TestActorDispatcherDefaultPoolFootprint(t *testing.T) {
	t.Setenv("IM_ACTOR_CALLBACK_WORKERS", "")
	t.Setenv("IM_ACTOR_EXECUTOR_WORKERS", "")
	baseline := runtime.NumGoroutine()
	dispatcher := NewActorDispatcher(nil)
	t.Cleanup(func() {
		dispatcher.timer.Stop()
		dispatcher.callbackPool.Close()
		dispatcher.executorCommonPool.Close()
	})

	deadline := time.Now().Add(time.Second)
	wantAtLeast := defaultCallbackPoolSize + defaultExecutorCommonPoolSize
	for runtime.NumGoroutine()-baseline < wantAtLeast && time.Now().Before(deadline) {
		runtime.Gosched()
	}
	growth := runtime.NumGoroutine() - baseline
	t.Logf("default dispatcher goroutines: baseline=%d growth=%d", baseline, growth)
	if growth > wantAtLeast+16 {
		t.Fatalf("default dispatcher created too many goroutines: workers=%d growth=%d", wantAtLeast, growth)
	}
}

// TestActorExecutorSaturationIsGoroutineBounded simulates a slow downstream
// dependency while a burst fills the actor queue. Queued work must stay in the
// bounded channel instead of creating one blocked goroutine per request.
func TestActorExecutorSaturationIsGoroutineBounded(t *testing.T) {
	release := make(chan struct{})
	const burst = 2000
	var completed sync.WaitGroup
	completed.Add(burst + 1)
	executor := NewActorExecutor(1, func() IUntypedActor {
		return &blockingLoadActor{release: release, completed: &completed}
	})

	executor.wraperChan <- wraper{}
	time.Sleep(20 * time.Millisecond)

	baseline := runtime.NumGoroutine()
	for i := 0; i < burst; i++ {
		executor.wraperChan <- wraper{}
	}
	time.Sleep(100 * time.Millisecond)
	peak := runtime.NumGoroutine()
	t.Logf("saturation goroutines: baseline=%d peak=%d growth=%d burst=%d", baseline, peak, peak-baseline, burst)

	close(release)
	if growth := peak - baseline; growth > 16 {
		t.Fatalf("goroutine count grew with queued requests: baseline=%d peak=%d growth=%d", baseline, peak, growth)
	}
	drained := make(chan struct{})
	go func() {
		completed.Wait()
		close(drained)
	}()
	select {
	case <-drained:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out draining saturated actor queue")
	}
}
