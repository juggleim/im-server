package dbs

import (
	"bytes"
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/friends/storages/models"
)

type FriendRelDao struct {
	ID       int64  `gorm:"primary_key"`
	UserId   string `gorm:"user_id"`
	FriendId string `gorm:"friend_id"`
	OrderTag string `gorm:"order_tag"`
	AppKey   string `gorm:"app_key"`
}

func (rel FriendRelDao) TableName() string {
	return "friendrels"
}

func (rel FriendRelDao) Upsert(item models.FriendRel) error {
	sql := fmt.Sprintf("INSERT IGNORE INTO %s (app_key,user_id,friend_id,order_tag)VALUES(?,?,?,?)", rel.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.UserId, item.FriendId, item.OrderTag).Error
}

func (rel FriendRelDao) BatchUpsert(items []models.FriendRel) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT IGNORE INTO %s (app_key,user_id,friend_id,order_tag)VALUES", rel.TableName())
	buffer.WriteString(sql)
	length := len(items)
	params := []interface{}{}
	for i, item := range items {
		if i == length-1 {
			buffer.WriteString("(?,?,?,?)")
		} else {
			buffer.WriteString("(?,?,?,?),")
		}
		params = append(params, item.AppKey, item.UserId, item.FriendId, item.OrderTag)
	}
	return dbcommons.GetDb().Exec(buffer.String(), params...).Error
}

func (rel FriendRelDao) QueryFriendRels(appkey, userId string, startId, limit int64) ([]*models.FriendRel, error) {
	var items []*FriendRelDao
	params := []interface{}{}
	condition := "app_key=?"
	params = append(params, appkey)
	if userId != "" {
		condition = condition + " and user_id=?"
		params = append(params, userId)
	}
	condition = condition + " and id>?"
	params = append(params, startId)
	err := dbcommons.GetDb().Where(condition, params...).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.FriendRel{}
	for _, rel := range items {
		ret = append(ret, &models.FriendRel{
			ID:       rel.ID,
			AppKey:   rel.AppKey,
			UserId:   rel.UserId,
			FriendId: rel.FriendId,
			OrderTag: rel.OrderTag,
		})
	}
	return ret, nil
}

func (rel FriendRelDao) QueryFriendRelsWithPage(appkey, userId string, orderTag string, page, size int64) ([]*models.FriendRel, error) {
	var items []*FriendRelDao
	params := []interface{}{}
	condition := "app_key=? and user_id=?"
	params = append(params, appkey, userId)
	if orderTag != "" {
		condition = condition + " and order_tag>=?"
		params = append(params, orderTag)
	}
	err := dbcommons.GetDb().Where(condition, params...).Order("order_tag asc").Offset((page - 1) * size).Limit(size).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.FriendRel{}
	for _, rel := range items {
		ret = append(ret, &models.FriendRel{
			ID:       rel.ID,
			AppKey:   rel.AppKey,
			UserId:   rel.UserId,
			FriendId: rel.FriendId,
			OrderTag: rel.OrderTag,
		})
	}
	return ret, nil
}

func (rel FriendRelDao) BatchDelete(appkey, userId string, friendIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and friend_id in (?)", appkey, userId, friendIds).Delete(&FriendRelDao{}).Error
}

func (rel FriendRelDao) QueryFriendRelsByFriendIds(appkey, userId string, friendIds []string) ([]*models.FriendRel, error) {
	var items []*FriendRelDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and friend_id in (?)", appkey, userId, friendIds).Order("id asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.FriendRel{}
	for _, rel := range items {
		ret = append(ret, &models.FriendRel{
			ID:       rel.ID,
			AppKey:   rel.AppKey,
			UserId:   rel.UserId,
			FriendId: rel.FriendId,
			OrderTag: rel.OrderTag,
		})
	}
	return ret, nil
}

func (rel FriendRelDao) UpdateOrderTag(appkey, friendId string, orderTag string) error {
	return dbcommons.GetDb().Model(&FriendRelDao{}).Where("app_key=? and friend_id=?", appkey, friendId).Update("order_tag", orderTag).Error
}
