package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"time"
)

type AppExtDao struct {
	ID           int64     `gorm:"primary_key"`
	AppKey       string    `gorm:"app_key"`
	AppItemKey   string    `gorm:"app_item_key"`
	AppItemValue string    `gorm:"app_item_value"`
	UpdatedTime  time.Time `gorm:"updated_time"`
}

func (appExt AppExtDao) TableName() string {
	return "appexts"
}

func (appExt AppExtDao) FindListByAppkey(appkey string) []*AppExtDao {
	var list []*AppExtDao
	dbcommons.GetDb().Where("app_key=?", appkey).Find(&list)
	return list
}

func (appExt AppExtDao) Find(appkey string, itemKey string) (*AppExtDao, error) {
	var item AppExtDao
	err := dbcommons.GetDb().Where("app_key=? and app_item_key=?", appkey, itemKey).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (appExt AppExtDao) FindByItemKeys(appkey string, itemKeys []string) ([]*AppExtDao, error) {
	var list []*AppExtDao
	err := dbcommons.GetDb().Where("app_key=? and app_item_key in(?)", appkey, itemKeys).Find(&list).Error
	return list, err
}

func (appExt AppExtDao) CreateOrUpdate(appkey string, fieldKey, fieldValue string) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,app_item_key,app_item_value)VALUES(?,?,?) ON DUPLICATE KEY UPDATE app_item_value=?", appExt.TableName()), appkey, fieldKey, fieldValue, fieldValue).Error
}
