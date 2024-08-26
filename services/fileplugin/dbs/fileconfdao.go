package dbs

import (
	"im-server/commons/dbcommons"

	"github.com/jinzhu/gorm"
)

type FileConfDao struct {
	ID      int64  `gorm:"primary_key" json:"id,omitempty"`
	AppKey  string `gorm:"column:app_key" json:"app_key,omitempty"`
	Channel string `gorm:"column:channel" json:"channel,omitempty"`
	Conf    string `gorm:"column:conf" json:"conf,omitempty"`
	Enable  int    `gorm:"column:enable" json:"enable,omitempty"`
}

func (file FileConfDao) TableName() string {
	return "fileconfs"
}

func (file FileConfDao) Upsert(item FileConfDao) error {
	return dbcommons.GetDb().Exec("INSERT INTO fileconfs (app_key, channel, conf) VALUES (?,?,?) ON DUPLICATE KEY UPDATE conf = ?",
		item.AppKey, item.Channel, item.Conf, item.Conf).Error
}

func (file FileConfDao) FindEnableFileConf(appkey string) (*FileConfDao, error) {
	var item FileConfDao
	err := dbcommons.GetDb().Where("app_key=? and enable=1", appkey).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (file FileConfDao) FindFileConfs(appkey string) ([]*FileConfDao, error) {
	var items []*FileConfDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Find(&items).Error
	return items, err
}

func (file FileConfDao) FindFileConf(appkey, channel string) (*FileConfDao, error) {
	var item FileConfDao
	err := dbcommons.GetDb().Where("app_key=? and channel=?", appkey, channel).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (file FileConfDao) UpdateEnable(appkey, channel string) error {
	err := dbcommons.GetDb().Model(&file).
		Where("app_key = ?", appkey).
		UpdateColumn("enable", gorm.Expr("CASE WHEN channel = ? THEN 1 ELSE 0 END", channel)).Error
	return err
}
