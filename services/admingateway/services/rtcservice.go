package services

import (
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/logs"
)

func SetZegoConf(appkey string, req *commonservices.ZegoConfigObj) AdminErrorCode {
	dao := dbs.AppExtDao{}
	err := dao.CreateOrUpdate(appkey, "zego_config", tools.ToJson(req))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return AdminErrorCode_Success
}

func GetZegoConf(appkey string) (AdminErrorCode, *commonservices.ZegoConfigObj) {
	zegoConf := &commonservices.ZegoConfigObj{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "zego_config")
	if err == nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), zegoConf)
	} else {
		logs.NewLogEntity().Error(err.Error())
	}

	return AdminErrorCode_Success, zegoConf
}

func SetAgoraConf(appkey string, req *commonservices.AgoraConfigObj) AdminErrorCode {
	dao := dbs.AppExtDao{}
	err := dao.CreateOrUpdate(appkey, "agora_config", tools.ToJson(req))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return AdminErrorCode_Success
}

func GetAgoraConf(appkey string) (AdminErrorCode, *commonservices.AgoraConfigObj) {
	agroaConf := &commonservices.AgoraConfigObj{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "agora_config")
	if err == nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), agroaConf)
	} else {
		logs.NewLogEntity().Error(err.Error())
	}

	return AdminErrorCode_Success, agroaConf
}

func SetLivekitConf(appkey string, req *commonservices.LivekitConfigObj) AdminErrorCode {
	dao := dbs.AppExtDao{}
	err := dao.CreateOrUpdate(appkey, "livekit_config", tools.ToJson(req))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return AdminErrorCode_Success
}

func GetLivekitConf(appkey string) (AdminErrorCode, *commonservices.LivekitConfigObj) {
	livekitConf := &commonservices.LivekitConfigObj{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "livekit_config")
	if err == nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), livekitConf)
	} else {
		logs.NewLogEntity().Error(err.Error())
	}

	return AdminErrorCode_Success, livekitConf
}
