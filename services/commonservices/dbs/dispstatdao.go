package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
)

type DispStatDao struct {
	ID          int64  `gorm:"primary_key"`
	ChannelType int    `gorm:"channel_type"`
	TimeMark    int64  `gorm:"time_mark"`
	Count       int64  `gorm:"count"`
	AppKey      string `gorm:"app_key"`
}

func (stat DispStatDao) TableName() string {
	return "dispstats"
}
func (stat DispStatDao) Create(item DispStatDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}
func (stat DispStatDao) IncrByStep(appkey string, channelType int, timeMark, step int64) error {
	sql := fmt.Sprintf("insert into %s (channel_type,time_mark,count,app_key)values(?,?,?,?) ON DUPLICATE KEY UPDATE count=count+?", stat.TableName())
	err := dbcommons.GetDb().Exec(sql, channelType, timeMark, step, appkey, step).Error
	return err
}
