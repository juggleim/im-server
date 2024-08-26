package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"sort"
)

type SystemHisMsgDao struct {
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
	IsRead      int    `gorm:"is_read"`
	AppKey      string `gorm:"app_key"`
}

func (msg SystemHisMsgDao) TableName() string {
	return "s_hismsgs"
}
func (msg SystemHisMsgDao) SaveSystemHisMsg(item models.SystemHisMsg) error {
	sMsg := SystemHisMsgDao{
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
		IsRead:      item.IsRead,
	}
	err := dbcommons.GetDb().Create(&sMsg).Error
	return err
}

func (msg SystemHisMsgDao) QryLatestMsgSeqNo(appkey, converId string) int64 {
	var items []*SystemHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return items[0].MsgSeqNo
	}
	return 0
}

func (msg SystemHisMsgDao) QryLatestMsg(appkey, converId string) (*models.SystemHisMsg, error) {
	var items []*SystemHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return dbMsg2SysMsg(items[0]), nil
	}
	return nil, err
}

func (msg SystemHisMsgDao) QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*models.SystemHisMsg, error) {
	var items []*SystemHisMsgDao
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
	retItems := []*models.SystemHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2SysMsg(dbMsg))
	}
	return retItems, err
}

func (msg SystemHisMsgDao) FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*models.SystemHisMsg, error) {
	var items []*SystemHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and send_time>? and msg_id in (?)", appkey, converId, cleanTime, msgIds).Order("send_time asc").Find(&items).Error
	retItems := []*models.SystemHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2SysMsg(dbMsg))
	}
	return retItems, err
}

func dbMsg2SysMsg(dbMsg *SystemHisMsgDao) *models.SystemHisMsg {
	return &models.SystemHisMsg{
		HisMsg: models.HisMsg{
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
		},
		IsRead: dbMsg.IsRead,
	}
}
