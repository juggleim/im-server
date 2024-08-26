package dbs

import "im-server/commons/dbcommons"

type SubRelationDao struct {
	ID         int64  `gorm:"primary_key"`
	UserId     string `gorm:"user_id"`
	Subscriber string `gorm:"subscriber"`
	// CreatedTime time.Time `gorm:"created_time"`
	AppKey string `gorm:"app_key"`
}

func (sub SubRelationDao) TableName() string {
	return "subrelations"
}

func (sub SubRelationDao) Create(item SubRelationDao) error {
	return dbcommons.GetDb().Create(&item).Error
}

func (sub SubRelationDao) QrySubscribers(appkey, userId string, startId int64, count int32) ([]*SubRelationDao, error) {
	var items []*SubRelationDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and id>?", appkey, userId, startId).Order("id asc").Limit(count).Find(&items).Error
	return items, err
}

func (sub SubRelationDao) Delete(appkey, userId, subscriber string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and subscriber=?", appkey, userId, subscriber).Delete(&SubRelationDao{}).Error
}
