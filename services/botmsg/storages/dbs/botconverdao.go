package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/botmsg/storages/models"
	"time"
)

type BotConverDao struct {
	ID          int64     `gorm:"primary_key"`
	AppKey      string    `gorm:"app_key"`
	ConverType  int       `gorm:"conver_type"`
	ConverKey   string    `gorm:"conver_key"`
	ConverId    string    `gorm:"conver_id"`
	UpdatedTime time.Time `gorm:"updated_time"`
}

func (ext BotConverDao) TableName() string {
	return "botconvers"
}

func (ext BotConverDao) Upsert(item models.BotConver) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,conver_type,conver_key,conver_id)VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE conver_id=VALUES(conver_id)", ext.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.ConverType, item.ConverKey, item.ConverId).Error
}

func (ext BotConverDao) Find(appkey string, converType models.BotConverType, converKey string) (*models.BotConver, error) {
	var item BotConverDao
	err := dbcommons.GetDb().Where("app_key=? and conver_type=? and conver_key=?", appkey, converType, converKey).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.BotConver{
		AppKey:      item.AppKey,
		ConverType:  models.BotConverType(item.ConverType),
		ConverKey:   item.ConverKey,
		ConverId:    item.ConverId,
		UpdatedTime: item.UpdatedTime.UnixMilli(),
	}, nil
}
