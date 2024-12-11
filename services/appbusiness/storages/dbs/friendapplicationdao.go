package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/appbusiness/storages/models"
	"time"
)

type FriendApplicationDao struct {
	ID          int64  `gorm:"primary_key"`
	RecipientId string `gorm:"recipient_id"`
	SponsorId   string `gorm:"sponsor_id"`
	ApplyTime   int64  `gorm:"apply_time"`
	Status      int    `gorm:"status"`
	AppKey      string `gorm:"app_key"`
}

func (apply FriendApplicationDao) TableName() string {
	return "friendapplications"
}

func (apply FriendApplicationDao) Upsert(item models.FriendApplication) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,recipient_id,sponsor_id,apply_time,status)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE apply_time=VALUES(apply_time),status=VALUES(status)", apply.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.RecipientId, item.SponsorId, item.ApplyTime, int(item.Status)).Error
}

func (apply FriendApplicationDao) QueryPendingApplications(appkey, recipientId string, startTime, count int64, isPositive bool) ([]*models.FriendApplication, error) {
	var items []*FriendApplicationDao
	params := []interface{}{}
	condition := "app_key=? and recipient_id=?"
	params = append(params, appkey)
	params = append(params, recipientId)
	orderStr := "apply_time desc"
	if isPositive {
		orderStr = "apply_time asc"
		condition = condition + " and apply_time>?"
	} else {
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
		condition = condition + " and apply_time<?"
	}
	params = append(params, startTime)
	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.FriendApplication{}
	for _, app := range items {
		ret = append(ret, &models.FriendApplication{
			RecipientId: app.RecipientId,
			SponsorId:   app.SponsorId,
			ApplyTime:   app.ApplyTime,
			Status:      models.FriendApplicationStatus(app.Status),
			AppKey:      app.AppKey,
		})
	}
	return ret, nil
}

func (apply FriendApplicationDao) QueryMyApplications(appkey, sponsorId string, startTime, count int64, isPositive bool) ([]*models.FriendApplication, error) {
	var items []*FriendApplicationDao
	params := []interface{}{}
	condition := "app_key=? and sponsor_id=?"
	params = append(params, appkey)
	params = append(params, sponsorId)
	orderStr := "apply_time desc"
	if isPositive {
		orderStr = "apply_time asc"
		condition = condition + " and apply_time>?"
	} else {
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
		condition = condition + " and apply_time<?"
	}
	params = append(params, startTime)
	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.FriendApplication{}
	for _, app := range items {
		ret = append(ret, &models.FriendApplication{
			RecipientId: app.RecipientId,
			SponsorId:   app.SponsorId,
			ApplyTime:   app.ApplyTime,
			Status:      models.FriendApplicationStatus(app.Status),
			AppKey:      app.AppKey,
		})
	}
	return ret, nil
}
