package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/historymsg/storages/models"
)

type GroupDelHisMsgDao struct {
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

func (msg GroupDelHisMsgDao) TableName() string {
	return "g_delhismsgs"
}

func (msg GroupDelHisMsgDao) Create(item models.GroupDelHisMsg) error {
	add := GroupDelHisMsgDao{
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

func (msg GroupDelHisMsgDao) BatchCreate(items []models.GroupDelHisMsg) error {
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

func (msg GroupDelHisMsgDao) QryDelHisMsgs(appkey, userId, targetId, subChannel string, startTime int64, count int32, isPositive bool) ([]*models.GroupDelHisMsg, error) {
	var items []*GroupDelHisMsgDao
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
	retItems := []*models.GroupDelHisMsg{}
	for _, item := range items {
		retItems = append(retItems, &models.GroupDelHisMsg{
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

func (msg GroupDelHisMsgDao) QryDelHisMsgsByMsgIds(appkey, userId, targetId, subChannel string, msgIds []string) ([]*models.GroupDelHisMsg, error) {
	var items []*GroupDelHisMsgDao
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
	retItems := []*models.GroupDelHisMsg{}
	for _, item := range items {
		retItems = append(retItems, &models.GroupDelHisMsg{
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

func (msg GroupDelHisMsgDao) ExistDelHisMsg(appkey, userId, targetId, subChannel string) bool {
	var items []*GroupDelHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and sub_channel=?", appkey, userId, targetId, subChannel).Order("msg_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) <= 0 {
		return false
	}
	return true
}
