package caches

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrStaleWhileRevalidateNotEnabled     = errors.New("stale-while-revalidate is not enabled")
	ErrStaleWhileRevalidateAlreadyEnabled = errors.New("stale-while-revalidate is already enabled")
	ErrStaleWhileRevalidateQueueFull      = errors.New("stale-while-revalidate load queue is full")
)

// StaleWhileRevalidateLoader loads one cache value. found=false represents a
// value that no longer exists and removes an unchanged stale entry.
type StaleWhileRevalidateLoader func(ctx context.Context, key interface{}) (value interface{}, found bool, err error)

// StaleWhileRevalidateApplyGuard serializes replacement of a loaded value
// with business mutations for the same key. The returned function releases
// the guard. It should only protect the short cache replacement step; the DB
// load always runs before the guard is acquired.
type StaleWhileRevalidateApplyGuard func(key interface{}) (release func())

// StaleWhileRevalidateOptions controls the optional stale-while-revalidate
// read path. SoftTTL marks a value stale. HardTTL stops serving it without a
// successful synchronous load, unless ServeStaleOnError is enabled.
//
// MaxConcurrentLoads limits loader executions. MaxPendingLoads bounds loads
// waiting for an execution slot. A stale read never waits for queue capacity;
// it returns the stale value and skips that refresh when the queue is full.
type StaleWhileRevalidateOptions struct {
	SoftTTL            time.Duration
	HardTTL            time.Duration
	RefreshJitter      time.Duration
	LoadTimeout        time.Duration
	MaxConcurrentLoads int
	MaxPendingLoads    int
	ServeStaleOnError  bool
	ApplyGuard         StaleWhileRevalidateApplyGuard
}

// StaleWhileRevalidateStats is a process-local snapshot for one cache.
type StaleWhileRevalidateStats struct {
	FreshHits             uint64 `json:"fresh_hits"`
	StaleHits             uint64 `json:"stale_hits"`
	Misses                uint64 `json:"misses"`
	HardExpired           uint64 `json:"hard_expired"`
	LoadsStarted          uint64 `json:"loads_started"`
	LoadsShared           uint64 `json:"loads_shared"`
	LoadsSucceeded        uint64 `json:"loads_succeeded"`
	LoadsNotFound         uint64 `json:"loads_not_found"`
	LoadsFailed           uint64 `json:"loads_failed"`
	RefreshQueueFull      uint64 `json:"refresh_queue_full"`
	ReplacementsDiscarded uint64 `json:"replacements_discarded"`
	StaleServedOnError    uint64 `json:"stale_served_on_error"`
	InFlight              int    `json:"in_flight"`
	Running               int    `json:"running"`
	Pending               int    `json:"pending"`
	MaxConcurrentLoads    int    `json:"max_concurrent_loads"`
	MaxPendingLoads       int    `json:"max_pending_loads"`
}

type staleWhileRevalidateCounters struct {
	freshHits             atomic.Uint64
	staleHits             atomic.Uint64
	misses                atomic.Uint64
	hardExpired           atomic.Uint64
	loadsStarted          atomic.Uint64
	loadsShared           atomic.Uint64
	loadsSucceeded        atomic.Uint64
	loadsNotFound         atomic.Uint64
	loadsFailed           atomic.Uint64
	refreshQueueFull      atomic.Uint64
	replacementsDiscarded atomic.Uint64
	staleServedOnError    atomic.Uint64
}

type staleWhileRevalidateState struct {
	options StaleWhileRevalidateOptions
	loader  StaleWhileRevalidateLoader

	// capacity bounds running and pending load goroutines; workers bounds the
	// subset that may execute the loader concurrently.
	capacity chan struct{}
	workers  chan struct{}

	lock     sync.Mutex
	inFlight map[interface{}]*staleWhileRevalidateCall
	counters staleWhileRevalidateCounters
}

type staleWhileRevalidateCall struct {
	done       chan struct{}
	key        interface{}
	expected   *lruCacheItem
	purgeEpoch uint64
	baseCtx    context.Context

	value interface{}
	found bool
	err   error
}

type staleWhileRevalidateFreshness uint8

const (
	staleWhileRevalidateFresh staleWhileRevalidateFreshness = iota
	staleWhileRevalidateStale
	staleWhileRevalidateHardExpired
)

