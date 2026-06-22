package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
)

type MsgRealtimeStatDao struct {
	ID          int64  `gorm:"primary_key"`
	StatType    int    `gorm:"stat_type"`
	ChannelType int    `gorm:"channel_type"`
	TimeMark    int64  `gorm:"time_mark"`
	Count       int64  `gorm:"count"`
	AppKey      string `gorm:"app_key"`
}

func (stat MsgRealtimeStatDao) TableName() string {
	return "msgrealtimestats"
}

func (stat MsgRealtimeStatDao) IncrByStep(appkey string, statType int, channelType int, timeMark, step int64) error {
	sql := fmt.Sprintf("insert into %s (stat_type,channel_type,time_mark,count,app_key)values(?,?,?,?,?) ON DUPLICATE KEY UPDATE count=count+?", stat.TableName())
	return dbcommons.GetDb().Exec(sql, statType, channelType, timeMark, step, appkey, step).Error
}

func (stat MsgRealtimeStatDao) DeleteBeforeTimeMarkBatch(timeMark int64, limit int) (int64, error) {
	result := dbcommons.GetDb().Where("time_mark < ?", timeMark).Limit(limit).Delete(&MsgRealtimeStatDao{})
	return result.RowsAffected, result.Error
}
