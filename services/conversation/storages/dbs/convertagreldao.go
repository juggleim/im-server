package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"
	"time"
)

type ConverTagRelDao struct {
	ID          int64     `gorm:"primary_key"`
	UserId      string    `gorm:"user_id"`
	Tag         string    `gorm:"tag"`
	TargetId    string    `gorm:"target_id"`
	ChannelType int       `gorm:"channel_type"`
	SubChannel  string    `gorm:"sub_channel"`
	CreatedTime time.Time `gorm:"created_time"`
	AppKey      string    `gorm:"app_key"`
}

func (utag *ConverTagRelDao) TableName() string {
	return "convertagrels"
}

func (utag *ConverTagRelDao) Create(item models.ConverTagRel) error {
	add := ConverTagRelDao{
		UserId:      item.UserId,
		Tag:         item.Tag,
		TargetId:    item.TargetId,
		ChannelType: int(item.ChannelType),
		SubChannel:  item.SubChannel,
		AppKey:      item.AppKey,
	}
	return dbcommons.GetDb().Create(&add).Error
}

func (utag *ConverTagRelDao) BatchCreate(items []models.ConverTagRel) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT IGNORE INTO %s (`user_id`,`tag`,`target_id`,`channel_type`,`sub_channel`,`app_key`) VALUES ", utag.TableName())
	buffer.WriteString(sql)
	params := []interface{}{}
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?,?,?,?);")
		} else {
			buffer.WriteString("(?,?,?,?,?,?),")
		}
		params = append(params, item.UserId, item.Tag, item.TargetId, item.ChannelType, item.SubChannel, item.AppKey)
	}
	return dbcommons.GetDb().Exec(buffer.String(), params...).Error
}

func (utag *ConverTagRelDao) Delete(appkey, userId, tag, targetId, subChannel string, channelType pbobjs.ChannelType) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and tag=? and target_id=? and channel_type=? and sub_channel=?", appkey, userId, tag, targetId, channelType, subChannel).Delete(&ConverTagRelDao{}).Error
}

func (utag *ConverTagRelDao) BatchDelete(appkey, userId, tag string, convers []models.TargetConver) error {
	if len(convers) <= 0 {
		return nil
	}
	condition := "app_key=? and user_id=? and tag=?"
	params := []interface{}{}
	params = append(params, appkey)
	params = append(params, userId)
	params = append(params, tag)
	condition = condition + " and ("
	for i, conver := range convers {
		if i == len(convers)-1 {
			condition = condition + "(target_id=? and channel_type=? and sub_channel=?)"
		} else {
			condition = condition + "(target_id=? and channel_type=? and sub_channel=?) or "
		}
		params = append(params, conver.TargetId)
		params = append(params, conver.ChannelType)
		params = append(params, conver.SubChannel)
	}
	condition = condition + ")"
	return dbcommons.GetDb().Where(condition, params...).Delete(&ConverTagRelDao{}).Error
}

func (utag *ConverTagRelDao) DeleteByTag(appkey, userId, tag string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and tag=?").Delete(&ConverTagRelDao{}).Error
}
