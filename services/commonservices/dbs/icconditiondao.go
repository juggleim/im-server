package dbs

import "im-server/commons/dbcommons"

type IcConditionDao struct {
	ID            int64  `gorm:"primary_key"`
	InterceptorId int64  `gorm:"interceptor_id"`
	ChannelType   string `gorm:"channel_type"`
	MsgType       string `gorm:"msg_type"`
	SenderId      string `gorm:"sender_id"`
	ReceiverId    string `gorm:"receiver_id"`
	AppKey        string `gorm:"app_key"`
}

func (condition IcConditionDao) TableName() string {
	return "ic_conditions"
}

func (condition IcConditionDao) Create(item IcConditionDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (condition IcConditionDao) QryConditions(appkey string, interId int64) ([]*IcConditionDao, error) {
	var items []*IcConditionDao
	err := dbcommons.GetDb().Where("app_key=? and interceptor_id=?", appkey, interId).Find(&items).Error
	return items, err
}
