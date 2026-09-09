package apnspush

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sideshow/apns2"
)

type testClock struct{ value atomic.Int64 }

func (c *testClock) now() time.Time          { return time.Unix(0, c.value.Load()) }
func (c *testClock) advance(d time.Duration) { c.value.Add(int64(d)) }

type configSource struct {
	mu    sync.Mutex
	cfg   *Config
	err   error
	loads atomic.Int64
}

func (s *configSource) load(_ context.Context, appKey, packageName string) (*Config, error) {
	s.loads.Add(1)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cfg == nil {
		return nil, s.err
	}
	copy := *s.cfg
	copy.AppKey, copy.Package = appKey, packageName
	return &copy, s.err
}
func (s *configSource) change(fn func(*Config), err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if fn != nil {
		fn(s.cfg)
	}
	s.err = err
}

func fakeClients(transport *fakeTransport) *Clients {
	client := &apns2.Client{Host: apns2.HostDevelopment, HTTPClient: &http.Client{Transport: transport}}
	return &Clients{client, client}
}

func testManager(t *testing.T, source *configSource, factory func(*Config) (*Clients, error), options ...Option) (*Manager, *testClock) {
	t.Helper()
	clock := &testClock{}
	clock.value.Store(time.Now().UnixNano())
	defaults := []Option{func(m *Manager) { m.now = clock.now; m.factory = factory }}
	m := NewManager(source.load, append(defaults, options...)...)
	t.Cleanup(m.Close)
	return m, clock
}

func send(m *Manager, appKey, packageName string) error {
	_, err := m.Send(context.Background(), appKey, packageName, &apns2.Notification{Payload: `{}`}, false, "aabb", "")
	return err
}

func cachedEntry(m *Manager, appKey, packageName string) *entry {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.entries[cacheKey{appKey, packageName}]
}

func TestManagerRefreshRotationRemovalRecovery(t *testing.T) {
	cfg, keyBytes, _ := testP8(t)
	source := &configSource{cfg: cfg}
	factory := Factory{}
	var builds atomic.Int64
	var transports []*fakeTransport
	m, clock := testManager(t, source, func(c *Config) (*Clients, error) {
		builds.Add(1)
		clients, err := factory.NewClients(c)
		if err == nil {
			transport := &fakeTransport{}
			transports = append(transports, transport)
			clients.ApnsClient.HTTPClient.Transport = transport
		}
		return clients, err
	})
	if source.loads.Load() != 0 {
		t.Fatal("constructor accessed database")
	}
	if err := send(m, cfg.AppKey, cfg.Package); err != nil {
		t.Fatal(err)
	}
	original := cachedEntry(m, cfg.AppKey, cfg.Package)
	jwt := original.clients.ApnsClient.Token.GenerateIfExpired()
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); err != nil {
		t.Fatal(err)
	}
	if cachedEntry(m, cfg.AppKey, cfg.Package) != original || original.clients.ApnsClient.Token.GenerateIfExpired() != jwt || builds.Load() != 1 || source.loads.Load() != 2 {
		t.Fatal("unchanged refresh recreated client/JWT")
	}
	if !bytes.Equal(cfg.P8PrivateKey, keyBytes) {
		t.Fatal("refresh mutated the loader-owned key")
	}

	// A legacy writer changes content without bumping ConfigVersion.
	source.change(func(c *Config) { c.IsProduct = 1 }, nil)
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); err != nil {
		t.Fatal(err)
	}
	rotated := cachedEntry(m, cfg.AppKey, cfg.Package)
	if rotated == original || rotated.clients.ApnsClient.Host != apns2.HostProduction || transports[0].closes.Load() != 1 {
		t.Fatal("legacy content change ignored")
	}
	source.change(func(c *Config) { c.ConfigVersion++ }, nil)
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); err != nil {
		t.Fatal(err)
	}
	if cachedEntry(m, cfg.AppKey, cfg.Package) == rotated {
		t.Fatal("version change ignored")
	}

	before := cachedEntry(m, cfg.AppKey, cfg.Package)
	lastTransport := transports[len(transports)-1]
	source.change(nil, errors.New("SQL contains secret credentials"))
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); !errors.Is(err, ErrLoad) || strings.Contains(err.Error(), "secret") {
		t.Fatalf("transient error: %v", err)
	}
	loadCount, sendCount := source.loads.Load(), lastTransport.calls.Load()
	if cachedEntry(m, cfg.AppKey, cfg.Package) != before || lastTransport.closes.Load() != 0 {
		t.Fatal("transient failure retired healthy client")
	}
	if err := send(m, cfg.AppKey, cfg.Package); !errors.Is(err, ErrLoad) || source.loads.Load() != loadCount || lastTransport.calls.Load() != sendCount {
		t.Fatal("retry delay served stale client or retried database")
	}
	source.change(nil, nil)
	clock.advance(5 * time.Second)
	if err := send(m, cfg.AppKey, cfg.Package); err != nil {
		t.Fatal(err)
	}
	if cachedEntry(m, cfg.AppKey, cfg.Package) != before {
		t.Fatal("recovery recreated unchanged client")
	}

	source.change(func(c *Config) { c.P8KeyID = "invalid" }, nil)
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); !errors.Is(err, ErrConfiguration) {
		t.Fatalf("invalid rotation accepted: %v", err)
	}
	if lastTransport.closes.Load() != 1 || cachedEntry(m, cfg.AppKey, cfg.Package).clients != nil {
		t.Fatal("invalid rotation retained old client")
	}
	source.change(func(c *Config) { c.P8KeyID = "ABCDEFGHIJ" }, nil)
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); err != nil {
		t.Fatal(err)
	}
	lastTransport = transports[len(transports)-1]
	source.change(nil, fmt.Errorf("dao: %w", ErrNotFound))
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); !errors.Is(err, ErrNotFound) {
		t.Fatalf("absence not distinguished: %v", err)
	}
	if lastTransport.closes.Load() != 1 || cachedEntry(m, cfg.AppKey, cfg.Package).clients != nil {
		t.Fatal("deletion retained old client")
	}
}

