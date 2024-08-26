package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"time"
)

type UserExtDao struct {
	ID          int64     `gorm:"primary_key"`
	UserId      string    `gorm:"user_id"`
	ItemKey     string    `gorm:"item_key"`
	ItemValue   string    `gorm:"item_value"`
	ItemType    int       `gorm:"item_type"`
	UpdatedTime time.Time `gorm:"updated_time"`
	AppKey      string    `gorm:"app_key"`
}

func (ext UserExtDao) TableName() string {
	return "userexts"
}

func (ext UserExtDao) Upsert(item UserExtDao) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,user_id,item_key,item_value,item_type)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=?", ext.TableName()), item.AppKey, item.UserId, item.ItemKey, item.ItemValue, item.ItemType, item.ItemValue).Error
}

func (ext UserExtDao) QryExtFields(appkey, userId string) ([]*UserExtDao, error) {
	var items []*UserExtDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Find(&items).Error
	return items, err
}
