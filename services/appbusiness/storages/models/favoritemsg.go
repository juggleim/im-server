package models

import "time"

type FavoriteMsg struct {
	ID          int64
	UserId      string
	SenderId    string
	ReceiverId  string
	ChannelType int32
	MsgId       string
	MsgTime     int64
	MsgType     string
	MsgContent  string
	CreatedTime time.Time
	AppKey      string
}

type IFavoriteMsgStorage interface {
	Create(item FavoriteMsg) error
	QueryFavoriteMsgs(appkey, userId string, startId, limit int64) ([]*FavoriteMsg, error)
}
