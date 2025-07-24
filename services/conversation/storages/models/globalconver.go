package models

import (
	"im-server/commons/pbdefines/pbobjs"
)

type GlobalConver struct {
	Id          int64
	ConverId    string
	SenderId    string
	TargetId    string
	ChannelType pbobjs.ChannelType
	SubChannel  string
	UpdatedTime int64
	AppKey      string
}

type IGlobalConverStorage interface {
	UpsertConversation(item GlobalConver) error
	QryConversations(appkey, targetId, subChannel string, channelType pbobjs.ChannelType, startTime int64, count int32, isPositiveOder bool, excludeUserIds []string) ([]*GlobalConver, error)
}
