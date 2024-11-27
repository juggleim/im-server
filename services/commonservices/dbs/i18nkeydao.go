package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
)

type I18nKeyDao struct {
	ID     int64  `gorm:"primary_key"`
	Lang   string `gorm:"lang"`
	Key    string `gorm:"key"`
	Value  string `gorm:"value"`
	AppKey string `gorm:"app_key"`
}

func (i18n I18nKeyDao) TableName() string {
	return "i18nkeys"
}

func (i18n I18nKeyDao) Upsert(item I18nKeyDao) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,lang,key,value)VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE value=VALUES(value)", i18n.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.Lang, item.Key, item.Value).Error
}

func (i18n I18nKeyDao) Delete(appkey, lang, key string) error {
	return dbcommons.GetDb().Where("app_key=? and lang=? and key=?", appkey, lang, key).Delete(&i18n).Error
}

func (i18n I18nKeyDao) Query(appkey string, startId, limit int64) ([]*I18nKeyDao, error) {
	var list []*I18nKeyDao
	err := dbcommons.GetDb().Where("app_key=? and id>?", appkey, startId).Order("id asc").Limit(limit).Find(&list).Error
	return list, err
}
