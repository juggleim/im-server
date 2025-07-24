package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
)

type HisMsgConverCleanTimeDao struct {
	ID          int64  `gorm:"primary_key"`
	ConverId    string `gorm:"conver_id"`
	SubChannel  string `gorm:"sub_channel"`
	ChannelType int    `gorm:"channel_type"`
	CleanTime   int64  `gorm:"clean_time"`
	AppKey      string `gorm:"app_key"`
}

func (msg HisMsgConverCleanTimeDao) TableName() string {
	return "convercleantimes"
}

func (msg HisMsgConverCleanTimeDao) UpsertDestroyTime(item models.HisMsgConverCleanTime) error {
	return dbcommons.GetDb().Exec("INSERT INTO convercleantimes(app_key,conver_id,channel_type,sub_channel,clean_time)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE clean_time=?",
		item.AppKey, item.ConverId, item.ChannelType, item.SubChannel, item.CleanTime, item.CleanTime).Error
}

func (msg HisMsgConverCleanTimeDao) FindOne(appkey, converId, subChannel string, channelType pbobjs.ChannelType) (*models.HisMsgConverCleanTime, error) {
	var item HisMsgConverCleanTimeDao
	err := dbcommons.GetDb().Where("app_key=? and conver_id=?and sub_channel=? and channel_type=?", appkey, converId, subChannel, channelType).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.HisMsgConverCleanTime{
		ConverId:    item.ConverId,
		ChannelType: pbobjs.ChannelType(item.ChannelType),
		SubChannel:  item.SubChannel,
		CleanTime:   item.CleanTime,
		AppKey:      item.AppKey,
	}, nil
}
