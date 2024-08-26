package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/services/chatroom/storages/models"
	"time"
)

type ChatroomMemberDao struct {
	ID          int64     `gorm:"primary_key"`
	ChatId      string    `gorm:"chat_id"`
	MemberId    string    `gorm:"member_id"`
	CreatedTime time.Time `gorm:"created_time"`
	AppKey      string    `gorm:"app_key"`
}

func (member *ChatroomMemberDao) TableName() string {
	return "chatroommembers"
}

func (member *ChatroomMemberDao) Create(item models.ChatroomMember) error {
	add := ChatroomMemberDao{
		ChatId:      item.ChatId,
		MemberId:    item.MemberId,
		CreatedTime: time.Now(),
		AppKey:      item.AppKey,
	}
	return dbcommons.GetDb().Create(&add).Error
}

func (member *ChatroomMemberDao) FindById(appkey, chatId, memberId string) (*models.ChatroomMember, error) {
	var item ChatroomMemberDao
	err := dbcommons.GetDb().Where("app_key=? and chat_id=? and member_id=?", appkey, chatId, memberId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.ChatroomMember{
		ChatId:      item.ChatId,
		MemberId:    item.MemberId,
		CreatedTime: item.CreatedTime,
		AppKey:      item.AppKey,
	}, nil
}

func (member *ChatroomMemberDao) DeleteMember(appkey, chatId, memberId string) error {
	return dbcommons.GetDb().Where("app_key=? and chat_id=? and member_id=?", appkey, chatId, memberId).Delete(&ChatroomMemberDao{}).Error
}

func (member *ChatroomMemberDao) ClearMembers(appkey, chatId string) error {
	return dbcommons.GetDb().Where("app_key=? and chat_id=?", appkey, chatId).Delete(&ChatroomMemberDao{}).Error
}

func (member *ChatroomMemberDao) QryMembers(appkey, chatId string, isPositive bool, startId, limit int64) ([]*models.ChatroomMember, error) {
	params := []interface{}{}
	condition := "app_key=? and chat_id=?"
	params = append(params, appkey)
	params = append(params, chatId)

	orderStr := "id asc"
	if !isPositive {
		orderStr = "id desc"
		condition = condition + " and id<?"
		params = append(params, startId)
	} else {
		condition = condition + " and id>?"
		params = append(params, startId)
	}
	var items []*ChatroomMemberDao
	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	retItems := []*models.ChatroomMember{}
	for _, item := range items {
		retItems = append(retItems, &models.ChatroomMember{
			ID:          item.ID,
			ChatId:      item.ChatId,
			MemberId:    item.MemberId,
			CreatedTime: item.CreatedTime,
			AppKey:      item.AppKey,
		})
	}
	return retItems, nil
}
