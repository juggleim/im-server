package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/historymsg/storages/models"
)

type PrivateDelHisMsgDao struct {
	ID            int64  `gorm:"primary_key"`
	UserId        string `gorm:"user_id"`
	TargetId      string `gorm:"target_id"`
	SubChannel    string `gorm:"sub_channel"`
	MsgId         string `gorm:"msg_id"`
	MsgTime       int64  `gorm:"msg_time"`
	MsgSeq        int64  `gorm:"msg_seq"`
	EffectiveTime int64  `gorm:"effective_time"`
	AppKey        string `gorm:"app_key"`
}

func (msg PrivateDelHisMsgDao) TableName() string {
	return "p_delhismsgs"
}

func (msg PrivateDelHisMsgDao) Create(item models.PrivateDelHisMsg) error {
	add := PrivateDelHisMsgDao{
		UserId:        item.UserId,
		TargetId:      item.TargetId,
		SubChannel:    item.SubChannel,
		MsgId:         item.MsgId,
		MsgTime:       item.MsgTime,
		MsgSeq:        item.MsgSeq,
		EffectiveTime: item.EffectiveTime,
		AppKey:        item.AppKey,
	}
	err := dbcommons.GetDb().Create(&add).Error
	return err
}

func (msg PrivateDelHisMsgDao) BatchCreate(items []models.PrivateDelHisMsg) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into %s (`user_id`,`target_id`,`sub_channel`,`msg_id`,`msg_time`,`msg_seq`,`effective_time`,`app_key`)values", msg.TableName())
	params := []interface{}{}

	buffer.WriteString(sql)
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?,?,?,?,?,?);")
		} else {
			buffer.WriteString("(?,?,?,?,?,?,?,?),")
		}
		params = append(params, item.UserId, item.TargetId, item.SubChannel, item.MsgId, item.MsgTime, item.MsgSeq, item.EffectiveTime, item.AppKey)
	}

	err := dbcommons.GetDb().Exec(buffer.String(), params...).Error
	return err
}

func (msg PrivateDelHisMsgDao) QryDelHisMsgs(appkey, userId, targetId, subChannel string, startTime int64, count int32, isPositive bool) ([]*models.PrivateDelHisMsg, error) {
	var items []*PrivateDelHisMsgDao
	params := []interface{}{}
	condition := "app_key=? and user_id=? and target_id=? and sub_channel=?"
	params = append(params, appkey)
	params = append(params, userId)
	params = append(params, targetId)
	params = append(params, subChannel)
	orderStr := "msg_time desc"
	if isPositive {
		condition = condition + " and msg_time>?"
	} else {
		condition = condition + " and msg_time<?"
	}
	params = append(params, startTime)
	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return nil, err
	}
	retItems := []*models.PrivateDelHisMsg{}
	for _, item := range items {
		retItems = append(retItems, &models.PrivateDelHisMsg{
			UserId:        item.UserId,
			TargetId:      item.TargetId,
			SubChannel:    item.SubChannel,
			MsgId:         item.MsgId,
			MsgTime:       item.MsgTime,
			MsgSeq:        item.MsgSeq,
			EffectiveTime: item.EffectiveTime,
			AppKey:        item.AppKey,
		})
	}
	return retItems, err
}

func (msg PrivateDelHisMsgDao) QryDelHisMsgsByMsgIds(appkey, userId, targetId, subChannel string, msgIds []string) ([]*models.PrivateDelHisMsg, error) {
	var items []*PrivateDelHisMsgDao
	params := []interface{}{}
	condition := "app_key=? and user_id=? and target_id=? and sub_channel=? and msg_id in (?)"
	params = append(params, appkey)
	params = append(params, userId)
	params = append(params, targetId)
	params = append(params, subChannel)
	params = append(params, msgIds)
	err := dbcommons.GetDb().Where(condition, params...).Find(&items).Error
	if err != nil {
		return nil, err
	}
	retItems := []*models.PrivateDelHisMsg{}
	for _, item := range items {
		retItems = append(retItems, &models.PrivateDelHisMsg{
			UserId:        item.UserId,
			TargetId:      item.TargetId,
			SubChannel:    item.SubChannel,
			MsgId:         item.MsgId,
			MsgTime:       item.MsgTime,
			MsgSeq:        item.MsgSeq,
			EffectiveTime: item.EffectiveTime,
			AppKey:        item.AppKey,
		})
	}
	return retItems, err
}
