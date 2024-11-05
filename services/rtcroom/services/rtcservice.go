package services

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"

	"github.com/zegoim/zego_server_assistant/token/go/src/token04"
)

func GenerateAuth(appkey, userId string, rtcChannel pbobjs.RtcChannel) (errs.IMErrorCode, *pbobjs.RtcAuth) {
	appInfo, exist := commonservices.GetAppInfo(appkey)
	if !exist {
		return errs.IMErrorCode_RTCROOM_RTCAUTHFAILED, nil
	}
	if rtcChannel == pbobjs.RtcChannel_Zego {
		if appInfo.ZegoConfigObj == nil {
			return errs.IMErrorCode_RTCROOM_RTCAUTHFAILED, nil
		}
		var appId uint32 = uint32(appInfo.ZegoConfigObj.AppId)
		serverSecret := appInfo.ZegoConfigObj.Secret
		var effectiveTimeInSeconds int64 = 3600 // token 的有效时长，单位：秒
		var payload string = ""                 // token业务认证扩展，基础鉴权token此处填空
		//生成token
		token, err := token04.GenerateToken04(appId, userId, serverSecret, effectiveTimeInSeconds, payload)
		if err != nil {
			return errs.IMErrorCode_RTCROOM_RTCAUTHFAILED, nil
		}
		return errs.IMErrorCode_SUCCESS, &pbobjs.RtcAuth{
			ZegoAuth: &pbobjs.ZegoAuth{
				Token: token,
			},
		}
	} else {
		return errs.IMErrorCode_RTCROOM_RTCAUTHFAILED, nil
	}
}
