package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/rtcroom/storages/models"
	"time"
)

type RtcRoomDao struct {
	ID           int64     `gorm:"primary_key"`
	RoomId       string    `gorm:"room_id"`
	RoomType     int       `gorm:"room_type"`
	RtcChannel   int       `gorm:"rtc_channel"`
	RtcMediaType int       `gorm:"rtc_media_type"`
	OwnerId      string    `gorm:"owner_id"`
	Ext          string    `gorm:"ext"`
	CreatedTime  time.Time `gorm:"created_time"`
	AcceptedTime int64     `gorm:"accepted_time"`
	AppKey       string    `gorm:"app_key"`
}

func (room *RtcRoomDao) TableName() string {
	return "rtcrooms"
}

func (room *RtcRoomDao) Create(item models.RtcRoom) error {
	add := RtcRoomDao{
		RoomId:       item.RoomId,
		RoomType:     int(item.RoomType),
		RtcChannel:   int(item.RtcChannel),
		RtcMediaType: int(item.RtcMediaType),
		OwnerId:      item.OwnerId,
		Ext:          item.Ext,
		CreatedTime:  time.Now(),
		AcceptedTime: item.AcceptedTime,
		AppKey:       item.AppKey,
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
		RoomId:       item.RoomId,
		RoomType:     pbobjs.RtcRoomType(item.RoomType),
		RtcChannel:   pbobjs.RtcChannel(item.RtcChannel),
		RtcMediaType: pbobjs.RtcMediaType(item.RtcMediaType),
		OwnerId:      item.OwnerId,
		Ext:          item.Ext,
		CreatedTime:  item.CreatedTime.UnixMilli(),
		AcceptedTime: item.AcceptedTime,
		AppKey:       item.AppKey,
	}, nil
}

func (room *RtcRoomDao) UpdateAcceptedTime(appkey, roomId string, acceptedTime int64) error {
	return dbcommons.GetDb().Model(room).Where("app_key=? and room_id=?", appkey, roomId).Update("accepted_time", acceptedTime).Error
}

func (room *RtcRoomDao) Delete(appkey, roomId string) error {
	return dbcommons.GetDb().Where("app_key=? and room_id=?", appkey, roomId).Delete(&RtcRoomDao{}).Error
}
