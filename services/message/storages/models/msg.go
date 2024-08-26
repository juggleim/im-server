package models

import (
	"im-server/commons/pbdefines/pbobjs"
)

type Msg struct {
	UserId      string
	SendTime    int64
	MsgId       string
	ChannelType pbobjs.ChannelType
	MsgBody     []byte
	AppKey      string
	TargetId    string
	MsgType     string
	UniqTag     string
}

type IMsgStorage interface {
	SaveMsg(msg Msg) error
	UpsertMsg(item Msg) error
	QryMsgsBaseTime(appkey, userId string, start int64, count int) ([]*Msg, error)
	DelMsgsBaseTime(appkey string, start int64) error
}
