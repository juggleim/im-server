package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"
	"time"

	"github.com/juggleim/commons/dbcommons"
)

type ConverConfDao struct {
	ID          int64     `gorm:"primary_key"`
	ConverId    string    `gorm:"conver_id"`
	ConverType  int32     `gorm:"conver_type"`
	SubChannel  string    `gorm:"sub_channel"`
	ItemKey     string    `gorm:"item_key"`
	ItemValue   string    `gorm:"item_value"`
	ItemType    int32     `gorm:"item_type"`
	UpdatedTime time.Time `gorm:"updated_time"`
	CreatedTime time.Time `gorm:"created_time"`
	AppKey      string    `gorm:"app_key"`
}

func (conf ConverConfDao) TableName() string {
	return "converconfs"
}

func (conf ConverConfDao) Upsert(item models.ConverConf) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,conver_id,conver_type,sub_channel,item_key,item_value,item_type)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=VALUES(item_value)", conf.TableName()), item.AppKey, item.ConverId, item.ConverType, item.SubChannel, item.ItemKey, item.ItemValue, item.ItemType).Error
}

func (conf ConverConfDao) BatchUpsert(items []models.ConverConf) error {
	if len(items) <= 0 {
		return nil
	}
	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT INTO %s (app_key,conver_id,conver_type,sub_channel,item_key,item_value,item_type)VALUES", conf.TableName())
	buffer.WriteString(sql)
	params := []interface{}{}
	for i, item := range items {
		if i == len(items)-1 {
			buffer.WriteString("(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=VALUES(item_value)")
		} else {
			buffer.WriteString("(?,?,?,?,?,?,?),")
		}
		params = append(params, item.AppKey, item.ConverId, item.ConverType, item.SubChannel, item.ItemKey, item.ItemValue, item.ItemType)
	}
	return dbcommons.GetDb().Exec(buffer.String(), params...).Error
}

func (conf ConverConfDao) Delete(appkey, converId string, channelType pbobjs.ChannelType, subChannel string, itemKey string) error {
	return dbcommons.GetDb().Where("app_key=? and conver_id=? and conver_type=? and sub_channel=? and item_key=?", appkey, converId, channelType, subChannel, itemKey).Delete(&ConverConfDao{}).Error
}

func (conf ConverConfDao) UpdateTime(appkey, converId string, channelType pbobjs.ChannelType, subChannel, itemKey string, t time.Time) error {
	return dbcommons.GetDb().Model(&ConverConfDao{}).Where("app_key=? and conver_id=? and conver_type=? and sub_channel=? and item_key=?", appkey, converId, channelType, subChannel, itemKey).Update("updated_time", t).Error
}

func (conf ConverConfDao) QryConverConfs(appkey, converId, subChannel string, converType int32) (map[string]*models.ConverConf, error) {
	var items []*ConverConfDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=? and conver_type=? and sub_channel=?", appkey, converId, converType, subChannel).Find(&items).Error
	ret := map[string]*models.ConverConf{}
	for _, item := range items {
		ret[item.ItemKey] = &models.ConverConf{
			ID:          item.ID,
			ConverId:    item.ConverId,
			ConverType:  pbobjs.ChannelType(item.ConverType),
			SubChannel:  item.SubChannel,
			ItemKey:     item.ItemKey,
			ItemValue:   item.ItemValue,
			ItemType:    item.ItemType,
			UpdatedTime: item.UpdatedTime.UnixMilli(),
			CreatedTime: item.CreatedTime.UnixMilli(),
			AppKey:      item.AppKey,
		}
	}
	return ret, err
}
