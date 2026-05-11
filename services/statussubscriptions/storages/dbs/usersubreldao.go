package dbs

import (
	"time"

	"im-server/commons/dbcommons"
	"im-server/services/statussubscriptions/storages/models"

	"gorm.io/gorm"
)

type UserSubRelDao struct {
	ID                 int64     `gorm:"primary_key"`
	UserId             string    `gorm:"column:user_id"`
	SubscriberId       string    `gorm:"column:subscriber_id"`
	SubscriberDeviceId string    `gorm:"column:subscriber_device_id"`
	CreatedTime        time.Time `gorm:"column:created_time"`
	AppKey             string    `gorm:"column:app_key"`
}

func (UserSubRelDao) TableName() string {
	return "usersubrels"
}

func rowToModel(d *UserSubRelDao) *models.UserSubRel {
	if d == nil {
		return nil
	}
	return &models.UserSubRel{
		ID:                 d.ID,
		UserId:             d.UserId,
		SubscriberId:       d.SubscriberId,
		SubscriberDeviceId: d.SubscriberDeviceId,
		CreatedTime:        d.CreatedTime.UnixMilli(),
		AppKey:             d.AppKey,
	}
}

func (sub UserSubRelDao) BatchCreate(items []*models.UserSubRel) error {
	if len(items) == 0 {
		return nil
	}
	return dbcommons.GetDb().Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if item == nil {
				continue
			}
			// MySQL NO_ZERO_DATE / strict mode rejects GORM zero Time → '0000-00-00'
			createdAt := time.Now()
			if item.CreatedTime > 0 {
				createdAt = time.UnixMilli(item.CreatedTime)
			}
			row := UserSubRelDao{
				AppKey:             item.AppKey,
				UserId:             item.UserId,
				SubscriberId:       item.SubscriberId,
				SubscriberDeviceId: item.SubscriberDeviceId,
				CreatedTime:        createdAt,
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
			item.ID = row.ID
			if !row.CreatedTime.IsZero() {
				item.CreatedTime = row.CreatedTime.UnixMilli()
			}
		}
		return nil
	})
}

func (sub UserSubRelDao) QryBySubscriber(appkey, subscriberId, deviceId string, limit int) ([]*models.UserSubRel, error) {
	var items []*UserSubRelDao
	err := dbcommons.GetDb().
		Where("app_key = ? AND subscriber_id = ? AND subscriber_device_id = ?", appkey, subscriberId, deviceId).
		Order("id desc").
		Limit(limit).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	out := make([]*models.UserSubRel, 0, len(items))
	for _, it := range items {
		out = append(out, rowToModel(it))
	}
	return out, nil
}

func (sub UserSubRelDao) QryByUserID(appkey, targetUserId string, afterID int64, limit int) ([]*models.UserSubRel, error) {
	if limit <= 0 {
		limit = 10000
	}
	var items []*UserSubRelDao
	err := dbcommons.GetDb().
		Where("app_key = ? AND user_id = ? AND id > ?", appkey, targetUserId, afterID).
		Order("id ASC").
		Limit(limit).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	out := make([]*models.UserSubRel, 0, len(items))
	for _, it := range items {
		out = append(out, rowToModel(it))
	}
	return out, nil
}

func (sub UserSubRelDao) Delete(appkey, userId, subscriberId, deviceId string) error {
	return dbcommons.GetDb().
		Where("app_key = ? AND user_id = ? AND subscriber_id = ? AND subscriber_device_id = ?", appkey, userId, subscriberId, deviceId).
		Delete(&UserSubRelDao{}).Error
}

func (sub UserSubRelDao) DeleteByRelIDs(relIds []int64) error {
	if len(relIds) == 0 {
		return nil
	}
	return dbcommons.GetDb().Where("id IN ?", relIds).Delete(&UserSubRelDao{}).Error
}
