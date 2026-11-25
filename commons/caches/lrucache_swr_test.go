package caches

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetOrLoadRequiresEnable(t *testing.T) {
	cache := NewLruCache("swr_disabled", 10, nil)
	if _, _, err := cache.GetOrLoad(context.Background(), "key"); !errors.Is(err, ErrStaleWhileRevalidateNotEnabled) {
		t.Fatalf("expected not-enabled error, got %v", err)
	}
}

func TestStaleWhileRevalidateReturnsStaleAndRefreshesAsynchronously(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	cache := newStaleWhileRevalidateTestCache(t, "swr_async", 2, 4, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		close(started)
		<-release
		return "new", true, nil
	})
	cache.Add("key", "old")
	ageStaleWhileRevalidateTestEntry(t, cache, "key", 2*time.Second)

	value, found, err := cache.GetOrLoad(context.Background(), "key")
	if err != nil || !found || value != "old" {
		t.Fatalf("expected immediate stale value, value=%v found=%v err=%v", value, found, err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("asynchronous refresh did not start")
	}
	close(release)
	waitForStaleWhileRevalidateValue(t, cache, "key", "new")

	stats, ok := cache.StaleWhileRevalidateStats()
	if !ok || stats.StaleHits != 1 || stats.LoadsStarted != 1 || stats.LoadsSucceeded != 1 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestStaleWhileRevalidateCoalescesLoadsForSameKey(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var loads atomic.Int32
	cache := newStaleWhileRevalidateTestCache(t, "swr_coalesce", 4, 8, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		if loads.Add(1) == 1 {
			close(started)
		}
		<-release
		return "new", true, nil
	})
	cache.Add("key", "old")
	ageStaleWhileRevalidateTestEntry(t, cache, "key", 2*time.Second)

	value, found, err := cache.GetOrLoad(context.Background(), "key")
	if err != nil || !found || value != "old" {
		t.Fatalf("unexpected first stale read: value=%v found=%v err=%v", value, found, err)
	}
	<-started

	const readers = 32
	var wg sync.WaitGroup
	wg.Add(readers)
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			value, found, err := cache.GetOrLoad(context.Background(), "key")
			if err != nil || !found || value != "old" {
				t.Errorf("unexpected shared stale read: value=%v found=%v err=%v", value, found, err)
			}
		}()
	}
	wg.Wait()
	close(release)
	waitForStaleWhileRevalidateValue(t, cache, "key", "new")
	if loads.Load() != 1 {
		t.Fatalf("expected one loader execution, got %d", loads.Load())
	}
}

func TestStaleWhileRevalidateMissWaitsAndCoalesces(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var loads atomic.Int32
	cache := newStaleWhileRevalidateTestCache(t, "swr_miss", 2, 4, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		if loads.Add(1) == 1 {
			close(started)
		}
		<-release
		return "loaded", true, nil
	})

	const readers = 16
	type result struct {
		value interface{}
		found bool
		err   error
	}
	results := make(chan result, readers)
	for i := 0; i < readers; i++ {
		go func() {
			value, found, err := cache.GetOrLoad(context.Background(), "missing")
			results <- result{value: value, found: found, err: err}
		}()
	}
	<-started
	close(release)
	for i := 0; i < readers; i++ {
		got := <-results
		if got.err != nil || !got.found || got.value != "loaded" {
			t.Fatalf("unexpected miss result: %+v", got)
		}
	}
	if loads.Load() != 1 {
		t.Fatalf("expected one loader execution, got %d", loads.Load())
	}
}

func TestStaleWhileRevalidateHardExpirationWaitsForLoad(t *testing.T) {
	cache := newStaleWhileRevalidateTestCache(t, "swr_hard", 1, 1, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		return "new", true, nil
	})
	cache.Add("key", "old")
	ageStaleWhileRevalidateTestEntry(t, cache, "key", 2*time.Hour)

	value, found, err := cache.GetOrLoad(context.Background(), "key")
	if err != nil || !found || value != "new" {
		t.Fatalf("expected hard-expired load, value=%v found=%v err=%v", value, found, err)
	}
}

func TestStaleWhileRevalidateServesHardExpiredValueOnLoadError(t *testing.T) {
	cache := NewLruCache("swr_error", 10, nil)
	errLoad := errors.New("db unavailable")
	err := cache.EnableStaleWhileRevalidate(StaleWhileRevalidateOptions{
		SoftTTL:            time.Second,
		HardTTL:            time.Minute,
		LoadTimeout:        time.Second,
		MaxConcurrentLoads: 1,
		ServeStaleOnError:  true,
	}, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		return nil, false, errLoad
	})
	if err != nil {
		t.Fatal(err)
	}
	cache.Add("key", "old")
	ageStaleWhileRevalidateTestEntry(t, cache, "key", 2*time.Hour)

	value, found, err := cache.GetOrLoad(context.Background(), "key")
	if err != nil || !found || value != "old" {
		t.Fatalf("expected stale fallback, value=%v found=%v err=%v", value, found, err)
	}
	stats, _ := cache.StaleWhileRevalidateStats()
	if stats.LoadsFailed != 1 || stats.StaleServedOnError != 1 {
		t.Fatalf("unexpected error stats: %+v", stats)
	}
}

