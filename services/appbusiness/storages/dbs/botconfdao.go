package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/appbusiness/storages/models"
	"im-server/services/commonservices"
)

type BotConfDao struct {
	ID          int64  `gorm:"primary_key"`
	BotId       string `gorm:"bot_id"`
	Nickname    string `gorm:"nickname"`
	BotPortrait string `gorm:"bot_portrait"`
	Description string `gorm:"description"`
	BotType     int    `gorm:"bot_type"`
	BotConf     string `gorm:"bot_conf"`
	Status      int    `gorm:"status"`
	AppKey      string `gorm:"app_key"`
}

func (conf BotConfDao) TableName() string {
	return "botconfs"
}

func (conf BotConfDao) Upsert(item models.BotConf) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,bot_id,nickname,bot_portrait,description,bot_type,bot_conf)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE nickname=VALUES(nickname), bot_portrait=VALUES(bot_portrait), description=VALUES(descriptiono), bot_type=VALUES(bot_type), bot_conf=VALUES(bot_conf)", conf.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.BotId, item.Nickname, item.BotPortrait, item.Description, item.BotType, item.BotConf).Error
}

func (conf BotConfDao) FindById(appkey, botId string) (*models.BotConf, error) {
	var item BotConfDao
	err := dbcommons.GetDb().Where("app_key=? and bot_id=?", appkey, botId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.BotConf{
		AppKey:      item.AppKey,
		BotId:       item.BotId,
		Nickname:    item.Nickname,
		BotPortrait: item.BotPortrait,
		Description: item.Description,
		BotType:     commonservices.BotType(item.BotType),
		BotConf:     item.BotConf,
		Status:      models.BotStatus(item.Status),
	}, err
}

func (conf BotConfDao) QryBotConfs(appkey string, startId, limit int64) ([]*models.BotConf, error) {
	var items []*BotConfDao
	err := dbcommons.GetDb().Where("app_key=? and id>?", appkey, startId).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.BotConf{}
	for _, item := range items {
		ret = append(ret, &models.BotConf{
			ID:          item.ID,
			AppKey:      item.AppKey,
			BotId:       item.BotId,
			Nickname:    item.Nickname,
			BotPortrait: item.BotPortrait,
			Description: item.Description,
			BotType:     commonservices.BotType(item.BotType),
			BotConf:     item.BotConf,
			Status:      models.BotStatus(item.Status),
		})
	}
	return ret, nil
}

func (conf BotConfDao) QryBotConfsWithStatus(appkey string, status models.BotStatus, startId, limit int64) ([]*models.BotConf, error) {
	var items []*BotConfDao
	err := dbcommons.GetDb().Where("app_key=? and id>? and status=?", appkey, startId, status).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.BotConf{}
	for _, item := range items {
		ret = append(ret, &models.BotConf{
			ID:          item.ID,
			AppKey:      item.AppKey,
			BotId:       item.BotId,
			Nickname:    item.Nickname,
			BotPortrait: item.BotPortrait,
			Description: item.Description,
			BotType:     commonservices.BotType(item.BotType),
			BotConf:     item.BotConf,
			Status:      models.BotStatus(item.Status),
		})
	}
	return ret, nil
}

func (conf BotConfDao) UpdateStatus(appkey, botId string, status models.BotStatus) error {
	return dbcommons.GetDb().Model(&BotConfDao{}).Where("app_key=? and bot_id=?", appkey, botId).Update("status", status).Error
}
