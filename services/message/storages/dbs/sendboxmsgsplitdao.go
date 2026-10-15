package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/storages/models"
)

// SendboxMsgSplitDao keeps the partitioned MySQL sendbox implementation
// separate from SendboxMsgDao so the legacy sendboxmsgs table remains intact.
type SendboxMsgSplitDao struct {
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

func (msg *SendboxMsgSplitDao) TableName() string {
	return "sendboxmsgs_split"
}

func (msg *SendboxMsgSplitDao) SaveMsg(item models.Msg) error {
	daoItem := SendboxMsgSplitDao{
		UserId:      item.UserId,
		SendTime:    item.SendTime,
		MsgId:       item.MsgId,
		ChannelType: int(item.ChannelType),
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
		TargetId:    item.TargetId,
		MsgType:     item.MsgType,
	}
	return dbcommons.GetDb().Create(&daoItem).Error
}

func (msg *SendboxMsgSplitDao) UpsertMsg(item models.Msg) error {
	return msg.SaveMsg(item)
}

func (msg *SendboxMsgSplitDao) QryMsgsBaseTime(appkey, userId string, start int64, count int) ([]*models.Msg, error) {
	var items []*SendboxMsgSplitDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and send_time>?", appkey, userId, start).Order("send_time asc").Limit(count).Find(&items).Error
	if err != nil {
		return []*models.Msg{}, err
	}
	sendboxMsgs := make([]*models.Msg, 0, len(items))
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
	return sendboxMsgs, nil
}
