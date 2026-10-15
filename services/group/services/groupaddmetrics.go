package services

import (
	"context"
	"expvar"
	"sync/atomic"
	"time"

	"im-server/commons/errs"
	"im-server/services/commonservices/logs"

	"golang.org/x/time/rate"
)

const groupAddSlowThreshold = 500 * time.Millisecond

type groupAddTiming struct {
	count     atomic.Uint64
	sumMicros atomic.Uint64
	maxMicros atomic.Uint64
	over100ms atomic.Uint64
	over1s    atomic.Uint64
}

func (metric *groupAddTiming) observe(duration time.Duration) {
	micros := uint64(duration.Microseconds())
	metric.count.Add(1)
	metric.sumMicros.Add(micros)
	if duration >= 100*time.Millisecond {
		metric.over100ms.Add(1)
	}
	if duration >= time.Second {
		metric.over1s.Add(1)
	}
	for old := metric.maxMicros.Load(); micros > old; old = metric.maxMicros.Load() {
		if metric.maxMicros.CompareAndSwap(old, micros) {
			break
		}
	}
}

func (metric *groupAddTiming) snapshot() map[string]uint64 {
	return map[string]uint64{
		"count": metric.count.Load(), "sum_us": metric.sumMicros.Load(),
		"max_us": metric.maxMicros.Load(), "over_100ms": metric.over100ms.Load(),
		"over_1s": metric.over1s.Load(),
	}
}

type groupAddPhaseDurations struct {
	groupInfo             time.Duration
	memberCache           time.Duration
	memberCacheLoad       time.Duration
	memberCacheSharedWait time.Duration
	memberStateLockWait   time.Duration
	memberLockWait        time.Duration
	memberFilter          time.Duration
	coldCountDB           time.Duration
	coldExistingDB        time.Duration
	cacheUpdate           time.Duration
	memberDB              time.Duration
	postPersistSchedule   time.Duration
	conversationSchedule  time.Duration
	coldPath              bool
	coldPathFailed        bool
}

var groupAddMetrics = struct {
	total                 groupAddTiming
	groupInfo             groupAddTiming
	memberCache           groupAddTiming
	memberCacheLoad       groupAddTiming
	memberCacheSharedWait groupAddTiming
	memberStateLockWait   groupAddTiming
	memberLockWait        groupAddTiming
	memberFilter          groupAddTiming
	coldCountDB           groupAddTiming
	coldExistingDB        groupAddTiming
	cacheUpdate           groupAddTiming
	memberDB              groupAddTiming
	snapshot              groupAddTiming
	memberExtDB           groupAddTiming
	postPersistSchedule   groupAddTiming
	conversationSchedule  groupAddTiming
	requests              atomic.Uint64
	slowRequests          atomic.Uint64
	failedRequests        atomic.Uint64
	memberCacheMisses     atomic.Uint64
	coldPathRequests      atomic.Uint64
	coldPathFailures      atomic.Uint64
	snapshotRuns          atomic.Uint64
	requestedMembers      atomic.Uint64
	addedMembers          atomic.Uint64
	maxGroupMembers       atomic.Uint64
	maxRequestedMembers   atomic.Uint64
	maxAddedMembers       atomic.Uint64
}{}

var groupAddSlowLogLimiter = rate.NewLimiter(rate.Limit(5), 5)

func updateGroupAddMax(metric *atomic.Uint64, value int) {
	if value < 0 {
		return
	}
	n := uint64(value)
	for old := metric.Load(); n > old; old = metric.Load() {
		if metric.CompareAndSwap(old, n) {
			break
		}
	}
}

