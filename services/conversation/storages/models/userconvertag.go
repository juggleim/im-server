package models

import "im-server/commons/pbdefines/pbobjs"

type UserConverTag struct {
	UserId      string
	Tag         string
	TagName     string
	CreatedTime int64
	AppKey      string
}

type IUserConverTagStorage interface {
	Upsert(item UserConverTag) error
	Delete(appkey, userId, tag string) error
	QryTags(appkey, userId string) ([]*UserConverTag, error)
	QryTagsByConver(appkey, userId, targetId string, channelType pbobjs.ChannelType) ([]*UserConverTag, error)
}

type ConverTagRel struct {
	UserId      string
	Tag         string
	TargetId    string
	ChannelType pbobjs.ChannelType
	CreatedTime int64
	AppKey      string
}

type TargetConver struct {
	TargetId    string
	ChannelType pbobjs.ChannelType
}

type IConverTagRelStorage interface {
	Create(item ConverTagRel) error
	BatchCreate(items []ConverTagRel) error
	Delete(appkey, userId, tag, targetId string, channelType pbobjs.ChannelType) error
	BatchDelete(appkey, userId, tag string, convers []TargetConver) error
	DeleteByTag(appkey, userId, tag string) error
}
