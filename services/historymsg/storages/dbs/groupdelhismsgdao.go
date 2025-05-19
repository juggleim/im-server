package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/historymsg/storages/models"
)

type GroupDelHisMsgDao struct {
	ID       int64  `gorm:"primary_key"`
	UserId   string `gorm:"user_id"`
	TargetId string `gorm:"target_id"`
	MsgId    string `gorm:"msg_id"`
	MsgTime  int64  `gorm:"msg_time"`
	MsgSeq   int64  `gorm:"msg_seq"`
	AppKey   string `gorm:"app_key"`
}

func (msg GroupDelHisMsgDao) TableName() string {
	return "g_delhismsgs"
}

func (msg GroupDelHisMsgDao) Create(item models.GroupDelHisMsg) error {
	add := GroupDelHisMsgDao{
		UserId:   item.UserId,
		TargetId: item.TargetId,
		MsgId:    item.MsgId,
		MsgTime:  item.MsgTime,
		MsgSeq:   item.MsgSeq,
		AppKey:   item.AppKey,
	}
	err := dbcommons.GetDb().Create(&add).Error
	return err
}

func (msg GroupDelHisMsgDao) BatchCreate(items []models.GroupDelHisMsg) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into %s (`user_id`,`target_id`,`msg_id`,`msg_time`,`msg_seq`,`app_key`)values", msg.TableName())
	params := []interface{}{}

	buffer.WriteString(sql)
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?,?,?,?);")
		} else {
			buffer.WriteString("(?,?,?,?,?,?),")
		}
		params = append(params, item.UserId, item.TargetId, item.MsgId, item.MsgTime, item.MsgSeq, item.AppKey)
	}

	err := dbcommons.GetDb().Exec(buffer.String(), params...).Error
	return err
}

func (msg GroupDelHisMsgDao) QryDelHisMsgs(appkey, userId, targetId string, startTime int64, count int32, isPositive bool) ([]*models.GroupDelHisMsg, error) {
	var items []*GroupDelHisMsgDao
	params := []interface{}{}
	condition := "app_key=? and user_id=? and target_id=?"
	params = append(params, appkey)
	params = append(params, userId)
	params = append(params, targetId)
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
			UserId:   item.UserId,
			TargetId: item.TargetId,
			MsgId:    item.MsgId,
			MsgTime:  item.MsgTime,
			MsgSeq:   item.MsgSeq,
			AppKey:   item.AppKey,
		})
	}
	return retItems, err
}

func (msg GroupDelHisMsgDao) ExistDelHisMsg(appkey, userId, targetId string) bool {
	var items []*GroupDelHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=?", appkey, userId, targetId).Order("msg_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) <= 0 {
		return false
	}
	return true
}
