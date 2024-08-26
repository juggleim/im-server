package models

import (
	"im-server/commons/pbdefines/pbobjs"
	"time"
)

type Chatroom struct {
	ChatId      string
	ChatName    string
	IsMute      int
	CreatedTime time.Time
	UpdatedTime time.Time
	AppKey      string
}

type IChatroomStorage interface {
	Create(item Chatroom) error
	FindById(appkey, chatId string) (*Chatroom, error)
	Delete(appkey, chatId string) error
	UpdateMute(appkey, chatId string, isMute int) error
}

type ChatroomMember struct {
	ID          int64
	ChatId      string
	MemberId    string
	CreatedTime time.Time
	AppKey      string
}

type IChatroomMemberStorage interface {
	Create(item ChatroomMember) error
	FindById(appkey, chatId, memberId string) (*ChatroomMember, error)
	DeleteMember(appkey, chatId, memberId string) error
	ClearMembers(appkey, chatId string) error
	QryMembers(appkey, chatId string, isPositive bool, startId, limit int64) ([]*ChatroomMember, error)
}

type ChatroomExt struct {
	ChatId    string
	ItemKey   string
	ItemValue string
	ItemType  int
	ItemTime  int64
	AppKey    string
	MemberId  string
	IsDelete  int
}

type IChatroomExtStorage interface {
	Upsert(item ChatroomExt) error
	QryExts(appkey, chatId string) ([]*ChatroomExt, error)
	ClearExts(appkey, chatId string) error
	DeleteExt(appkey, chatId, key string) error
}

type ChatroomBanUser struct {
	ID          int64
	ChatId      string
	BanType     pbobjs.ChrmBanType
	MemberId    string
	AppKey      string
	CreatedTime int64
}

type IChatroomBanUserStorage interface {
	Create(item ChatroomBanUser) error
	DelBanUser(appkey, chatId, memberId string, banType pbobjs.ChrmBanType) error
	QryBanUsers(appkey, chatId string, banType pbobjs.ChrmBanType, startId, limit int64) ([]*ChatroomBanUser, error)
	ClearBanUsers(appkey, chatId string) error
}
