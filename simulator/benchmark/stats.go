package benchmark

import (
	"math"
	"sort"
	"sync"
	"time"
)

type LatencySummary struct {
	Count  int     `json:"count"`
	MinMS  float64 `json:"min_ms"`
	P50MS  float64 `json:"p50_ms"`
	P95MS  float64 `json:"p95_ms"`
	P99MS  float64 `json:"p99_ms"`
	MaxMS  float64 `json:"max_ms"`
	MeanMS float64 `json:"mean_ms"`
}

type MetricSnapshot struct {
	Attempted        int64            `json:"attempted"`
	Succeeded        int64            `json:"succeeded"`
	Failed           int64            `json:"failed"`
	ThroughputPerSec float64          `json:"throughput_per_second"`
	Latency          LatencySummary   `json:"latency"`
	Errors           map[string]int64 `json:"errors,omitempty"`
}

type metricRecorder struct {
	mu        sync.Mutex
	attempted int64
	succeeded int64
	errors    map[string]int64
	latencies []time.Duration
}

func newMetricRecorder() *metricRecorder {
	return &metricRecorder{errors: make(map[string]int64)}
}

func (r *metricRecorder) recordSuccess(latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.attempted++
	r.succeeded++
	r.latencies = append(r.latencies, latency)
}

func (r *metricRecorder) recordFailure(code string, latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.attempted++
	r.errors[code]++
}

func (r *metricRecorder) recordObservation(latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.attempted++
	r.succeeded++
	r.latencies = append(r.latencies, latency)
}

func (r *metricRecorder) snapshot(window time.Duration) MetricSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	errors := make(map[string]int64, len(r.errors))
	for key, value := range r.errors {
		errors[key] = value
	}
	throughput := 0.0
	if window > 0 {
		throughput = float64(r.succeeded) / window.Seconds()
	}
	return MetricSnapshot{
		Attempted:        r.attempted,
		Succeeded:        r.succeeded,
		Failed:           r.attempted - r.succeeded,
		ThroughputPerSec: throughput,
		Latency:          summarizeLatencies(r.latencies),
		Errors:           errors,
	}
}

func summarizeLatencies(values []time.Duration) LatencySummary {
	if len(values) == 0 {
		return LatencySummary{}
	}
	ordered := append([]time.Duration(nil), values...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i] < ordered[j] })
	var total time.Duration
	for _, value := range ordered {
		total += value
	}
	toMS := func(value time.Duration) float64 {
		return float64(value) / float64(time.Millisecond)
	}
	return LatencySummary{
		Count:  len(ordered),
		MinMS:  toMS(ordered[0]),
		P50MS:  toMS(nearestRank(ordered, 0.50)),
		P95MS:  toMS(nearestRank(ordered, 0.95)),
		P99MS:  toMS(nearestRank(ordered, 0.99)),
		MaxMS:  toMS(ordered[len(ordered)-1]),
		MeanMS: toMS(total / time.Duration(len(ordered))),
	}
}

func nearestRank(values []time.Duration, percentile float64) time.Duration {
	if len(values) == 0 {
		return 0
	}
	index := int(math.Ceil(percentile*float64(len(values)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(values) {
		index = len(values) - 1
	}
	return values[index]
}
