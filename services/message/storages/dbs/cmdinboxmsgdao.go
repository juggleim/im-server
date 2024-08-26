package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/storages/models"
)

type CmdInboxMsgDao struct {
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

func (msg *CmdInboxMsgDao) TableName() string {
	return "cmdinboxmsgs"
}

func (msg *CmdInboxMsgDao) SaveMsg(item models.Msg) error {
	daoItem := CmdInboxMsgDao{
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

func (msg *CmdInboxMsgDao) UpsertMsg(item models.Msg) error {
	if item.UniqTag == "" {
		return msg.SaveMsg(item)
	}
	sql := fmt.Sprintf("INSERT INTO %s (user_id,send_time,msg_id,channel_type,msg_body,app_key,target_id,msg_type,uniq_tag)VALUES(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE send_time=?, msg_id=?, msg_body=?", msg.TableName())
	return dbcommons.GetDb().Exec(sql, item.UserId, item.SendTime, item.MsgId, item.ChannelType, item.MsgBody, item.AppKey, item.TargetId, item.MsgType, item.UniqTag, item.SendTime, item.MsgId, item.MsgBody).Error
}

func (msg *CmdInboxMsgDao) QryMsgsBaseTime(appkey, userId string, start int64, count int) ([]*models.Msg, error) {
	var items []*CmdInboxMsgDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and send_time>?", appkey, userId, start).Order("send_time asc").Limit(count).Find(&items).Error
	if err != nil {
		return []*models.Msg{}, err
	}
	cmdMsgs := []*models.Msg{}
	for _, item := range items {
		cmdMsgs = append(cmdMsgs, &models.Msg{
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
	return cmdMsgs, nil
}

func (msg *CmdInboxMsgDao) DelMsgsBaseTime(appkey string, start int64) error {
	return dbcommons.GetDb().Where("app_key=? and send_time<?", appkey, start).Delete(&CmdInboxMsgDao{}).Error
}

func (msg *CmdInboxMsgDao) QryBaseTime(limit, offset int64) ([]*CmdInboxMsgDao, error) {
	var items []*CmdInboxMsgDao
	err := dbcommons.GetDb().Order("id asc").Limit(limit).Offset(offset).Find(&items).Error
	return items, err
}

func (msg *CmdInboxMsgDao) DelBaseTime(id int64) error {
	return dbcommons.GetDb().Debug().Where("id<?", id).Delete(&CmdInboxMsgDao{}).Error
}
