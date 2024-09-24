package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"sort"
	"time"
)

type GroupHisMsgDao struct {
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
	IsExt       int    `gorm:"is_ext"`
	IsReaction  int    `gorm:"is_reaction"`

	MemberCount int `gorm:"member_count"`
	ReadCount   int `gorm:"read_count"`
	IsDelete    int `gorm:"is_delete"`
}

func (msg GroupHisMsgDao) TableName() string {
	return "g_hismsgs"
}
func (msg GroupHisMsgDao) SaveGroupHisMsg(item models.GroupHisMsg) error {
	gMsg := GroupHisMsgDao{
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
		IsExt:       item.IsExt,
		IsReaction:  item.IsReaction,
		MemberCount: item.MemberCount,
		ReadCount:   item.ReadCount,
		IsDelete:    item.IsDelete,
	}
	err := dbcommons.GetDb().Create(&gMsg).Error
	return err
}

func (msg GroupHisMsgDao) QryLatestMsgSeqNo(appkey, converId string) int64 {
	var items []*GroupHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return items[0].MsgSeqNo
	}
	return 0
}

func (msg GroupHisMsgDao) QryLatestMsg(appkey, converId string) (*models.GroupHisMsg, error) {
	var items []*GroupHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return dbMsg2GrpMsg(items[0]), nil
	}
	return nil, err
}

func (msg GroupHisMsgDao) QryHisMsgsExcludeDel(appkey, converId, userId, targetId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*models.GroupHisMsg, error) {
	var items []*GroupHisMsgDao
	params := []interface{}{}
	condition := "app_key=? and conver_id=?"
	params = append(params, appkey)
	params = append(params, converId)

	condition = condition + fmt.Sprintf(" and msg_id not in (select msg_id from %s where app_key=? and user_id=? and target_id=?)", (&GroupDelHisMsgDao{}).TableName())
	params = append(params, appkey)
	params = append(params, userId)
	params = append(params, targetId)

	orderStr := "send_time desc"
	start := startTime
	if isPositiveOrder {
		orderStr = "send_time asc"
		if start < cleanTime {
			start = cleanTime
		}
		condition = condition + " and send_time>?"
		params = append(params, start)
	} else {
		if start <= 0 {
			start = time.Now().UnixMilli()
		}
		condition = condition + " and send_time<?"
		params = append(params, start)
		if cleanTime > 0 {
			condition = condition + " and send_time>?"
			params = append(params, cleanTime)
		}
	}
	if len(msgTypes) > 0 {
		condition = condition + " and msg_type in (?)"
		params = append(params, msgTypes)
	}
	condition = condition + " and is_delete=0"
	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if !isPositiveOrder {
		sort.Slice(items, func(i, j int) bool {
			return items[i].SendTime < items[j].SendTime
		})
	}
	retItems := []*models.GroupHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2GrpMsg(dbMsg))
	}
	return retItems, err
}

func (msg GroupHisMsgDao) QryHisMsgs(appkey, converId string, startTime, endTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string, excludeMsgIds []string) ([]*models.GroupHisMsg, error) {
	var items []*GroupHisMsgDao

	params := []interface{}{}
	condition := "app_key=? and conver_id=?"
	params = append(params, appkey)
	params = append(params, converId)

	orderStr := "send_time desc"
	start := startTime
	end := endTime
	if isPositiveOrder {
		orderStr = "send_time asc"
		if start < cleanTime {
			start = cleanTime
		}
		condition = condition + " and send_time>?"
		params = append(params, start)
		if end > 0 {
			condition = condition + " and send_time<?"
			params = append(params, end)
		}
	} else {
		if start <= 0 {
			start = time.Now().UnixMilli()
		}
		if end < cleanTime {
			end = cleanTime
		}
		condition = condition + " and send_time<? and send_time>?"
		params = append(params, start)
		params = append(params, end)
	}

	if len(excludeMsgIds) > 0 {
		condition = condition + " and msg_id not in (?)"
		params = append(params, excludeMsgIds)
	}
	if len(msgTypes) > 0 {
		condition = condition + " and msg_type in (?)"
		params = append(params, msgTypes)
	}
	condition = condition + " and is_delete=0"

	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if !isPositiveOrder {
		sort.Slice(items, func(i, j int) bool {
			return items[i].SendTime < items[j].SendTime
		})
	}
	retItems := []*models.GroupHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2GrpMsg(dbMsg))
	}
	return retItems, err
}

func (msg GroupHisMsgDao) UpdateMsgBody(appkey, conver_id, msgId, msgType string, msgBody []byte) error {
	upd := map[string]interface{}{}
	upd["msg_body"] = msgBody
	upd["msg_type"] = msgType
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id=?", appkey, conver_id, msgId).Update(upd).Error
}

func (msg GroupHisMsgDao) UpdateReadCount(appkey, converId, msgId string, readCount int) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id=? and read_count<?", appkey, converId, msgId, readCount).Update("read_count", readCount).Error
}

func (msg GroupHisMsgDao) FindById(appkey, conver_id, msgId string) (*models.GroupHisMsg, error) {
	var item GroupHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and msg_id=?", appkey, conver_id, msgId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return dbMsg2GrpMsg(&item), nil
}

func (msg GroupHisMsgDao) FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*models.GroupHisMsg, error) {
	var items []*GroupHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and send_time>? and msg_id in (?)", appkey, converId, cleanTime, msgIds).Order("send_time asc").Find(&items).Error

	retItems := []*models.GroupHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2GrpMsg(dbMsg))
	}
	return retItems, err
}

func (msg GroupHisMsgDao) DelMsgs(appkey, converId string, msgIds []string) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id in (?)", appkey, converId, msgIds).Update("is_delete", 1).Error
}

func (msg GroupHisMsgDao) UpdateMsgExtState(appkey, converId, msgId string, isExt int) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id=?", appkey, converId, msgId).Update("is_ext", isExt).Error
}

func (msg GroupHisMsgDao) UpdateMsgReactionState(appkey, converId, msgId string, isReaction int) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id=?", appkey, converId, msgId).Update("is_reaction", isReaction).Error
}

// TODO need batch delete
func (msg GroupHisMsgDao) DelSomeoneMsgsBaseTime(appkey, converId string, cleanTime int64, senderId string) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sender_id=? and send_time<?", appkey, converId, senderId, cleanTime).Update("is_delete", 1).Error
}

func dbMsg2GrpMsg(dbMsg *GroupHisMsgDao) *models.GroupHisMsg {
	return &models.GroupHisMsg{
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
			IsExt:       dbMsg.IsExt,
			IsReaction:  dbMsg.IsReaction,
			IsDelete:    dbMsg.IsDelete,
		},
		MemberCount: dbMsg.MemberCount,
		ReadCount:   dbMsg.ReadCount,
	}
}