func TestConcurrentLoadsAndRefresh(t *testing.T) {
	source := &configSource{cfg: &Config{}}
	var builds atomic.Int64
	m, clock := testManager(t, source, func(*Config) (*Clients, error) { builds.Add(1); return fakeClients(&fakeTransport{}), nil })
	burst := func() {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := send(m, "tenant", "com.example"); err != nil {
					t.Error(err)
				}
			}()
		}
		wg.Wait()
	}
	burst()
	clock.advance(DefaultRefreshInterval)
	burst()
	if source.loads.Load() != 2 || builds.Load() != 1 {
		t.Fatalf("loads=%d builds=%d", source.loads.Load(), builds.Load())
	}
}

func TestStructuredKeysAndTenantIsolation(t *testing.T) {
	source := &configSource{cfg: &Config{}}
	m, _ := testManager(t, source, func(*Config) (*Clients, error) { return fakeClients(&fakeTransport{}), nil })
	keys := []cacheKey{{"a-b", "c"}, {"a", "b-c"}, {"other", "c"}, {"a-b", "different"}}
	seen := map[*apns2.Client]bool{}
	for _, key := range keys {
		if err := send(m, key.appKey, key.packageName); err != nil {
			t.Fatal(err)
		}
		client := cachedEntry(m, key.appKey, key.packageName).clients.ApnsClient
		if seen[client] {
			t.Fatal("client shared across tenant/bundle")
		}
		seen[client] = true
	}
	if source.loads.Load() != int64(len(keys)) {
		t.Fatal("structured key collision")
	}
}

func TestIdentityMismatchRetiresClient(t *testing.T) {
	var wrong atomic.Bool
	transport := &fakeTransport{}
	clock := &testClock{}
	m := NewManager(func(_ context.Context, appKey, packageName string) (*Config, error) {
		if wrong.Load() {
			appKey = "other"
		}
		return &Config{AppKey: appKey, Package: packageName}, nil
	}, func(m *Manager) {
		m.now = clock.now
		m.factory = func(*Config) (*Clients, error) { return fakeClients(transport), nil }
	})
	defer m.Close()
	if err := send(m, "tenant", "bundle"); err != nil {
		t.Fatal(err)
	}
	wrong.Store(true)
	clock.advance(DefaultRefreshInterval)
	if err := send(m, "tenant", "bundle"); !errors.Is(err, ErrConfiguration) {
		t.Fatal("cross-tenant loader result accepted")
	}
	if transport.closes.Load() != 1 || transport.calls.Load() != 1 {
		t.Fatal("mismatched identity reused old client")
	}
}

