package caches

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Clear must evict every buffered item (invoking onEvict for each) and leave the
// cache empty. This is what lets graceful shutdown flush pending connect-log
// items to the database instead of losing them when the process exits before
// their normal lifetime elapses.
func TestEphemeralCache_ClearEvictsAllPendingItems(t *testing.T) {
	var mu sync.Mutex
	evicted := map[interface{}]interface{}{}
	// Long lifetime so nothing evicts on its own during the test.
	c := NewEphemeralCache(time.Second, time.Hour, func(key, value interface{}) {
		mu.Lock()
		evicted[key] = value
		mu.Unlock()
	})

	keys := []string{"a", "b", "c"}
	for _, k := range keys {
		c.Add(k, "v-"+k)
	}
	if l, _ := c.Len(); l != len(keys) {
		t.Fatalf("expected %d items before Clear, got %d", len(keys), l)
	}

	c.Clear()

	if l, x := c.Len(); l != 0 || x != 0 {
		t.Fatalf("cache should be empty after Clear, got items=%d evictList=%d", l, x)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(evicted) != len(keys) {
		t.Fatalf("Clear should have evicted %d items, got %d", len(keys), len(evicted))
	}
	for _, k := range keys {
		if evicted[k] != "v-"+k {
			t.Fatalf("missing/incorrect eviction for key %q: %v", k, evicted[k])
		}
	}
}

func TestEphemeralCache_ClearOnEmptyIsNoop(t *testing.T) {
	calls := 0
	c := NewEphemeralCache(time.Second, time.Hour, func(key, value interface{}) { calls++ })
	c.Clear()
	if calls != 0 {
		t.Fatalf("Clear on empty cache should not invoke onEvict, got %d calls", calls)
	}
	if l, _ := c.Len(); l != 0 {
		t.Fatalf("empty cache should stay empty")
	}
}

func TestEphemeralCache_ConcurrentMutationAndExpiration(t *testing.T) {
	var evicted atomic.Int64
	c := NewEphemeralCache(time.Millisecond, 3*time.Millisecond, func(key, value interface{}) {
		evicted.Add(1)
	})

	const (
		workerCount = 8
		operations  = 300
	)
	var workers sync.WaitGroup
	for worker := 0; worker < workerCount; worker++ {
		workers.Add(1)
		go func(worker int) {
			defer workers.Done()
			for operation := 0; operation < operations; operation++ {
				key := fmt.Sprintf("%d-%d", worker, operation)
				c.Add(key, operation)
				c.Upsert(key, func(oldVal interface{}) interface{} {
					return operation + 1
				})
				if operation%3 == 0 {
					c.Remove(key)
				}
				if operation%25 == 0 {
					c.Len()
				}
			}
		}(worker)
	}

	workers.Wait()
	time.Sleep(10 * time.Millisecond)
	c.Clear()

	if items, entries := c.Len(); items != 0 || entries != 0 {
		t.Fatalf("cache should be empty after final Clear, got items=%d evictList=%d", items, entries)
	}
	if evicted.Load() == 0 {
		t.Fatal("expected concurrent removal or expiration to evict items")
	}
}

func TestEphemeralCacheRemoveCallsOnEvictOutsideLock(t *testing.T) {
	callbackDone := make(chan struct{})
	var cache *EphemeralCache
	cache = NewEphemeralCache(time.Hour, time.Hour, func(key, value interface{}) {
		cache.Add("from-callback", "value")
		close(callbackDone)
	})
	defer cache.timer.Stop()

	cache.Add("removed", "value")
	go cache.Remove("removed")

	select {
	case <-callbackDone:
	case <-time.After(time.Second):
		t.Fatal("onEvict callback was blocked while re-entering the cache")
	}

	items, entries := cache.Len()
	if items != 1 || entries != 1 {
		t.Fatalf("unexpected cache size after callback: items=%d entries=%d", items, entries)
	}
}

func TestEphemeralCacheSlowOnEvictDoesNotBlockAdd(t *testing.T) {
	callbackStarted := make(chan struct{})
	releaseCallback := make(chan struct{})
	removeDone := make(chan struct{})
	cache := NewEphemeralCache(time.Hour, time.Hour, func(key, value interface{}) {
		close(callbackStarted)
		<-releaseCallback
	})
	defer cache.timer.Stop()

	cache.Add("removed", "value")
	go func() {
		cache.Remove("removed")
		close(removeDone)
	}()

	select {
	case <-callbackStarted:
	case <-time.After(time.Second):
		t.Fatal("onEvict callback did not start")
	}

	addDone := make(chan struct{})
	go func() {
		cache.Add("new", "value")
		close(addDone)
	}()

	select {
	case <-addDone:
	case <-time.After(time.Second):
		t.Fatal("cache Add was blocked by onEvict callback")
	}

	close(releaseCallback)
	select {
	case <-removeDone:
	case <-time.After(time.Second):
		t.Fatal("Remove did not finish after onEvict callback was released")
	}
}

func TestEphemeralCacheExpirationCallsOnEvictOutsideLock(t *testing.T) {
	callbackDone := make(chan struct{})
	var cache *EphemeralCache
	cache = NewEphemeralCache(time.Hour, time.Hour, func(key, value interface{}) {
		cache.Len()
		close(callbackDone)
	})
	defer cache.timer.Stop()

	cache.Add("expired", "value")
	cache.cleanOldestByCreatedTime(time.Now().Add(time.Second).UnixMilli())

	select {
	case <-callbackDone:
	case <-time.After(time.Second):
		t.Fatal("expiration onEvict callback was blocked while re-entering the cache")
	}
}
