package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
)

type DownStatDao struct {
	ID          int64  `gorm:"primary_key"`
	ChannelType int    `gorm:"channel_type"`
	TimeMark    int64  `gorm:"time_mark"`
	Count       int64  `gorm:"count"`
	AppKey      string `gorm:"app_key"`
}

func (stat DownStatDao) TableName() string {
	return "downstats"
}
func (stat DownStatDao) Create(item DownStatDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}
func (stat DownStatDao) IncrByStep(appkey string, channelType int, timeMark, step int64) error {
	sql := fmt.Sprintf("insert into %s (channel_type,time_mark,count,app_key)values(?,?,?,?) ON DUPLICATE KEY UPDATE count=count+?", stat.TableName())
	err := dbcommons.GetDb().Exec(sql, channelType, timeMark, step, appkey, step).Error
	return err
}
