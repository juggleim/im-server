package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"sort"
	"strings"
	"time"
)

type PrivateHisMsgDao struct {
	ID          int64  `gorm:"primary_key"`
	ConverId    string `gorm:"conver_id"`
	SenderId    string `gorm:"sender_id"`
	ReceiverId  string `gorm:"receiver_id"`
	MsgType     string `gorm:"msg_type"`
	ChannelType int    `gorm:"channel_type"`
	SendTime    int64  `gorm:"send_time"`
	MsgId       string `gorm:"msg_id"`
	MsgSeqNo    int64  `gorm:"msg_seq_no"`
	MsgBody     []byte `gorm:"msg_body"`
	IsRead      int    `gorm:"is_read"`
	AppKey      string `gorm:"app_key"`
	IsDelete    int    `gorm:"is_delete"`
	IsExt       int    `gorm:"is_ext"`
	IsExset     int    `gorm:"is_exset"`
	MsgExt      []byte `hbase:"msg_ext"`
	MsgExset    []byte `hbase:"msg_exset"`
}

func (msg PrivateHisMsgDao) TableName() string {
	return "p_hismsgs"
}

func (msg PrivateHisMsgDao) SavePrivateHisMsg(item models.PrivateHisMsg) error {
	pMsg := PrivateHisMsgDao{
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
		IsExset:     item.IsExset,
		MsgExt:      item.MsgExt,
		MsgExset:    item.MsgExset,
		IsDelete:    item.IsDelete,
		IsRead:      item.IsRead,
	}
	err := dbcommons.GetDb().Create(&pMsg).Error
	return err
}

func (msg PrivateHisMsgDao) FindById(appkey, conver_id, msgId string) (*models.PrivateHisMsg, error) {
	var item PrivateHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and msg_id=?", appkey, conver_id, msgId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return dbMsg2PrivateMsg(&item), nil
}

func (msg PrivateHisMsgDao) UpdateMsgBody(appkey, conver_id, msgId, msgType string, msgBody []byte) error {
	upd := map[string]interface{}{}
	upd["msg_body"] = msgBody
	upd["msg_type"] = msgType
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id=?", appkey, conver_id, msgId).Update(upd).Error
}

func (msg PrivateHisMsgDao) QryLatestMsgSeqNo(appkey, converId string) int64 {
	var items []*PrivateHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return items[0].MsgSeqNo
	}
	return 0
}

func (msg PrivateHisMsgDao) QryLatestMsg(appkey, converId string) (*models.PrivateHisMsg, error) {
	var items []*PrivateHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?", appkey, converId).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return dbMsg2PrivateMsg(items[0]), nil
	}
	return nil, err
}

func (msg PrivateHisMsgDao) QryHisMsgsExcludeDel(appkey, converId, userId, targetId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*models.PrivateHisMsg, error) {
	var items []*PrivateHisMsgDao
	params := []interface{}{}
	condition := "app_key=? and conver_id=?"
	params = append(params, appkey)
	params = append(params, converId)

	hismsgsTable := msg.TableName()
	delHismsgsTable := (&PrivateDelHisMsgDao{}).TableName()
	condition = condition + fmt.Sprintf(" and not exists (select msg_id from %s where app_key=? and user_id=? and target_id=? and %s.msg_id=%s.msg_id)", delHismsgsTable, hismsgsTable, delHismsgsTable)
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
	retItems := []*models.PrivateHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2PrivateMsg(dbMsg))
	}
	return retItems, err
}

func (msg PrivateHisMsgDao) QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string, excludeMsgIds []string) ([]*models.PrivateHisMsg, error) {
	var items []*PrivateHisMsgDao
	params := []interface{}{}
	condition := "app_key=? and conver_id=?"
	params = append(params, appkey)
	params = append(params, converId)
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
	retItems := []*models.PrivateHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2PrivateMsg(dbMsg))
	}
	return retItems, err
}

func (msg PrivateHisMsgDao) MarkReadByMsgIds(appkey, converId string, msgIds []string) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id in (?)", appkey, converId, msgIds).Update("is_read", 1).Error
}

func (msg PrivateHisMsgDao) MarkReadByScope(appkey, converId string, start, end int64) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_index>=? and msg_index<=?", appkey, converId, start, end).Update("is_read", 1).Error
}

func (msg PrivateHisMsgDao) FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*models.PrivateHisMsg, error) {
	var items []*PrivateHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and send_time>? and msg_id in (?)", appkey, converId, cleanTime, msgIds).Order("send_time asc").Find(&items).Error

	retItems := []*models.PrivateHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2PrivateMsg(dbMsg))
	}

	return retItems, err
}

func (msg PrivateHisMsgDao) FindByConvers(appkey string, convers []models.ConverItem) ([]*models.PrivateHisMsg, error) {
	length := len(convers)
	if length <= 0 {
		return []*models.PrivateHisMsg{}, nil
	}
	var items []*PrivateHisMsgDao
	var sqlBuilder strings.Builder
	params := []interface{}{}
	sqlBuilder.WriteString("app_key=? and (")
	params = append(params, appkey)
	for i, conver := range convers {
		if i == length-1 {
			sqlBuilder.WriteString("(conver_id=? and msg_id=?)")
		} else {
			sqlBuilder.WriteString("(conver_id=? and msg_id=?) or ")
		}
		params = append(params, conver.ConverId)
		params = append(params, conver.MsgId)
	}
	sqlBuilder.WriteString(")")
	err := dbcommons.GetDb().Where(sqlBuilder.String(), params...).Find(&items).Error
	retItems := []*models.PrivateHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2PrivateMsg(dbMsg))
	}
	return retItems, err
}

func (msg PrivateHisMsgDao) DelMsgs(appkey, converId string, msgIds []string) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id in (?)", appkey, converId, msgIds).Update("is_delete", 1).Error
}

func (msg PrivateHisMsgDao) UpdateMsgExtState(appkey, converId, msgId string, isExt int) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id=?", appkey, converId, msgId).Update("is_ext", isExt).Error
}

func (msg PrivateHisMsgDao) UpdateMsgExt(appkey, converId, msgId string, ext []byte) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id=?", appkey, converId, msgId).Update("msg_ext", ext).Error
}

func (msg PrivateHisMsgDao) UpdateMsgExsetState(appkey, converId, msgId string, isExset int) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id=?", appkey, converId, msgId).Update("is_exset", isExset).Error
}

func (msg PrivateHisMsgDao) UpdateMsgExset(appkey, converId, msgId string, ext []byte) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and msg_id=?", appkey, converId, msgId).Update("msg_exset", ext).Error
}

// TODO need batch delete
func (msg PrivateHisMsgDao) DelSomeoneMsgsBaseTime(appkey, converId string, cleanTime int64, senderId string) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sender_id=? and send_time<?", appkey, converId, senderId, cleanTime).Update("is_delete", 1).Error
}

func dbMsg2PrivateMsg(dbMsg *PrivateHisMsgDao) *models.PrivateHisMsg {
	return &models.PrivateHisMsg{
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
			IsExset:     dbMsg.IsExset,
			MsgExt:      dbMsg.MsgExt,
			MsgExset:    dbMsg.MsgExset,
			IsDelete:    dbMsg.IsDelete,
		},
		IsRead: dbMsg.IsRead,
	}
}
