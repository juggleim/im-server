package models

import "im-server/commons/pbdefines/pbobjs"

type MentionMsg struct {
	UserId      string
	TargetId    string
	ChannelType pbobjs.ChannelType
	SenderId    string
	MentionType pbobjs.MentionType
	MsgId       string
	MsgTime     int64
	MsgIndex    int64
	IsRead      int
	AppKey      string
}

type IMentionMsgStorage interface {
	SaveMentionMsg(msg MentionMsg) error
	QryMentionMsgs(appkey, userId, targetId string, channelType pbobjs.ChannelType, startTime int64, count int, isPositiveOrder bool, startIndex int64) ([]*MentionMsg, error)
	QryUnreadMentionMsgs(appkey, userId, targetId string, channelType pbobjs.ChannelType, startTime int64, count int, isPositiveOrder bool) ([]*MentionMsg, error)
	QryMentionSenderIdsBaseIndex(appkey, userId, targetId string, channelType pbobjs.ChannelType, startIndex int64, count int) ([]*MentionMsg, error)
	MarkRead(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgIds []string) error
	DelMentionMsgs(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgIds []string) error
	DelMentionMsg(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgId string) error
}
