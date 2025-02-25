package services

import (
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/dbs"
)

type TranslateConf struct {
	AppKey string                          `json:"app_key"`
	Conf   *commonservices.TransEngineConf `json:"conf"`
}

func SetTranslateConf(appkey string, req *commonservices.TransEngineConf) AdminErrorCode {
	dao := dbs.AppExtDao{}
	dao.CreateOrUpdate(appkey, "trans_engine_conf", tools.ToJson(req))
	return AdminErrorCode_Success
}

func GetTranslateConf(appkey string) (AdminErrorCode, *commonservices.TransEngineConf) {
	transConf := &commonservices.TransEngineConf{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "trans_engine_conf")
	if err == nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), transConf)
	}
	return AdminErrorCode_Success, transConf
}
