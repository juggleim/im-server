package services

import (
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/dbs"
)

type RtcConfReq struct {
	AppKey string   `json:"app_key"`
	Conf   *RtcConf `json:"conf"`
}
type RtcConf struct {
	ZegoConf *commonservices.ZegoConfigObj `json:"zego_conf,omitempty"`
}

func SetRtcConf(appkey string, req *RtcConf) AdminErrorCode {
	dao := dbs.AppExtDao{}
	if req.ZegoConf != nil {
		dao.CreateOrUpdate(appkey, "zego_config", tools.ToJson(req.ZegoConf))
	}
	return AdminErrorCode_Success
}

func GetRtcConf(appkey string) (AdminErrorCode, *RtcConf) {
	dao := dbs.AppExtDao{}
	exts, err := dao.FindByItemKeys(appkey, []string{"zego_config"})
	if err != nil {
		return AdminErrorCode_Default, nil
	}
	ret := &RtcConf{}
	for _, ext := range exts {
		if ext.AppItemKey == "zego_config" {
			zegoConf := &commonservices.ZegoConfigObj{}
			err = tools.JsonUnMarshal([]byte(ext.AppItemValue), zegoConf)
			if err == nil {
				ret.ZegoConf = zegoConf
			}
		}
	}
	return AdminErrorCode_Success, ret
}
