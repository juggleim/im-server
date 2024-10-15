package dbs

import (
	"im-server/commons/dbcommons"
	"time"
)

type ClientLogState int

var (
	ClientLogState_Default      ClientLogState = 0
	ClientLogState_SendOK       ClientLogState = 1
	ClientLogState_SendFail     ClientLogState = 2
	ClientLogState_Uploaded     ClientLogState = 3
	ClientLogState_UploadFailed ClientLogState = 4
	ClientLogState_NoLog        ClientLogState = 5
)

type ClientLogDao struct {
	ID          int64          `gorm:"primary_key"`
	AppKey      string         `gorm:"app_key"`
	UserId      string         `gorm:"user_id"`
	CreatedTime time.Time      `gorm:"created_time"`
	Start       int64          `gorm:"start"`
	End         int64          `gorm:"end"`
	Log         []byte         `gorm:"log"`
	State       ClientLogState `gorm:"state"`
	Platform    string         `gorm:"platform"`
	DeviceId    string         `gorm:"device_id"`
	LogUrl      string         `gorm:"log_url"`
	TraceId     string         `gorm:"trace_id"`
	MsgId       string         `gorm:"msg_id"`
	FailReason  string         `gorm:"fail_reason"`
	Description string         `gorm:"description"`
}

func (log *ClientLogDao) TableName() string {
	return "clientlogs"
}

func (log *ClientLogDao) Create(item ClientLogDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (log *ClientLogDao) FindById(appkey string, id int64) *ClientLogDao {
	var item ClientLogDao
	err := dbcommons.GetDb().Where("id=? and app_key=?", id, appkey).Take(&item).Error
	if err != nil {
		return nil
	}
	return &item
}

func (log *ClientLogDao) QryLogs(appkey, userId string, start, end, startId, limit int64) ([]*ClientLogDao, error) {
	var items []*ClientLogDao

	params := []interface{}{}
	condition := "app_key=?"
	params = append(params, appkey)

	if userId != "" {
		condition = condition + " and user_id=?"
		params = append(params, userId)
	}

	if start > 0 {
		condition = condition + " and created_time>?"
		params = append(params, time.UnixMilli(start))
	}
	if end > 0 {
		condition = condition + " and created_time<?"
		params = append(params, time.UnixMilli(end))
	}

	condition = condition + " and id>?"
	params = append(params, startId)

	err := dbcommons.GetDb().Where(condition, params...).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (log *ClientLogDao) Update(appkey, msgId string, data []byte, state ClientLogState) error {
	upd := map[string]interface{}{}
	upd["log"] = data
	upd["state"] = state
	return dbcommons.GetDb().Model(&ClientLogDao{}).Where("app_key=? and msg_id=?", appkey, msgId).Update(upd).Error
}

func (log *ClientLogDao) UpdateLogUrl(appkey, msgId string, logUrl string, state ClientLogState) error {
	upd := map[string]interface{}{}
	upd["log_url"] = logUrl
	upd["state"] = state
	return dbcommons.GetDb().Model(&ClientLogDao{}).Where("app_key=? and msg_id=? and state!=?", appkey, msgId, ClientLogState_Uploaded).Update(upd).Error
}
