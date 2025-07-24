package models

import (
	"im-server/commons/pbdefines/pbobjs"
	"time"
)

type TopMsg struct {
	ID          int64
	ConverId    string
	SubChannel  string
	ChannelType pbobjs.ChannelType
	MsgId       string
	UserId      string
	CreatedTime time.Time
	AppKey      string
}

type ITopMsgStorage interface {
	Upsert(item TopMsg) error
	FindTopMsg(appkey, converId, subChannel string, channelType pbobjs.ChannelType) (*TopMsg, error)
	DelTopMsg(appkey, converId, subChannel string, channelType pbobjs.ChannelType, msgId string) error
}
