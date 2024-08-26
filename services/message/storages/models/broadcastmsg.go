package models

import "im-server/commons/pbdefines/pbobjs"

type BrdInboxMsgMsg struct {
	SenderId    string
	SendTime    int64
	MsgId       string
	ChannelType pbobjs.ChannelType
	MsgBody     []byte
	AppKey      string
}
type IBroadcastMsgStorage interface {
	SaveMsg(msg BrdInboxMsgMsg) error
	QryMsgsBaseTime(appkey string, start int64, count int) ([]*BrdInboxMsgMsg, error)
	QryLatestMsg(appkey string, count int) ([]*BrdInboxMsgMsg, error)
	DelMsgsBaseTime(appkey string, start int64) error
}
