package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/usermanager/storages/models"
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

func (ext UserExtDao) Upsert(item models.UserExt) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,user_id,item_key,item_value,item_type)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=?", ext.TableName()), item.AppKey, item.UserId, item.ItemKey, item.ItemValue, item.ItemType, item.ItemValue).Error
}

func (ext UserExtDao) BatchUpsert(items []models.UserExt) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT INTO %s (app_key,user_id,item_key,item_value,item_type)VALUES", ext.TableName())
	params := []interface{}{}
	buffer.WriteString(sql)
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=VALUES(item_value);")
		} else {
			buffer.WriteString("(?,?,?,?,?),")
		}
		params = append(params, item.AppKey, item.UserId, item.ItemKey, item.ItemValue, item.ItemType)
	}
	return dbcommons.GetDb().Exec(buffer.String(), params...).Error
}

func (ext UserExtDao) BatchDelete(appkey, itemKey string, userIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and item_key=? and user_id in (?)", appkey, itemKey, userIds).Delete(&UserExtDao{}).Error
}

func (ext UserExtDao) QryExtFields(appkey, userId string) ([]*models.UserExt, error) {
	var items []*UserExtDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Find(&items).Error
	ret := []*models.UserExt{}
	for _, item := range items {
		ret = append(ret, &models.UserExt{
			ID:          item.ID,
			UserId:      item.UserId,
			ItemKey:     item.ItemKey,
			ItemValue:   item.ItemValue,
			ItemType:    item.ItemType,
			UpdatedTime: item.UpdatedTime,
			AppKey:      item.AppKey,
		})
	}
	return ret, err
}

func (ext UserExtDao) QryExtFieldsByItemKeys(appkey, userId string, itemKeys []string) (map[string]*models.UserExt, error) {
	var items []*UserExtDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and item_key in (?)", appkey, userId, itemKeys).Find(&items).Error
	ret := map[string]*models.UserExt{}
	for _, item := range items {
		ret[item.ItemKey] = &models.UserExt{
			ID:          item.ID,
			UserId:      item.UserId,
			ItemKey:     item.ItemKey,
			ItemValue:   item.ItemValue,
			ItemType:    item.ItemType,
			UpdatedTime: item.UpdatedTime,
			AppKey:      item.AppKey,
		}
	}
	return ret, err
}

func (ext UserExtDao) QryExtsBaseItemKey(appkey, itemKey string, startId, limit int64) ([]*models.UserExt, error) {
	var items []*UserExtDao
	err := dbcommons.GetDb().Where("app_key=? and item_key=? and id>?", appkey, itemKey, startId).Order("id asc").Limit(limit).Find(&items).Error
	ret := []*models.UserExt{}
	for _, item := range items {
		ret = append(ret, &models.UserExt{
			ID:          item.ID,
			UserId:      item.UserId,
			ItemKey:     item.ItemKey,
			ItemValue:   item.ItemValue,
			ItemType:    item.ItemType,
			UpdatedTime: item.UpdatedTime,
			AppKey:      item.AppKey,
		})
	}
	return ret, err
}
