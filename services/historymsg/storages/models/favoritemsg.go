package models

import (
	"im-server/commons/pbdefines/pbobjs"
	"time"
)

type FavoriteMsg struct {
	ID          int64
	UserId      string
	SenderId    string
	ReceiverId  string
	ChannelType pbobjs.ChannelType
	MsgId       string
	MsgTime     int64
	MsgType     string
	MsgBody     []byte
	CreatedTime time.Time
	AppKey      string
}

type IFavoriteMsgStorage interface {
	Create(item FavoriteMsg) error
	QueryFavoriteMsgs(appkey, userId string, startId, limit int64) ([]*FavoriteMsg, error)
	BatchDelete(appkey, userId string, msgIds []string) error
}
