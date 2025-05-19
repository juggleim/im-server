package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/historymsg/storages/models"
)

type PrivateDelHisMsgDao struct {
	ID       int64  `gorm:"primary_key"`
	UserId   string `gorm:"user_id"`
	TargetId string `gorm:"target_id"`
	MsgId    string `gorm:"msg_id"`
	MsgTime  int64  `gorm:"msg_time"`
	MsgSeq   int64  `gorm:"msg_seq"`
	AppKey   string `gorm:"app_key"`
}

func (msg PrivateDelHisMsgDao) TableName() string {
	return "p_delhismsgs"
}

func (msg PrivateDelHisMsgDao) Create(item models.PrivateDelHisMsg) error {
	add := PrivateDelHisMsgDao{
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

func (msg PrivateDelHisMsgDao) BatchCreate(items []models.PrivateDelHisMsg) error {
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

func (msg PrivateDelHisMsgDao) QryDelHisMsgs(appkey, userId, targetId string, startTime int64, count int32, isPositive bool) ([]*models.PrivateDelHisMsg, error) {
	var items []*PrivateDelHisMsgDao
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
	retItems := []*models.PrivateDelHisMsg{}
	for _, item := range items {
		retItems = append(retItems, &models.PrivateDelHisMsg{
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