func TestIdleAndCapacityEviction(t *testing.T) {
	source := &configSource{cfg: &Config{}}
	transports := map[string]*fakeTransport{}
	m, clock := testManager(t, source, func(c *Config) (*Clients, error) {
		transport := &fakeTransport{}
		transports[c.Package] = transport
		return fakeClients(transport), nil
	}, func(m *Manager) { m.capacity = 2 })
	for _, bundle := range []string{"a", "b", "a", "c"} {
		if err := send(m, "tenant", bundle); err != nil {
			t.Fatal(err)
		}
	}
	if cachedEntry(m, "tenant", "b") != nil || transports["b"].closes.Load() != 1 || cachedEntry(m, "tenant", "a") == nil {
		t.Fatal("LRU capacity eviction incorrect")
	}
	clock.advance(DefaultIdleTTL - time.Second)
	m.cleanup()
	if cachedEntry(m, "tenant", "a") == nil {
		t.Fatal("idle cleanup too early")
	}
	clock.advance(time.Second)
	m.cleanup()
	if cachedEntry(m, "tenant", "a") != nil || cachedEntry(m, "tenant", "c") != nil || transports["a"].closes.Load() != 1 || transports["c"].closes.Load() != 1 {
		t.Fatal("idle cleanup failed")
	}
}

func TestPeriodicCleanup(t *testing.T) {
	source := &configSource{cfg: &Config{}}
	transport := &fakeTransport{}
	m, clock := testManager(t, source, func(*Config) (*Clients, error) { return fakeClients(transport), nil }, func(m *Manager) { m.cleanupInterval = time.Millisecond })
	if err := send(m, "tenant", "bundle"); err != nil {
		t.Fatal(err)
	}
	clock.advance(DefaultIdleTTL)
	deadline := time.Now().Add(time.Second)
	for transport.closes.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if transport.closes.Load() != 1 {
		t.Fatal("background idle cleanup did not run")
	}
	m.Close()
	select {
	case <-m.done:
	default:
		t.Fatal("cleanup goroutine not stopped")
	}
}

func TestEvictionDoesNotCancelActiveRequest(t *testing.T) {
	for _, kind := range []string{"capacity", "idle", "rotation", "deletion"} {
		t.Run(kind, func(t *testing.T) {
			started, finish := make(chan struct{}), make(chan struct{})
			active := &fakeTransport{fn: func(r *http.Request) (*http.Response, error) {
				close(started)
				select {
				case <-finish:
					return reply(200, ""), nil
				case <-r.Context().Done():
					return nil, r.Context().Err()
				}
			}}
			source := &configSource{cfg: &Config{}}
			var builds atomic.Int64
			m, clock := testManager(t, source, func(*Config) (*Clients, error) {
				if builds.Add(1) == 1 {
					return fakeClients(active), nil
				}
				return fakeClients(&fakeTransport{}), nil
			}, func(m *Manager) { m.capacity = 1 })
			result := make(chan error, 1)
			go func() { result <- send(m, "tenant", "a") }()
			<-started
			switch kind {
			case "capacity":
				if err := send(m, "tenant", "b"); err != nil {
					t.Fatal(err)
				}
			case "idle":
				clock.advance(DefaultIdleTTL)
				m.cleanup()
			case "rotation":
				source.change(func(c *Config) { c.ConfigVersion++ }, nil)
				clock.advance(DefaultRefreshInterval)
				if err := send(m, "tenant", "a"); err != nil {
					t.Fatal(err)
				}
			case "deletion":
				source.change(nil, ErrNotFound)
				clock.advance(DefaultRefreshInterval)
				if err := send(m, "tenant", "a"); !errors.Is(err, ErrNotFound) {
					t.Fatal(err)
				}
			}
			select {
			case err := <-result:
				t.Fatalf("eviction canceled active request: %v", err)
			default:
			}
			close(finish)
			if err := <-result; err != nil {
				t.Fatal(err)
			}
			if active.closes.Load() != 1 {
				t.Fatal("retired active client not closed after completion")
			}
		})
	}
}

