package models

import (
	"im-server/commons/pbdefines/pbobjs"
)

type Conversation struct {
	AppKey      string
	UserId      string
	TargetId    string
	ChannelType pbobjs.ChannelType

	SortTime int64
	SyncTime int64

	LatestMsgId          string
	LatestMsg            []byte
	LatestUnreadMsgIndex int64

	LatestReadMsgIndex int64
	LatestReadMsgId    string
	LatestReadMsgTime  int64

	IsTop          int
	TopUpdatedTime int64
	UndisturbType  int32

	UnreadTag  int
	ConverExts *pbobjs.ConverExts

	IsDeleted int

	//mention info
	MentionInfo *ConverMentionInfo
}

type ConverMentionInfo struct {
	IsMentioned     bool
	MentionMsgCount int
	SenderIds       []string
	MentionMsgs     []*MentionMsg
}

type IConversationStorage interface {
	FindOne(appkey, userId, targetId string, channelType pbobjs.ChannelType) (*Conversation, error)
	UpsertConversation(conversation Conversation) error
	Upsert(item Conversation) error
	QryConvers(appkey, userId string, startTime int64, count int32) ([]*Conversation, error)
	ClearTotalUnreadCount(appkey, userId string) error
}
