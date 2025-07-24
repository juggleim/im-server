package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"sort"
)

type MergedMsgDao struct {
	ID          int64  `gorm:"primary_key"`
	ParentMsgId string `gorm:"parent_msg_id"`
	FromId      string `gorm:"from_id"`
	TargetId    string `gorm:"target_id"`
	ChannelType int    `gorm:"channel_type"`
	SubChannel  string `gorm:"sub_channel"`
	MsgId       string `gorm:"msg_id"`
	MsgTime     int64  `gorm:"msg_time"`
	MsgBody     []byte `gorm:"msg_body"`
	AppKey      string `gorm:"app_key"`
}

func (msg MergedMsgDao) TableName() string {
	return "mergedmsgs"
}
func (msg MergedMsgDao) SaveMergedMsg(item models.MergedMsg) error {
	mMsg := MergedMsgDao{
		ParentMsgId: item.ParentMsgId,
		FromId:      item.FromId,
		TargetId:    item.TargetId,
		ChannelType: int(item.ChannelType),
		SubChannel:  item.SubChannel,
		MsgId:       item.MsgId,
		MsgTime:     item.MsgTime,
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
	}
	err := dbcommons.GetDb().Create(&mMsg).Error
	return err
}
func (msg MergedMsgDao) BatchSaveMergedMsgs(items []models.MergedMsg) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into %s (`parent_msg_id`,`from_id`,`target_id`,`channel_type`,`sub_channel`,`msg_id`,`msg_time`,`msg_body`,`app_key`)values ", msg.TableName())
	buffer.WriteString(sql)
	vals := []interface{}{}
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?,?,?,?,?,?,?);")
		} else {
			buffer.WriteString("(?,?,?,?,?,?,?,?,?),")
		}
		vals = append(vals, item.ParentMsgId, item.FromId, item.TargetId, item.ChannelType, item.SubChannel, item.MsgId, item.MsgTime, item.MsgBody, item.AppKey)
	}
	err := dbcommons.GetDb().Exec(buffer.String(), vals...).Error
	return err
}
func (msg MergedMsgDao) QryMergedMsgs(appkey, parentMsgId string, startTime int64, count int32, isPositiveOrder bool) ([]*models.MergedMsg, error) {
	var items []*MergedMsgDao
	condition := "msg_time<?"
	orderStr := "msg_time desc"
	if isPositiveOrder {
		condition = "msg_time>?"
		orderStr = "msg_time asc"
	}
	err := dbcommons.GetDb().Where("app_key=? and parent_msg_id=? and "+condition, appkey, parentMsgId, startTime).Order(orderStr).Limit(count).Find(&items).Error
	if !isPositiveOrder {
		sort.Slice(items, func(i, j int) bool {
			return items[i].MsgTime < items[j].MsgTime
		})
	}
	retItems := []*models.MergedMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, &models.MergedMsg{
			ParentMsgId: dbMsg.ParentMsgId,
			FromId:      dbMsg.FromId,
			TargetId:    dbMsg.TargetId,
			ChannelType: pbobjs.ChannelType(dbMsg.ChannelType),
			SubChannel:  dbMsg.SubChannel,
			MsgId:       dbMsg.MsgId,
			MsgTime:     dbMsg.MsgTime,
			MsgBody:     dbMsg.MsgBody,
			AppKey:      dbMsg.AppKey,
		})
	}
	return retItems, err
}
