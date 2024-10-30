package dbs

import (
	"bytes"
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

func (ext UserExtDao) BatchUpsert(items []UserExtDao) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT INTO %s (app_key,user_id,item_key,item_value,item_type)VALUES", ext.TableName())
	buffer.WriteString(sql)
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s',%d);", item.AppKey, item.UserId, item.ItemKey, item.ItemValue, item.ItemType))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s',%d),", item.AppKey, item.UserId, item.ItemKey, item.ItemValue, item.ItemType))
		}
	}
	return dbcommons.GetDb().Exec(buffer.String()).Error
}

func (ext UserExtDao) BatchDelete(appkey, itemKey string, userIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and item_key=? and user_id in (?)", appkey, itemKey, userIds).Delete(&UserExtDao{}).Error
}

func (ext UserExtDao) QryExtFields(appkey, userId string) ([]*UserExtDao, error) {
	var items []*UserExtDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Find(&items).Error
	return items, err
}

func (ext UserExtDao) QryExtsBaseItemKey(appkey, itemKey string, startId, limit int64) ([]*UserExtDao, error) {
	var items []*UserExtDao
	err := dbcommons.GetDb().Where("app_key=? and item_key=? and id>?", appkey, itemKey, startId).Order("id asc").Limit(limit).Find(&items).Error
	return items, err
}
