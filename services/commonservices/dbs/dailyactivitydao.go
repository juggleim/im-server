package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"time"
)

type DailyActivityDao struct {
	ID          int64     `gorm:"primary_key"`
	TimeMark    int64     `gorm:"time_mark"`
	Count       int64     `gorm:"count"`
	AppKey      string    `gorm:"app_key"`
	CreatedTime time.Time `gorm:"created_time"`
}

func (stat DailyActivityDao) TableName() string {
	return "dailyactivities"
}

func (stat DailyActivityDao) IncrUpsert(item DailyActivityDao) error {
	sql := fmt.Sprintf("insert into %s (time_mark,count,app_key)values(?,?,?) ON DUPLICATE KEY UPDATE count=count+VALUES(count)", stat.TableName())
	return dbcommons.GetDb().Exec(sql, item.TimeMark, item.Count, item.AppKey).Error
}
