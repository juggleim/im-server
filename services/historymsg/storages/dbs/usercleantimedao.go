package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
)

type HisMsgUserCleanTimeDao struct {
	ID          int64  `gorm:"primary_key"`
	UserId      string `gorm:"user_id"`
	TargetId    string `gorm:"target_id"`
	ChannelType int    `gorm:"channel_type"`
	CleanTime   int64  `gorm:"clean_time"`
	AppKey      string `gorm:"app_key"`
}

func (msg HisMsgUserCleanTimeDao) TableName() string {
	return "usercleantimes"
}

func (msg HisMsgUserCleanTimeDao) UpsertCleanTime(item models.HisMsgUserCleanTime) error {
	return dbcommons.GetDb().Exec("INSERT INTO usercleantimes (app_key,user_id,target_id,channel_type,clean_time)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE clean_time=?",
		item.AppKey, item.UserId, item.TargetId, item.ChannelType, item.CleanTime, item.CleanTime).Error
}

func (msg HisMsgUserCleanTimeDao) FindOne(appkey, userId, targetId string, channelType pbobjs.ChannelType) (*models.HisMsgUserCleanTime, error) {
	var item HisMsgUserCleanTimeDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=?", appkey, userId, targetId, channelType).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.HisMsgUserCleanTime{
		UserId:      item.UserId,
		TargetId:    item.TargetId,
		ChannelType: pbobjs.ChannelType(item.ChannelType),
		CleanTime:   item.CleanTime,
		AppKey:      item.AppKey,
	}, nil
}
