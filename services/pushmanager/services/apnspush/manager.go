package apnspush

import (
	"container/list"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"hash/maphash"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sideshow/apns2"
)

const (
	DefaultCapacity        = 10000
	DefaultRefreshInterval = 5 * time.Minute
	DefaultIdleTTL         = time.Hour
	DefaultPushTimeout     = 10 * time.Second
	MaxPushTimeoutSeconds  = 120
)

type cacheKey struct{ appKey, packageName string }

// No Config, passwords or PEM buffers survive in cache entries.
// Only the SDK client retains the parsed signing key required for sends.
type entry struct {
	key       cacheKey
	clients   *Clients
	digest    [sha256.Size]byte
	version   int64
	nextCheck time.Time
	lastUsed  time.Time
	err       error
	refs      int
	retired   bool
	lru       *list.Element
}

type Manager struct {
	mu      sync.Mutex
	entries map[cacheKey]*entry
	lru     list.List
	closed  bool
	stripes [128]chan struct{}
	seed    maphash.Seed
	loader  func(context.Context, string, string) (*Config, error)
	factory func(*Config) (*Clients, error)
	ctx     context.Context
	cancel  context.CancelFunc
	done    chan struct{}
	once    sync.Once

	capacity        int
	refreshInterval time.Duration
	idleTTL         time.Duration
	retryDelay      time.Duration
	cleanupInterval time.Duration
	timeout         time.Duration
	now             func() time.Time
}

type Option func(*Manager)

// NewManager starts idle cleanup but never calls loader. The loader must honor
// ctx, return a stable snapshot with matching AppKey/Package, and map absence to
// ErrNotFound. Other loader errors are treated as transient and sanitized.
func NewManager(loader func(context.Context, string, string) (*Config, error), options ...Option) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{
		entries: make(map[cacheKey]*entry), seed: maphash.MakeSeed(),
		loader: loader, factory: (Factory{}).NewClients,
		ctx: ctx, cancel: cancel, done: make(chan struct{}),
		capacity: DefaultCapacity, refreshInterval: DefaultRefreshInterval,
		idleTTL: DefaultIdleTTL, retryDelay: 5 * time.Second,
		cleanupInterval: time.Minute, timeout: PushTimeoutFromEnv(), now: time.Now,
	}
	for i := range m.stripes {
		m.stripes[i] = make(chan struct{}, 1)
	}
	for _, option := range options {
		option(m)
	}
	go m.cleanupLoop()
	return m
}

// PushTimeoutFromEnv accepts integers from 1 through 120 seconds; missing or
// invalid IM_APNS_PUSH_TIMEOUT_SECONDS falls back to ten seconds.
func PushTimeoutFromEnv() time.Duration {
	seconds, err := strconv.Atoi(os.Getenv("IM_APNS_PUSH_TIMEOUT_SECONDS"))
	if err != nil || seconds < 1 || seconds > MaxPushTimeoutSeconds {
		return DefaultPushTimeout
	}
	return time.Duration(seconds) * time.Second
}

// Send bounds the entire operation, including loading and waiting for another
// load. Close cancels it; eviction does not. APNs responses are never retried.
func (m *Manager) Send(ctx context.Context, appKey, packageName string, notification *apns2.Notification, isVoip bool, ordinaryToken, voipToken string) (*apns2.Response, error) {
	m.mu.Lock()
	closed := m.closed
	m.mu.Unlock()
	if closed {
		return nil, ErrClosed
	}
	if appKey == "" || !validTopic(packageName) || notification == nil {
		return nil, ErrConfiguration
	}
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	stop := context.AfterFunc(m.ctx, cancel)
	defer func() { stop(); cancel() }()
	e, err := m.acquire(ctx, cacheKey{appKey, packageName})
	if err != nil {
		return nil, err
	}
	defer m.release(e)
	client, routed, err := Route(e.clients, packageName, notification, isVoip, ordinaryToken, voipToken)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	response, err := client.PushWithContext(ctx, routed)
	if err != nil {
		return response, &sendError{cause: err}
	}
	return response, nil
}

// cached is called with mu held; a hit reserves the client until release.
func (m *Manager) cached(key cacheKey, now time.Time) (*entry, error, bool) {
	if m.closed {
		return nil, ErrClosed, true
	}
	e := m.entries[key]
	if e == nil {
		return nil, nil, false
	}
	e.lastUsed = now
	m.lru.MoveToFront(e.lru)
	if !now.Before(e.nextCheck) {
		return nil, nil, false
	}
	if e.err != nil {
		return nil, e.err, true
	}
	e.refs++
	return e, nil, true
}

