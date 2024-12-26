package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/services/appbusiness/storages/models"
)

type QrCodeRecordDao struct {
	ID          int64  `gorm:"primary_key"`
	CodeId      string `gorm:"code_id"`
	Status      int    `gorm:"status"`
	CreatedTime int64  `gorm:"created_time"`
	UserId      string `gorm:"user_id"`
	AppKey      string `gorm:"app_key"`
}

func (record QrCodeRecordDao) TableName() string {
	return "qrcoderecords"
}

func (record QrCodeRecordDao) Create(item models.QrCodeRecord) error {
	return dbcommons.GetDb().Create(&QrCodeRecordDao{
		CodeId:      item.CodeId,
		CreatedTime: item.CreatedTime,
		AppKey:      item.AppKey,
	}).Error
}

func (record QrCodeRecordDao) FindById(appkey, codeId string) (*models.QrCodeRecord, error) {
	var item QrCodeRecordDao
	err := dbcommons.GetDb().Where("app_key=? and code_id=?", appkey, codeId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.QrCodeRecord{
		CodeId:      item.CodeId,
		Status:      models.QrCodeRecordStatus(item.Status),
		CreatedTime: item.CreatedTime,
		UserId:      item.UserId,
		AppKey:      item.AppKey,
	}, nil
}

func (record QrCodeRecordDao) UpdateStatus(appkey, codeId string, status models.QrCodeRecordStatus, userId string) error {
	upd := map[string]interface{}{}
	upd["user_id"] = userId
	upd["status"] = status
	return dbcommons.GetDb().Model(&QrCodeRecordDao{}).Where("app_key=? and code_id=?", appkey, codeId).Update(upd).Error
}
