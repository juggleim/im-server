package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"time"
)

type GroupExtDao struct {
	ID          int64     `gorm:"primary_key"`
	GroupId     string    `gorm:"group_id"`
	ItemKey     string    `gorm:"item_key"`
	ItemValue   string    `gorm:"item_value"`
	ItemType    int       `gorm:"item_type"`
	UpdatedTime time.Time `gorm:"updated_time"`
	AppKey      string    `gorm:"app_key"`
}

func (ext GroupExtDao) TableName() string {
	return "groupinfoexts"
}

func (ext GroupExtDao) QryExtFields(appkey, groupId string) ([]*GroupExtDao, error) {
	var items []*GroupExtDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Find(&items).Error
	return items, err
}

func (ext GroupExtDao) Upsert(item GroupExtDao) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,group_id,item_key,item_value,item_type)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=?", ext.TableName()), item.AppKey, item.GroupId, item.ItemKey, item.ItemValue, item.ItemType, item.ItemValue).Error
}
