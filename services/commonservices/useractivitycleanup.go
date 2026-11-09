package commonservices

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/logs"
	"time"
)

const (
	userActivityCleanupRouteMethod      = "connect"
	userActivityCleanupGlobalRouteKey   = "useractivity-cleanup"
	userActivityCleanupRetainDays       = 3
	msgRealtimeStatCleanupRetainDays    = 3
	userActivityCleanupScheduleHour     = 5
	userActivityCleanupDeleteBatchSize  = 10000
	userActivityCleanupDeleteBatchPause = 200 * time.Millisecond
	connectCountCompactScanBatchSize    = 10000
	connectCountCompactDeleteBatchSize  = 10000
	connectCountCompactBatchPause       = 200 * time.Millisecond
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
	if canRunGlobalUserActivityMaintenance() {
		if err := archiveUserActivityDays(userActivityDao, now); err != nil {
			logs.NewLogEntity().Errorf("archive useractivities failed, err:%v", err)
		}
	}
	cleanupRetainedMsgRealtimeStats(dbs.MsgRealtimeStatDao{}, now)
	compactRetainedConnectCounts(dbs.ConnectCountDao{}, now)
}

type userActivityArchive interface {
	DaysBefore(today int64) ([]int64, error)
	ArchiveDay(day int64, removeDetails bool) error
}

func archiveUserActivityDays(storage userActivityArchive, now time.Time) error {
	today := dayTimeMark(now)
	days, err := storage.DaysBefore(today)
	if err != nil {
		return err
	}
	for _, day := range days {
		// Only closed UTC days are eligible, regardless of malformed legacy data.
		if day <= 0 || day >= today || day%oneDay != 0 {
			continue
		}
		if err := storage.ArchiveDay(day, day < userActivityRetentionCutoff(now)); err != nil {
			return fmt.Errorf("archive activity day %d: %w", day, err)
		}
	}
	return nil
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
