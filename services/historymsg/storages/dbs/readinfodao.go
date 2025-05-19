package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"time"
)

type ReadInfoDao struct {
	ID          int64     `gorm:"primary_key"`
	AppKey      string    `gorm:"app_key"`
	MsgId       string    `gorm:"msg_id"`
	ChannelType int       `gorm:"channel_type"`
	GroupId     string    `gorm:"group_id"`
	MemberId    string    `gorm:"member_id"`
	CreatedTime time.Time `gorm:"created_time"`
}

func (info ReadInfoDao) TableName() string {
	return "readinfos"
}

func (info ReadInfoDao) Create(item models.ReadInfo) error {
	var addTime time.Time
	if item.CreatedTime > 0 {
		addTime = time.UnixMilli(item.CreatedTime)
	} else {
		addTime = time.Now()
	}
	add := ReadInfoDao{
		AppKey:      item.AppKey,
		MsgId:       item.MsgId,
		ChannelType: int(item.ChannelType),
		GroupId:     item.GroupId,
		MemberId:    item.MemberId,
		CreatedTime: addTime,
	}
	err := dbcommons.GetDb().Create(&add).Error
	return err
}

func (info ReadInfoDao) BatchCreate(items []models.ReadInfo) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("insert into %s (`app_key`,`msg_id`,`channel_type`,`group_id`,`member_id`)values", info.TableName())
	params := []interface{}{}

	buffer.WriteString(sql)
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?,?,?);")
		} else {
			buffer.WriteString("(?,?,?,?,?),")
		}
		params = append(params, item.AppKey, item.MsgId, item.ChannelType, item.GroupId, item.MemberId)
	}
	err := dbcommons.GetDb().Exec(buffer.String(), params...).Error
	return err
}

func (info ReadInfoDao) QryReadInfosByMsgId(appkey, groupId string, channelType pbobjs.ChannelType, msgId string, startId, limit int64) ([]*models.ReadInfo, error) {
	var items []*ReadInfoDao
	err := dbcommons.GetDb().Where("app_key=? and channel_type=? and group_id=? and msg_id=? and id>?", appkey, channelType, groupId, msgId, startId).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	retItems := []*models.ReadInfo{}
	for _, item := range items {
		retItems = append(retItems, &models.ReadInfo{
			AppKey:      item.AppKey,
			MsgId:       item.MsgId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			GroupId:     item.GroupId,
			MemberId:    item.MemberId,
			CreatedTime: item.CreatedTime.UnixMilli(),
		})
	}
	return retItems, err
}

func (info ReadInfoDao) CountReadInfosByMsgId(appkey, groupId string, channelType pbobjs.ChannelType, msgId string) int32 {
	var count int32
	err := dbcommons.GetDb().Model(&ReadInfoDao{}).Where("app_key=? and channel_type=? and group_id=? and msg_id=?", appkey, channelType, groupId, msgId).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}

func (info ReadInfoDao) CheckMsgsRead(appkey, groupId, memberId string, channelType pbobjs.ChannelType, msgIds []string) (map[string]bool, error) {
	ret := map[string]bool{}
	for _, msgId := range msgIds {
		ret[msgId] = false
	}
	var items []*ReadInfoDao
	err := dbcommons.GetDb().Where("app_key=? and channel_type=? and group_id=? and member_id=? and msg_id in (?)", appkey, channelType, groupId, memberId, msgIds).Find(&items).Error
	if err != nil {
		return ret, err
	}
	for _, item := range items {
		if _, exist := ret[item.MsgId]; exist {
			ret[item.MsgId] = true
		}
	}
	return ret, nil
}
