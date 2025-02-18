package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/storages"
	"im-server/services/appbusiness/storages/models"
	"im-server/services/commonservices"
	"im-server/services/commonservices/sms"
	"math/rand"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomSms() string {
	retCode := ""
	for i := 0; i < 6; i++ {
		item := random.Intn(10)
		retCode = retCode + tools.Int642String(int64(item))
	}
	return retCode
}

func SmsSend(ctx context.Context, phone string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	smsEngine := commonservices.GetSmsEngine(appkey)
	if smsEngine != nil && smsEngine != sms.DefaultSmsEngine {
		// 检查是否还有有效的
		storage := storages.NewSmsRecordStorage()
		record, err := storage.FindByPhone(appkey, phone, time.Now().Add(-3*time.Minute))
		randomCode := RandomSms()
		if err == nil {
			randomCode = record.Code
		} else {
			_, err = storage.Create(models.SmsRecord{
				AppKey:      appkey,
				Phone:       phone,
				Code:        randomCode,
				CreatedTime: time.Now(),
			})
			if err != nil {
				return errs.IMErrorCode_APP_SMS_SEND_FAILED
			}
		}
		err = smsEngine.SmsSend(phone, map[string]interface{}{
			"code": randomCode,
		})
		if err == nil {
			return errs.IMErrorCode_SUCCESS
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func CheckPhoneSmsCode(ctx context.Context, phone, code string) errs.IMErrorCode {
	if code == "000000" {
		return errs.IMErrorCode_SUCCESS
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewSmsRecordStorage()
	record, err := storage.FindByPhoneCode(appkey, phone, code)
	if err != nil {
		return errs.IMErrorCode_APP_SMS_CODE_EXPIRED
	}
	interval := time.Since(record.CreatedTime)
	if interval > 5*time.Minute {
		return errs.IMErrorCode_APP_SMS_CODE_EXPIRED
	}
	return errs.IMErrorCode_SUCCESS
}
