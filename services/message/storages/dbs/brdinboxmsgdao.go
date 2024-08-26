package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/storages/models"
	"sort"
)

type BrdInboxMsgDao struct {
	ID          int64  `gorm:"primary_key"`
	SenderId    string `gorm:"sender_id"`
	MsgId       string `gorm:"msg_id"`
	SendTime    int64  `gorm:"send_time"`
	ChannelType int    `gorm:"channel_type"`
	MsgBody     []byte `gorm:"msg_body"`
	AppKey      string `gorm:"app_key"`
}

func (msg *BrdInboxMsgDao) TableName() string {
	return "brdinboxmsgs"
}

func (msg *BrdInboxMsgDao) SaveMsg(item models.BrdInboxMsgMsg) error {
	err := dbcommons.GetDb().Create(&BrdInboxMsgDao{
		SenderId:    item.SenderId,
		SendTime:    item.SendTime,
		MsgId:       item.MsgId,
		ChannelType: int(item.ChannelType),
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
	}).Error
	return err
}

func (msg *BrdInboxMsgDao) QryMsgsBaseTime(appkey string, start int64, count int) ([]*models.BrdInboxMsgMsg, error) {
	var items []*BrdInboxMsgDao
	err := dbcommons.GetDb().Where("app_key=? and send_time>?", appkey, start).Order("send_time asc").Limit(count).Find(&items).Error
	if err != nil {
		return []*models.BrdInboxMsgMsg{}, err
	}
	msgs := []*models.BrdInboxMsgMsg{}
	for _, item := range items {
		msgs = append(msgs, &models.BrdInboxMsgMsg{
			SenderId:    item.SenderId,
			SendTime:    item.SendTime,
			MsgId:       item.MsgId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			MsgBody:     item.MsgBody,
			AppKey:      item.AppKey,
		})
	}
	return msgs, nil
}

func (msg *BrdInboxMsgDao) QryLatestMsg(appkey string, count int) ([]*models.BrdInboxMsgMsg, error) {
	var items []*BrdInboxMsgDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Order("send_time desc").Limit(count).Find(&items).Error
	if err != nil || len(items) <= 0 {
		return []*models.BrdInboxMsgMsg{}, err
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SendTime < items[j].SendTime
	})
	ret := []*models.BrdInboxMsgMsg{}
	for _, item := range items {
		ret = append(ret, &models.BrdInboxMsgMsg{
			SenderId:    item.SenderId,
			SendTime:    item.SendTime,
			MsgId:       item.MsgId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			MsgBody:     item.MsgBody,
			AppKey:      item.AppKey,
		})
	}
	return ret, nil
}

func (msg *BrdInboxMsgDao) DelMsgsBaseTime(appkey string, start int64) error {
	return dbcommons.GetDb().Where("app_key=? and send_time<?", appkey, start).Delete(&BrdInboxMsgDao{}).Error
}
