package services

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"time"

	agoraTokenBuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"github.com/livekit/protocol/auth"
	"github.com/zegoim/zego_server_assistant/token/go/src/token04"
)

func GenerateAuth(appkey, userId, roomId string, rtcChannel pbobjs.RtcChannel) (errs.IMErrorCode, *pbobjs.RtcAuth) {
	appInfo, exist := commonservices.GetAppInfo(appkey)
	if !exist {
		return errs.IMErrorCode_CONNECT_APP_NOT_EXISTED, nil
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
	} else if rtcChannel == pbobjs.RtcChannel_Agora {
		if appInfo.AgoraConfigObj == nil {
			return errs.IMErrorCode_RTCROOM_RTCAUTHFAILED, nil
		}
		appId := appInfo.AgoraConfigObj.AppId
		appCertificate := appInfo.AgoraConfigObj.AppCertificate
		tokenExpirationInSeconds := uint32(3600)
		privilegeExpirationInSeconds := uint32(3600)
		// joinChannelPrivilegeExpireInSeconds := uint32(3600)
		// pubAudioPrivilegeExpireInSeconds := uint32(3600)
		// pubVideoPrivilegeExpireInSeconds := uint32(3600)
		// pubDataStreamPrivilegeExpireInSeconds := uint32(3600)
		result, err := agoraTokenBuilder.BuildTokenWithUserAccount(appId, appCertificate, roomId, userId, agoraTokenBuilder.RolePublisher, tokenExpirationInSeconds, privilegeExpirationInSeconds)
		if err != nil {
			return errs.IMErrorCode_RTCROOM_RTCAUTHFAILED, nil
		}
		return errs.IMErrorCode_SUCCESS, &pbobjs.RtcAuth{
			AgoraAuth: &pbobjs.AgoraAuth{
				Token: result,
			},
		}
	} else if rtcChannel == pbobjs.RtcChannel_LivekitRtc {
		if appInfo.LivekitConfigObj == nil {
			return errs.IMErrorCode_RTCROOM_RTCAUTHFAILED, nil
		}
		at := auth.NewAccessToken(appInfo.LivekitConfigObj.AppKey, appInfo.LivekitConfigObj.AppSecret)
		grant := &auth.VideoGrant{
			RoomJoin: true,
			Room:     roomId,
		}
		at.SetVideoGrant(grant).SetIdentity(userId).SetValidFor(time.Hour)
		token, err := at.ToJWT()
		if err != nil {
			return errs.IMErrorCode_RTCROOM_RTCAUTHFAILED, nil
		}
		return errs.IMErrorCode_SUCCESS, &pbobjs.RtcAuth{
			LivekitRtcAuth: &pbobjs.LivekitRtcAuth{
				Token:      token,
				ServiceUrl: appInfo.LivekitConfigObj.ServiceUrl,
			},
		}
	} else {
		return errs.IMErrorCode_RTCROOM_RTCAUTHFAILED, nil
	}
}
