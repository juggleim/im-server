package services

import (
	"expvar"
	"sync/atomic"
	"time"
)

// Dispatch timings are exposed through /debug/vars without a label per user or pool.
// Values are cumulative; compute rates and approximate latency quantiles from successive snapshots.
// Bucket boundaries in milliseconds: <=1, <=5, <=10, <=50, <=100,
// <=500, <=1000, <=5000, >5000.
type dispatchTiming struct {
	count     atomic.Uint64
	sumMicros atomic.Uint64
	maxMicros atomic.Uint64
	over100ms atomic.Uint64
	over1s    atomic.Uint64
	buckets   [9]atomic.Uint64
}

// Target-count bucket boundaries: <=1, <=100, <=500, <=1000, <=5000, >5000.
type dispatchTargetMetric struct {
	count   atomic.Uint64
	sum     atomic.Uint64
	max     atomic.Uint64
	buckets [6]atomic.Uint64
}

func (metric *dispatchTargetMetric) observe(targetCount int) {
	if targetCount < 0 {
		targetCount = 0
	}
	n := uint64(targetCount)
	metric.count.Add(1)
	metric.sum.Add(n)
	for old := metric.max.Load(); n > old; old = metric.max.Load() {
		if metric.max.CompareAndSwap(old, n) {
			break
		}
	}
	for i, bound := range [...]uint64{1, 100, 500, 1000, 5000} {
		if n <= bound {
			metric.buckets[i].Add(1)
			return
		}
	}
	metric.buckets[5].Add(1)
}

func (metric *dispatchTargetMetric) snapshot() map[string]interface{} {
	buckets := make([]uint64, len(metric.buckets))
	for i := range buckets {
		buckets[i] = metric.buckets[i].Load()
	}
	return map[string]interface{}{
		"count": metric.count.Load(), "sum": metric.sum.Load(),
		"max": metric.max.Load(), "buckets": buckets,
	}
}

func (metric *dispatchTiming) observe(duration time.Duration) {
	micros := uint64(duration.Microseconds())
	metric.count.Add(1)
	metric.sumMicros.Add(micros)
	if duration >= 100*time.Millisecond {
		metric.over100ms.Add(1)
	}
	if duration >= time.Second {
		metric.over1s.Add(1)
	}
	metric.updateMax(micros)
	for i, bound := range [...]time.Duration{
		time.Millisecond, 5 * time.Millisecond, 10 * time.Millisecond,
		50 * time.Millisecond, 100 * time.Millisecond, 500 * time.Millisecond,
		time.Second, 5 * time.Second,
	} {
		if duration <= bound {
			metric.buckets[i].Add(1)
			return
		}
	}
	metric.buckets[8].Add(1)
}

func (metric *dispatchTiming) updateMax(micros uint64) {
	for old := metric.maxMicros.Load(); micros > old; old = metric.maxMicros.Load() {
		if metric.maxMicros.CompareAndSwap(old, micros) {
			break
		}
	}
}

func (metric *dispatchTiming) snapshot() map[string]interface{} {
	buckets := make([]uint64, len(metric.buckets))
	for i := range buckets {
		buckets[i] = metric.buckets[i].Load()
	}
	return map[string]interface{}{
		"count": metric.count.Load(), "sum_us": metric.sumMicros.Load(),
		"max_us": metric.maxMicros.Load(), "over_100ms": metric.over100ms.Load(),
		"over_1s": metric.over1s.Load(), "buckets": buckets,
	}
}

var dispatchTimings = struct {
	actorSubmit        dispatchTiming
	fanoutSubmit       dispatchTiming
	fanoutQueueWait    dispatchTiming
	fanoutTask         dispatchTiming
	preheat            dispatchTiming
	fanout             dispatchTiming
	submit             dispatchTiming
	queueWait          dispatchTiming
	task               dispatchTiming
	conversation       dispatchTiming
	userInfo           dispatchTiming
	inbox              dispatchTiming
	downstream         dispatchTiming
	maxQueueLen        atomic.Uint64
	queueFullOnSubmit  atomic.Uint64
	fanoutFullOnSubmit atomic.Uint64
}{}

var dispatchTargets dispatchTargetMetric

func observeDispatchQueueLen(length int) {
	n := uint64(length)
	for old := dispatchTimings.maxQueueLen.Load(); n > old; old = dispatchTimings.maxQueueLen.Load() {
		if dispatchTimings.maxQueueLen.CompareAndSwap(old, n) {
			break
		}
	}
}

func init() {
	expvar.Publish("msg_dispatch", expvar.Func(func() interface{} {
		return map[string]interface{}{
			"actor_submit":          dispatchTimings.actorSubmit.snapshot(),
			"fanout_submit":         dispatchTimings.fanoutSubmit.snapshot(),
			"fanout_queue_wait":     dispatchTimings.fanoutQueueWait.snapshot(),
			"fanout_task":           dispatchTimings.fanoutTask.snapshot(),
			"preheat":               dispatchTimings.preheat.snapshot(),
			"fanout":                dispatchTimings.fanout.snapshot(),
			"targets_per_call":      dispatchTargets.snapshot(),
			"submit":                dispatchTimings.submit.snapshot(),
			"queue_wait":            dispatchTimings.queueWait.snapshot(),
			"task":                  dispatchTimings.task.snapshot(),
			"conversation":          dispatchTimings.conversation.snapshot(),
			"user_info":             dispatchTimings.userInfo.snapshot(),
			"inbox":                 dispatchTimings.inbox.snapshot(),
			"downstream":            dispatchTimings.downstream.snapshot(),
			"max_queue_len_ever":    dispatchTimings.maxQueueLen.Load(),
			"queue_full_on_submit":  dispatchTimings.queueFullOnSubmit.Load(),
			"fanout_full_on_submit": dispatchTimings.fanoutFullOnSubmit.Load(),
			"queues":                MsgSinglePools.QueueStats(),
			"fanout_scheduler":      msgDispatchFanoutScheduler.Stats(),
		}
	}))
}