func TestStaleWhileRevalidateDoesNotOverwriteConcurrentAdd(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	cache := newStaleWhileRevalidateTestCache(t, "swr_cas", 1, 1, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		close(started)
		<-release
		return "loaded", true, nil
	})
	cache.Add("key", "old")
	ageStaleWhileRevalidateTestEntry(t, cache, "key", 2*time.Second)
	if value, found, err := cache.GetOrLoad(context.Background(), "key"); err != nil || !found || value != "old" {
		t.Fatalf("unexpected stale read: value=%v found=%v err=%v", value, found, err)
	}
	<-started
	cache.Add("key", "manual")
	close(release)
	waitForStaleWhileRevalidateIdle(t, cache)
	value, found := cache.Peek("key")
	if !found || value != "manual" {
		t.Fatalf("refresh overwrote concurrent add: value=%v found=%v", value, found)
	}
	stats, _ := cache.StaleWhileRevalidateStats()
	if stats.ReplacementsDiscarded != 1 {
		t.Fatalf("expected discarded replacement, got %+v", stats)
	}
}

func TestStaleWhileRevalidatePurgePreventsRepopulation(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	cache := newStaleWhileRevalidateTestCache(t, "swr_purge", 1, 1, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		close(started)
		<-release
		return "loaded", true, nil
	})
	result := make(chan error, 1)
	go func() {
		value, found, err := cache.GetOrLoad(context.Background(), "key")
		if err == nil && (found || value != nil) {
			err = errors.New("purged load unexpectedly populated cache")
		}
		result <- err
	}()
	<-started
	cache.Purge()
	close(release)
	if err := <-result; err != nil {
		t.Fatal(err)
	}
	if _, found := cache.Peek("key"); found {
		t.Fatal("load repopulated cache after purge")
	}
}

func TestStaleWhileRevalidateBoundsBackgroundRefreshes(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	cache := newStaleWhileRevalidateTestCache(t, "swr_bounded", 1, 0, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		if key == "first" {
			close(started)
			<-release
		}
		return key.(string) + "-new", true, nil
	})
	cache.Add("first", "first-old")
	cache.Add("second", "second-old")
	ageStaleWhileRevalidateTestEntry(t, cache, "first", 2*time.Second)
	ageStaleWhileRevalidateTestEntry(t, cache, "second", 2*time.Second)

	cache.GetOrLoad(context.Background(), "first")
	<-started
	value, found, err := cache.GetOrLoad(context.Background(), "second")
	if err != nil || !found || value != "second-old" {
		t.Fatalf("expected second stale value, value=%v found=%v err=%v", value, found, err)
	}
	stats, _ := cache.StaleWhileRevalidateStats()
	if stats.RefreshQueueFull != 1 {
		t.Fatalf("expected one dropped refresh, got %+v", stats)
	}
	close(release)
}

func newStaleWhileRevalidateTestCache(t *testing.T, name string, maxConcurrent, maxPending int, loader StaleWhileRevalidateLoader) *LruCache {
	t.Helper()
	cache := NewLruCache(name, 100, nil)
	err := cache.EnableStaleWhileRevalidate(StaleWhileRevalidateOptions{
		SoftTTL:            time.Second,
		HardTTL:            time.Hour,
		LoadTimeout:        time.Second,
		MaxConcurrentLoads: maxConcurrent,
		MaxPendingLoads:    maxPending,
	}, loader)
	if err != nil {
		t.Fatal(err)
	}
	return cache
}

func ageStaleWhileRevalidateTestEntry(t *testing.T, cache *LruCache, key interface{}, age time.Duration) {
	t.Helper()
	cache.lock.Lock()
	defer cache.lock.Unlock()
	raw, found := cache.lru.Peek(key)
	if !found {
		t.Fatalf("missing cache entry %v", key)
	}
	item := raw.(*lruCacheItem)
	item.addedTime = time.Now().Add(-age).UnixMilli()
	item.swrSoftExpireAt = 0
	item.swrHardExpireAt = 0
}

func waitForStaleWhileRevalidateValue(t *testing.T, cache *LruCache, key, expected interface{}) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if value, found := cache.Peek(key); found && value == expected {
			return
		}
		time.Sleep(time.Millisecond)
	}
	value, found := cache.Peek(key)
	t.Fatalf("cache value did not become %v: value=%v found=%v", expected, value, found)
}

func waitForStaleWhileRevalidateIdle(t *testing.T, cache *LruCache) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		stats, _ := cache.StaleWhileRevalidateStats()
		if stats.InFlight == 0 {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("stale-while-revalidate loader did not become idle")
}
