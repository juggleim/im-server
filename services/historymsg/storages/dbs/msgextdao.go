package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/services/historymsg/storages/models"
	"time"
)

type MsgExtDao struct {
	ID          int64     `gorm:"primary_key"`
	AppKey      string    `gorm:"app_key"`
	MsgId       string    `gorm:"msg_id"`
	Key         string    `gorm:"key"`
	Value       string    `gorm:"value"`
	CreatedTime time.Time `gorm:"created_time"`
}

func (ext MsgExtDao) TableName() string {
	return "msgexts"
}

func (ext MsgExtDao) Upsert(item models.MsgExt) error {
	return dbcommons.GetDb().Exec("INSERT INTO msgexts (app_key,msg_id,key,value,created_time)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE value=?,created_time=?",
		item.AppKey, item.MsgId, item.Key, item.Value, time.UnixMilli(item.CreatedTime), item.Value, time.UnixMilli(item.CreatedTime)).Error
}

func (exset MsgExtDao) Delete(appkey, msgId, key string) error {
	return dbcommons.GetDb().Where("app_key=? and msg_id=? and key=?", appkey, msgId, key).Delete(&MsgExtDao{}).Error
}

func (ext MsgExtDao) QryExtsByMsgIds(appkey string, msgIds []string) ([]*models.MsgExt, error) {
	var items []*MsgExtDao
	err := dbcommons.GetDb().Where("app_key=? and msg_id in (?)", appkey, msgIds).Order("created_time asc").Find(&items).Error
	retItems := []*models.MsgExt{}
	for _, ext := range items {
		retItems = append(retItems, &models.MsgExt{
			AppKey:      appkey,
			MsgId:       ext.MsgId,
			Key:         ext.Key,
			Value:       ext.Value,
			CreatedTime: ext.CreatedTime.UnixMilli(),
		})
	}
	return retItems, err
}
