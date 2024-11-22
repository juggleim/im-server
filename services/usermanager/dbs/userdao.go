package dbs

import (
	"fmt"
	"time"

	"im-server/commons/dbcommons"

	"github.com/jinzhu/gorm"
)

type UserDao struct {
	ID           int64     `gorm:"primary_key"`
	UserType     int       `gorm:"user_type"`
	UserId       string    `gorm:"user_id"`
	Nickname     string    `gorm:"nickname"`
	UserPortrait string    `gorm:"user_portrait"`
	CreatedTime  time.Time `gorm:"created_time"`
	UpdatedTime  time.Time `gorm:"updated_time"`
	AppKey       string    `gorm:"app_key"`
}

func (user UserDao) TableName() string {
	return "users"
}
func (user UserDao) FindByUserId(appkey, userId string) *UserDao {
	var item UserDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Take(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return &item
}

func (user UserDao) Create(item UserDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (user UserDao) Upsert(item UserDao) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (user_id,user_type,nickname,user_portrait,created_time,updated_time,app_key)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE nickname=?,user_portrait=?,updated_time=?", user.TableName()), item.UserId, item.UserType, item.Nickname, item.UserPortrait, item.CreatedTime, item.UpdatedTime, item.AppKey, item.Nickname, item.UserPortrait, item.UpdatedTime).Error
}

func (user UserDao) Update(appkey, userId, nickname, userPortrait string) error {
	upd := map[string]interface{}{}
	if nickname != "" {
		upd["nickname"] = nickname
	}
	if userPortrait != "" {
		upd["user_portrait"] = userPortrait
	}
	if len(upd) > 0 {
		upd["updated_time"] = time.Now()
	} else {
		return fmt.Errorf("do nothing")
	}
	err := dbcommons.GetDb().Model(&UserDao{}).Where("app_key=? and user_id=?", appkey, userId).Update(upd).Error
	return err
}

func (user UserDao) Count(appkey string) int {
	var count int
	err := dbcommons.GetDb().Model(&UserDao{}).Where("app_key=?", appkey).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}

func (user UserDao) CountByTime(appkey string, start, end int64) int64 {
	var count int64
	err := dbcommons.GetDb().Model(&UserDao{}).Where("app_key=? and created_time>=? and created_time<=?", appkey, time.UnixMilli(start), time.UnixMilli(end)).Count(&count).Error
	if err != nil {
		return count
	}
	return count
}
