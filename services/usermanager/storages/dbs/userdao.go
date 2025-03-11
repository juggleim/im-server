package dbs

import (
	"fmt"
	"strings"
	"time"

	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/usermanager/storages/models"
)

type UserDao struct {
	ID           int64     `gorm:"primary_key"`
	UserType     int       `gorm:"user_type"`
	UserId       string    `gorm:"user_id"`
	Nickname     string    `gorm:"nickname"`
	UserPortrait string    `gorm:"user_portrait"`
	Pinyin       string    `gorm:"pinyin"`
	Phone        string    `gorm:"phone"`
	Email        string    `gorm:"email"`
	LoginAccount string    `gorm:"login_account"`
	LoginPass    string    `gorm:"login_pass"`
	CreatedTime  time.Time `gorm:"created_time"`
	UpdatedTime  time.Time `gorm:"updated_time"`
	AppKey       string    `gorm:"app_key"`
}

func (user UserDao) TableName() string {
	return "users"
}

func (user UserDao) FindByUserId(appkey, userId string) (*models.User, error) {
	var item UserDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:           item.ID,
		UserId:       item.UserId,
		Nickname:     item.Nickname,
		UserPortrait: item.UserPortrait,
		Pinyin:       item.Pinyin,
		UserType:     pbobjs.UserType(item.UserType),
		Phone:        item.Phone,
		Email:        item.Email,
		LoginAccount: item.LoginAccount,
		UpdatedTime:  item.UpdatedTime,
		AppKey:       item.AppKey,
	}, nil
}

func (user UserDao) FindByPhone(appkey, phone string) (*models.User, error) {
	var item UserDao
	err := dbcommons.GetDb().Where("app_key=? and phone=?", appkey, phone).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:           item.ID,
		UserId:       item.UserId,
		Nickname:     item.Nickname,
		UserPortrait: item.UserPortrait,
		Pinyin:       item.Pinyin,
		UserType:     pbobjs.UserType(item.UserType),
		Phone:        item.Phone,
		Email:        item.Email,
		LoginAccount: item.LoginAccount,
		UpdatedTime:  item.UpdatedTime,
		AppKey:       item.AppKey,
	}, nil
}

func (user UserDao) FindByEmail(appkey, email string) (*models.User, error) {
	var item UserDao
	err := dbcommons.GetDb().Where("app_key=? and email=?", appkey, email).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:           item.ID,
		UserId:       item.UserId,
		Nickname:     item.Nickname,
		UserPortrait: item.UserPortrait,
		Pinyin:       item.Pinyin,
		UserType:     pbobjs.UserType(item.UserType),
		Phone:        item.Phone,
		Email:        item.Email,
		LoginAccount: item.LoginAccount,
		UpdatedTime:  item.UpdatedTime,
		AppKey:       item.AppKey,
	}, nil
}

func (user UserDao) Create(item models.User) error {
	var sqlBuilder strings.Builder
	params := []interface{}{}
	sqlBuilder.WriteString("INSERT INTO ")
	sqlBuilder.WriteString(user.TableName())
	sqlBuilder.WriteString(" (app_key,user_id,user_type,nickname,user_portrait,pinyin")
	params = append(params, item.AppKey)
	params = append(params, item.UserId)
	params = append(params, item.UserType)
	params = append(params, item.Nickname)
	params = append(params, item.UserPortrait)
	params = append(params, item.Pinyin)
	if item.Phone != "" {
		sqlBuilder.WriteString(",phone")
		params = append(params, item.Phone)
	}
	if item.Email != "" {
		sqlBuilder.WriteString(",email")
		params = append(params, item.Email)
	}
	if item.LoginAccount != "" {
		sqlBuilder.WriteString(",login_account,login_pass")
		params = append(params, item.LoginAccount)
		params = append(params, item.LoginPass)
	}
	sqlBuilder.WriteString(")VALUES(")
	marks := []string{}
	for range params {
		marks = append(marks, "?")
	}
	sqlBuilder.WriteString(strings.Join(marks, ","))
	sqlBuilder.WriteString(")")
	err := dbcommons.GetDb().Exec(sqlBuilder.String(), params...).Error
	return err
}

func (user UserDao) Upsert(item models.User) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (user_id,user_type,nickname,user_portrait,app_key)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE nickname=?,user_portrait=?", user.TableName()), item.UserId, item.UserType, item.Nickname, item.UserPortrait, item.AppKey, item.Nickname, item.UserPortrait).Error
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
