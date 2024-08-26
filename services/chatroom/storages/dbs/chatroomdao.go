package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/services/chatroom/storages/models"
	"time"
)

type ChatroomDao struct {
	ID          int64     `gorm:"primary_key"`
	ChatId      string    `gorm:"chat_id"`
	ChatName    string    `gorm:"chat_name"`
	IsMute      int       `gorm:"is_mute"`
	CreatedTime time.Time `gorm:"created_time"`
	UpdatedTime time.Time `gorm:"updated_time"`
	AppKey      string    `gorm:"app_key"`
}

func (chat *ChatroomDao) TableName() string {
	return "chatroominfos"
}

func (chat *ChatroomDao) Create(item models.Chatroom) error {
	add := ChatroomDao{
		ChatId:      item.ChatId,
		ChatName:    item.ChatName,
		IsMute:      item.IsMute,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		AppKey:      item.AppKey,
	}
	return dbcommons.GetDb().Create(&add).Error
}

func (chat *ChatroomDao) FindById(appkey, chatId string) (*models.Chatroom, error) {
	var item ChatroomDao
	err := dbcommons.GetDb().Where("app_key=? and chat_id=?", appkey, chatId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.Chatroom{
		ChatId:      item.ChatId,
		ChatName:    item.ChatName,
		IsMute:      item.IsMute,
		CreatedTime: item.CreatedTime,
		UpdatedTime: item.UpdatedTime,
		AppKey:      item.AppKey,
	}, nil
}

func (chat *ChatroomDao) Delete(appkey, chatId string) error {
	return dbcommons.GetDb().Where("app_key=? and chat_id=?", appkey, chatId).Delete(&ChatroomDao{}).Error
}

func (chat *ChatroomDao) UpdateMute(appkey, chatId string, isMute int) error {
	return dbcommons.GetDb().Model(&ChatroomDao{}).Where("app_key=? and chat_id=?", appkey, chatId).Update("is_mute", isMute).Error
}
