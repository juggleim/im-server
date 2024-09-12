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

func Online(ctx imcontext.WsHandleContext, ext string) {
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
}