func TestSendTimeoutCancellationAndClose(t *testing.T) {
	for _, kind := range []string{"timeout", "cancel", "close"} {
		t.Run(kind, func(t *testing.T) {
			started := make(chan struct{})
			transport := &fakeTransport{fn: func(r *http.Request) (*http.Response, error) {
				close(started)
				<-r.Context().Done()
				return nil, r.Context().Err()
			}}
			source := &configSource{cfg: &Config{}}
			m, _ := testManager(t, source, func(*Config) (*Clients, error) { return fakeClients(transport), nil }, func(m *Manager) {
				if kind == "timeout" {
					m.timeout = 20 * time.Millisecond
				}
			})
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			result := make(chan error, 1)
			go func() {
				_, err := m.Send(ctx, "tenant", "bundle", &apns2.Notification{Payload: `{}`}, false, "aabb", "")
				result <- err
			}()
			<-started
			if kind == "cancel" {
				cancel()
			}
			if kind == "close" {
				m.Close()
			}
			want := context.Canceled
			if kind == "timeout" {
				want = context.DeadlineExceeded
			}
			select {
			case err := <-result:
				if !errors.Is(err, want) || strings.Contains(err.Error(), "aabb") {
					t.Fatalf("cancellation lost or token leaked: %v", err)
				}
			case <-time.After(time.Second):
				t.Fatal("send did not stop")
			}
			m.Close()
			if err := send(m, "tenant", "bundle"); !errors.Is(err, ErrClosed) {
				t.Fatal("send after Close accepted")
			}
			if transport.closes.Load() != 1 {
				t.Fatal("shutdown did not close distinct client once")
			}
		})
	}
}

func TestCanceledLoadAndStripedWait(t *testing.T) {
	started := make(chan struct{})
	var loads atomic.Int64
	m := NewManager(func(ctx context.Context, _, _ string) (*Config, error) {
		loads.Add(1)
		close(started)
		<-ctx.Done()
		return nil, ctx.Err()
	})
	defer m.Close()
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := m.Send(ctx, "tenant", "bundle", &apns2.Notification{}, false, "aabb", "")
		result <- err
	}()
	<-started
	waitCtx, stop := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer stop()
	if _, err := m.Send(waitCtx, "tenant", "bundle", &apns2.Notification{}, false, "aabb", ""); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("striped wait ignored deadline: %v", err)
	}
	if loads.Load() != 1 {
		t.Fatal("same-key loads overlap")
	}
	cancel()
	if err := <-result; !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if cachedEntry(m, "tenant", "bundle") != nil {
		t.Fatal("canceled loader poisoned cache")
	}
}

func TestNoGlobalLockDuringLoad(t *testing.T) {
	started, release := make(chan struct{}), make(chan struct{})
	m := NewManager(func(ctx context.Context, appKey, bundle string) (*Config, error) {
		if bundle == "blocked" {
			close(started)
			select {
			case <-release:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		return &Config{AppKey: appKey, Package: bundle}, nil
	}, func(m *Manager) {
		m.factory = func(*Config) (*Clients, error) { return fakeClients(&fakeTransport{}), nil }
	})
	defer m.Close()
	if err := send(m, "tenant", "cached"); err != nil {
		t.Fatal(err)
	}
	result := make(chan error, 1)
	go func() { result <- send(m, "tenant", "blocked") }()
	<-started
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := m.Send(ctx, "tenant", "cached", &apns2.Notification{}, false, "aabb", ""); err != nil {
		t.Fatal("database load blocked cached send:", err)
	}
	close(release)
	if err := <-result; err != nil {
		t.Fatal(err)
	}
}

func TestCloseDuringLoadAndBuild(t *testing.T) {
	for _, phase := range []string{"load", "build"} {
		t.Run(phase, func(t *testing.T) {
			started, release := make(chan struct{}), make(chan struct{})
			transport := &fakeTransport{}
			m := NewManager(func(ctx context.Context, a, b string) (*Config, error) {
				if phase == "load" {
					close(started)
					<-ctx.Done()
					return nil, ctx.Err()
				}
				return &Config{AppKey: a, Package: b}, nil
			}, func(m *Manager) {
				m.factory = func(*Config) (*Clients, error) { close(started); <-release; return fakeClients(transport), nil }
			})
			defer m.Close()
			result := make(chan error, 1)
			go func() { result <- send(m, "tenant", "bundle") }()
			<-started
			m.Close()
			close(release)
			if err := <-result; !errors.Is(err, context.Canceled) && !errors.Is(err, ErrClosed) {
				t.Fatalf("shutdown result: %v", err)
			}
			if phase == "build" && transport.closes.Load() != 1 {
				t.Fatal("client built during Close leaked")
			}
			if transport.calls.Load() != 0 || cachedEntry(m, "tenant", "bundle") != nil {
				t.Fatal("client published after Close")
			}
		})
	}
}

func TestSendCloseRace(t *testing.T) {
	for iteration := 0; iteration < 20; iteration++ {
		cfg, _, _ := testP8(t)
		source := &configSource{cfg: cfg}
		factory := Factory{}
		m, _ := testManager(t, source, func(c *Config) (*Clients, error) {
			clients, err := factory.NewClients(c)
			if err == nil {
				clients.ApnsClient.HTTPClient.Transport = &fakeTransport{}
			}
			return clients, err
		})
		if err := send(m, cfg.AppKey, cfg.Package); err != nil {
			t.Fatal(err)
		}
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if i%11 == 0 {
					m.Close()
					return
				}
				_, err := m.Send(context.Background(), cfg.AppKey, cfg.Package, &apns2.Notification{Payload: `{}`}, i%2 == 0, "aabb", "ccdd")
				if err != nil && !errors.Is(err, ErrClosed) && !errors.Is(err, context.Canceled) {
					t.Error(err)
				}
			}(i)
		}
		wg.Wait()
		m.Close()
	}
}