// EnableStaleWhileRevalidate enables the optional read-through path used by
// GetOrLoad. Existing Get, Peek, and GetByCreator behavior is unchanged.
//
// Values should be treated as immutable while this mode is enabled. Add,
// Remove, eviction, and Purge are detected while a refresh is in flight, but
// mutations made inside a cached pointer cannot be observed by the cache.
func (c *LruCache) EnableStaleWhileRevalidate(options StaleWhileRevalidateOptions, loader StaleWhileRevalidateLoader) error {
	if loader == nil {
		return errors.New("stale-while-revalidate loader is nil")
	}
	if options.SoftTTL < time.Millisecond {
		return errors.New("stale-while-revalidate SoftTTL must be at least 1ms")
	}
	if options.RefreshJitter < 0 {
		return errors.New("stale-while-revalidate RefreshJitter cannot be negative")
	}
	if options.HardTTL <= options.SoftTTL+options.RefreshJitter {
		return errors.New("stale-while-revalidate HardTTL must be greater than SoftTTL plus RefreshJitter")
	}
	if options.LoadTimeout <= 0 {
		options.LoadTimeout = 10 * time.Second
	}
	if options.MaxConcurrentLoads <= 0 {
		options.MaxConcurrentLoads = 4
	}
	if options.MaxPendingLoads < 0 {
		return errors.New("stale-while-revalidate MaxPendingLoads cannot be negative")
	}

	state := &staleWhileRevalidateState{
		options:  options,
		loader:   loader,
		capacity: make(chan struct{}, options.MaxConcurrentLoads+options.MaxPendingLoads),
		workers:  make(chan struct{}, options.MaxConcurrentLoads),
		inFlight: make(map[interface{}]*staleWhileRevalidateCall),
	}
	c.swrLock.Lock()
	defer c.swrLock.Unlock()
	if c.swr != nil {
		return ErrStaleWhileRevalidateAlreadyEnabled
	}
	c.swr = state
	return nil
}

// GetOrLoad returns a fresh cached value, serves a stale value while scheduling
// an asynchronous refresh, or synchronously waits for the loader on a miss or
// hard expiration. Concurrent loads for the same key are coalesced per cache.
func (c *LruCache) GetOrLoad(ctx context.Context, key interface{}) (interface{}, bool, error) {
	state := c.staleWhileRevalidateState()
	if state == nil {
		return nil, false, ErrStaleWhileRevalidateNotEnabled
	}
	if ctx == nil {
		ctx = context.Background()
	}

	for {
		item, purgeEpoch, exists, freshness := c.getStaleWhileRevalidateItem(key, state)
		if exists {
			switch freshness {
			case staleWhileRevalidateFresh:
				state.counters.freshHits.Add(1)
				return item.value, true, nil
			case staleWhileRevalidateStale:
				state.counters.staleHits.Add(1)
				state.startLoad(c, ctx, key, item, purgeEpoch, false)
				return item.value, true, nil
			case staleWhileRevalidateHardExpired:
				state.counters.hardExpired.Add(1)
			}
		} else {
			state.counters.misses.Add(1)
		}

		call, err := state.startLoad(c, ctx, key, item, purgeEpoch, true)
		if err != nil {
			return nil, false, err
		}
		select {
		case <-call.done:
			if call.err != nil {
				if exists && state.options.ServeStaleOnError {
					state.counters.staleServedOnError.Add(1)
					return item.value, true, nil
				}
				return nil, false, call.err
			}
			return call.value, call.found, nil
		case <-ctx.Done():
			return nil, false, ctx.Err()
		}
	}
}

func (c *LruCache) staleWhileRevalidateState() *staleWhileRevalidateState {
	c.swrLock.RLock()
	defer c.swrLock.RUnlock()
	return c.swr
}

