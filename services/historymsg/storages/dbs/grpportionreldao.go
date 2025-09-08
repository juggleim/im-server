package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/historymsg/storages/models"
	"sort"
	"time"
)

type GroupPortionRelDao struct {
	ID          int64  `gorm:"primary_key"`
	ConverId    string `gorm:"conver_id"`
	ChannelType int32  `gorm:"channel_type"`
	SubChannel  string `gorm:"sub_channel"`
	UserId      string `gorm:"user_id"`
	MsgId       string `gorm:"msg_id"`
	MsgTime     int64  `gorm:"msg_time"`
	AppKey      string `gorm:"app_key"`
}

func (rel GroupPortionRelDao) TableName() string {
	return "g_portionrels"
}

func (rel GroupPortionRelDao) Upsert(item models.GroupPortionRel) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT IGNORE INTO %s (app_key,conver_id,channel_type,sub_channel,user_id,msg_id,msg_time)VALUES(?,?,?,?,?,?,?) ", rel.TableName()), item.AppKey, item.ConverId, item.ChannelType, item.SubChannel, item.UserId, item.MsgId, item.MsgTime).Error
}

func (rel GroupPortionRelDao) BatchUpsert(items []models.GroupPortionRel) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT IGNORE INTO %s (app_key,conver_id,channel_type,sub_channel,user_id,msg_id,msg_time)VALUES", rel.TableName())
	buffer.WriteString(sql)
	params := []interface{}{}
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?,?,?,?,?)")
		} else {
			buffer.WriteString("(?,?,?,?,?,?,?),")
		}
		params = append(params, item.AppKey, item.ConverId, item.ChannelType, item.SubChannel, item.UserId, item.MsgId, item.MsgTime)
	}
	return dbcommons.GetDb().Exec(buffer.String(), params...).Error
}

func (rel GroupPortionRelDao) Delete(item models.GroupPortionRel) error {
	return dbcommons.GetDb().Where("app_key=? and conver_id=? and channel_type=? and sub_channel=? and user_id=? and msg_id=?", item.AppKey, item.ConverId, item.ChannelType, item.SubChannel, item.UserId, item.MsgId).Error
}

func (rel GroupPortionRelDao) QryPortionMsgs(appkey, userId, converId, subChannel string, startTime int64, count int32, isPositive bool, cleanTime int64) ([]*models.GroupHisMsg, error) {
	var items []*GroupHisMsgDao
	params := []interface{}{}
	sql := fmt.Sprintf("select his.* from %s as portions left join %s as his on his.app_key=portions.app_key and his.conver_id=portions.conver_id and his.sub_channel=portions.sub_channel and his.msg_id=portions.msg_id where portions.app_key=? and portions.conver_id=? and portions.sub_channel=? and portions.user_id=?", rel.TableName(), (&GroupHisMsgDao{}).TableName())
	params = append(params, appkey)
	params = append(params, converId)
	params = append(params, subChannel)
	params = append(params, userId)
	orderStr := "portions.msg_time desc"
	start := startTime
	if isPositive {
		orderStr = "portions.msg_time asc"
		if start < cleanTime {
			start = cleanTime
		}
		sql = sql + " and portions.msg_time>?"
		params = append(params, start)
	} else {
		if start <= 0 {
			start = time.Now().UnixMilli()
		}
		sql = sql + " and portions.msg_time<?"
		params = append(params, start)
		if cleanTime > 0 {
			sql = sql + " and portions.msg_time>?"
			params = append(params, cleanTime)
		}
	}
	err := dbcommons.GetDb().Raw(sql, params...).Order(orderStr).Limit(count).Find(&items).Error
	if !isPositive {
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
