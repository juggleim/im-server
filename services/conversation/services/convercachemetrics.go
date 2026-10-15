package services

import (
	"expvar"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Scan in small batches so /debug/vars never walks the entire cache or holds
// a user's lock for the duration of an HTTP request. These figures describe
// one completed scan, not an atomic snapshot of a changing cache.
const (
	converCacheScanBatchSize  = 64
	converCacheScanInterval   = 100 * time.Millisecond
	converCacheRescanInterval = 5 * time.Minute
	converMentionSampleSize   = 1024
)

type userConverCacheStats struct {
	CompletedAt                  string `json:"completed_at"`
	ScanDurationMS               int64  `json:"scan_duration_ms"`
	UsersScanned                 int    `json:"users_scanned"`
	UsersSkipped                 int    `json:"users_skipped"`
	ConversationsTotal           int64  `json:"conversations_total"`
	ConversationsP50             int    `json:"conversations_per_user_p50"`
	ConversationsP95             int    `json:"conversations_per_user_p95"`
	ConversationsP99             int    `json:"conversations_per_user_p99"`
	ConversationsMax             int    `json:"conversations_per_user_max"`
	UsersOver3000                int    `json:"users_over_3000"`
	UsersOver3100                int    `json:"users_over_3100"`
	ConversationsInUsersOver3000 int64  `json:"conversations_in_users_over_3000"`
	ConversationsInUsersOver3100 int64  `json:"conversations_in_users_over_3100"`
	MentionUsersSampled          int    `json:"mention_users_sampled"`
	MentionMessagesSampled       int64  `json:"mention_messages_sampled"`
	UsersWithMentionInSample     int    `json:"users_with_mention_in_sample"`
	MentionMessagesSampleMax     int    `json:"mention_messages_per_user_sample_max"`
}

var converCacheMetrics struct {
	sync.RWMutex
	last     userConverCacheStats
	progress int
	target   int
}

func init() {
	expvar.Publish("user_conver_cache", expvar.Func(func() interface{} {
		converCacheMetrics.RLock()
		last, progress, target := converCacheMetrics.last, converCacheMetrics.progress, converCacheMetrics.target
		converCacheMetrics.RUnlock()
		return struct {
			UsersCachedNow int                  `json:"users_cached_now"`
			Goroutines     int                  `json:"goroutines"`
			ScanProgress   int                  `json:"scan_progress"`
			ScanTarget     int                  `json:"scan_target"`
			LastComplete   userConverCacheStats `json:"last_complete"`
		}{userConverCache.Len(), runtime.NumGoroutine(), progress, target, last}
	}))
}

func scanUserConverCachePeriodically() {
	for {
		// Snapshot keys first, then fetch entry pointers one small batch at a
		// time. Retaining all values for the whole scan would delay GC on an
		// already memory-constrained node.
		keys := userConverCache.Keys()
		started := time.Now()
		counts := make([]int, 0, len(keys))
		stats := userConverCacheStats{}
		mentionSampleStride := (len(keys) + converMentionSampleSize - 1) / converMentionSampleSize
		if mentionSampleStride == 0 {
			mentionSampleStride = 1
		}
		converCacheMetrics.Lock()
		converCacheMetrics.progress, converCacheMetrics.target = 0, len(keys)
		converCacheMetrics.Unlock()
		for start := 0; start < len(keys); start += converCacheScanBatchSize {
			end := start + converCacheScanBatchSize
			if end > len(keys) {
				end = len(keys)
			}
			entries := userConverCache.SnapshotForKeys(keys[start:end])
			for offset, entry := range entries {
				if !entry.Found {
					stats.UsersSkipped++
					continue
				}
				cacheKey, ok := entry.Key.(string)
				if !ok {
					stats.UsersSkipped++
					continue
				}
				uc, ok := entry.Value.(*UserConversations)
				if !ok || uc == nil {
					stats.UsersSkipped++
					continue
				}
				lock := userLocks.GetLocks(cacheKey)
				lock.RLock()
				count := len(uc.ConverItemMap)
				// Inspect mention objects for a bounded, evenly spaced sample.
				// A full walk can be expensive on a 100k-user cache.
				if (start+offset)%mentionSampleStride == 0 {
					mentionCount := 0
					for _, conver := range uc.ConverItemMap {
						if conver != nil && conver.MentionInfo != nil {
							mentionCount += len(conver.MentionInfo.MentionMsgs)
						}
					}
					stats.MentionUsersSampled++
					stats.MentionMessagesSampled += int64(mentionCount)
					if mentionCount > 0 {
						stats.UsersWithMentionInSample++
					}
					if mentionCount > stats.MentionMessagesSampleMax {
						stats.MentionMessagesSampleMax = mentionCount
					}
				}
				lock.RUnlock()
				counts = append(counts, count)
				stats.ConversationsTotal += int64(count)
				if count > 3000 {
					stats.UsersOver3000++
					stats.ConversationsInUsersOver3000 += int64(count)
				}
				if count > 3100 {
					stats.UsersOver3100++
					stats.ConversationsInUsersOver3100 += int64(count)
				}
			}
			entries = nil
			converCacheMetrics.Lock()
			converCacheMetrics.progress = end
			converCacheMetrics.Unlock()
			time.Sleep(converCacheScanInterval)
		}
		stats.UsersScanned = len(counts)
		if len(counts) > 0 {
			sort.Ints(counts)
			stats.ConversationsP50 = cachePercentile(counts, 50)
			stats.ConversationsP95 = cachePercentile(counts, 95)
			stats.ConversationsP99 = cachePercentile(counts, 99)
			stats.ConversationsMax = counts[len(counts)-1]
		}
		stats.ScanDurationMS = time.Since(started).Milliseconds()
		stats.CompletedAt = time.Now().Format(time.RFC3339)
		converCacheMetrics.Lock()
		converCacheMetrics.last = stats
		converCacheMetrics.progress, converCacheMetrics.target = 0, 0
		converCacheMetrics.Unlock()
		time.Sleep(converCacheRescanInterval)
	}
}

// cachePercentile returns the nearest-rank percentile of sorted counts.
func cachePercentile(sorted []int, percentile int) int {
	if len(sorted) == 0 {
		return 0
	}
	index := (len(sorted)*percentile+99)/100 - 1
	return sorted[index]
}
