package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"sort"
)

type BrdCastHisMsgDao struct {
	ID          int64  `gorm:"primary_key"`
	ConverId    string `gorm:"conver_id"`
	SenderId    string `gorm:"sender_id"`
	ChannelType int    `gorm:"channel_type"`
	MsgType     string `gorm:"msg_type"`
	MsgId       string `gorm:"msg_id"`
	SendTime    int64  `gorm:"send_time"`
	MsgSeqNo    int64  `gorm:"msg_seq_no"`
	MsgBody     []byte `gorm:"msg_body"`
	AppKey      string `gorm:"app_key"`
}

func (msg BrdCastHisMsgDao) TableName() string {
	return "bc_hismsgs"
}

func (msg BrdCastHisMsgDao) SaveBrdCastHisMsg(item models.BrdCastHisMsg) error {
	bMsg := BrdCastHisMsgDao{
		ConverId:    item.ConverId,
		SenderId:    item.SenderId,
		ChannelType: int(item.ChannelType),
		MsgType:     item.MsgType,
		MsgId:       item.MsgId,
		SendTime:    item.SendTime,
		MsgSeqNo:    item.MsgSeqNo,
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
	}
	return dbcommons.GetDb().Create(&bMsg).Error
}

func (msg BrdCastHisMsgDao) QryLatestMsgSeqNo(appkey, converId string) int64 {
	var items []*BrdCastHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return items[0].MsgSeqNo
	}
	return 0
}

func (msg BrdCastHisMsgDao) QryLatestMsg(appkey, converId string) (*models.BrdCastHisMsg, error) {
	var items []*BrdCastHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return dbMsg2BrdMsg(items[0]), nil
	}
	return nil, err
}

func (msg BrdCastHisMsgDao) QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*models.BrdCastHisMsg, error) {
	var items []*BrdCastHisMsgDao
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
	retItems := []*models.BrdCastHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2BrdMsg(dbMsg))
	}
	return retItems, err
}

func (msg BrdCastHisMsgDao) FindById(appkey, conver_id, msgId string) (*models.BrdCastHisMsg, error) {
	var item BrdCastHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and msg_id=?", appkey, conver_id, msgId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return dbMsg2BrdMsg(&item), nil
}

func (msg BrdCastHisMsgDao) FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*models.BrdCastHisMsg, error) {
	var items []*BrdCastHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and send_time>? and msg_id in (?)", appkey, converId, cleanTime, msgIds).Order("send_time asc").Find(&items).Error
	retItems := []*models.BrdCastHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2BrdMsg(dbMsg))
	}
	return retItems, err
}

func dbMsg2BrdMsg(dbMsg *BrdCastHisMsgDao) *models.BrdCastHisMsg {
	return &models.BrdCastHisMsg{
		ConverId:    dbMsg.ConverId,
		SenderId:    dbMsg.SenderId,
		ChannelType: pbobjs.ChannelType(dbMsg.ChannelType),
		MsgType:     dbMsg.MsgType,
		MsgId:       dbMsg.MsgId,
		SendTime:    dbMsg.SendTime,
		MsgSeqNo:    dbMsg.MsgSeqNo,
		MsgBody:     dbMsg.MsgBody,
		AppKey:      dbMsg.AppKey,
	}
}
