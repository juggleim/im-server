package models

import "time"

type SmsRecord struct {
	Phone       string
	Code        string
	CreatedTime time.Time
	AppKey      string
}

type ISmsRecordStorage interface {
	Create(s SmsRecord) (int64, error)
	FindByPhoneCode(appkey, phone, code string) (*SmsRecord, error)
	FindByPhone(appkey, phone string, startTime time.Time) (*SmsRecord, error)
}