func (c *LruCache) getStaleWhileRevalidateItem(key interface{}, state *staleWhileRevalidateState) (*lruCacheItem, uint64, bool, staleWhileRevalidateFreshness) {
	c.lock.Lock()
	defer c.lock.Unlock()
	raw, exists := c.lru.Get(key)
	if !exists {
		return nil, c.purgeEpoch, false, staleWhileRevalidateHardExpired
	}
	item := raw.(*lruCacheItem)
	now := time.Now().UnixMilli()
	item.lastReadTime = now
	if item.swrSoftExpireAt == 0 {
		jitter := int64(0)
		if maxJitterMillis := state.options.RefreshJitter.Milliseconds(); maxJitterMillis > 0 {
			jitter = rand.Int63n(maxJitterMillis + 1)
		}
		item.swrSoftExpireAt = item.addedTime + state.options.SoftTTL.Milliseconds() + jitter
		item.swrHardExpireAt = item.addedTime + state.options.HardTTL.Milliseconds()
	}
	if now < item.swrSoftExpireAt {
		return item, c.purgeEpoch, true, staleWhileRevalidateFresh
	}
	if now < item.swrHardExpireAt {
		return item, c.purgeEpoch, true, staleWhileRevalidateStale
	}
	return item, c.purgeEpoch, true, staleWhileRevalidateHardExpired
}

func (state *staleWhileRevalidateState) startLoad(cache *LruCache, ctx context.Context, key interface{}, expected *lruCacheItem, purgeEpoch uint64, waitForCapacity bool) (*staleWhileRevalidateCall, error) {
	state.lock.Lock()
	if current := state.inFlight[key]; current != nil {
		state.counters.loadsShared.Add(1)
		state.lock.Unlock()
		return current, nil
	}
	// A previous coalesced load or a business write may have replaced the
	// value after GetOrLoad inspected it but before it reached this lock. Avoid
	// issuing another load from that obsolete observation.
	if waitForCapacity {
		value, found, unchanged := cache.staleWhileRevalidateExpectedCurrent(key, expected, purgeEpoch)
		if !unchanged {
			call := &staleWhileRevalidateCall{
				done:  make(chan struct{}),
				key:   key,
				value: value,
				found: found,
			}
			close(call.done)
			state.lock.Unlock()
			return call, nil
		}
	}
	call := &staleWhileRevalidateCall{
		done:       make(chan struct{}),
		key:        key,
		expected:   expected,
		purgeEpoch: purgeEpoch,
		baseCtx:    ctx,
	}
	state.inFlight[key] = call
	state.lock.Unlock()

	if waitForCapacity {
		select {
		case state.capacity <- struct{}{}:
		case <-ctx.Done():
			state.failUnscheduledCall(call, ctx.Err())
			return nil, ctx.Err()
		}
	} else {
		select {
		case state.capacity <- struct{}{}:
		default:
			state.counters.refreshQueueFull.Add(1)
			state.failUnscheduledCall(call, ErrStaleWhileRevalidateQueueFull)
			return nil, ErrStaleWhileRevalidateQueueFull
		}
	}

	state.counters.loadsStarted.Add(1)
	go state.executeLoad(cache, call)
	return call, nil
}

func (c *LruCache) staleWhileRevalidateExpectedCurrent(key interface{}, expected *lruCacheItem, purgeEpoch uint64) (interface{}, bool, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	raw, exists := c.lru.Peek(key)
	if c.purgeEpoch != purgeEpoch {
		value, found, _ := staleWhileRevalidateCurrentValue(raw, exists)
		return value, found, false
	}
	if expected == nil {
		if !exists {
			return nil, false, true
		}
		return raw.(*lruCacheItem).value, true, false
	}
	if exists && raw.(*lruCacheItem) == expected {
		return expected.value, true, true
	}
	value, found, _ := staleWhileRevalidateCurrentValue(raw, exists)
	return value, found, false
}

func (state *staleWhileRevalidateState) failUnscheduledCall(call *staleWhileRevalidateCall, err error) {
	state.lock.Lock()
	if state.inFlight[call.key] == call {
		delete(state.inFlight, call.key)
		call.err = err
		close(call.done)
	}
	state.lock.Unlock()
}

func (state *staleWhileRevalidateState) executeLoad(cache *LruCache, call *staleWhileRevalidateCall) {
	state.workers <- struct{}{}
	defer func() {
		<-state.workers
		<-state.capacity
	}()

	baseCtx := call.baseCtx
	if baseCtx == nil {
		baseCtx = context.Background()
	}
	loadCtx, cancel := context.WithTimeout(context.WithoutCancel(baseCtx), state.options.LoadTimeout)
	value, found, err := state.loadSafely(loadCtx, call.key)
	cancel()

	if err == nil {
		value, found, discarded, applyErr := state.applyLoadSafely(cache, call, value, found)
		if applyErr != nil {
			call.err = applyErr
			state.counters.loadsFailed.Add(1)
			state.finishLoad(call)
			return
		}
		if discarded {
			state.counters.replacementsDiscarded.Add(1)
		}
		call.value = value
		call.found = found
		if found {
			state.counters.loadsSucceeded.Add(1)
		} else {
			state.counters.loadsNotFound.Add(1)
		}
	} else {
		call.err = err
		state.counters.loadsFailed.Add(1)
	}

	state.finishLoad(call)
}

