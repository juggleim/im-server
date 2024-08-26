package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/storages/models"
)

type SendboxMsgDao struct {
	ID          int64  `gorm:"primary_key"`
	UserId      string `gorm:"user_id"`
	SendTime    int64  `gorm:"send_time"`
	MsgId       string `gorm:"msg_id"`
	ChannelType int    `gorm:"channel_type"`
	MsgBody     []byte `gorm:"msg_body"`
	AppKey      string `gorm:"app_key"`
	TargetId    string `gorm:"target_id"`
	MsgType     string `gorm:"msg_type"`
}

func (msg *SendboxMsgDao) TableName() string {
	return "sendboxmsgs"
}
func (msg *SendboxMsgDao) SaveMsg(item models.Msg) error {
	err := dbcommons.GetDb().Create(&SendboxMsgDao{
		UserId:      item.UserId,
		SendTime:    item.SendTime,
		MsgId:       item.MsgId,
		ChannelType: int(item.ChannelType),
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
		TargetId:    item.TargetId,
		MsgType:     item.MsgType,
	}).Error
	return err
}

func (msg *SendboxMsgDao) UpsertMsg(item models.Msg) error {
	return msg.SaveMsg(item)
}

func (msg *SendboxMsgDao) QryMsgsBaseTime(appkey, userId string, start int64, count int) ([]*models.Msg, error) {
	var items []*SendboxMsgDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and send_time>?", appkey, userId, start).Order("send_time asc").Limit(count).Find(&items).Error
	if err != nil {
		return []*models.Msg{}, err
	}
	sendboxMsgs := []*models.Msg{}
	for _, item := range items {
		sendboxMsgs = append(sendboxMsgs, &models.Msg{
			UserId:      item.UserId,
			SendTime:    item.SendTime,
			MsgId:       item.MsgId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			MsgBody:     item.MsgBody,
			AppKey:      item.AppKey,
			TargetId:    item.TargetId,
			MsgType:     item.MsgType,
		})
	}
	return sendboxMsgs, err
}

func (msg *SendboxMsgDao) DelMsgsBaseTime(appkey string, start int64) error {
	return dbcommons.GetDb().Where("app_key=? and send_time<?", appkey, start).Delete(&SendboxMsgDao{}).Error
}

func (msg *SendboxMsgDao) QryBaseTime(limit, offset int64) ([]*SendboxMsgDao, error) {
	var items []*SendboxMsgDao
	err := dbcommons.GetDb().Order("id asc").Limit(limit).Offset(offset).Find(&items).Error
	return items, err
}

func (msg *SendboxMsgDao) DelBaseTime(id int64) error {
	return dbcommons.GetDb().Debug().Where("id<?", id).Delete(&SendboxMsgDao{}).Error
}
