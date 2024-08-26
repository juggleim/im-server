package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/services/historymsg/storages/models"
)

type MsgExtDao struct {
	ID     int64  `gorm:"primary_key"`
	AppKey string `gorm:"app_key"`
	MsgId  string `gorm:"msg_id"`
	Key    string `gorm:"key"`
	Value  string `gorm:"value"`
}

func (ext MsgExtDao) TableName() string {
	return "msgexts"
}

func (ext MsgExtDao) Upsert(item models.MsgExt) error {
	return dbcommons.GetDb().Exec("INSERT INTO msgexts (app_key,msg_id,key,value)VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE value=?",
		item.AppKey, item.MsgId, item.Key, item.Value, item.Value).Error
}
