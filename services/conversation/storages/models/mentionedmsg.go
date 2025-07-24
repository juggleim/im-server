package models

import "im-server/commons/pbdefines/pbobjs"

type MentionMsg struct {
	UserId      string
	TargetId    string
	ChannelType pbobjs.ChannelType
	SubChannel  string
	SenderId    string
	MentionType pbobjs.MentionType
	MsgId       string
	MsgTime     int64
	MsgIndex    int64
	IsRead      int
	AppKey      string
}

type ConverItem struct {
	TargetId    string
	ChannelType pbobjs.ChannelType
	MsgIndex    int64
	SubChannel  string
}

type IMentionMsgStorage interface {
	SaveMentionMsg(msg MentionMsg) error
	QryMentionMsgs(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, startTime int64, count int, isPositiveOrder bool, startIndex int64, cleanTime int64) ([]*MentionMsg, error)
	QryUnreadMentionMsgs(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, startTime int64, count int, isPositiveOrder bool, cleanTime int64) ([]*MentionMsg, error)
	QryMentionSenderIdsBaseIndex(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, startIndex int64, count int) ([]*MentionMsg, error)
	BatchQryMentionSenderIdsBaseIndex(appkey, userId string, convers []ConverItem) ([]*MentionMsg, error)
	MarkRead(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, msgIds []string) error
	DelMentionMsgs(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, msgIds []string) error
	DelMentionMsg(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, msgId string) error
	CleanMentionMsgsBaseIndex(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, msgIndex int64) error
	CleanMentionMsgsBaseUserId(appkey, userId string) error
	DelOnlyByMsgIds(appkey string, msgIds []string) error
}
