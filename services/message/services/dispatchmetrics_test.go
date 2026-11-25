package services

import (
	"encoding/json"
	"expvar"
	"testing"
	"time"
)

func TestDispatchTimingBucketsAndExpvar(t *testing.T) {
	var metric dispatchTiming
	metric.observe(2 * time.Millisecond)
	metric.observe(150 * time.Millisecond)
	metric.observe(6 * time.Second)
	got := metric.snapshot()
	if got["count"] != uint64(3) || got["max_us"] != uint64(6_000_000) {
		t.Fatalf("unexpected timing snapshot: %v", got)
	}
	buckets := got["buckets"].([]uint64)
	if buckets[1] != 1 || buckets[5] != 1 || buckets[8] != 1 {
		t.Fatalf("unexpected timing buckets: %v", buckets)
	}
	if !json.Valid([]byte(expvar.Get("msg_dispatch").String())) {
		t.Fatal("msg_dispatch expvar value is not valid JSON")
	}
}

func TestDispatchTargetMetricBuckets(t *testing.T) {
	var metric dispatchTargetMetric
	for _, targetCount := range []int{1, 100, 500, 1000, 5000, 5001} {
		metric.observe(targetCount)
	}
	got := metric.snapshot()
	if got["count"] != uint64(6) || got["sum"] != uint64(11602) || got["max"] != uint64(5001) {
		t.Fatalf("unexpected target-count snapshot: %v", got)
	}
	buckets := got["buckets"].([]uint64)
	for i, count := range buckets {
		if count != 1 {
			t.Fatalf("bucket %d count = %d, want 1; all buckets: %v", i, count, buckets)
		}
	}
}