func (m *Manager) acquire(ctx context.Context, key cacheKey) (*entry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.Lock()
	e, err, hit := m.cached(key, m.now())
	m.mu.Unlock()
	if hit {
		return e, err
	}
	stripe := m.stripes[maphash.Comparable(m.seed, key)%uint64(len(m.stripes))]
	select {
	case stripe <- struct{}{}:
		defer func() { <-stripe }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.Lock()
	e, err, hit = m.cached(key, m.now())
	m.mu.Unlock()
	if hit {
		return e, err
	}

	var cfg *Config
	if m.loader == nil {
		err = ErrConfiguration
	} else {
		cfg, err = m.loader(ctx, key.appKey, key.packageName)
	}
	// A canceled loader must not poison a healthy entry for other callers.
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if err == nil && cfg == nil {
		err = ErrNotFound
	}
	var digest [sha256.Size]byte
	if err == nil {
		if cfg.AppKey != key.appKey || cfg.Package != key.packageName {
			err = ErrConfiguration
		} else {
			digest = configDigest(cfg)
		}
	}
	m.mu.Lock()
	old := m.entries[key]
	if m.closed {
		m.mu.Unlock()
		return nil, ErrClosed
	}
	if err == nil && old != nil && old.clients != nil && old.digest == digest && old.version == cfg.ConfigVersion {
		old.nextCheck, old.lastUsed, old.err = m.now().Add(m.refreshInterval), m.now(), nil
		old.refs++
		m.lru.MoveToFront(old.lru)
		m.mu.Unlock()
		return old, nil
	}
	m.mu.Unlock()

	var clients *Clients
	if err == nil {
		clients, err = m.factory(cfg)
	}
	// Never keep arbitrary DB or injected helper errors in the cache.
	transient := false
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			err = ErrNotFound
		case errors.Is(err, ErrConfiguration):
			err = ErrConfiguration
		default:
			err, transient = ErrLoad, true
		}
	}
	m.mu.Lock()
	if m.closed || ctx.Err() != nil {
		closed := m.closed
		m.mu.Unlock()
		clients.Close()
		if closed {
			return nil, ErrClosed
		}
		return nil, ctx.Err()
	}
	now := m.now()
	old = m.entries[key]
	if transient && old != nil {
		old.err, old.nextCheck, old.lastUsed = err, now.Add(m.retryDelay), now
		m.lru.MoveToFront(old.lru)
		m.mu.Unlock()
		return nil, err
	}
	var retired []*Clients
	if old != nil {
		retired = append(retired, m.retire(old))
	}
	for len(m.entries) >= m.capacity {
		retired = append(retired, m.retire(m.lru.Back().Value.(*entry)))
	}
	e = &entry{key: key, clients: clients, digest: digest, lastUsed: now, err: err, nextCheck: now.Add(m.refreshInterval)}
	if cfg != nil {
		e.version = cfg.ConfigVersion
	}
	if transient {
		e.nextCheck = now.Add(m.retryDelay)
	}
	if err == nil {
		e.refs = 1
	}
	e.lru = m.lru.PushFront(e)
	m.entries[key] = e
	m.mu.Unlock()
	for _, client := range retired {
		client.Close()
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

// retire runs under mu. Active users own the final close, so even HTTP/2
// connections that become idle after eviction are eventually released.
func (m *Manager) retire(e *entry) *Clients {
	delete(m.entries, e.key)
	m.lru.Remove(e.lru)
	e.retired = true
	if e.refs == 0 {
		return e.clients
	}
	return nil
}

func (m *Manager) release(e *entry) {
	m.mu.Lock()
	e.refs--
	closeClients := e.retired && e.refs == 0
	m.mu.Unlock()
	if closeClients {
		e.clients.Close()
	}
}

func (m *Manager) cleanupLoop() {
	defer close(m.done)
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.cleanup()
		}
	}
}

func (m *Manager) cleanup() {
	var retired []*Clients
	m.mu.Lock()
	now := m.now()
	for _, e := range m.entries {
		if now.Sub(e.lastUsed) >= m.idleTTL {
			retired = append(retired, m.retire(e))
		}
	}
	m.mu.Unlock()
	for _, clients := range retired {
		clients.Close()
	}
}

// Close is idempotent and safe concurrently with Send. It waits for cleanup to
// stop, not for arbitrary loader code to return. In-flight requests are canceled
// and their resources are closed when they unwind.
func (m *Manager) Close() {
	m.once.Do(func() {
		var retired []*Clients
		m.mu.Lock()
		m.closed = true
		m.cancel()
		for _, e := range m.entries {
			retired = append(retired, m.retire(e))
		}
		m.mu.Unlock()
		for _, clients := range retired {
			clients.Close()
		}
		<-m.done
	})
}

func configDigest(cfg *Config) [sha256.Size]byte {
	h := sha256.New()
	// Length-prefix every field: concatenation must not create credential or
	// identity collisions.
	var size [8]byte
	for _, value := range [][]byte{
		[]byte(cfg.AppKey), []byte(cfg.Package), []byte(strconv.Itoa(cfg.IsProduct)),
		[]byte(cfg.P8KeyID), []byte(cfg.P8TeamID), cfg.P8PrivateKey,
	} {
		binary.BigEndian.PutUint64(size[:], uint64(len(value)))
		h.Write(size[:])
		h.Write(value)
	}
	var digest [sha256.Size]byte
	h.Sum(digest[:0])
	return digest
}
