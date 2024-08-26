package dbs

import (
	"im-server/commons/dbcommons"
	"time"
)

type BlockDao struct {
	ID          int64     `gorm:"primary_key"`
	UserId      string    `gorm:"user_id"`
	BlockUserId string    `gorm:"block_user_id"`
	CreatedTime time.Time `gorm:"created_time"`
	AppKey      string    `gorm:"app_key"`
}

func (block *BlockDao) TableName() string {
	return "blocks"
}

func (block BlockDao) Create(item BlockDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (block BlockDao) DelBlockUser(appkey, userId, blockUserId string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and block_user_id=?", appkey, userId, blockUserId).Delete(&BlockDao{}).Error
}

func (block BlockDao) BatchDelBlockUsers(appkey, userId string, blockUserIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and block_user_id in (?)", appkey, userId, blockUserIds).Delete(&BlockDao{}).Error
}

func (block BlockDao) QryBlockUsers(appkey, userId string, limit, startId int64) ([]*BlockDao, error) {
	var items []*BlockDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and id>?", appkey, userId, startId).Order("id asc").Limit(limit).Find(&items).Error
	return items, err
}
