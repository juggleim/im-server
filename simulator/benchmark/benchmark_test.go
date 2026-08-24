package benchmark

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestSummarizeLatenciesUsesNearestRank(t *testing.T) {
	values := make([]time.Duration, 100)
	for index := range values {
		values[index] = time.Duration(index+1) * time.Millisecond
	}
	summary := summarizeLatencies(values)
	if summary.Count != 100 || summary.P50MS != 50 || summary.P95MS != 95 || summary.P99MS != 99 {
		t.Fatalf("unexpected latency summary: %+v", summary)
	}
	if summary.MinMS != 1 || summary.MaxMS != 100 || summary.MeanMS != 50.5 {
		t.Fatalf("unexpected latency bounds or mean: %+v", summary)
	}
}

func TestMetricRecorderSeparatesFailures(t *testing.T) {
	recorder := newMetricRecorder()
	recorder.recordSuccess(10 * time.Millisecond)
	recorder.recordFailure("timeout", 20*time.Millisecond)
	snapshot := recorder.snapshot(time.Second)
	if snapshot.Attempted != 2 || snapshot.Succeeded != 1 || snapshot.Failed != 1 {
		t.Fatalf("unexpected counters: %+v", snapshot)
	}
	if snapshot.Errors["timeout"] != 1 || snapshot.ThroughputPerSec != 1 {
		t.Fatalf("unexpected errors or throughput: %+v", snapshot)
	}
}

func TestConfigRejectsRemoteTargetsByDefault(t *testing.T) {
	config := DefaultConfig()
	config.AppKey = "test"
	config.AppSecret = "secret"
	config.WSURL = "wss://example.com"
	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "not loopback") {
		t.Fatalf("expected non-loopback validation error, got %v", err)
	}
	config.AllowNonLoopback = true
	if err := config.Validate(); err != nil {
		t.Fatalf("allow non-loopback should permit isolated remote targets: %v", err)
	}
}

func TestMakePayloadHasRequestedMinimumSizeAndRoundTrips(t *testing.T) {
	sentAt := time.Unix(123, 456)
	data, err := makePayload("run", "measure", 42, sentAt, 512)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 512 {
		t.Fatalf("payload length %d is less than requested 512", len(data))
	}
	var payload benchmarkPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.BenchmarkID != "run" || payload.Phase != "measure" || payload.Sequence != 42 || payload.SentAtNS != sentAt.UnixNano() {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestSignatureMatchesServerContract(t *testing.T) {
	got := signature("appsecret", "nonce", "1672568121910")
	const want = "2e639ae3600a48b6c595495dcf61fb88b76f485b"
	if got != want {
		t.Fatalf("signature mismatch: got %s, want %s", got, want)
	}
}