func TestNoRetryOrResignOnAPNsRejection(t *testing.T) {
	cfg, _, _ := testP8(t)
	source := &configSource{cfg: cfg}
	var clients *Clients
	transport := &fakeTransport{fn: func(*http.Request) (*http.Response, error) {
		return reply(403, `{"reason":"ExpiredProviderToken"}`), nil
	}}
	m, _ := testManager(t, source, func(c *Config) (*Clients, error) {
		var err error
		clients, err = (Factory{}).NewClients(c)
		if err == nil {
			clients.ApnsClient.HTTPClient.Transport = transport
		}
		return clients, err
	})
	var first string
	for i := 0; i < 3; i++ {
		response, err := m.Send(context.Background(), cfg.AppKey, cfg.Package, &apns2.Notification{}, false, "aabb", "")
		if err != nil || Classify(response, err) != ClassAuthentication {
			t.Fatal("APNs rejection not preserved")
		}
		bearer := clients.ApnsClient.Token.GenerateIfExpired()
		if i == 0 {
			first = bearer
		} else if first != bearer {
			t.Fatal("auth rejection regenerated token")
		}
	}
	if transport.calls.Load() != 3 || source.loads.Load() != 1 {
		t.Fatal("APNs rejection triggered automatic retry")
	}
}

func TestDigestIncludesAllCredentialFields(t *testing.T) {
	base := &Config{AppKey: "tenant", Package: "bundle", P8KeyID: "key", P8TeamID: "team", P8PrivateKey: []byte("raw-key")}
	original := configDigest(base)
	for _, change := range []func(*Config){
		func(c *Config) { c.AppKey = "other" },
		func(c *Config) { c.Package = "other" },
		func(c *Config) { c.IsProduct = 1 },
		func(c *Config) { c.P8KeyID = "other" },
		func(c *Config) { c.P8TeamID = "other" },
		func(c *Config) { c.P8PrivateKey = []byte("other") },
	} {
		copy := *base
		change(&copy)
		if configDigest(&copy) == original {
			t.Fatal("credential content excluded from digest")
		}
	}
	a, b := *base, *base
	a.P8KeyID, a.P8TeamID = "ab", "c"
	b.P8KeyID, b.P8TeamID = "a", "bc"
	if configDigest(&a) == configDigest(&b) {
		t.Fatal("digest field boundary collision")
	}
}

