package dbs

import (
	"errors"
	"im-server/commons/dbcommons"
	"time"

	"gorm.io/gorm"
)

type BindDeviceDao struct {
	ID            int64     `gorm:"primary_key"`
	AppKey        string    `gorm:"app_key"`
	UserId        string    `gorm:"user_id"`
	Platform      string    `gorm:"platform"`
	DeviceId      string    `gorm:"device_id"`
	DeviceCompany string    `gorm:"device_company"`
	DeviceModel   string    `gorm:"device_model"`
	CreatedTime   time.Time `gorm:"created_time"`
}

func (dev BindDeviceDao) TableName() string {
	return "binddevices"
}

func (dev BindDeviceDao) Upsert(item BindDeviceDao) error {
	err := dbcommons.GetDb().Exec("INSERT INTO binddevices (app_key, user_id, platform, device_id, device_company, device_model)VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE platform=?, device_company=?, device_model=?",
		item.AppKey, item.UserId, item.Platform, item.DeviceId, item.DeviceCompany, item.DeviceModel, item.Platform, item.DeviceCompany, item.DeviceModel).Error
	return err
}

func (dev BindDeviceDao) FindByUserId(appkey, userId string) ([]*BindDeviceDao, error) {
	var items []*BindDeviceDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Limit(100).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (dev BindDeviceDao) FindByUserAndDevice(appkey, userId, deviceId string) (*BindDeviceDao, error) {
	item := &BindDeviceDao{}
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and device_id=?", appkey, userId, deviceId).First(item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (dev BindDeviceDao) DelByUserAndDevice(appkey, userId, deviceId string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and device_id=?", appkey, userId, deviceId).Delete(&BindDeviceDao{}).Error
}
