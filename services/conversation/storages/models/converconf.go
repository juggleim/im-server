package models

import (
	"im-server/commons/pbdefines/pbobjs"
	"time"
)

// ConverConf is a configuration item attached to a conversation.
type ConverConf struct {
	ID          int64
	ConverId    string
	ConverType  pbobjs.ChannelType
	SubChannel  string
	ItemKey     string
	ItemValue   string
	ItemType    int
	CreatedTime time.Time
	UpdatedTime time.Time
	AppKey      string
}

type IConverConfStorage interface {
	Upsert(item ConverConf) error
	BatchUpsert(items []ConverConf) error
	Delete(appkey, converId string, converType pbobjs.ChannelType, subChannel, itemKey string) error
	Find(appkey, converId string, converType pbobjs.ChannelType, subChannel, itemKey string) (*ConverConf, error)
	QryConfs(appkey, converId string, converType pbobjs.ChannelType, subChannel string) ([]*ConverConf, error)
	QryConfsByItemKeys(appkey, converId string, converType pbobjs.ChannelType, subChannel string, itemKeys []string) (map[string]*ConverConf, error)
}
