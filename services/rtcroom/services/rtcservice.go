package services

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"

	"github.com/zegoim/zego_server_assistant/token/go/src/token04"
)

func GenerateAuth(appkey, userId string, rtcChannel pbobjs.RtcChannel) (errs.IMErrorCode, *pbobjs.RtcAuth) {
	var appId uint32 = 1881186044
	serverSecret := "66c19f780de15fb0d4d5fa35f5f77ab2"
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
}
