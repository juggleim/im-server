package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"
	"time"
)

type ConverConfDao struct {
	ID          int64     `gorm:"primaryKey;column:id"`
	ConverId    string    `gorm:"column:conver_id"`
	ConverType  int       `gorm:"column:conver_type"`
	SubChannel  string    `gorm:"column:sub_channel"`
	ItemKey     string    `gorm:"column:item_key"`
	ItemValue   string    `gorm:"column:item_value"`
	ItemType    int       `gorm:"column:item_type"`
	CreatedTime time.Time `gorm:"column:created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time"`
	AppKey      string    `gorm:"column:app_key"`
}

func (conf ConverConfDao) TableName() string {
	return "converconfs"
}

func (conf ConverConfDao) Upsert(item models.ConverConf) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,conver_id,conver_type,sub_channel,item_key,item_value,item_type) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=VALUES(item_value),item_type=VALUES(item_type)", conf.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.ConverId, item.ConverType, item.SubChannel, item.ItemKey, item.ItemValue, item.ItemType).Error
}

func (conf ConverConfDao) BatchUpsert(items []models.ConverConf) error {
	if len(items) == 0 {
		return nil
	}

	var sqlBuilder bytes.Buffer
	sqlBuilder.WriteString(fmt.Sprintf("INSERT INTO %s (app_key,conver_id,conver_type,sub_channel,item_key,item_value,item_type) VALUES ", conf.TableName()))
	params := make([]interface{}, 0, len(items)*7)
	for i, item := range items {
		if i > 0 {
			sqlBuilder.WriteByte(',')
		}
		sqlBuilder.WriteString("(?,?,?,?,?,?,?)")
		params = append(params, item.AppKey, item.ConverId, item.ConverType, item.SubChannel, item.ItemKey, item.ItemValue, item.ItemType)
	}
	sqlBuilder.WriteString(" ON DUPLICATE KEY UPDATE item_value=VALUES(item_value),item_type=VALUES(item_type)")
	return dbcommons.GetDb().Exec(sqlBuilder.String(), params...).Error
}

func (conf ConverConfDao) Delete(appkey, converId string, converType pbobjs.ChannelType, subChannel, itemKey string) error {
	return dbcommons.GetDb().Where(
		"app_key=? and conver_id=? and conver_type=? and sub_channel=? and item_key=?",
		appkey, converId, converType, subChannel, itemKey,
	).Delete(&ConverConfDao{}).Error
}

func (conf ConverConfDao) Find(appkey, converId string, converType pbobjs.ChannelType, subChannel, itemKey string) (*models.ConverConf, error) {
	item := &ConverConfDao{}
	err := dbcommons.GetDb().Where(
		"app_key=? and conver_id=? and conver_type=? and sub_channel=? and item_key=?",
		appkey, converId, converType, subChannel, itemKey,
	).Take(item).Error
	if err != nil {
		return nil, err
	}
	return item.toModel(), nil
}

func (conf ConverConfDao) QryConfs(appkey, converId string, converType pbobjs.ChannelType, subChannel string) ([]*models.ConverConf, error) {
	var items []*ConverConfDao
	err := dbcommons.GetDb().Where(
		"app_key=? and conver_id=? and conver_type=? and sub_channel=?",
		appkey, converId, converType, subChannel,
	).Find(&items).Error
	return toConverConfModels(items), err
}

func (conf ConverConfDao) QryConfsByItemKeys(appkey, converId string, converType pbobjs.ChannelType, subChannel string, itemKeys []string) (map[string]*models.ConverConf, error) {
	ret := make(map[string]*models.ConverConf, len(itemKeys))
	if len(itemKeys) == 0 {
		return ret, nil
	}

	var items []*ConverConfDao
	err := dbcommons.GetDb().Where(
		"app_key=? and conver_id=? and conver_type=? and sub_channel=? and item_key in (?)",
		appkey, converId, converType, subChannel, itemKeys,
	).Find(&items).Error
	if err != nil {
		return ret, err
	}
	for _, item := range items {
		ret[item.ItemKey] = item.toModel()
	}
	return ret, nil
}

func (conf *ConverConfDao) toModel() *models.ConverConf {
	return &models.ConverConf{
		ID:          conf.ID,
		ConverId:    conf.ConverId,
		ConverType:  pbobjs.ChannelType(conf.ConverType),
		SubChannel:  conf.SubChannel,
		ItemKey:     conf.ItemKey,
		ItemValue:   conf.ItemValue,
		ItemType:    conf.ItemType,
		CreatedTime: conf.CreatedTime,
		UpdatedTime: conf.UpdatedTime,
		AppKey:      conf.AppKey,
	}
}

func toConverConfModels(items []*ConverConfDao) []*models.ConverConf {
	ret := make([]*models.ConverConf, 0, len(items))
	for _, item := range items {
		ret = append(ret, item.toModel())
	}
	return ret
}
