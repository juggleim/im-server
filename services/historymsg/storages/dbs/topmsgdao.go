package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"time"
)

type TopMsgDao struct {
	ID          int64     `gorm:"primary_key"`
	ConverId    string    `gorm:"conver_id"`
	SubChannel  string    `gorm:"sub_channel"`
	ChannelType int       `gorm:"channel_type"`
	MsgId       string    `gorm:"msg_id"`
	UserId      string    `gorm:"user_id"`
	CreatedTime time.Time `gorm:"created_time"`
	AppKey      string    `gorm:"app_key"`
}

func (msg TopMsgDao) TableName() string {
	return "topmsgs"
}

func (msg TopMsgDao) Upsert(item models.TopMsg) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,conver_id,sub_channel,channel_type,msg_id,user_id)VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE msg_id=VALUES(msg_id),user_id=VALUES(user_id)", msg.TableName()), item.AppKey, item.ConverId, item.SubChannel, item.ChannelType, item.MsgId, item.UserId).Error
}

func (msg TopMsgDao) FindTopMsg(appkey, converId, subChannel string, channelType pbobjs.ChannelType) (*models.TopMsg, error) {
	var item TopMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and sub_channel=? and channel_type=?", appkey, converId, subChannel, channelType).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.TopMsg{
		ID:          item.ID,
		ConverId:    item.ConverId,
		SubChannel:  item.SubChannel,
		ChannelType: pbobjs.ChannelType(item.ChannelType),
		MsgId:       item.MsgId,
		UserId:      item.UserId,
		CreatedTime: item.CreatedTime,
		AppKey:      item.AppKey,
	}, nil
}

func (msg TopMsgDao) DelTopMsg(appkey, converId, subChannel string, channelType pbobjs.ChannelType, msgId string) error {
	return dbcommons.GetDb().Where("app_key=? and conver_id=? and sub_channel=? and channel_type=? and msg_id=?", appkey, converId, subChannel, channelType, msgId).Delete(&TopMsgDao{}).Error
}
