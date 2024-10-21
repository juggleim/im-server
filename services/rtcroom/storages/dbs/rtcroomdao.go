package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/rtcroom/storages/models"
)

type RtcRoomDao struct {
	ID       int64  `gorm:"primary_key"`
	RoomId   string `gorm:"room_id"`
	RoomType int    `gorm:"room_type"`
	OwnerId  string `gorm:"owner_id"`
	AppKey   string `gorm:"app_key"`
}

func (room *RtcRoomDao) TableName() string {
	return "rtcrooms"
}

func (room *RtcRoomDao) Create(item models.RtcRoom) error {
	add := RtcRoomDao{
		RoomId:   item.RoomId,
		RoomType: int(item.RoomType),
		OwnerId:  item.OwnerId,
		AppKey:   item.AppKey,
	}
	return dbcommons.GetDb().Create(&add).Error
}

func (room *RtcRoomDao) FindById(appkey, roomId string) (*models.RtcRoom, error) {
	var item RtcRoomDao
	err := dbcommons.GetDb().Where("app_key=? and room_id=?", appkey, roomId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.RtcRoom{
		RoomId:   item.RoomId,
		RoomType: pbobjs.RtcRoomType(item.RoomType),
		OwnerId:  item.OwnerId,
		AppKey:   item.AppKey,
	}, nil
}

func (room *RtcRoomDao) Delete(appkey, roomId string) error {
	return dbcommons.GetDb().Where("app_key=? and room_id=?", appkey, roomId).Delete(&RtcRoomDao{}).Error
}
