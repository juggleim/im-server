// Package apnspush owns APNs authentication, routing and client lifetimes.
// It does not perform database access or logging.
package apnspush

import (
	"bytes"
	"crypto/elliptic"
	"fmt"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
	"golang.org/x/net/http2"
)

// Config is a loader-owned P8 snapshot. It must not be mutated while being loaded.
// IsProduct > 0 selects production.
type Config struct {
	AppKey        string
	Package       string
	IsProduct     int
	P8KeyID       string
	P8TeamID      string
	P8PrivateKey  []byte
	ConfigVersion int64
}

// Clients are immutable after publication, except for SDK-managed token state.
// P8 uses the same client for both routes. Do not mutate the shared token.
type Clients struct {
	ApnsClient     *apns2.Client
	ApnsVoipClient *apns2.Client
}

// Close closes idle connections, once per distinct client. It does not cancel
// active requests. Managers also close retired clients after their last request.
func (c *Clients) Close() {
	if c == nil {
		return
	}
	if c.ApnsClient != nil && c.ApnsClient.HTTPClient != nil {
		c.ApnsClient.HTTPClient.CloseIdleConnections()
	}
	if c.ApnsVoipClient != nil && c.ApnsVoipClient != c.ApnsClient && c.ApnsVoipClient.HTTPClient != nil {
		c.ApnsVoipClient.HTTPClient.CloseIdleConnections()
	}
}

type Factory struct{}

// NewClients validates and constructs clients without contacting Apple. Errors
// deliberately omit underlying parser errors and credential values.
func (f Factory) NewClients(cfg *Config) (*Clients, error) {
	if cfg == nil || cfg.AppKey == "" || !validTopic(cfg.Package) {
		return nil, fmt.Errorf("%w: missing or invalid identity", ErrConfiguration)
	}
	// Clear only our copy; the loader may reuse its snapshot on refresh.
	keyBytes := bytes.Clone(cfg.P8PrivateKey)
	defer clear(keyBytes)
	if validateCredentials(keyBytes, cfg.P8KeyID, cfg.P8TeamID) != nil {
		return nil, fmt.Errorf("%w: invalid P8 credentials", ErrConfiguration)
	}
	key, err := token.AuthKeyFromBytes(bytes.TrimSpace(keyBytes))
	if err != nil || key == nil || key.Curve != elliptic.P256() {
		return nil, fmt.Errorf("%w: P8 requires PKCS#8 ECDSA P-256", ErrConfiguration)
	}
	providerToken := &token.Token{AuthKey: key, KeyID: cfg.P8KeyID, TeamID: cfg.P8TeamID}
	// Generate is not independently synchronized in apns2 v0.23.0. Pre-sign
	// before publication, then leave all refreshes to GenerateIfExpired.
	if ok, err := providerToken.Generate(); err != nil || !ok {
		return nil, fmt.Errorf("%w: P8 signing failed", ErrConfiguration)
	}
	client := apns2.NewTokenClient(providerToken)
	clients := &Clients{ApnsClient: client, ApnsVoipClient: client}
	for _, client := range []*apns2.Client{clients.ApnsClient, clients.ApnsVoipClient} {
		if client != nil {
			// The SDK installs a context-free TLS dialer. Use HTTP/2's native
			// context-aware TLS dial path so cancellation also bounds handshakes.
			client.HTTPClient.Transport.(*http2.Transport).DialTLS = nil
			if cfg.IsProduct > 0 {
				client.Production()
			} else {
				client.Development()
			}
		}
	}
	return clients, nil
}
