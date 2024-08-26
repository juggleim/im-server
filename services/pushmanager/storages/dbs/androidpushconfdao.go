package dbs

import "im-server/commons/dbcommons"

type AndroidPushConfDao struct {
	AppKey      string `gorm:"app_key"`
	PushChannel string `gorm:"push_channel"`
	Package     string `gorm:"package"`
	PushConf    string `gorm:"push_conf"`
}

func (conf AndroidPushConfDao) TableName() string {
	return "androidpushconfs"
}

func (conf AndroidPushConfDao) Upsert(item AndroidPushConfDao) error {
	err := dbcommons.GetDb().Exec("INSERT INTO androidpushconfs (app_key,push_channel,package,push_conf)VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE package=?,push_conf=?",
		item.AppKey, item.PushChannel, item.Package, item.PushConf, item.Package, item.PushConf).Error
	return err
}

func (conf AndroidPushConfDao) Create(item AndroidPushConfDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (conf AndroidPushConfDao) Find(appkey, pushChannel string) (*AndroidPushConfDao, error) {
	var item AndroidPushConfDao
	err := dbcommons.GetDb().Where("app_key=? and push_channel=?", appkey, pushChannel).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (conf AndroidPushConfDao) FindByChannels(appkey, packageName string) ([]*AndroidPushConfDao, error) {
	var list []*AndroidPushConfDao
	err := dbcommons.GetDb().Where("app_key=? and package=?", appkey, packageName).Find(&list).Error
	return list, err
}
