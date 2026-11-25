package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/storages/models"
)

type CmdSendboxMsgSplitDao struct {
	ID           int64  `gorm:"primary_key"`
	UserId       string `gorm:"user_id"`
	SendTime     int64  `gorm:"send_time"`
	PartitionDay int64  `gorm:"partition_day"`
	MsgId        string `gorm:"msg_id"`
	ChannelType  int    `gorm:"channel_type"`
	MsgBody      []byte `gorm:"msg_body"`
	AppKey       string `gorm:"app_key"`
	TargetId     string `gorm:"target_id"`
	MsgType      string `gorm:"msg_type"`
}

func (msg *CmdSendboxMsgSplitDao) TableName() string {
	return "cmdsendboxmsgs_split"
}

func (msg *CmdSendboxMsgSplitDao) SaveMsg(item models.Msg) error {
	return dbcommons.GetDb().Create(&CmdSendboxMsgSplitDao{
		UserId:       item.UserId,
		SendTime:     item.SendTime,
		PartitionDay: alignCmdMsgSplitPartitionStart(item.SendTime),
		MsgId:        item.MsgId,
		ChannelType:  int(item.ChannelType),
		MsgBody:      item.MsgBody,
		AppKey:       item.AppKey,
		TargetId:     item.TargetId,
		MsgType:      item.MsgType,
	}).Error
}

func (msg *CmdSendboxMsgSplitDao) UpsertMsg(item models.Msg) error {
	if item.UniqTag == "" {
		return msg.SaveMsg(item)
	}
	return upsertCmdMsgSplit(msg.TableName(), item)
}

func (msg *CmdSendboxMsgSplitDao) QryMsgsBaseTime(appkey, userId string, start int64, count int) ([]*models.Msg, error) {
	legacyMsgs, err := (&CmdSendboxMsgDao{}).QryMsgsBaseTime(appkey, userId, start, count)
	if err != nil {
		return []*models.Msg{}, err
	}
	splitMsgs, err := msg.qrySplitMsgsBaseTime(appkey, userId, start, count)
	if err != nil {
		return []*models.Msg{}, err
	}
	return mergeCmdCompatMsgs(legacyMsgs, splitMsgs, count), nil
}

func (msg *CmdSendboxMsgSplitDao) qrySplitMsgsBaseTime(appkey, userId string, start int64, count int) ([]*models.Msg, error) {
	var items []*CmdSendboxMsgSplitDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and send_time>? and partition_day>=?", appkey, userId, start, alignCmdMsgSplitPartitionStart(start)).Order("send_time asc").Limit(count).Find(&items).Error
	if err != nil {
		return []*models.Msg{}, err
	}
	cmdMsgs := make([]*models.Msg, 0, len(items))
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
