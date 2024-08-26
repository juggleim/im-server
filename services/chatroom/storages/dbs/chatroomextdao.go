package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/chatroom/storages/models"
)

type ChatroomExtDao struct {
	ID        int64  `gorm:"primary_key"`
	ChatId    string `gorm:"chat_id"`
	ItemKey   string `gorm:"item_key"`
	ItemValue string `gorm:"item_value"`
	ItemType  int    `gorm:"item_type"`
	ItemTime  int64  `gorm:"item_time"`
	AppKey    string `gorm:"app_key"`
	MemberId  string `gorm:"member_id"`
	IsDelete  int    `gorm:"is_delete"`
}

func (ext ChatroomExtDao) TableName() string {
	return "chatroomexts"
}

func (ext ChatroomExtDao) Upsert(item models.ChatroomExt) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,chat_id,item_key,item_value,item_type,member_id,item_time)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=?,member_id=?,item_time=?", ext.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.ChatId, item.ItemKey, item.ItemValue, item.ItemType, item.MemberId, item.ItemTime, item.ItemValue, item.MemberId, item.ItemTime).Error
}

func (ext ChatroomExtDao) QryExts(appkey, chatId string) ([]*models.ChatroomExt, error) {
	var items []*ChatroomExtDao
	err := dbcommons.GetDb().Where("app_key=? and chat_id=?", appkey, chatId).Find(&items).Error
	if err != nil {
		return nil, err
	}
	retItems := []*models.ChatroomExt{}
	for _, item := range items {
		retItems = append(retItems, &models.ChatroomExt{
			ChatId:    item.ChatId,
			ItemKey:   item.ItemKey,
			ItemValue: item.ItemValue,
			ItemType:  item.ItemType,
			ItemTime:  item.ItemTime,
			AppKey:    item.AppKey,
			MemberId:  item.MemberId,
			IsDelete:  item.IsDelete,
		})
	}
	return retItems, nil
}

func (ext ChatroomExtDao) ClearExts(appkey, chatId string) error {
	return dbcommons.GetDb().Where("app_key=? and chat_id=?", appkey, chatId).Delete(&ChatroomExtDao{}).Error
}

func (ext ChatroomExtDao) DeleteExt(appkey, chatId, key string) error {
	return dbcommons.GetDb().Model(&ChatroomExtDao{}).Where("app_key=? and chat_id=? and item_key=?", appkey, chatId, key).Update("is_delete", 1).Error
}
