package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/services/botmsg/storages/models"
)

type AssistantDao struct {
	ID          int64  `gorm:"primary_key"`
	AssistantId string `gorm:"assistant_id"`
	OwnerId     string `gorm:"owner_id"`
	Nickname    string `gorm:"nickname"`
	Portrait    string `gorm:"portrait"`
	Description string `gorm:"description"`
	BotType     int    `gorm:"bot_type"`
	BotConf     string `gorm:"bot_conf"`
	Status      int    `gorm:"status"`
	AppKey      string `gorm:"app_key"`
}

func (assis AssistantDao) TableName() string {
	return "assistants"
}

func (assis AssistantDao) Create(item models.Assistant) error {
	return dbcommons.GetDb().Create(&AssistantDao{
		AssistantId: item.AssistantId,
		OwnerId:     item.OwnerId,
		Nickname:    item.Nickname,
		Portrait:    item.Portrait,
		Description: item.Description,
		BotType:     int(item.BotType),
		BotConf:     item.BotConf,
		Status:      item.Status,
		AppKey:      item.AppKey,
	}).Error
}

func (assis AssistantDao) FindByAssistantId(appkey, assistantId string) (*models.Assistant, error) {
	var item AssistantDao
	err := dbcommons.GetDb().Where("app_key=? and assistant_id=?", appkey, assistantId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.Assistant{
		ID:          item.ID,
		AssistantId: item.AssistantId,
		OwnerId:     item.OwnerId,
		Nickname:    item.Nickname,
		Portrait:    item.Portrait,
		Description: item.Description,
		BotType:     models.BotType(item.BotType),
		BotConf:     item.BotConf,
		Status:      item.Status,
		AppKey:      item.AppKey,
	}, nil
}
