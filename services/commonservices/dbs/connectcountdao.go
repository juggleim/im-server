package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
)

type ConnectCountDao struct {
	ID          int64  `gorm:"primary_key"`
	ConnectType int    `gorm:"connect_type"`
	TimeMark    int64  `gorm:"time_mark"`
	Count       int64  `gorm:"count"`
	AppKey      string `gorm:"app_key"`
}

func (count ConnectCountDao) TableName() string {
	return "connectcounts"
}

func (count ConnectCountDao) IncrByStep(appkey string, connectType int, timeMark, step int64) error {
	sql := fmt.Sprintf("insert into %s (connect_type,time_mark,count,app_key)values(?,?,?,?) ON DUPLICATE KEY UPDATE count=count+?", count.TableName())
	return dbcommons.GetDb().Exec(sql, connectType, timeMark, step, appkey, step).Error
}

func (count ConnectCountDao) QryStats(appkey string, connectType int, start, end int64) []*ConnectCountDao {
	var items []*ConnectCountDao
	err := dbcommons.GetDb().Where("app_key=? and connect_type=? and time_mark>=? and time_mark<=?", appkey, connectType, start, end).Limit(1000).Find(&items).Error
	if err == nil {
		return items
	}
	return []*ConnectCountDao{}
}

func (count ConnectCountDao) MaxByTime(appkey string, connectType int, start, end int64) *ConnectCountDao {
	var item ConnectCountDao
	err := dbcommons.GetDb().Where("app_key=? and connect_type=? and time_mark>=? and time_mark<=?", appkey, connectType, start, end).Order("count desc").Limit(1).Find(&item)
	if err != nil {
		return nil
	}
	return &item
}

func (count ConnectCountDao) ScanByTimeRangeAfterCursor(start, end, lastTimeMark, lastID int64, limit int) ([]ConnectCountDao, error) {
	items := []ConnectCountDao{}
	sql := fmt.Sprintf("select id,connect_type,time_mark,count,app_key from %s where time_mark>=? and time_mark<? and (time_mark>? or (time_mark=? and id>?)) order by time_mark asc,id asc limit ?", count.TableName())
	err := dbcommons.GetDb().Raw(sql, start, end, lastTimeMark, lastTimeMark, lastID, limit).Scan(&items).Error
	return items, err
}

func (count ConnectCountDao) DeleteByTimeRangeExceptKeepersBatch(start, end int64, keeperIDs []int64, limit int) (int64, error) {
	if len(keeperIDs) == 0 {
		return 0, nil
	}
	sql := fmt.Sprintf("delete from %s where time_mark>=? and time_mark<? and id not in (?) limit ?", count.TableName())
	tx := dbcommons.GetDb().Exec(sql, start, end, keeperIDs, limit)
	return tx.RowsAffected, tx.Error
}
