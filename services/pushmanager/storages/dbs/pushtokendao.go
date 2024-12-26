package dbs

import (
	"im-server/commons/dbcommons"
)

type PushTokenDao struct {
	UserId      string `gorm:"user_id"`
	DeviceId    string `gorm:"device_id"`
	Platform    string `gorm:"platform"`
	PushChannel string `gorm:"push_channel"`
	Package     string `gorm:"package"`
	PushToken   string `gorm:"push_token"`
	VoipToken   string `gorm:"voip_token"`
	//CreatedTime time.Time `gorm:"created_time"`
	//UpdatedTme  time.Time `gorm:"updated_time"`
	AppKey string `gorm:"app_key"`
}

func (token PushTokenDao) TableName() string {
	return "pushtokens"
}

func (token PushTokenDao) Upsert(item PushTokenDao) error {
	err := dbcommons.GetDb().Exec("INSERT INTO pushtokens (app_key,user_id,device_id,platform,push_channel,package,push_token,voip_token)VALUES(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE device_id=?,platform=?,push_channel=?,package=?,push_token=?,voip_token=?",
		item.AppKey, item.UserId, item.DeviceId, item.Platform, item.PushChannel, item.Package, item.PushToken, item.VoipToken, item.DeviceId, item.Platform, item.PushChannel, item.Package, item.PushToken, item.VoipToken).Error
	return err
}

func (token PushTokenDao) FindByUserId(appkey, userId string) (*PushTokenDao, error) {
	var item PushTokenDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (token PushTokenDao) DeleteByDeviceId(appkey, deviceId, exceptUserId string) error {
	return dbcommons.GetDb().Where("app_key=? and device_id=? and user_id!=?", appkey, deviceId, exceptUserId).Delete(&PushTokenDao{}).Error
}

func (token PushTokenDao) GetUserWithToken(appkey string, pushToken string) (*PushTokenDao, error) {
	var item PushTokenDao
	err := dbcommons.GetDb().Where("app_key=? and push_token=?", appkey, pushToken).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (token PushTokenDao) QueryByDeviceId(appkey, deviceId string) ([]*PushTokenDao, error) {
	var items []*PushTokenDao
	err := dbcommons.GetDb().Where("app_key=? and device_id=?", appkey, deviceId).Find(&items).Error
	return items, err
}

func (token PushTokenDao) DeleteByUserId(appkey, userId string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Delete(&PushTokenDao{}).Error
}
