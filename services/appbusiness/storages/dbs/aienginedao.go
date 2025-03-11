package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/services/appbusiness/storages/models"
)

type AiEngineDao struct {
	ID         int64  `gorm:"primary_key"`
	EngineType int    `gorm:"engine_type"`
	EngineConf string `gorm:"engine_conf"`
	Status     int    `gorm:"status"`
	AppKey     string `gorm:"app_key"`
}

func (eng AiEngineDao) TableName() string {
	return "ai_engines"
}

func (eng AiEngineDao) Create(item models.AiEngine) error {
	return dbcommons.GetDb().Create(&AiEngineDao{
		EngineType: int(item.EngineType),
		EngineConf: item.EngineConf,
		Status:     item.Status,
		AppKey:     item.AppKey,
	}).Error
}

func (eng AiEngineDao) FindById(appkey string, id int64) (*models.AiEngine, error) {
	var item AiEngineDao
	err := dbcommons.GetDb().Where("app_key=? and id=?", appkey, id).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.AiEngine{
		ID:         item.ID,
		EngineType: models.EngineType(item.EngineType),
		EngineConf: item.EngineConf,
		Status:     item.Status,
		AppKey:     item.AppKey,
	}, nil
}

func (eng AiEngineDao) FindEnableAiEngine(appkey string) (*models.AiEngine, error) {
	var item AiEngineDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.AiEngine{
		ID:         item.ID,
		EngineType: models.EngineType(item.EngineType),
		EngineConf: item.EngineConf,
		Status:     item.Status,
		AppKey:     item.AppKey,
	}, nil
}
