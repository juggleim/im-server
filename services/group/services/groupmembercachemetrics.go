package services

import (
	"context"
	"expvar"
	"sync/atomic"
	"time"

	"im-server/services/commonservices/logs"

	"golang.org/x/time/rate"
)

const groupMemberCacheSlowLoadThreshold = 500 * time.Millisecond

type groupMemberCacheLoadStats struct {
	cacheHits       atomic.Uint64
	cacheMisses     atomic.Uint64
	leaderLoads     atomic.Uint64
	recheckHits     atomic.Uint64
	sharedWaits     atomic.Uint64
	loadFailures    atomic.Uint64
	dbPages         atomic.Uint64
	dbRows          atomic.Uint64
	maxPagesPerLoad atomic.Uint64
	maxRowsPerLoad  atomic.Uint64
	stateLockWait   groupAddTiming
	sharedWait      groupAddTiming
	load            groupAddTiming
	dbQuery         groupAddTiming
	mapBuild        groupAddTiming
}

var groupMemberCacheMetrics groupMemberCacheLoadStats
var groupMemberCacheSlowLogLimiter = rate.NewLimiter(rate.Limit(5), 5)

func updateGroupMemberCacheMax(metric *atomic.Uint64, value int) {
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

func observeGroupMemberDBPage(queryDuration, buildDuration time.Duration, rows int, failed bool) {
	groupMemberCacheMetrics.dbQuery.observe(queryDuration)
	if failed {
		groupMemberCacheMetrics.loadFailures.Add(1)
		return
	}
	groupMemberCacheMetrics.mapBuild.observe(buildDuration)
	groupMemberCacheMetrics.dbPages.Add(1)
	groupMemberCacheMetrics.dbRows.Add(uint64(rows))
}

func observeGroupMemberDBLoad(pages, rows int) {
	updateGroupMemberCacheMax(&groupMemberCacheMetrics.maxPagesPerLoad, pages)
	updateGroupMemberCacheMax(&groupMemberCacheMetrics.maxRowsPerLoad, rows)
}

func logSlowGroupMemberDBLoad(ctx context.Context, groupID string, total, query, build time.Duration, pages, rows int) {
	if total < groupMemberCacheSlowLoadThreshold || !groupMemberCacheSlowLogLimiter.Allow() {
		return
	}
	logs.WithContext(ctx).Warnf("[GROUP_MEMBERS_CACHE_LOAD_SLOW] group_id:%s\tmember_count:%d\tpages:%d\ttotal:%v\tdb_query:%v\tmap_build:%v",
		groupID, rows, pages, total, query, build)
}

func init() {
	expvar.Publish("group_members_cache", expvar.Func(func() interface{} {
		swrStats, _ := groupMembersCache.StaleWhileRevalidateStats()
		return map[string]interface{}{
			"cache_size":             groupMembersCache.Len(),
			"cache_hits":             groupMemberCacheMetrics.cacheHits.Load(),
			"cache_misses":           groupMemberCacheMetrics.cacheMisses.Load(),
			"leader_loads":           groupMemberCacheMetrics.leaderLoads.Load(),
			"recheck_hits":           groupMemberCacheMetrics.recheckHits.Load(),
			"shared_waits":           groupMemberCacheMetrics.sharedWaits.Load(),
			"load_failures":          groupMemberCacheMetrics.loadFailures.Load(),
			"db_pages":               groupMemberCacheMetrics.dbPages.Load(),
			"db_rows":                groupMemberCacheMetrics.dbRows.Load(),
			"max_pages_per_load":     groupMemberCacheMetrics.maxPagesPerLoad.Load(),
			"max_rows_per_load":      groupMemberCacheMetrics.maxRowsPerLoad.Load(),
			"state_lock_wait":        groupMemberCacheMetrics.stateLockWait.snapshot(),
			"shared_wait":            groupMemberCacheMetrics.sharedWait.snapshot(),
			"load":                   groupMemberCacheMetrics.load.snapshot(),
			"db_query":               groupMemberCacheMetrics.dbQuery.snapshot(),
			"map_build":              groupMemberCacheMetrics.mapBuild.snapshot(),
			"active_state_locks":     groupMemberStateLocks.size(),
			"stale_while_revalidate": swrStats,
		}
	}))
}
