package services

import (
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
)

func RegPushToken(ctx imcontext.WsHandleContext, appkey, userId, deviceId, platformStr, pushChannelStr, packageName, pushToken string) {
	if deviceId != "" && platformStr != "" && pushChannelStr != "" && packageName != "" && pushToken != "" {
		platform := commonservices.Str2Platform(platformStr)
		pushChannel := commonservices.Str2PushChannel(pushChannelStr)

		if platform == pbobjs.Platform_Android || platform == pbobjs.Platform_iOS {
			req := &pbobjs.RegPushTokenReq{
				DeviceId:    deviceId,
				Platform:    platform,
				PushChannel: pushChannel,
				PushToken:   pushToken,
				PackageName: packageName,
			}
			data, _ := tools.PbMarshal(req)
			bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
				RpcMsgType:   pbobjs.RpcMsgType_QueryAck,
				AppKey:       appkey,
				Session:      imcontext.GetConnSession(ctx),
				Method:       "reg_push_token",
				RequesterId:  userId,
				ReqIndex:     0,
				Qos:          int32(codec.QoS_NoAck),
				AppDataBytes: data,
				TargetId:     userId,
				TerminalNum:  0,
			})
		}
	}
}
