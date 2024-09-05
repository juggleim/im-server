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
