package services

import (
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/dbs"
)

type SmsConf struct {
	AppKey string                        `json:"app_key"`
	Conf   *commonservices.SmsEngineConf `json:"conf"`
}

func SetSmsConf(appkey string, req *commonservices.SmsEngineConf) AdminErrorCode {
	dao := dbs.AppExtDao{}
	dao.CreateOrUpdate(appkey, "sms_engine_conf", tools.ToJson(req))
	return AdminErrorCode_Success
}

func GetSmsConf(appkey string) (AdminErrorCode, *commonservices.SmsEngineConf) {
	smsConf := &commonservices.SmsEngineConf{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "sms_engine_conf")
	if err == nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), smsConf)
	}
	return AdminErrorCode_Success, smsConf
}
