package commonservices

import (
	"im-server/commons/bases"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/logs"
	"time"
)

const (
	userActivityCleanupRouteMethod         = "connect"
	userActivityCleanupGlobalRouteKey      = "useractivity-cleanup"
	userActivityCleanupRetainDays          = 3
	msgRealtimeStatCleanupRetainDays       = 3
	userActivityCleanupScheduleHour        = 5
	userActivityCleanupDeleteBatchSize     = 10000
	userActivityCleanupAggregateBatchSize  = 10000
	userActivityCleanupDeleteBatchPause    = 200 * time.Millisecond
	userActivityCleanupAggregateBatchPause = 200 * time.Millisecond
	connectCountCompactScanBatchSize       = 10000
	connectCountCompactDeleteBatchSize     = 10000
	connectCountCompactBatchPause          = 200 * time.Millisecond
)

var (
	userActivityCleanupTimer *time.Timer
	userActivityCleanupStop  chan struct{}
)

func StartUserActivityCleanupTask() {
	StopUserActivityCleanupTask()

	userActivityCleanupStop = make(chan struct{})
	go runUserActivityCleanupLoop(userActivityCleanupStop)
}

func StopUserActivityCleanupTask() {
	if userActivityCleanupTimer != nil {
		userActivityCleanupTimer.Stop()
		userActivityCleanupTimer = nil
	}
	if userActivityCleanupStop != nil {
		close(userActivityCleanupStop)
		userActivityCleanupStop = nil
	}
}

func runUserActivityCleanupLoop(stop <-chan struct{}) {
	for {
		nextRun := nextUserActivityCleanupTime(time.Now())
		userActivityCleanupTimer = time.NewTimer(time.Until(nextRun))
		select {
		case <-userActivityCleanupTimer.C:
			cleanupExpiredUserActivities(nextRun)
		case <-stop:
			return
		}
	}
}

func cleanupExpiredUserActivities(now time.Time) {
	userActivityDao := dbs.UserActivityDao{}
	aggregatePreviousDayActivities(userActivityDao, now)
	cleanupRetainedUserActivities(userActivityDao, now)
	cleanupRetainedMsgRealtimeStats(dbs.MsgRealtimeStatDao{}, now)
	compactRetainedConnectCounts(dbs.ConnectCountDao{}, now)
}

func aggregatePreviousDayActivities(userActivityDao dbs.UserActivityDao, now time.Time) {
	if !canRunGlobalUserActivityMaintenance() {
		return
	}
	previousDayMark := userActivityPreviousDayTimeMark(now)
	var lastID int64
	for {
		rows, err := userActivityDao.ScanByTimeMarkAfterID(previousDayMark, lastID, userActivityCleanupAggregateBatchSize)
		if err != nil {
			logs.NewLogEntity().Errorf("scan previous day useractivities failed, time_mark:%d, last_id:%d, err:%v", previousDayMark, lastID, err)
			return
		}
		if len(rows) == 0 {
			return
		}

		aggregates := aggregateUserActivityScanRows(rows)
		if err := incrUpsertDailyActivityBuckets(aggregates); err != nil {
			logs.NewLogEntity().Errorf("persist dailyactivities failed, time_mark:%d, last_id:%d, err:%v", previousDayMark, lastID, err)
			return
		}

		lastID = rows[len(rows)-1].ID
		if len(rows) < userActivityCleanupAggregateBatchSize {
			return
		}
		time.Sleep(userActivityCleanupAggregateBatchPause)
	}
}

type dailyActivityBucketKey struct {
	appKey   string
	timeMark int64
}

func aggregateUserActivityScanRows(rows []dbs.UserActivityScanRow) map[dailyActivityBucketKey]int64 {
	aggregates := map[dailyActivityBucketKey]int64{}
	for _, row := range rows {
		if row.AppKey == "" {
			continue
		}
		key := dailyActivityBucketKey{
			appKey:   row.AppKey,
			timeMark: row.TimeMark,
		}
		aggregates[key]++
	}
	return aggregates
}

func incrUpsertDailyActivityBuckets(aggregates map[dailyActivityBucketKey]int64) error {
	if len(aggregates) == 0 {
		return nil
	}

	dailyActivityDao := dbs.DailyActivityDao{}
	for key, count := range aggregates {
		if count <= 0 || key.appKey == "" {
			continue
		}
		if err := dailyActivityDao.IncrUpsert(dbs.DailyActivityDao{
			AppKey:   key.appKey,
			TimeMark: key.timeMark,
			Count:    count,
		}); err != nil {
			return err
		}
	}
	return nil
}

func cleanupRetainedUserActivities(userActivityDao dbs.UserActivityDao, now time.Time) {
	if !canRunGlobalUserActivityMaintenance() {
		return
	}
	cutoff := userActivityRetentionCutoff(now)
	for {
		affected, err := userActivityDao.DeleteBeforeTimeMarkBatch(cutoff, userActivityCleanupDeleteBatchSize)
		if err != nil {
			logs.NewLogEntity().Errorf("cleanup useractivities failed, cutoff:%d, err:%v", cutoff, err)
			return
		}
		if affected == 0 {
			return
		}
		time.Sleep(userActivityCleanupDeleteBatchPause)
	}
}

