package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/services/appbusiness/storages/models"
	"time"
)

type PromptDao struct {
	ID          int64     `gorm:"primary_key"`
	UserId      string    `gorm:"user_id"`
	Prompts     string    `gorm:"prompts"`
	CreatedTime time.Time `gorm:"created_time"`
	AppKey      string    `gorm:"app_key"`
}

func (pro PromptDao) TableName() string {
	return "assistant_prompts"
}

func (pro PromptDao) Create(prompt models.Prompt) error {
	return dbcommons.GetDb().Create(&PromptDao{
		UserId:      prompt.UserId,
		Prompts:     prompt.Prompts,
		CreatedTime: time.Now(),
		AppKey:      prompt.AppKey,
	}).Error
}

func (pro PromptDao) UpdatePrompts(appkey, userId string, id int64, prompts string) error {
	return dbcommons.GetDb().Model(&PromptDao{}).Where("app_key=? and id=? and user_id=?", appkey, id, userId).Update("prompts", prompts).Error
}

func (pro PromptDao) DelPrompts(appkey, userId string, id int64) error {
	return dbcommons.GetDb().Where("app_key=? and id=? and user_id=?", appkey, id, userId).Delete(&PromptDao{}).Error
}

func (pro PromptDao) BatchDelPrompts(appkey, userId string, ids []int64) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and id in (?)", appkey, userId, ids).Delete(&PromptDao{}).Error
}

func (pro PromptDao) FindPrompt(appkey, userId string, id int64) (*models.Prompt, error) {
	var item PromptDao
	err := dbcommons.GetDb().Where("app_key=? and id=? and user_id=?", appkey, id, userId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.Prompt{
		ID:          item.ID,
		UserId:      item.UserId,
		Prompts:     item.Prompts,
		CreatedTime: item.CreatedTime.UnixMilli(),
		AppKey:      item.AppKey,
	}, nil
}

func (pro PromptDao) QryPrompts(appkey, userId string, limit int64, startId int64) ([]*models.Prompt, error) {
	var items []*PromptDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and id<?", appkey, userId, startId).Order("id desc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.Prompt{}
	for _, item := range items {
		ret = append(ret, &models.Prompt{
			ID:          item.ID,
			UserId:      item.UserId,
			Prompts:     item.Prompts,
			CreatedTime: item.CreatedTime.UnixMilli(),
			AppKey:      item.AppKey,
		})
	}
	return ret, nil
}
