package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/services/historymsg/storages/models"
	"time"
)

type MsgExSetDao struct {
	ID          int64     `gorm:"primary_key"`
	AppKey      string    `gorm:"app_key"`
	MsgId       string    `gorm:"msg_id"`
	Key         string    `gorm:"key"`
	Item        string    `gorm:"item"`
	UserId      string    `gorm:"user_id"`
	CreatedTime time.Time `gorm:"created_time"`
}

func (exset MsgExSetDao) TableName() string {
	return "msgexsets"
}

func (exset MsgExSetDao) Create(item models.MsgExSet) error {
	add := MsgExSetDao{
		AppKey:      item.AppKey,
		MsgId:       item.MsgId,
		Key:         item.Key,
		Item:        item.Item,
		UserId:      item.UserId,
		CreatedTime: time.UnixMilli(item.CreatedTime),
	}
	return dbcommons.GetDb().Create(&add).Error
}

func (exset MsgExSetDao) Delete(appkey, msgId, key, item string) error {
	return dbcommons.GetDb().Where("app_key=? and msg_id=? and `key`=? and item=?", appkey, msgId, key, item).Delete(&MsgExSetDao{}).Error
}

func (exset MsgExSetDao) QryExtsByMsgIds(appkey string, msgIds []string) ([]*models.MsgExSet, error) {
	var items []*MsgExSetDao
	err := dbcommons.GetDb().Where("app_key=? and msg_id in (?)", appkey, msgIds).Order("created_time asc").Find(&items).Error
	retItems := []*models.MsgExSet{}
	for _, ext := range items {
		retItems = append(retItems, &models.MsgExSet{
			AppKey:      appkey,
			MsgId:       ext.MsgId,
			Key:         ext.Key,
			Item:        ext.Item,
			UserId:      ext.UserId,
			CreatedTime: ext.CreatedTime.UnixMilli(),
		})
	}
	return retItems, err
}