func (state *staleWhileRevalidateState) loadSafely(ctx context.Context, key interface{}) (value interface{}, found bool, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			value = nil
			found = false
			err = fmt.Errorf("stale-while-revalidate loader panic: %v", recovered)
		}
	}()
	return state.loader(ctx, key)
}

func (state *staleWhileRevalidateState) applyLoadSafely(cache *LruCache, call *staleWhileRevalidateCall, value interface{}, found bool) (appliedValue interface{}, appliedFound bool, discarded bool, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			appliedValue = nil
			appliedFound = false
			discarded = false
			err = fmt.Errorf("stale-while-revalidate apply guard panic: %v", recovered)
		}
	}()
	if state.options.ApplyGuard != nil {
		if release := state.options.ApplyGuard(call.key); release != nil {
			defer release()
		}
	}
	appliedValue, appliedFound, discarded = cache.applyStaleWhileRevalidateLoad(call, value, found)
	return appliedValue, appliedFound, discarded, nil
}

func (state *staleWhileRevalidateState) finishLoad(call *staleWhileRevalidateCall) {
	state.lock.Lock()
	if state.inFlight[call.key] == call {
		delete(state.inFlight, call.key)
	}
	close(call.done)
	state.lock.Unlock()
}

func (c *LruCache) applyStaleWhileRevalidateLoad(call *staleWhileRevalidateCall, value interface{}, found bool) (interface{}, bool, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	currentRaw, currentExists := c.lru.Peek(call.key)
	if c.purgeEpoch != call.purgeEpoch {
		return staleWhileRevalidateCurrentValue(currentRaw, currentExists)
	}
	if call.expected == nil {
		if currentExists {
			current := currentRaw.(*lruCacheItem)
			return current.value, true, true
		}
	} else if !currentExists || currentRaw.(*lruCacheItem) != call.expected {
		value, exists, _ := staleWhileRevalidateCurrentValue(currentRaw, currentExists)
		return value, exists, true
	}

	if !found {
		if currentExists {
			c.lru.Remove(call.key)
		}
		return nil, false, false
	}
	c.innerAdd(call.key, value)
	return value, true, false
}

func staleWhileRevalidateCurrentValue(raw interface{}, exists bool) (interface{}, bool, bool) {
	if !exists {
		return nil, false, true
	}
	return raw.(*lruCacheItem).value, true, true
}

// StaleWhileRevalidateStats returns the optional loader's current metrics.
func (c *LruCache) StaleWhileRevalidateStats() (StaleWhileRevalidateStats, bool) {
	state := c.staleWhileRevalidateState()
	if state == nil {
		return StaleWhileRevalidateStats{}, false
	}
	state.lock.Lock()
	inFlight := len(state.inFlight)
	state.lock.Unlock()
	running := len(state.workers)
	pending := len(state.capacity) - running
	if pending < 0 {
		pending = 0
	}
	return StaleWhileRevalidateStats{
		FreshHits:             state.counters.freshHits.Load(),
		StaleHits:             state.counters.staleHits.Load(),
		Misses:                state.counters.misses.Load(),
		HardExpired:           state.counters.hardExpired.Load(),
		LoadsStarted:          state.counters.loadsStarted.Load(),
		LoadsShared:           state.counters.loadsShared.Load(),
		LoadsSucceeded:        state.counters.loadsSucceeded.Load(),
		LoadsNotFound:         state.counters.loadsNotFound.Load(),
		LoadsFailed:           state.counters.loadsFailed.Load(),
		RefreshQueueFull:      state.counters.refreshQueueFull.Load(),
		ReplacementsDiscarded: state.counters.replacementsDiscarded.Load(),
		StaleServedOnError:    state.counters.staleServedOnError.Load(),
		InFlight:              inFlight,
		Running:               running,
		Pending:               pending,
		MaxConcurrentLoads:    state.options.MaxConcurrentLoads,
		MaxPendingLoads:       state.options.MaxPendingLoads,
	}, true
}
