package benchmark

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type Scenario string

const (
	ScenarioPrivate Scenario = "private"
	ScenarioGroup   Scenario = "group"
)

type Config struct {
	Scenario           Scenario
	WSURL              string
	APIURL             string
	AppKey             string
	AppSecret          string
	Clients            int
	GroupSenders       int
	Rate               int
	Warmup             time.Duration
	Duration           time.Duration
	DeliveryGrace      time.Duration
	PayloadBytes       int
	SetupConcurrency   int
	ConnectConcurrency int
	StoreMessages      bool
	CountMessages      bool
	AllowNonLoopback   bool
	OutputPath         string
	EnvironmentLabel   string
	DatabaseLabel      string
	ServerCommit       string
}

func DefaultConfig() Config {
	return Config{
		Scenario:           ScenarioPrivate,
		WSURL:              "ws://127.0.0.1:9003",
		APIURL:             "http://127.0.0.1:9001/apigateway",
		Clients:            50,
		GroupSenders:       2,
		Rate:               200,
		Warmup:             10 * time.Second,
		Duration:           30 * time.Second,
		DeliveryGrace:      3 * time.Second,
		PayloadBytes:       256,
		SetupConcurrency:   10,
		ConnectConcurrency: 20,
		StoreMessages:      true,
		CountMessages:      true,
		EnvironmentLabel:   "unspecified",
		DatabaseLabel:      "unspecified",
	}
}

func (c Config) Validate() error {
	if c.Scenario != ScenarioPrivate && c.Scenario != ScenarioGroup {
		return fmt.Errorf("scenario must be %q or %q", ScenarioPrivate, ScenarioGroup)
	}
	if c.AppKey == "" {
		return fmt.Errorf("app key is required")
	}
	if c.AppSecret == "" {
		return fmt.Errorf("app secret is required; set JIM_BENCH_APP_SECRET")
	}
	if c.Clients < 2 {
		return fmt.Errorf("clients must be at least 2")
	}
	if c.GroupSenders < 1 || c.GroupSenders > c.Clients {
		return fmt.Errorf("group senders must be between 1 and clients")
	}
	if c.Rate < 1 {
		return fmt.Errorf("rate must be at least 1 message per second")
	}
	if c.Warmup < 0 || c.Duration <= 0 || c.DeliveryGrace < 0 {
		return fmt.Errorf("warmup and delivery grace cannot be negative, and duration must be positive")
	}
	if c.PayloadBytes < 128 {
		return fmt.Errorf("payload bytes must be at least 128")
	}
	if c.SetupConcurrency < 1 || c.ConnectConcurrency < 1 {
		return fmt.Errorf("setup and connect concurrency must be positive")
	}
	for name, rawURL := range map[string]string{"WebSocket": c.WSURL, "API": c.APIURL} {
		u, err := url.Parse(rawURL)
		if err != nil || u.Hostname() == "" {
			return fmt.Errorf("invalid %s URL %q", name, rawURL)
		}
		if !c.AllowNonLoopback && !isLoopbackHost(u.Hostname()) {
			return fmt.Errorf("%s URL %q is not loopback; pass --allow-non-loopback only for an isolated benchmark environment", name, rawURL)
		}
	}
	return nil
}

func isLoopbackHost(host string) bool {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
