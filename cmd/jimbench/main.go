package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"im-server/simulator/benchmark"
)

func main() {
	config := benchmark.DefaultConfig()
	var scenario string
	flag.StringVar(&scenario, "scenario", string(config.Scenario), "workload scenario: private or group")
	flag.StringVar(&config.WSURL, "ws-url", config.WSURL, "IM WebSocket base URL")
	flag.StringVar(&config.APIURL, "api-url", config.APIURL, "signed server API base URL")
	flag.StringVar(&config.AppKey, "app-key", os.Getenv("JIM_BENCH_APP_KEY"), "benchmark application key")
	flag.IntVar(&config.Clients, "clients", config.Clients, "number of connected clients")
	flag.IntVar(&config.GroupSenders, "group-senders", config.GroupSenders, "number of group-message senders")
	flag.IntVar(&config.Rate, "rate", config.Rate, "aggregate target messages per second")
	flag.DurationVar(&config.Warmup, "warmup", config.Warmup, "warm-up duration excluded from measurements")
	flag.DurationVar(&config.Duration, "duration", config.Duration, "measurement duration")
	flag.DurationVar(&config.DeliveryGrace, "delivery-grace", config.DeliveryGrace, "time to collect late deliveries after sending stops")
	flag.IntVar(&config.PayloadBytes, "payload-bytes", config.PayloadBytes, "minimum JSON message payload size")
	flag.IntVar(&config.SetupConcurrency, "setup-concurrency", config.SetupConcurrency, "concurrent user registration requests")
	flag.IntVar(&config.ConnectConcurrency, "connect-concurrency", config.ConnectConcurrency, "concurrent WebSocket connection attempts")
	flag.BoolVar(&config.StoreMessages, "store-messages", config.StoreMessages, "set the stored-message flag")
	flag.BoolVar(&config.CountMessages, "count-messages", config.CountMessages, "set the counted-message flag")
	flag.BoolVar(&config.AllowNonLoopback, "allow-non-loopback", false, "allow an isolated non-loopback target (never use against production)")
	flag.StringVar(&config.OutputPath, "output", "", "JSON result path (default: benchmark-results/<run>-<scenario>.json)")
	flag.StringVar(&config.EnvironmentLabel, "environment", envOr("JIM_BENCH_ENVIRONMENT", config.EnvironmentLabel), "human-readable deployment/configuration label")
	flag.StringVar(&config.DatabaseLabel, "database", envOr("JIM_BENCH_DATABASE", config.DatabaseLabel), "database engine and configuration label")
	flag.Parse()

	config.Scenario = benchmark.Scenario(scenario)
	config.AppSecret = os.Getenv("JIM_BENCH_APP_SECRET")
	config.ServerCommit = os.Getenv("JIM_BENCH_SERVER_COMMIT")
	runner, err := benchmark.NewRunner(config)
	if err != nil {
		exitf("invalid benchmark configuration: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("Preparing %s benchmark: clients=%d rate=%d/s warmup=%s duration=%s payload>=%d bytes\n",
		config.Scenario, config.Clients, config.Rate, config.Warmup, config.Duration, config.PayloadBytes)
	report, err := runner.Run(ctx)
	if err != nil {
		exitf("benchmark failed: %v", err)
	}
	path, err := benchmark.WriteReport(config.OutputPath, report)
	if err != nil {
		exitf("save benchmark report: %v", err)
	}

	fmt.Printf("Benchmark complete.\n")
	fmt.Printf("  Connections: %d/%d, P95 %.2f ms\n", report.Connections.Succeeded, report.Connections.Attempted, report.Connections.Latency.P95MS)
	fmt.Printf("  ACKs: %d succeeded, %d failed, %.2f/s, P95 %.2f ms, P99 %.2f ms\n",
		report.MessageAcknowledgements.Succeeded, report.MessageAcknowledgements.Failed,
		report.MessageAcknowledgements.ThroughputPerSec, report.MessageAcknowledgements.Latency.P95MS,
		report.MessageAcknowledgements.Latency.P99MS)
	fmt.Printf("  Deliveries: %d observed, %.2f/s, P95 %.2f ms, P99 %.2f ms\n",
		report.Deliveries.Succeeded, report.Deliveries.ThroughputPerSec,
		report.Deliveries.Latency.P95MS, report.Deliveries.Latency.P99MS)
	fmt.Printf("  Result: %s\n", path)
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
