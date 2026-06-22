package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
)

type UserActivityDao struct {
	ID       int64  `gorm:"primary_key"`
	UserId   int    `gorm:"user_id"`
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
	err := dbcommons.GetDb().Model(&UserActivityDao{}).Where("app_key=? and time_mark=?", appkey, timeMark).Count(&count).Error
	if err != nil {
		return count
	}
	return count
}

func (stat UserActivityDao) ScanByTimeMarkAfterID(timeMark, lastID int64, limit int) ([]UserActivityScanRow, error) {
	rows := []UserActivityScanRow{}
	err := dbcommons.GetDb().Table(stat.TableName()).
		Select("id, app_key, time_mark").
		Where("time_mark = ? and id > ?", timeMark, lastID).
		Order("id asc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (stat UserActivityDao) DeleteBeforeTimeMarkBatch(timeMark int64, limit int) (int64, error) {
	result := dbcommons.GetDb().Where("time_mark < ?", timeMark).Limit(limit).Delete(&UserActivityDao{})
	return result.RowsAffected, result.Error
}
