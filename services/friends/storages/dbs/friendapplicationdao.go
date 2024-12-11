package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/friends/storages/models"
)

type FriendApplicationDao struct {
	ID        int64  `gorm:"primary_key"`
	UserId    string `gorm:"user_id"`
	SponsorId string `gorm:"sponsor_id"`
	ApplyTime int64  `gorm:"apply_time"`
	Status    int    `gorm:"status"`
	AppKey    string `gorm:"app_key"`
}

func (apply FriendApplicationDao) TableName() string {
	return "friendapplications"
}

func (apply FriendApplicationDao) Upsert(item models.FriendApplication) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,user_id,sponsor_id,apply_time,status)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE apply_time=VALUES(apply_time),status=VALUES(status)", apply.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.UserId, item.SponsorId, item.ApplyTime, item.Status).Error
}

func (apply FriendApplicationDao) QueryApplications(appkey, userId string, startTime, count int64) ([]*models.FriendApplication, error) {
	// var items []*FriendApplicationDao
	// err := dbcommons.GetDb().Where("app_key=? and user_id=? and apply_time>?")
	return nil, nil
}

func (apply FriendApplicationDao) QueryApplicationsBaseSponsor(appkey, sponsorId string, startTime, count int64) ([]*models.FriendApplication, error) {
	return nil, nil
}
