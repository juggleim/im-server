package dbs

import (
	"im-server/commons/dbcommons"
	"time"
)

var (
	InterceptorType_Custom int = 0
	InterceptorType_Baidu  int = 1
	InterceptorType_Yidun  int = 2
)

type InterceptorDao struct {
	ID              int64     `gorm:"primary_key"`
	Name            string    `gorm:"name"`
	Sort            int       `gorm:"sort"`
	RequestUrl      string    `gorm:"request_url"`
	RequestTemplate string    `gorm:"request_template"`
	SuccTemplate    string    `gorm:"succ_template"`
	IsAsync         int       `gorm:"is_async"`
	CreatedTime     time.Time `gorm:"created_time"`
	UpdatedTime     time.Time `gorm:"updated_time"`
	AppKey          string    `gorm:"app_key"`
	Conf            string    `gorm:"conf"`
	InterceptType   int       `gorm:"intercept_type"`
}

func (interceptor InterceptorDao) TableName() string {
	return "interceptors"
}

func (interceptor InterceptorDao) Create(item InterceptorDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (interceptor InterceptorDao) QryInterceptors(appkey string) ([]*InterceptorDao, error) {
	var items []*InterceptorDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Order("sort asc").Find(&items).Error
	return items, err
}
