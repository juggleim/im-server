package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
)

type MsgStatDao struct {
	ID          int64  `gorm:"primary_key"`
	StatType    int    `gorm:"stat_type"`
	ChannelType int    `gorm:"channel_type"`
	TimeMark    int64  `gorm:"time_mark"`
	Count       int64  `gorm:"count"`
	AppKey      string `gorm:"app_key"`
}

func (stat MsgStatDao) TableName() string {
	return "msgstats"
}
func (stat MsgStatDao) Create(item MsgStatDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}
func (stat MsgStatDao) IncrByStep(appkey string, statType int, channelType int, timeMark, step int64) error {
	sql := fmt.Sprintf("insert into %s (stat_type,channel_type,time_mark,count,app_key)values(?,?,?,?,?) ON DUPLICATE KEY UPDATE count=count+?", stat.TableName())
	err := dbcommons.GetDb().Exec(sql, statType, channelType, timeMark, step, appkey, step).Error
	return err
}
func (stat MsgStatDao) QryStats(appkey string, statType int, channelType int, start int64, end int64) []*MsgStatDao {
	var items []*MsgStatDao
	err := dbcommons.GetDb().Where("app_key=? and stat_type=? and channel_type=? and time_mark>=? and time_mark<=?", appkey, statType, channelType, start, end).Limit(1000).Find(&items).Error
	if err == nil {
		return items
	}
	return []*MsgStatDao{}
}
