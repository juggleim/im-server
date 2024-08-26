package models

import (
	"im-server/commons/pbdefines/pbobjs"
)

type Conversation struct {
	UserId               string
	TargetId             string
	SortTime             int64
	ChannelType          pbobjs.ChannelType
	LatestMsgId          string
	LatestMsg            []byte
	LatestUnreadMsgIndex int64
	LatestReadMsgIndex   int64
	LatestReadMsgId      string
	LatestReadMsgTime    int64
	IsTop                int
	TopUpdatedTime       int64
	UndisturbType        int32
	IsDeleted            int
	SyncTime             int64
	UnreadTag            int
	AppKey               string
}

type IConversationStorage interface {
	FindOne(appkey, userId, targetId string, channelType pbobjs.ChannelType) (*Conversation, error)
	UpsertConversation(conversation Conversation) error
	QryConversations(appkey, userId, targetId string, channelType pbobjs.ChannelType, startTime int64, count int32, isPositiveOrder bool) ([]*Conversation, error)
	DelConversation(appkey, userId, targetId string, channelType pbobjs.ChannelType) error
	UpdateLatestReadMsgIndex(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgIndex int64, readMsgId string, readMsgTime int64) (int64, error)
	UpdateIsTopState(appkey, userId, targetId string, channelType pbobjs.ChannelType, isTop int, optTime int64) (int64, error)
	TotalUnreadCount(appkey, userId string, channelTypes []pbobjs.ChannelType, excludeConvers []*pbobjs.SimpleConversation) int64
	ClearTotalUnreadCount(appkey, userId string) error
	QryTopConvers(appkey, userId string, startTime, limit int64) ([]*Conversation, error)
	SyncConversations(appkey, userId string, startTime int64, count int32) ([]*Conversation, error)
	FindUndisturb(appkey, userId, targetId string, channelType pbobjs.ChannelType) (*Conversation, error)
	UpdateUndisturbType(appkey, userId, targetId string, channelType pbobjs.ChannelType, undisturbType int32) (int64, error)
	FindUnreadIndex(appkey, userId, targetId string, channelType pbobjs.ChannelType) (*Conversation, error)
	UpdateLatestMsgBody(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgId string, msgBs []byte) error
	UpdateUnreadTag(appkey, userId, targetId string, channelType pbobjs.ChannelType) (int64, error)
}
