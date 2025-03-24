package models

import (
	"im-server/commons/pbdefines/pbobjs"
	"time"
)

type User struct {
	ID           int64
	UserId       string
	Nickname     string
	UserPortrait string
	Pinyin       string
	UserType     pbobjs.UserType
	Phone        string
	Email        string
	LoginAccount string
	LoginPass    string
	Status       int
	UpdatedTime  time.Time
	AppKey       string
}

type IUserStorage interface {
	Create(item User) error
	Upsert(item User) error
	FindByPhone(appkey, phone string) (*User, error)
	FindByEmail(appkey, email string) (*User, error)
	FindByUserId(appkey, userId string) (*User, error)
	Update(appkey, userId, nickname, userPortrait string) error
	Count(appkey string) int
	CountByTime(appkey string, start, end int64) int64
}

type UserExt struct {
	ID          int64
	UserId      string
	ItemKey     string
	ItemValue   string
	ItemType    int
	UpdatedTime time.Time
	AppKey      string
}

type IUserExtStorage interface {
	Upsert(item UserExt) error
	BatchUpsert(items []UserExt) error
	BatchDelete(appkey, itemKey string, userIds []string) error
	QryExtFields(appkey, userId string) ([]*UserExt, error)
	QryExtFieldsByItemKeys(appkey, userId string, itemKeys []string) (map[string]*UserExt, error)
	QryExtsBaseItemKey(appkey, itemKey string, startId, limit int64) ([]*UserExt, error)
}
