package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	commonDbs "im-server/services/commonservices/dbs"
)

type AccountAppRelDao struct {
	ID      int64  `gorm:"primary_key"`
	AppKey  string `gorm:"app_key"`
	Account string `gorm:"account"`
}

func (rel AccountAppRelDao) TableName() string {
	return "accountapprels"
}

func (rel AccountAppRelDao) Create(item AccountAppRelDao) error {
	return dbcommons.GetDb().Create(&item).Error
}

func (rel AccountAppRelDao) CheckExist(appkey string, account string) bool {
	var item AccountAppRelDao
	err := dbcommons.GetDb().Where("account=? and app_key=?", account, appkey).Take(&item).Error
	return err == nil
}

func (rel AccountAppRelDao) BatchDelete(account string, appkeys []string) error {
	return dbcommons.GetDb().Where("account=? and app_key in (?)", account, appkeys).Delete(&rel).Error
}

func (rel AccountAppRelDao) FindByAppkey(account string, appkey string) *commonDbs.AppInfoDao {
	var appItem commonDbs.AppInfoDao
	sql := fmt.Sprintf("select app.* from %s as rel left join %s as app on rel.app_key=app.app_key where rel.account=? and rel.app_key=?", rel.TableName(), commonDbs.AppInfoDao{}.TableName())
	err := dbcommons.GetDb().Raw(sql, account, appkey).Take(&appItem).Error
	if err != nil {
		return nil
	}
	return &appItem
}

func (rel AccountAppRelDao) QryApps(account string, limit int64, offset int64) ([]*commonDbs.AppInfoDao, error) {
	var list []*commonDbs.AppInfoDao
	sql := fmt.Sprintf("select app.* from %s as rel left join %s as app on rel.app_key=app.app_key where rel.account=? and app.id<?", rel.TableName(), commonDbs.AppInfoDao{}.TableName())
	err := dbcommons.GetDb().Raw(sql, account, offset).Order("app.id desc").Limit(limit).Find(&list).Error
	return list, err
}