func cleanupRetainedMsgRealtimeStats(msgRealtimeStatDao dbs.MsgRealtimeStatDao, now time.Time) {
	if !canRunGlobalUserActivityMaintenance() {
		return
	}
	cutoff := msgRealtimeStatRetentionCutoff(now)
	for {
		affected, err := msgRealtimeStatDao.DeleteBeforeTimeMarkBatch(cutoff, userActivityCleanupDeleteBatchSize)
		if err != nil {
			logs.NewLogEntity().Errorf("cleanup msgrealtimestats failed, cutoff:%d, err:%v", cutoff, err)
			return
		}
		if affected == 0 {
			return
		}
		time.Sleep(userActivityCleanupDeleteBatchPause)
	}
}

func compactRetainedConnectCounts(connectCountDao dbs.ConnectCountDao, now time.Time) {
	if !canRunGlobalUserActivityMaintenance() {
		return
	}
	start, end := connectCountCompactDayRange(now)
	keepers, err := findConnectCountKeepers(connectCountDao, start, end)
	if err != nil {
		logs.NewLogEntity().Errorf("find compact connectcounts keepers failed, start:%d, end:%d, err:%v", start, end, err)
		return
	}
	if len(keepers) == 0 {
		return
	}
	if err := deleteCompactConnectCountRows(connectCountDao, start, end, keepers); err != nil {
		logs.NewLogEntity().Errorf("delete compact connectcounts rows failed, start:%d, end:%d, err:%v", start, end, err)
	}
}

type connectCountKeeperKey struct {
	appKey      string
	connectType int
}

func findConnectCountKeepers(connectCountDao dbs.ConnectCountDao, start, end int64) (map[connectCountKeeperKey]int64, error) {
	keepers := map[connectCountKeeperKey]dbs.ConnectCountDao{}
	lastTimeMark := start
	var lastID int64
	for {
		rows, err := connectCountDao.ScanByTimeRangeAfterCursor(start, end, lastTimeMark, lastID, connectCountCompactScanBatchSize)
		if err != nil {
			return nil, err
		}
		if len(rows) == 0 {
			break
		}
		for i := range rows {
			row := rows[i]
			if row.AppKey == "" {
				continue
			}
			key := connectCountKeeperKey{
				appKey:      row.AppKey,
				connectType: row.ConnectType,
			}
			current, ok := keepers[key]
			if !ok || isBetterConnectCountKeeper(&row, &current) {
				keepers[key] = row
			}
		}
		last := rows[len(rows)-1]
		lastTimeMark = last.TimeMark
		lastID = last.ID
		if len(rows) < connectCountCompactScanBatchSize {
			break
		}
		time.Sleep(connectCountCompactBatchPause)
	}

	keeperIDs := map[connectCountKeeperKey]int64{}
	for key, keeper := range keepers {
		keeperIDs[key] = keeper.ID
	}
	return keeperIDs, nil
}

func deleteCompactConnectCountRows(connectCountDao dbs.ConnectCountDao, start, end int64, keepers map[connectCountKeeperKey]int64) error {
	keeperIDs := connectCountKeeperIDList(keepers)
	if len(keeperIDs) == 0 {
		return nil
	}
	for {
		affected, err := connectCountDao.DeleteByTimeRangeExceptKeepersBatch(start, end, keeperIDs, connectCountCompactDeleteBatchSize)
		if err != nil {
			return err
		}
		if affected == 0 {
			return nil
		}
		time.Sleep(connectCountCompactBatchPause)
	}
}

func connectCountKeeperIDList(keepers map[connectCountKeeperKey]int64) []int64 {
	ids := make([]int64, 0, len(keepers))
	for _, id := range keepers {
		ids = append(ids, id)
	}
	return ids
}

func connectCountCompactDayRange(now time.Time) (int64, int64) {
	end := userActivityRetentionCutoff(now)
	return end - oneDay, end
}

func connectCountCompactDayEnd(dayMark int64) int64 {
	return dayMark + oneDay
}

func isBetterConnectCountKeeper(candidate, current *dbs.ConnectCountDao) bool {
	if candidate == nil {
		return false
	}
	if current == nil {
		return true
	}
	if candidate.Count != current.Count {
		return candidate.Count > current.Count
	}
	return candidate.ID < current.ID
}

func canRunGlobalUserActivityMaintenance() bool {
	return canCleanupUserActivities(userActivityCleanupGlobalRouteKey)
}

func canCleanupUserActivities(appkey string) bool {
	cluster := bases.GetCluster()
	if cluster == nil || cluster.GetCurrentNode() == nil {
		return false
	}
	node := cluster.GetTargetNode(userActivityCleanupRouteMethod, appkey)
	if node == nil {
		return false
	}
	return node.Name == cluster.GetCurrentNode().Name
}

func userActivityPreviousDayTimeMark(now time.Time) int64 {
	return dayTimeMark(now) - oneDay
}

func userActivityRetentionCutoff(now time.Time) int64 {
	return dayTimeMark(now) - int64(userActivityCleanupRetainDays-1)*oneDay
}

func msgRealtimeStatRetentionCutoff(now time.Time) int64 {
	return dayTimeMark(now) - int64(msgRealtimeStatCleanupRetainDays-1)*oneDay
}

func dayTimeMark(now time.Time) int64 {
	return now.UnixMilli() / oneDay * oneDay
}

func nextUserActivityCleanupTime(now time.Time) time.Time {
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), userActivityCleanupScheduleHour, 0, 0, 0, now.Location())
	if !nextRun.After(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}
	return nextRun
}
