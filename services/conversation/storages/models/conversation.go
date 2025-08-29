package models

import (
	"im-server/commons/pbdefines/pbobjs"
	"time"
)

type Conversation struct {
	AppKey      string
	UserId      string
	TargetId    string
	SubChannel  string
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
	FindOne(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType) (*Conversation, error)
	UpsertConversation(conversation Conversation) error
	Upsert(item Conversation) error
	QryConvers(appkey, userId string, startTime int64, count int32) ([]*Conversation, error)
	ClearTotalUnreadCount(appkey, userId string) error
}

type ConverConfItemKey string

const (
	ConverConfItemKey_RtcRoomId ConverConfItemKey = "rtc_room_id"
)

type ConverConf struct {
	ID          int64
	ConverId    string
	ConverType  pbobjs.ChannelType
	SubChannel  string
	ItemKey     string
	ItemValue   string
	ItemType    int32
	UpdatedTime int64
	CreatedTime int64
	AppKey      string
}

type IConverConfStorage interface {
	Upsert(item ConverConf) error
	BatchUpsert(items []ConverConf) error
	Delete(appkey, converId string, channelType pbobjs.ChannelType, subChannel string, itemKey string) error
	UpdateTime(appkey, converId string, channelType pbobjs.ChannelType, subChannel string, itemKey string, t time.Time) error
	QryConverConfs(appkey, converId, subChannel string, converType int32) (map[string]*ConverConf, error)
}
