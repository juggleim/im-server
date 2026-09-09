package apnspush

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sideshow/apns2"
	"golang.org/x/net/http2"
)

type fakeTransport struct {
	fn     func(*http.Request) (*http.Response, error)
	calls  atomic.Int64
	closes atomic.Int64
}

func (f *fakeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	f.calls.Add(1)
	if f.fn != nil {
		return f.fn(r)
	}
	return reply(200, ""), nil
}
func (f *fakeTransport) CloseIdleConnections() { f.closes.Add(1) }

func reply(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Apns-Id": {"test-id"}}, Body: io.NopCloser(strings.NewReader(body))}
}

func testP8(t *testing.T) (*Config, []byte, *ecdsa.PrivateKey) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	privateKey := keyPEM(t, key)
	return &Config{AppKey: "tenant", Package: "com.example.app", P8KeyID: "ABCDEFGHIJ", P8TeamID: "1234567890", P8PrivateKey: bytes.Clone(privateKey), ConfigVersion: 1}, privateKey, key
}

func keyPEM(t *testing.T, key any) []byte {
	t.Helper()
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

func TestP8Protocol(t *testing.T) {
	for _, product := range []int{0, 1} {
		cfg, _, key := testP8(t)
		cfg.IsProduct = product
		clients, err := (Factory{}).NewClients(cfg)
		if err != nil {
			t.Fatal(err)
		}
		if clients.ApnsClient != clients.ApnsVoipClient {
			t.Fatal("P8 clients must alias")
		}
		client := clients.ApnsClient
		if client.HTTPClient.Transport.(*http2.Transport).DialTLS != nil {
			t.Fatal("context-free dialer retained")
		}
		wantHost := apns2.HostDevelopment
		if product > 0 {
			wantHost = apns2.HostProduction
		}
		if client.Host != wantHost {
			t.Fatalf("host = %s", client.Host)
		}
		initialJWT := client.Token.GenerateIfExpired()
		transport := &fakeTransport{}
		client.HTTPClient.Transport = transport
		for _, voip := range []bool{false, true} {
			n := &apns2.Notification{Payload: `{"aps":{"alert":"hello"}}`, Priority: 10, ApnsID: "request-id", CollapseID: "collapse", Expiration: time.Unix(1900000000, 0)}
			transport.fn = func(r *http.Request) (*http.Response, error) {
				wantTopic, wantType, wantToken := cfg.Package, "alert", "aabb"
				if voip {
					wantTopic, wantType, wantToken = cfg.Package+".voip", "voip", "ccdd"
				}
				if r.Method != "POST" || r.URL.Path != "/3/device/"+wantToken || r.Header.Get("apns-topic") != wantTopic || r.Header.Get("apns-push-type") != wantType {
					t.Error("incorrect APNs route or headers")
				}
				if r.Header.Get("apns-priority") != "10" || r.Header.Get("apns-id") != "request-id" || r.Header.Get("apns-collapse-id") != "collapse" || r.Header.Get("apns-expiration") != "1900000000" {
					t.Error("notification metadata lost")
				}
				body, _ := io.ReadAll(r.Body)
				if string(body) != n.Payload {
					t.Error("payload changed")
				}
				bearer := strings.TrimPrefix(r.Header.Get("authorization"), "bearer ")
				if bearer != initialJWT {
					t.Error("provider token unnecessarily regenerated")
				}
				parsed, err := jwt.Parse(bearer, func(j *jwt.Token) (any, error) {
					if j.Method != jwt.SigningMethodES256 || j.Header["kid"] != cfg.P8KeyID {
						t.Error("incorrect JWT algorithm or key id")
					}
					return &key.PublicKey, nil
				}, jwt.WithValidMethods([]string{"ES256"}))
				if err != nil || !parsed.Valid {
					t.Fatalf("JWT verification: %v", err)
				}
				claims := parsed.Claims.(jwt.MapClaims)
				if claims["iss"] != cfg.P8TeamID || time.Since(time.Unix(int64(claims["iat"].(float64)), 0)) > time.Minute {
					t.Error("incorrect JWT claims")
				}
				return reply(200, ""), nil
			}
			selected, routed, err := Route(clients, cfg.Package, n, voip, "aabb", "ccdd")
			if err != nil {
				t.Fatal(err)
			}
			if resp, err := selected.PushWithContext(context.Background(), routed); err != nil || !resp.Sent() {
				t.Fatalf("push failed: %v", err)
			}
			if n.Topic != "" || n.DeviceToken != "" || n.PushType != "" {
				t.Fatal("notification mutated")
			}
		}
		clients.Close()
		if transport.closes.Load() != 1 {
			t.Fatal("aliased client closed more than once")
		}
	}
}

func TestFactoryRejectsInvalidCredentials(t *testing.T) {
	cfg, valid, _ := testP8(t)
	p384, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	rsaKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	for _, input := range [][]byte{nil, []byte("secret-invalid-pem"), keyPEM(t, p384), keyPEM(t, rsaKey), bytes.Repeat([]byte("x"), 16385)} {
		cfg.P8PrivateKey = bytes.Clone(input)
		if clients, err := (Factory{}).NewClients(cfg); clients != nil || !errors.Is(err, ErrConfiguration) || strings.Contains(err.Error(), "secret-invalid-pem") {
			t.Fatalf("invalid key accepted or leaked: %v", err)
		}
		if !bytes.Equal(cfg.P8PrivateKey, input) {
			t.Fatal("invalid loader-owned key mutated")
		}
	}
	cfg.P8PrivateKey = bytes.Clone(valid)
	for _, mutate := range []func(*Config){
		func(c *Config) { c.P8KeyID = "" },
		func(c *Config) { c.P8TeamID = "" },
		func(c *Config) { c.Package = "bad\r\ntopic" },
	} {
		copy := *cfg
		mutate(&copy)
		if _, err := (Factory{}).NewClients(&copy); !errors.Is(err, ErrConfiguration) {
			t.Fatalf("invalid configuration accepted: %v", err)
		}
	}
	clients, err := (Factory{}).NewClients(cfg)
	if err != nil {
		t.Fatal(err)
	}
	clients.Close()
	if !bytes.Equal(cfg.P8PrivateKey, valid) {
		t.Fatal("loader-owned key mutated")
	}
}

func TestFactoryP8RawWhitespace(t *testing.T) {
	cfg, privateKey, key := testP8(t)
	original := append([]byte("\r\n  \t"), bytes.ReplaceAll(privateKey, []byte("\n"), []byte("\r\n"))...)
	original = append(original, []byte(" \t\r\n")...)
	cfg.P8PrivateKey = bytes.Clone(original)
	digest := configDigest(cfg)
	for range 2 {
		clients, err := (Factory{}).NewClients(cfg)
		if err != nil {
			t.Fatal(err)
		}
		bearer := clients.ApnsClient.Token.GenerateIfExpired()
		clients.Close()
		parsed, err := jwt.Parse(bearer, func(*jwt.Token) (any, error) {
			return &key.PublicKey, nil
		}, jwt.WithValidMethods([]string{"ES256"}))
		if err != nil || !parsed.Valid {
			t.Fatalf("whitespace-wrapped raw key failed signing: %v", err)
		}
		if !bytes.Equal(cfg.P8PrivateKey, original) || configDigest(cfg) != digest {
			t.Fatal("factory modified loader-owned raw bytes or digest")
		}
	}
}