func observeGroupAdd(
	ctx context.Context,
	groupID string,
	total time.Duration,
	phases groupAddPhaseDurations,
	cacheChecked bool,
	cacheHit bool,
	groupMemberCount int,
	requestedMemberCount int,
	addedMemberCount int,
	result errs.IMErrorCode,
) {
	groupAddMetrics.requests.Add(1)
	groupAddMetrics.total.observe(total)
	groupAddMetrics.groupInfo.observe(phases.groupInfo)
	groupAddMetrics.memberCache.observe(phases.memberCache)
	if phases.memberCacheLoad > 0 {
		groupAddMetrics.memberCacheLoad.observe(phases.memberCacheLoad)
	}
	if phases.memberCacheSharedWait > 0 {
		groupAddMetrics.memberCacheSharedWait.observe(phases.memberCacheSharedWait)
	}
	groupAddMetrics.memberStateLockWait.observe(phases.memberStateLockWait)
	groupAddMetrics.memberLockWait.observe(phases.memberLockWait)
	groupAddMetrics.memberFilter.observe(phases.memberFilter)
	if phases.coldPath {
		groupAddMetrics.coldPathRequests.Add(1)
		if phases.coldCountDB > 0 {
			groupAddMetrics.coldCountDB.observe(phases.coldCountDB)
		}
		if phases.coldExistingDB > 0 {
			groupAddMetrics.coldExistingDB.observe(phases.coldExistingDB)
		}
	}
	if phases.coldPathFailed {
		groupAddMetrics.coldPathFailures.Add(1)
	}
	groupAddMetrics.cacheUpdate.observe(phases.cacheUpdate)
	groupAddMetrics.memberDB.observe(phases.memberDB)
	groupAddMetrics.postPersistSchedule.observe(phases.postPersistSchedule)
	groupAddMetrics.conversationSchedule.observe(phases.conversationSchedule)
	if cacheChecked && !cacheHit {
		groupAddMetrics.memberCacheMisses.Add(1)
	}
	if result != errs.IMErrorCode_SUCCESS {
		groupAddMetrics.failedRequests.Add(1)
	}
	groupAddMetrics.requestedMembers.Add(uint64(requestedMemberCount))
	groupAddMetrics.addedMembers.Add(uint64(addedMemberCount))
	updateGroupAddMax(&groupAddMetrics.maxGroupMembers, groupMemberCount)
	updateGroupAddMax(&groupAddMetrics.maxRequestedMembers, requestedMemberCount)
	updateGroupAddMax(&groupAddMetrics.maxAddedMembers, addedMemberCount)

	if total < groupAddSlowThreshold {
		return
	}
	groupAddMetrics.slowRequests.Add(1)
	if groupAddSlowLogLimiter.Allow() {
		logs.WithContext(ctx).Warnf("[G_ADD_MEMBERS_SLOW] group_id:%s\tresult:%d\tgroup_member_count:%d\trequest_member_count:%d\tadd_member_count:%d\tmember_cache_hit:%t\tcold_path:%t\ttotal:%v\tgroup_info:%v\tmember_cache:%v\tmember_cache_load:%v\tmember_cache_shared_wait:%v\tmember_state_lock_wait:%v\tmember_lock_wait:%v\tmember_filter:%v\tcold_count_db:%v\tcold_existing_db:%v\tcache_update:%v\tmember_db:%v\tpost_persist_schedule:%v\tconversation_schedule:%v",
			groupID, result, groupMemberCount, requestedMemberCount, addedMemberCount, cacheHit, phases.coldPath,
			total, phases.groupInfo, phases.memberCache, phases.memberCacheLoad, phases.memberCacheSharedWait,
			phases.memberStateLockWait, phases.memberLockWait, phases.memberFilter,
			phases.coldCountDB, phases.coldExistingDB, phases.cacheUpdate,
			phases.memberDB, phases.postPersistSchedule, phases.conversationSchedule)
	}
}

func init() {
	expvar.Publish("g_add_members_business", expvar.Func(func() interface{} {
		return map[string]interface{}{
			"total":                    groupAddMetrics.total.snapshot(),
			"group_info":               groupAddMetrics.groupInfo.snapshot(),
			"member_cache":             groupAddMetrics.memberCache.snapshot(),
			"member_cache_load":        groupAddMetrics.memberCacheLoad.snapshot(),
			"member_cache_shared_wait": groupAddMetrics.memberCacheSharedWait.snapshot(),
			"member_state_lock_wait":   groupAddMetrics.memberStateLockWait.snapshot(),
			"member_lock_wait":         groupAddMetrics.memberLockWait.snapshot(),
			"member_filter":            groupAddMetrics.memberFilter.snapshot(),
			"cold_count_db":            groupAddMetrics.coldCountDB.snapshot(),
			"cold_existing_db":         groupAddMetrics.coldExistingDB.snapshot(),
			"cache_update":             groupAddMetrics.cacheUpdate.snapshot(),
			"member_db":                groupAddMetrics.memberDB.snapshot(),
			"snapshot":                 groupAddMetrics.snapshot.snapshot(),
			"member_ext_db":            groupAddMetrics.memberExtDB.snapshot(),
			"post_persist_schedule":    groupAddMetrics.postPersistSchedule.snapshot(),
			"conversation_schedule":    groupAddMetrics.conversationSchedule.snapshot(),
			"requests":                 groupAddMetrics.requests.Load(),
			"slow_requests":            groupAddMetrics.slowRequests.Load(),
			"failed_requests":          groupAddMetrics.failedRequests.Load(),
			"member_cache_misses":      groupAddMetrics.memberCacheMisses.Load(),
			"cold_path_requests":       groupAddMetrics.coldPathRequests.Load(),
			"cold_path_failures":       groupAddMetrics.coldPathFailures.Load(),
			"snapshot_runs":            groupAddMetrics.snapshotRuns.Load(),
			"requested_members":        groupAddMetrics.requestedMembers.Load(),
			"added_members":            groupAddMetrics.addedMembers.Load(),
			"max_group_members":        groupAddMetrics.maxGroupMembers.Load(),
			"max_requested_members":    groupAddMetrics.maxRequestedMembers.Load(),
			"max_added_members":        groupAddMetrics.maxAddedMembers.Load(),
			"post_persist":             groupMemberPostPersistStats(),
		}
	}))
}
