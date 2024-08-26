package dbs

import (
	"im-server/commons/dbcommons"

	"github.com/jinzhu/gorm"
)

type GlobalConfKey string

const (
	GlobalConfKey_NaviAddress    GlobalConfKey = "nav_address"
	GlobalConfKey_ConnectAddress GlobalConfKey = "connect_address"
	GlobalConfKey_ApiAddress     GlobalConfKey = "api_address"
)

type GlobalConfDao struct {
	ID        int64  `gorm:"primary_key"`
	ConfKey   string `gorm:"conf_key"`
	ConfValue string `gorm:"conf_value"`
}

func (conf GlobalConfDao) TableName() string {
	return "globalconfs"
}

func (conf GlobalConfDao) Create(item GlobalConfDao) error {
	return dbcommons.GetDb().Create(&item).Error
}

func (conf GlobalConfDao) FindByKey(key string) (*GlobalConfDao, error) {
	var item GlobalConfDao
	err := dbcommons.GetDb().Where("conf_key=?", key).Take(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	return &item, nil
}
