package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"time"
)

type FavoriteMsgDao struct {
	ID          int64     `gorm:"primary_key"`
	UserId      string    `gorm:"user_id"`
	SenderId    string    `gorm:"user_id"`
	ReceiverId  string    `gorm:"receiver_id"`
	ChannelType int32     `gorm:"channel_type"`
	MsgId       string    `gorm:"msg_id"`
	MsgTime     int64     `gorm:"msg_time"`
	MsgType     string    `gorm:"msg_type"`
	MsgBody     []byte    `gorm:"msg_body"`
	CreatedTime time.Time `gorm:"created_time"`
	AppKey      string    `gorm:"app_key"`
}

func (msg FavoriteMsgDao) TableName() string {
	return "favoritemsgs"
}

func (msg FavoriteMsgDao) Create(item models.FavoriteMsg) error {
	return dbcommons.GetDb().Create(&FavoriteMsgDao{
		UserId:      item.UserId,
		SenderId:    item.SenderId,
		ReceiverId:  item.ReceiverId,
		ChannelType: int32(item.ChannelType),
		MsgId:       item.MsgId,
		MsgTime:     item.MsgTime,
		MsgType:     item.MsgType,
		MsgBody:     item.MsgBody,
		CreatedTime: time.Now(),
		AppKey:      item.AppKey,
	}).Error
}

func (msg FavoriteMsgDao) QueryFavoriteMsgs(appkey, userId string, startId, limit int64) ([]*models.FavoriteMsg, error) {
	var items []*FavoriteMsgDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and id<?", appkey, userId, startId).Order("id desc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.FavoriteMsg{}
	for _, item := range items {
		ret = append(ret, &models.FavoriteMsg{
			ID:          item.ID,
			UserId:      item.UserId,
			SenderId:    item.SenderId,
			ReceiverId:  item.ReceiverId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			MsgId:       item.MsgId,
			MsgTime:     item.MsgTime,
			MsgType:     item.MsgType,
			MsgBody:     item.MsgBody,
			CreatedTime: item.CreatedTime,
			AppKey:      item.AppKey,
		})
	}
	return ret, nil
}
