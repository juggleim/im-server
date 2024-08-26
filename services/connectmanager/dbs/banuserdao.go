package dbs

import (
	"im-server/commons/dbcommons"
	"time"
)

type UserBanScope string

const (
	UserBanScopeDefault  UserBanScope = "default"
	UserBanScopePlatform UserBanScope = "platform"
	UserBanScopeDevice   UserBanScope = "device"
	UserBanScopeIp       UserBanScope = "ip"
)

type BanUserDao struct {
	ID          int64     `gorm:"primary_key"`
	UserId      string    `gorm:"user_id"`
	BanType     int       `gorm:"ban_type"`
	CreatedTime time.Time `gorm:"created_time"`
	EndTime     int64     `gorm:"end_time"`
	ScopeKey    string    `gorm:"scope_key"`
	ScopeValue  string    `gorm:"scope_value"`
	Ext         string    `gorm:"ext"`
	AppKey      string    `gorm:"app_key"`
}

func (user BanUserDao) TableName() string {
	return "banusers"
}

func (user BanUserDao) Upsert(item BanUserDao) error {
	err := dbcommons.GetDb().Exec("INSERT INTO banusers (user_id, ban_type, end_time, scope_key, scope_value, ext, app_key)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE end_time=?, ban_type=?, scope_value=?, ext=?",
		item.UserId, item.BanType, item.EndTime, item.ScopeKey, item.ScopeValue, item.Ext, item.AppKey, item.EndTime, item.BanType, item.ScopeValue, item.Ext).Error
	return err
}

func (user BanUserDao) FindById(appkey, userId string) ([]*BanUserDao, error) {
	var items []*BanUserDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (user BanUserDao) DelBanUser(appkey, userId, scopeKey string) error {
	if scopeKey == "" {
		return dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Delete(&BanUserDao{}).Error
	}
	return dbcommons.GetDb().Where("app_key=? and user_id=? and scope_key=?", appkey, userId, scopeKey).Delete(&BanUserDao{}).Error
}

func (user BanUserDao) QryBanUsers(appkey string, limit, startId int64) ([]*BanUserDao, error) {
	var items []*BanUserDao
	err := dbcommons.GetDb().Where("app_key=? and (ban_type=0 or (ban_type=1 and end_time>?)) and id>?", appkey, time.Now().UnixMilli(), startId).Order("id asc").Limit(limit).Find(&items).Error
	return items, err
}
