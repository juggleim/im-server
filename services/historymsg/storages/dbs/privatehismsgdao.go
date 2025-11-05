package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PrivateHisMsgDao struct {
	ID                int64  `gorm:"primary_key"`
	ConverId          string `gorm:"conver_id"`
	SubChannel        string `gorm:"sub_channel"`
	SenderId          string `gorm:"sender_id"`
	ReceiverId        string `gorm:"receiver_id"`
	MsgType           string `gorm:"msg_type"`
	ChannelType       int    `gorm:"channel_type"`
	SendTime          int64  `gorm:"send_time"`
	MsgId             string `gorm:"msg_id"`
	MsgSeqNo          int64  `gorm:"msg_seq_no"`
	MsgBody           []byte `gorm:"msg_body"`
	IsRead            int    `gorm:"is_read"`
	ReadTime          int64  `gorm:"read_time"`
	AppKey            string `gorm:"app_key"`
	IsDelete          int    `gorm:"is_delete"`
	IsExt             int    `gorm:"is_ext"`
	IsExset           int    `gorm:"is_exset"`
	MsgExt            []byte `hbase:"msg_ext"`
	MsgExset          []byte `hbase:"msg_exset"`
	DestroyTime       int64  `gorm:"destroy_time"`
	LifeTimeAfterRead int64  `gorm:"life_time_after_read"`
}

func (msg PrivateHisMsgDao) TableName() string {
	return "p_hismsgs"
}

func (msg PrivateHisMsgDao) SavePrivateHisMsg(item models.PrivateHisMsg) error {
	pMsg := PrivateHisMsgDao{
		ConverId:          item.ConverId,
		SubChannel:        item.SubChannel,
		SenderId:          item.SenderId,
		ReceiverId:        item.ReceiverId,
		ChannelType:       int(item.ChannelType),
		MsgType:           item.MsgType,
		MsgId:             item.MsgId,
		SendTime:          item.SendTime,
		MsgSeqNo:          item.MsgSeqNo,
		MsgBody:           item.MsgBody,
		AppKey:            item.AppKey,
		IsExt:             item.IsExt,
		IsExset:           item.IsExset,
		MsgExt:            item.MsgExt,
		MsgExset:          item.MsgExset,
		IsDelete:          item.IsDelete,
		IsRead:            item.IsRead,
		DestroyTime:       item.DestroyTime,
		LifeTimeAfterRead: item.LifeTimeAfterRead,
	}
	err := dbcommons.GetDb().Create(&pMsg).Error
	return err
}

func (msg PrivateHisMsgDao) FindById(appkey, conver_id, subChannel, msgId string) (*models.PrivateHisMsg, error) {
	var item PrivateHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and sub_channel=? and msg_id=?", appkey, conver_id, subChannel, msgId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return dbMsg2PrivateMsg(&item), nil
}

func (msg PrivateHisMsgDao) UpdateMsgBody(appkey, conver_id, subChannel, msgId, msgType string, msgBody []byte) error {
	upd := map[string]interface{}{}
	upd["msg_body"] = msgBody
	upd["msg_type"] = msgType
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_id=?", appkey, conver_id, subChannel, msgId).Update(upd).Error
}

func (msg PrivateHisMsgDao) QryLatestMsgSeqNo(appkey, converId, subChannel string) int64 {
	var items []*PrivateHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and sub_channel=?", appkey, converId, subChannel).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return items[0].MsgSeqNo
	}
	return 0
}

func (msg PrivateHisMsgDao) QryLatestMsg(appkey, converId, subChannel string) (*models.PrivateHisMsg, error) {
	var items []*PrivateHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and sub_channel=?", appkey, converId, subChannel).Order("send_time desc").Limit(1).Find(&items).Error
	if err == nil && len(items) > 0 {
		return dbMsg2PrivateMsg(items[0]), nil
	}
	return nil, err
}

