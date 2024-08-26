package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/storages/models"
)

type InboxMsgDao struct {
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

func (msg *InboxMsgDao) TableName() string {
	return "inboxmsgs"
}

func (msg *InboxMsgDao) SaveMsg(item models.Msg) error {
	daoItem := InboxMsgDao{
		UserId:      item.UserId,
		SendTime:    item.SendTime,
		MsgId:       item.MsgId,
		ChannelType: int(item.ChannelType),
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
		TargetId:    item.TargetId,
		MsgType:     item.MsgType,
	}
	err := dbcommons.GetDb().Create(&daoItem).Error
	return err
}

func (msg *InboxMsgDao) UpsertMsg(item models.Msg) error {
	return msg.SaveMsg(item)
}

func (msg *InboxMsgDao) QryMsgsBaseTime(appkey, userId string, start int64, count int) ([]*models.Msg, error) {
	var items []*InboxMsgDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and send_time>?", appkey, userId, start).Order("send_time asc").Limit(count).Find(&items).Error
	if err != nil {
		return []*models.Msg{}, err
	}
	inboxMsgs := []*models.Msg{}
	for _, item := range items {
		inboxMsgs = append(inboxMsgs, &models.Msg{
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
	return inboxMsgs, nil
}

func (msg *InboxMsgDao) DelMsgsBaseTime(appkey string, start int64) error {
	return dbcommons.GetDb().Where("app_key=? and send_time<?", appkey, start).Delete(&InboxMsgDao{}).Error
}

func (msg *InboxMsgDao) QryBaseTime(limit, offset int64) ([]*InboxMsgDao, error) {
	var items []*InboxMsgDao
	err := dbcommons.GetDb().Order("id asc").Limit(limit).Offset(offset).Find(&items).Error
	return items, err
}

func (msg *InboxMsgDao) DelBaseTime(id int64) error {
	return dbcommons.GetDb().Debug().Where("id<?", id).Delete(&InboxMsgDao{}).Error
}
