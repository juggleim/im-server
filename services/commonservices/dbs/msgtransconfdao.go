package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
)

type MsgTransConfDao struct {
	ID       int64  `gorm:"primary_key"`
	MsgType  string `gorm:"msg_type"`
	JsonPath string `gorm:"json_path"`
	AppKey   string `gorm:"app_key"`
}

func (conf MsgTransConfDao) TableName() string {
	return "msgtransconfs"
}

func (conf MsgTransConfDao) Upsert(item MsgTransConfDao) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,msg_type,json_path)VALUES(?,?,?) ON DUPLICATE KEY UPDATE json_path=VALUES(json_path)", conf.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.MsgType, item.JsonPath).Error
}

func (conf MsgTransConfDao) QueryConfs(appkey string, startId, limit int64) ([]*MsgTransConfDao, error) {
	var list []*MsgTransConfDao
	err := dbcommons.GetDb().Where("app_key=? and id>?", appkey, startId).Order("id asc").Limit(limit).Find(&list).Error
	return list, err
}
