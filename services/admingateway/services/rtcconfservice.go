package services

import (
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/logs"
)

type RtcConfReq struct {
	AppKey string   `json:"app_key"`
	Conf   *RtcConf `json:"conf"`
}
type RtcConf struct {
	ZegoConf    *commonservices.ZegoConfigObj    `json:"zego_conf,omitempty"`
	AgoraConf   *commonservices.AgoraConfigObj   `json:"agora_conf,omitempty"`
	LivekitConf *commonservices.LivekitConfigObj `json:"livekit_conf,omitempty"`
}

func SetRtcConf(appkey string, req *RtcConf) AdminErrorCode {
	dao := dbs.AppExtDao{}
	if req.ZegoConf != nil {
		err := dao.CreateOrUpdate(appkey, "zego_config", tools.ToJson(req.ZegoConf))
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	if req.AgoraConf != nil {
		err := dao.CreateOrUpdate(appkey, "agora_config", tools.ToJson(req.AgoraConf))
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	if req.LivekitConf != nil {
		err := dao.CreateOrUpdate(appkey, "livekit_config", tools.ToJson(req.LivekitConf))
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	return AdminErrorCode_Success
}

func GetRtcConf(appkey string) (AdminErrorCode, *RtcConf) {
	ret := &RtcConf{}
	dao := dbs.AppExtDao{}
	exts, err := dao.FindByItemKeys(appkey, []string{"zego_config,agora_config,livekit_config"})
	if err == nil {
		for _, ext := range exts {
			if ext.AppItemKey == "zego_config" {
				zegoConf := &commonservices.ZegoConfigObj{}
				err = tools.JsonUnMarshal([]byte(ext.AppItemValue), zegoConf)
				if err == nil {
					ret.ZegoConf = zegoConf
				} else {
					logs.NewLogEntity().Error(err.Error())
				}
			} else if ext.AppItemKey == "agora_config" {
				agoraConf := &commonservices.AgoraConfigObj{}
				err = tools.JsonUnMarshal([]byte(ext.AppItemValue), agoraConf)
				if err == nil {
					ret.AgoraConf = agoraConf
				} else {
					logs.NewLogEntity().Error(err.Error())
				}
			} else if ext.AppItemKey == "livekit_config" {
				livekitConf := &commonservices.LivekitConfigObj{}
				err = tools.JsonUnMarshal([]byte(ext.AppItemValue), livekitConf)
				if err == nil {
					ret.LivekitConf = livekitConf
				} else {
					logs.NewLogEntity().Error(err.Error())
				}
			}
		}
	}
	return AdminErrorCode_Success, ret
}
