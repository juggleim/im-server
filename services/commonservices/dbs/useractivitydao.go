package dbs

import (
	"fmt"
	"gorm.io/gorm"
	"im-server/commons/dbcommons"
)

type UserActivityDao struct {
	ID       int64  `gorm:"primary_key"`
	UserId   string `gorm:"user_id"`
	TimeMark int64  `gorm:"time_mark"`
	Count    int64  `gorm:"count"`
	AppKey   string `gorm:"app_key"`
}

type UserActivityScanRow struct {
	ID       int64  `gorm:"column:id"`
	AppKey   string `gorm:"column:app_key"`
	TimeMark int64  `gorm:"column:time_mark"`
}

func (stat UserActivityDao) TableName() string {
	return "useractivities"
}

func (stat UserActivityDao) Create(item UserActivityDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (stat UserActivityDao) IncrByStep(appkey, userId string, timeMark, step int64) error {
	sql := fmt.Sprintf("insert into %s (user_id,time_mark,count,app_key)values(?,?,?,?) ON DUPLICATE KEY UPDATE count=count+?", stat.TableName())
	err := dbcommons.GetDb().Exec(sql, userId, timeMark, step, appkey, step).Error
	return err
}

func (stat UserActivityDao) CountUserActivities(appkey string, timeMark int64) int64 {
	var count int64
	err := dbcommons.GetDb().Model(&UserActivityDao{}).Where("app_key=? and time_mark IN (?, ?)", appkey, timeMark, timeMark/1000).Distinct("user_id").Count(&count).Error
	if err != nil {
		return count
	}
	return count
}

func (stat UserActivityDao) QryLatestUserActivity(appkey, userId string) (*UserActivityDao, error) {
	var items []*UserActivityDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Order("time_mark desc").Limit(1).Table(stat.TableName()).Find(&items).Error
	if err == nil && len(items) > 0 {
		return items[0], nil
	}
	return nil, err
}

func (stat UserActivityDao) ScanByTimeMarkAfterID(timeMark, lastID int64, limit int) ([]UserActivityScanRow, error) {
	rows := []UserActivityScanRow{}
	err := dbcommons.GetDb().Table(stat.TableName()).
		Select("id, app_key, time_mark").
		Where("time_mark IN (?, ?) and id > ?", timeMark, timeMark/1000, lastID).
		Order("id asc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (stat UserActivityDao) DeleteBeforeTimeMarkBatch(timeMark int64, limit int) (int64, error) {
	// Keep recent second-based rows: comparing seconds directly to a
	// millisecond cutoff would otherwise delete even today's records.
	result := dbcommons.GetDb().Where("(time_mark >= ? AND time_mark < ?) OR (time_mark > 0 AND time_mark < ?)", activityMillisBoundary, timeMark, timeMark/1000).Limit(limit).Delete(&UserActivityDao{})
	return result.RowsAffected, result.Error
}

// Activity timestamps from the IM server are modern Unix dates. Values below
// this boundary are legacy seconds, not milliseconds (the boundary is 1973 in ms).
const activityMillisBoundary int64 = 100000000000

func (stat UserActivityDao) DaysBefore(today int64) ([]int64, error) {
	return activityDaysBefore(dbcommons.GetDb(), today)
}

func activityDaysBefore(db *gorm.DB, today int64) ([]int64, error) {
	var days []int64
	err := db.Table("useractivities").
		Select("DISTINCT CASE WHEN time_mark < ? THEN time_mark * 1000 ELSE time_mark END AS day_mark", activityMillisBoundary).
		Where("(time_mark > 0 AND time_mark < ?) OR (time_mark >= ? AND time_mark < ?)", today/1000, activityMillisBoundary, today).
		Order("day_mark asc").Pluck("day_mark", &days).Error
	return days, err
}

func (stat UserActivityDao) ArchiveDay(day int64, removeDetails bool) error {
	return archiveActivityDay(dbcommons.GetDb(), day, removeDetails)
}

func archiveActivityDay(db *gorm.DB, day int64, removeDetails bool) error {
	// A complete day is the atomic unit: a failed summary or delete rolls back
	// both. Replacing the snapshot makes retries and repeated maintenance safe.
	return db.Transaction(func(tx *gorm.DB) error {
		var counts []DailyActivityDao
		if err := tx.Table("useractivities").Select("app_key, COUNT(DISTINCT user_id) AS count").
			Where("time_mark IN (?, ?)", day, day/1000).Group("app_key").Scan(&counts).Error; err != nil {
			return err
		}
		for _, count := range counts {
			if count.AppKey == "" {
				continue
			}
			if err := tx.Exec("INSERT INTO dailyactivities (app_key,time_mark,count) VALUES (?,?,?) ON DUPLICATE KEY UPDATE count=VALUES(count)", count.AppKey, day, count.Count).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM dailyactivities WHERE app_key=? AND time_mark=?", count.AppKey, day/1000).Error; err != nil {
				return err
			}
		}
		if removeDetails {
			return tx.Where("time_mark IN (?, ?)", day, day/1000).Delete(&UserActivityDao{}).Error
		}
		return nil
	})
}
