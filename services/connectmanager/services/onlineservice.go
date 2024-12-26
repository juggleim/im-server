package services

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
	"im-server/services/logmanager"
	"time"
)

func RegPushToken(ctx imcontext.WsHandleContext, appkey, userId, deviceId, platformStr, pushChannelStr, packageName, pushToken, voipToken string) {
	if deviceId != "" && platformStr != "" && pushChannelStr != "" && packageName != "" && (pushToken != "" || voipToken != "") {
		platform := commonservices.Str2Platform(platformStr)
		pushChannel := commonservices.Str2PushChannel(pushChannelStr)

		if platform == pbobjs.Platform_Android || platform == pbobjs.Platform_iOS {
			req := &pbobjs.RegPushTokenReq{
				DeviceId:    deviceId,
				Platform:    platform,
				PushChannel: pushChannel,
				PushToken:   pushToken,
				VoipToken:   voipToken,
				PackageName: packageName,
			}
			data, _ := tools.PbMarshal(req)
			bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
				RpcMsgType:   pbobjs.RpcMsgType_QueryMsg,
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

func RemovePushToken(ctx imcontext.WsHandleContext) {
	appkey := imcontext.GetAppkey(ctx)
	userId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	req := &pbobjs.Nil{}
	data, _ := tools.PbMarshal(req)
	bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_QueryMsg,
		AppKey:       appkey,
		Session:      imcontext.GetConnSession(ctx),
		Method:       "remove_push_token",
		RequesterId:  userId,
		ReqIndex:     0,
		Qos:          int32(codec.QoS_NoAck),
		AppDataBytes: data,
		TargetId:     userId,
	})
}

func SetLanguage(ctx imcontext.WsHandleContext, language string) {
	if language == "" {
		return
	}
	appkey := imcontext.GetAppkey(ctx)
	userId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	req := &pbobjs.UserInfo{
		UserId: userId,
		Settings: []*pbobjs.KvItem{
			{
				Key:   string(commonservices.AttItemKey_Language),
				Value: language,
			},
		},
	}
	data, _ := tools.PbMarshal(req)
	bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_QueryMsg,
		AppKey:       appkey,
		Session:      imcontext.GetConnSession(ctx),
		Method:       "set_user_settings",
		RequesterId:  userId,
		TargetId:     userId,
		ReqIndex:     0,
		Qos:          int32(codec.QoS_NoAck),
		AppDataBytes: data,
	})
}

func Online(ctx imcontext.WsHandleContext, ext, language string, isBackend bool) {
	userId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	deviceId := imcontext.GetDeviceId(ctx)
	platform := imcontext.GetPlatform(ctx)
	clientIp := imcontext.GetContextAttrString(ctx, imcontext.StateKey_ClientIp)
	instanceId := imcontext.GetInstanceId(ctx)
	//online subscription
	onlineMsg := &pbobjs.OnlineOfflineMsg{
		Type:          pbobjs.OnlineType_Online,
		UserId:        userId,
		DeviceId:      deviceId,
		Platform:      platform,
		ClientIp:      clientIp,
		SessionId:     imcontext.GetConnSession(ctx),
		Timestamp:     time.Now().UnixMilli(),
		ConnectionExt: ext,
		InstanceId:    instanceId,
	}
	commonservices.SubOnlineEvent(imcontext.GetRpcContext(ctx), userId, onlineMsg)
	//set language for user
	SetLanguage(ctx, language)
	//close push switch
	if platform == string(commonservices.Platform_Android) || platform == string(commonservices.Platform_IOS) {
		var pushSwitch int32 = 0
		if isBackend {
			pushSwitch = 1
		}
		bases.AsyncRpcCall(imcontext.GetRpcContext(ctx), "push_switch", userId, &pbobjs.PushSwitch{
			Switch: pushSwitch,
		})
	}
}

func Offline(ctx imcontext.WsHandleContext, code errs.IMErrorCode) {
	rpcCtx := imcontext.GetRpcContext(ctx)
	//offline event
	userId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	deviceId := imcontext.GetDeviceId(ctx)
	platform := imcontext.GetPlatform(ctx)
	clientIp := imcontext.GetContextAttrString(ctx, imcontext.StateKey_ClientIp)
	instanceId := imcontext.GetInstanceId(ctx)
	commonservices.SubOfflineEvent(rpcCtx, userId, &pbobjs.OnlineOfflineMsg{
		Type:       pbobjs.OnlineType_Offline,
		UserId:     userId,
		DeviceId:   deviceId,
		Platform:   platform,
		ClientIp:   clientIp,
		SessionId:  imcontext.GetConnSession(ctx),
		Timestamp:  time.Now().UnixMilli(),
		InstanceId: instanceId,
	})
	//visual log
	logmanager.WriteConnectionLog(rpcCtx, &pbobjs.ConnectionLog{
		AppKey:  imcontext.GetAppkey(ctx),
		Session: imcontext.GetConnSession(ctx),
		Action:  string(imcontext.Action_Disconnect),
		Code:    int32(code),
	})
	//quit rtc room
	if code == errs.IMErrorCode_CONNECT_KICKED_OFF {
		bases.AsyncRpcCall(rpcCtx, "rtc_sync_user_quit", userId, &pbobjs.Nil{})
	}
}