func (msg PrivateHisMsgDao) QryHisMsgsExcludeDel(appkey, converId, subChannel, userId, targetId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*models.PrivateHisMsg, error) {
	curr := time.Now().UnixMilli()
	var items []*PrivateHisMsgDao
	params := []interface{}{}
	hismsgTableName := msg.TableName()
	delHismsgTableName := (&PrivateDelHisMsgDao{}).TableName()
	sql := fmt.Sprintf("select his.* from %s as his left join %s as delhis on his.app_key=delhis.app_key and delhis.user_id=? and delhis.target_id=? and his.sub_channel=delhis.sub_channel and his.msg_id=delhis.msg_id where his.app_key=? and his.conver_id=? and his.sub_channel=?", hismsgTableName, delHismsgTableName)
	params = append(params, userId)
	params = append(params, targetId)
	params = append(params, appkey)
	params = append(params, converId)
	params = append(params, subChannel)

	orderStr := "his.send_time desc"
	start := startTime
	if isPositiveOrder {
		orderStr = "his.send_time asc"
		if start < cleanTime {
			start = cleanTime
		}
		sql = sql + " and his.send_time>?"
		params = append(params, start)
	} else {
		if start <= 0 {
			start = curr
		}
		sql = sql + " and his.send_time<?"
		params = append(params, start)
		if cleanTime > 0 {
			sql = sql + " and his.send_time>?"
			params = append(params, cleanTime)
		}
	}
	if len(msgTypes) > 0 {
		sql = sql + " and his.msg_type in (?)"
		params = append(params, msgTypes)
	}
	sql = sql + " and his.is_delete=0 and delhis.msg_id is null"
	sql = sql + " and his.is_delete=0 and (his.destroy_time=0 or his.destroy_time>?) and delhis.msg_id is null"
	params = append(params, curr)
	err := dbcommons.GetDb().Raw(sql, params...).Order(orderStr).Limit(count).Find(&items).Error
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

func (msg PrivateHisMsgDao) QryHisMsgs(appkey, converId, subChannel string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string, excludeMsgIds []string) ([]*models.PrivateHisMsg, error) {
	curr := time.Now().UnixMilli()
	var items []*PrivateHisMsgDao
	params := []interface{}{}
	condition := "app_key=? and conver_id=? and sub_channel=?"
	params = append(params, appkey)
	params = append(params, converId)
	params = append(params, subChannel)
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
			start = curr
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
	condition = condition + " and is_delete=0 and (destroy_time=0 or destroy_time>?)"
	params = append(params, curr)

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

func (msg PrivateHisMsgDao) MarkReadByMsgIds(appkey, converId, subChannel string, msgIds []string) error {
	upd := map[string]interface{}{}
	upd["is_read"] = 1
	upd["read_time"] = time.Now().UnixMilli()
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_id in (?)", appkey, converId, subChannel, msgIds).Update(upd).Error
}

func (msg PrivateHisMsgDao) UpdateDestroyTimeAfterReadByMsgIds(appkey, converId, subChannel string, msgIds []string) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_id in (?) and life_time_after_read>0", appkey, converId, subChannel, msgIds).Update("destroy_time", gorm.Expr("(UNIX_TIMESTAMP(NOW(3))*1000)+life_time_after_read")).Error
}

func (msg PrivateHisMsgDao) MarkReadByScope(appkey, converId, subChannel string, start, end int64) error {
	upd := map[string]interface{}{}
	upd["is_read"] = 1
	upd["read_time"] = time.Now().UnixMilli()
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_index>=? and msg_index<=?", appkey, converId, subChannel, start, end).Update(upd).Error
}

func (msg PrivateHisMsgDao) UpdateDestroyTimeAfterReadByScope(appkey, converId, subChannel string, start, end int64) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_index>=? and msg_index<=? and life_time_after_read>0", appkey, converId, subChannel, start, end).Update("destroy_time", gorm.Expr("(UNIX_TIMESTAMP(NOW(3))*1000)+life_time_after_read")).Error
}