func TestP8KeyRotationWithoutVersionChange(t *testing.T) {
	cfg, _, oldKey := testP8(t)
	_, newPEM, newKey := testP8(t)
	source := &configSource{cfg: cfg}
	var clients []*Clients
	factory := Factory{}
	m, clock := testManager(t, source, func(c *Config) (*Clients, error) {
		client, err := factory.NewClients(c)
		if err == nil {
			client.ApnsClient.HTTPClient.Transport = &fakeTransport{}
			clients = append(clients, client)
		}
		return client, err
	})
	if err := send(m, cfg.AppKey, cfg.Package); err != nil {
		t.Fatal(err)
	}
	// Change only the original file bytes, not version or metadata.
	source.change(func(c *Config) { c.P8PrivateKey = bytes.Clone(newPEM) }, nil)
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); err != nil {
		t.Fatal(err)
	}
	if len(clients) != 2 {
		t.Fatal("key rotation not detected")
	}
	bearer := clients[1].ApnsClient.Token.GenerateIfExpired()
	parsed, err := jwt.Parse(bearer, func(*jwt.Token) (any, error) { return &newKey.PublicKey, nil }, jwt.WithValidMethods([]string{"ES256"}))
	if err != nil || !parsed.Valid || parsed.Header["kid"] != cfg.P8KeyID {
		t.Fatalf("rotated JWT invalid: %v", err)
	}
	if _, err := jwt.Parse(bearer, func(*jwt.Token) (any, error) { return &oldKey.PublicKey, nil }); err == nil {
		t.Fatal("rotated JWT still uses old key")
	}
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); err != nil {
		t.Fatal(err)
	}
	if len(clients) != 2 || clients[1].ApnsClient.Token.GenerateIfExpired() != bearer || !bytes.Equal(cfg.P8PrivateKey, newPEM) {
		t.Fatal("unchanged raw key not reused or loader snapshot mutated")
	}
	source.change(func(c *Config) { c.P8PrivateKey = nil }, nil)
	clock.advance(DefaultRefreshInterval)
	if err := send(m, cfg.AppKey, cfg.Package); !errors.Is(err, ErrConfiguration) {
		t.Fatalf("missing raw key retained old authentication: %v", err)
	}
	if cachedEntry(m, cfg.AppKey, cfg.Package).clients != nil || clients[1].ApnsClient.HTTPClient.Transport.(*fakeTransport).closes.Load() != 1 {
		t.Fatal("removed raw key did not retire the cached client")
	}
}

func TestNegativeEntriesAreBoundedAndRecover(t *testing.T) {
	source := &configSource{cfg: nil}
	var builds atomic.Int64
	m, clock := testManager(t, source, func(*Config) (*Clients, error) { builds.Add(1); return fakeClients(&fakeTransport{}), nil }, func(m *Manager) { m.capacity = 3 })
	for i := 0; i < 20; i++ {
		if err := send(m, "tenant", fmt.Sprintf("bundle%d", i)); !errors.Is(err, ErrNotFound) {
			t.Fatalf("nil row should mean absence: %v", err)
		}
	}
	m.mu.Lock()
	count := len(m.entries)
	m.mu.Unlock()
	if count != 3 || builds.Load() != 0 {
		t.Fatal("negative entries unbounded or constructed clients")
	}
	source.mu.Lock()
	source.cfg = &Config{}
	source.mu.Unlock()
	clock.advance(DefaultRefreshInterval)
	if err := send(m, "tenant", "bundle19"); err != nil {
		t.Fatal(err)
	}
	if builds.Load() != 1 {
		t.Fatal("negative cache failed to recover")
	}
}

func TestTransientInitialLoadRecovery(t *testing.T) {
	source := &configSource{cfg: &Config{}, err: errors.New("database unavailable")}
	m, clock := testManager(t, source, func(*Config) (*Clients, error) { return fakeClients(&fakeTransport{}), nil })
	for i := 0; i < 10; i++ {
		if err := send(m, "tenant", "bundle"); !errors.Is(err, ErrLoad) {
			t.Fatal(err)
		}
	}
	if source.loads.Load() != 1 {
		t.Fatal("initial transient failure retried without backoff")
	}
	source.change(nil, nil)
	clock.advance(5 * time.Second)
	if err := send(m, "tenant", "bundle"); err != nil {
		t.Fatal(err)
	}
	if source.loads.Load() != 2 {
		t.Fatal("initial transient failure did not recover")
	}
}

func TestInvalidInputNeverContactsAPNs(t *testing.T) {
	source := &configSource{cfg: &Config{}}
	transport := &fakeTransport{}
	m, _ := testManager(t, source, func(*Config) (*Clients, error) { return fakeClients(transport), nil })
	for _, tc := range []struct {
		token   string
		payload any
		class   ErrorClass
	}{
		{"", `{}`, ClassDevice},
		{"aa/bb", `{}`, ClassDevice},
		{"aabb", make(chan int), ClassConfiguration},
	} {
		resp, err := m.Send(context.Background(), "tenant", "bundle", &apns2.Notification{Payload: tc.payload}, false, tc.token, "")
		if err == nil || Classify(resp, err) != tc.class {
			t.Fatalf("invalid input classified incorrectly: %v", err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := m.Send(ctx, "tenant", "bundle", &apns2.Notification{}, false, "aabb", ""); !errors.Is(err, context.Canceled) {
		t.Fatal("pre-canceled request accepted")
	}
	if transport.calls.Load() != 0 {
		t.Fatal("invalid input contacted APNs")
	}
}
