package tools

import (
	"runtime"
	"testing"
	"time"
)

func TestBatchExecutorStopReleasesGoroutine(t *testing.T) {
	baseline := runtime.NumGoroutine()
	const count = 64
	executors := make([]*BatchExecutor, 0, count)
	for i := 0; i < count; i++ {
		executors = append(executors, NewBatchExecutor(10, time.Hour, func([]interface{}) {}))
	}

	for _, executor := range executors {
		executor.Stop()
	}

	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > baseline+4 && time.Now().Before(deadline) {
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
	if current := runtime.NumGoroutine(); current > baseline+4 {
		t.Fatalf("batch executor goroutines were not released: baseline=%d current=%d", baseline, current)
	}
}