func (msg PrivateHisMsgDao) FindByIds(appkey, converId, subChannel string, msgIds []string, cleanTime int64) ([]*models.PrivateHisMsg, error) {
	var items []*PrivateHisMsgDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and sub_channel=? and send_time>? and msg_id in (?)", appkey, converId, subChannel, cleanTime, msgIds).Order("send_time asc").Find(&items).Error

	retItems := []*models.PrivateHisMsg{}
	for _, dbMsg := range items {
		retItems = append(retItems, dbMsg2PrivateMsg(dbMsg))
	}

	return retItems, err
}

func (msg PrivateHisMsgDao) FindReadTimeByIds(appkey, converId, subChannel string, msgIds []string) ([]*models.PrivateHisMsg, error) {
	var items []*PrivateHisMsgDao
	err := dbcommons.GetDb().Select("id,msg_id,is_read,read_time").Where("app_key=? and conver_id=? and sub_channel=? and msg_id in (?)", appkey, converId, subChannel, msgIds).Find(&items).Error

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
			sqlBuilder.WriteString("(conver_id=? and sub_channel=? and msg_id=?)")
		} else {
			sqlBuilder.WriteString("(conver_id=? and sub_channel=? and msg_id=?) or ")
		}
		params = append(params, conver.ConverId)
		params = append(params, conver.SubChannel)
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

func (msg PrivateHisMsgDao) DelMsgs(appkey, converId, subChannel string, msgIds []string) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_id in (?)", appkey, converId, subChannel, msgIds).Update("is_delete", 1).Error
}

func (msg PrivateHisMsgDao) UpdateMsgExtState(appkey, converId, subChannel, msgId string, isExt int) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_id=?", appkey, converId, subChannel, msgId).Update("is_ext", isExt).Error
}

func (msg PrivateHisMsgDao) UpdateMsgExt(appkey, converId, subChannel, msgId string, ext []byte) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_id=?", appkey, converId, subChannel, msgId).Update("msg_ext", ext).Error
}

func (msg PrivateHisMsgDao) UpdateMsgExsetState(appkey, converId, subChannel, msgId string, isExset int) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_id=?", appkey, converId, subChannel, msgId).Update("is_exset", isExset).Error
}

func (msg PrivateHisMsgDao) UpdateMsgExset(appkey, converId, subChannel, msgId string, ext []byte) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and msg_id=?", appkey, converId, subChannel, msgId).Update("msg_exset", ext).Error
}

// TODO need batch delete
func (msg PrivateHisMsgDao) DelSomeoneMsgsBaseTime(appkey, converId, subChannel string, cleanTime int64, senderId string) error {
	return dbcommons.GetDb().Model(&msg).Where("app_key=? and conver_id=? and sub_channel=? and sender_id=? and send_time<?", appkey, converId, subChannel, senderId, cleanTime).Update("is_delete", 1).Error
}

func dbMsg2PrivateMsg(dbMsg *PrivateHisMsgDao) *models.PrivateHisMsg {
	return &models.PrivateHisMsg{
		HisMsg: models.HisMsg{
			ConverId:          dbMsg.ConverId,
			SubChannel:        dbMsg.SubChannel,
			SenderId:          dbMsg.SenderId,
			ReceiverId:        dbMsg.ReceiverId,
			ChannelType:       pbobjs.ChannelType(dbMsg.ChannelType),
			MsgType:           dbMsg.MsgType,
			MsgId:             dbMsg.MsgId,
			SendTime:          dbMsg.SendTime,
			MsgSeqNo:          dbMsg.MsgSeqNo,
			MsgBody:           dbMsg.MsgBody,
			AppKey:            dbMsg.AppKey,
			IsExt:             dbMsg.IsExt,
			IsExset:           dbMsg.IsExset,
			MsgExt:            dbMsg.MsgExt,
			MsgExset:          dbMsg.MsgExset,
			IsDelete:          dbMsg.IsDelete,
			DestroyTime:       dbMsg.DestroyTime,
			LifeTimeAfterRead: dbMsg.LifeTimeAfterRead,
		},
		IsRead:   dbMsg.IsRead,
		ReadTime: dbMsg.ReadTime,
	}
}
