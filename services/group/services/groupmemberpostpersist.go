package services

import (
	"context"
	"sync/atomic"
	"time"

	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/group/dbs"

	"golang.org/x/time/rate"
)

const (
	groupMemberPostPersistPoolCount     = 64
	groupMemberPostPersistQueueSize     = 500
	groupMemberPostPersistWarnThreshold = 500 * time.Millisecond
	groupMemberPostPersistMaxAttempts   = 3
)

// The queue is deliberately bounded and blocking. Saturation applies backpressure
// to g_add_members instead of silently dropping database work. Tasks for the same
// app/group always land on the same worker and therefore keep submission order.
var groupMemberPostPersistPools = tools.NewSinglePools(
	groupMemberPostPersistPoolCount,
	groupMemberPostPersistQueueSize,
	true,
)

var groupMemberPostPersistMetrics = struct {
	enqueue           groupAddTiming
	queueWait         groupAddTiming
	task              groupAddTiming
	snapshot          groupAddTiming
	memberExtDB       groupAddTiming
	submitted         atomic.Uint64
	completed         atomic.Uint64
	memberExtFailures atomic.Uint64
	memberExtRetries  atomic.Uint64
	queueFullOnSubmit atomic.Uint64
	maxQueueLen       atomic.Uint64
}{}

var groupMemberPostPersistLogLimiter = rate.NewLimiter(rate.Limit(5), 5)

func observeGroupMemberPostPersistQueueLen(length int) {
	n := uint64(length)
	for old := groupMemberPostPersistMetrics.maxQueueLen.Load(); n > old; old = groupMemberPostPersistMetrics.maxQueueLen.Load() {
		if groupMemberPostPersistMetrics.maxQueueLen.CompareAndSwap(old, n) {
			break
		}
	}
}

func submitGroupMemberPostPersist(
	ctx context.Context,
	appKey string,
	groupID string,
	memberExts []dbs.GroupMemberExtDao,
	generateSnapshot bool,
	snapshotMemberIDs []string,
) {
	poolKey := appKey + "\x00" + groupID
	poolIndex := groupMemberPostPersistPools.PoolIndex(poolKey)
	pool := groupMemberPostPersistPools.GetPool(poolKey)
	queueLenOnSubmit := pool.QueueLen()
	if queueLenOnSubmit == pool.QueueCapacity() {
		groupMemberPostPersistMetrics.queueFullOnSubmit.Add(1)
	}

	enqueued := time.Now()
	groupMemberPostPersistMetrics.submitted.Add(1)
	pool.Submit(func() {
		queueWait := time.Since(enqueued)
		groupMemberPostPersistMetrics.queueWait.observe(queueWait)
		if queueWait >= groupMemberPostPersistWarnThreshold && groupMemberPostPersistLogLimiter.Allow() {
			logs.WithContext(ctx).Warnf("[G_ADD_MEMBERS_POST_PERSIST_QUEUE_WAIT] group_id:%s\tpool_index:%d\tqueue_wait:%v\tqueue_len_on_submit:%d\tqueue_len_on_start:%d\tqueue_capacity:%d",
				groupID, poolIndex, queueWait, queueLenOnSubmit, pool.QueueLen(), pool.QueueCapacity())
		}

		taskStarted := time.Now()
		defer func() {
			taskDuration := time.Since(taskStarted)
			groupMemberPostPersistMetrics.task.observe(taskDuration)
			groupMemberPostPersistMetrics.completed.Add(1)
			if taskDuration >= groupMemberPostPersistWarnThreshold && groupMemberPostPersistLogLimiter.Allow() {
				logs.WithContext(ctx).Warnf("[G_ADD_MEMBERS_POST_PERSIST_SLOW] group_id:%s\ttask:%v", groupID, taskDuration)
			}
		}()

		if generateSnapshot {
			started := time.Now()
			GenerateGroupSnapshot(appKey, groupID, snapshotMemberIDs)
			duration := time.Since(started)
			groupMemberPostPersistMetrics.snapshot.observe(duration)
			groupAddMetrics.snapshot.observe(duration)
			groupAddMetrics.snapshotRuns.Add(1)
		}

		if len(memberExts) > 0 {
			started := time.Now()
			var err error
			for attempt := 1; attempt <= groupMemberPostPersistMaxAttempts; attempt++ {
				err = (dbs.GroupMemberExtDao{}).BatchCreate(memberExts)
				if err == nil {
					break
				}
				if attempt < groupMemberPostPersistMaxAttempts {
					groupMemberPostPersistMetrics.memberExtRetries.Add(1)
					time.Sleep(time.Duration(attempt) * 50 * time.Millisecond)
				}
			}
			duration := time.Since(started)
			groupMemberPostPersistMetrics.memberExtDB.observe(duration)
			groupAddMetrics.memberExtDB.observe(duration)
			if err != nil {
				groupMemberPostPersistMetrics.memberExtFailures.Add(1)
				logs.WithContext(ctx).Errorf("[G_ADD_MEMBERS_POST_PERSIST_FAILED] group_id:%s\terr:%v", groupID, err)
			}
		}
	})
	groupMemberPostPersistMetrics.enqueue.observe(time.Since(enqueued))
	observeGroupMemberPostPersistQueueLen(pool.QueueLen())
}

func groupMemberPostPersistStats() map[string]interface{} {
	return map[string]interface{}{
		"enqueue":              groupMemberPostPersistMetrics.enqueue.snapshot(),
		"queue_wait":           groupMemberPostPersistMetrics.queueWait.snapshot(),
		"task":                 groupMemberPostPersistMetrics.task.snapshot(),
		"snapshot":             groupMemberPostPersistMetrics.snapshot.snapshot(),
		"member_ext_db":        groupMemberPostPersistMetrics.memberExtDB.snapshot(),
		"submitted":            groupMemberPostPersistMetrics.submitted.Load(),
		"completed":            groupMemberPostPersistMetrics.completed.Load(),
		"member_ext_failures":  groupMemberPostPersistMetrics.memberExtFailures.Load(),
		"member_ext_retries":   groupMemberPostPersistMetrics.memberExtRetries.Load(),
		"queue_full_on_submit": groupMemberPostPersistMetrics.queueFullOnSubmit.Load(),
		"max_queue_len_ever":   groupMemberPostPersistMetrics.maxQueueLen.Load(),
		"queues":               groupMemberPostPersistPools.QueueStats(),
	}
}
