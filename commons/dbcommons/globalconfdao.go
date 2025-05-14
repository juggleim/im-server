package dbcommons

import (
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
	return GetDb().Create(&item).Error
}

func (conf GlobalConfDao) FindByKey(key string) (*GlobalConfDao, error) {
	var item GlobalConfDao
	err := GetDb().Where("conf_key=?", key).Take(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	return &item, nil
}

func (conf GlobalConfDao) UpdateValue(key, val string) error {
	return GetDb().Model(&GlobalConfDao{}).Where("conf_key=?", key).Update("conf_value", val).Error
}
