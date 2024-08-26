package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/storages/models"
	"time"
)

type ChatroomBanUserDao struct {
	ID          int64              `gorm:"primary_key"`
	ChatId      string             `gorm:"chat_id"`
	BanType     pbobjs.ChrmBanType `gorm:"ban_type"`
	MemberId    string             `gorm:"member_id"`
	CreatedTime time.Time          `gorm:"created_time"`
	AppKey      string             `gorm:"app_key"`
}

func (ban *ChatroomBanUserDao) TableName() string {
	return "chrmbanusers"
}

func (ban *ChatroomBanUserDao) Create(item models.ChatroomBanUser) error {
	add := ChatroomBanUserDao{
		ChatId:      item.ChatId,
		MemberId:    item.MemberId,
		BanType:     item.BanType,
		CreatedTime: time.Now(),
		AppKey:      item.AppKey,
	}
	return dbcommons.GetDb().Create(&add).Error
}

func (ban *ChatroomBanUserDao) DelBanUser(appkey, chatId, memberId string, banType pbobjs.ChrmBanType) error {
	return dbcommons.GetDb().Where("app_key=? and chat_id=? and ban_type=? and member_id=?", appkey, chatId, banType, memberId).Delete(&ChatroomBanUserDao{}).Error
}

func (ban *ChatroomBanUserDao) QryBanUsers(appkey, chatId string, banType pbobjs.ChrmBanType, startId, limit int64) ([]*models.ChatroomBanUser, error) {
	var items []*ChatroomBanUserDao
	err := dbcommons.GetDb().Where("app_key=? and chat_id=? and ban_type=? and id>?", appkey, chatId, banType, startId).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	retItems := []*models.ChatroomBanUser{}
	for _, item := range items {
		retItems = append(retItems, &models.ChatroomBanUser{
			ID:          item.ID,
			ChatId:      item.ChatId,
			BanType:     item.BanType,
			MemberId:    item.MemberId,
			AppKey:      item.AppKey,
			CreatedTime: item.CreatedTime.UnixMilli(),
		})
	}
	return retItems, nil
}

func (ban *ChatroomBanUserDao) ClearBanUsers(appkey, chatId string) error {
	return dbcommons.GetDb().Where("app_key=? and chat_id=?", appkey, chatId).Delete(&ChatroomBanUserDao{}).Error
}
