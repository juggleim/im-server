package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"sort"
)

type GrpCastHisMsgDao struct {
	ID          int64  `gorm:"primary_key"`
	ConverId    string `gorm:"conver_id"`
	SenderId    string `gorm:"sender_id"`
	ReceiverId  string `gorm:"receiver_id"`
	ChannelType int    `gorm:"channel_type"`
	MsgType     string `gorm:"msg_type"`
	MsgId       string `gorm:"msg_id"`
	SendTime    int64  `gorm:"send_time"`
	MsgSeqNo    int64  `gorm:"msg_seq_no"`
	MsgBody     []byte `gorm:"msg_body"`
	AppKey      string `gorm:"app_key"`
}

func (msg GrpCastHisMsgDao) TableName() string {
	return "gc_hismsgs"
}

func (msg GrpCastHisMsgDao) SaveGrpCastHisMsg(item models.GrpCastHisMsg) error {
	gcMsg := GrpCastHisMsgDao{
		ConverId:    item.ConverId,
		SenderId:    item.SenderId,
		ReceiverId:  item.ReceiverId,
		ChannelType: int(item.ChannelType),
		MsgType:     item.MsgType,
		MsgId:       item.MsgId,
		SendTime:    item.SendTime,
		MsgSeqNo:    item.MsgSeqNo,
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
	}
	return dbcommons.GetDb().Create(&gcMsg).Error
}

func (msg GrpCastHisMsgDao) QryLatestMsgSeqNo(appkey, converId string) int64 {
	var items []*GrpCastHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return items[0].MsgSeqNo
	}
	return 0
}

func (msg GrpCastHisMsgDao) QryLatestMsg(appkey, converId string) (*models.GrpCastHisMsg, error) {
	var items []*GrpCastHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return dbMsg2GrpCastMsg(items[0]), nil
	}
	return nil, err
}

func (msg GrpCastHisMsgDao) QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*models.GrpCastHisMsg, error) {
	var items []*GrpCastHisMsgDao
	condition := "send_time<? and send_time>?"
	orderStr := "send_time desc"
	if isPositiveOrder {
		condition = "send_time>? and send_time>?"
		orderStr = "send_time asc"
	}
	var err error
	if len(msgTypes) > 0 {
		condition = condition + " and msg_type in (?)"
		err = dbcommons.GetDb().Where("app_key=? and conver_id=? and "+condition, appkey, converId, startTime, cleanTime, msgTypes).Order(orderStr).Limit(count).Find(&items).Error
	} else {
		err = dbcommons.GetDb().Where("app_key=? and conver_id=? and "+condition, appkey, converId, startTime, cleanTime).Order(orderStr).Limit(count).Find(&items).Error
	}
	if !isPositiveOrder {
		sort.Slice(items, func(i, j int) bool {
			return items[i].SendTime < items[j].SendTime
		})
	}
	retItems := []*models.GrpCastHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2GrpCastMsg(dbMsg))
	}
	return retItems, err
}

func (msg GrpCastHisMsgDao) FindById(appkey, conver_id, msgId string) (*models.GrpCastHisMsg, error) {
	var item GrpCastHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and msg_id=?", appkey, conver_id, msgId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return dbMsg2GrpCastMsg(&item), nil
}

func (msg GrpCastHisMsgDao) FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*models.GrpCastHisMsg, error) {
	var items []*GrpCastHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and send_time>? and msg_id in (?)", appkey, converId, cleanTime, msgIds).Order("send_time asc").Find(&items).Error
	retItems := []*models.GrpCastHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2GrpCastMsg(dbMsg))
	}
	return retItems, err
}

func dbMsg2GrpCastMsg(dbMsg *GrpCastHisMsgDao) *models.GrpCastHisMsg {
	return &models.GrpCastHisMsg{
		ConverId:    dbMsg.ConverId,
		SenderId:    dbMsg.SenderId,
		ReceiverId:  dbMsg.ReceiverId,
		ChannelType: pbobjs.ChannelType(dbMsg.ChannelType),
		MsgType:     dbMsg.MsgType,
		MsgId:       dbMsg.MsgId,
		SendTime:    dbMsg.SendTime,
		MsgSeqNo:    dbMsg.MsgSeqNo,
		MsgBody:     dbMsg.MsgBody,
		AppKey:      dbMsg.AppKey,
	}
}
