package server

import (
	"context"
	"errors"
	"fmt"
	"im-server/services/commonservices"
	"im-server/services/logmanager"
	"time"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/logs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
	"im-server/services/connectmanager/services"
)

type ImListener interface {
	Create(ctx imcontext.WsHandleContext)
	Close(ctx imcontext.WsHandleContext)
	ExceptionCaught(ctx imcontext.WsHandleContext, e error)

	Connected(msg *codec.ConnectMsgBody, ctx imcontext.WsHandleContext)
	Diconnected(msg *codec.DisconnectMsgBody, ctx imcontext.WsHandleContext)
	PublishArrived(msg *codec.PublishMsgBody, qos int, ctx imcontext.WsHandleContext)
	PubAckArrived(msg *codec.PublishAckMsgBody, ctx imcontext.WsHandleContext)
	QueryArrived(msg *codec.QueryMsgBody, ctx imcontext.WsHandleContext)
	QueryConfirmArrived(msg *codec.QueryConfirmMsgBody, ctx imcontext.WsHandleContext)
	PingArrived(ctx imcontext.WsHandleContext)
}

type ImListenerImpl struct{}

func (*ImListenerImpl) Create(ctx imcontext.WsHandleContext) {
}
func (listener *ImListenerImpl) Close(ctx imcontext.WsHandleContext) {
	logs.Infof("1session:%s\taction:%s", imcontext.GetConnSession(ctx), imcontext.Action_Disconnect)
	services.RemoveFromContextCache(ctx)

	//subscriptions.PublishOnlineOfflineMsg(listener.context(ctx), userId, &pbobjs.OnlineOfflineMsg{
	//	Type:      pbobjs.OnlineType_Online,
	//	UserId:    userId,
	//	DeviceId:  msg.DeviceId,
	//	Platform:  msg.Platform,
	//	ClientIp:  msg.ClientIp,
	//	SessionId: utils.GetConnSession(ctx),
	//	Timestamp: time.Now().UnixMilli(),
	//})
}
func (listener *ImListenerImpl) ExceptionCaught(ctx imcontext.WsHandleContext, e error) {
	logs.Infof("2session:%s\taction:%s\terr:%v", imcontext.GetConnSession(ctx), imcontext.Action_Disconnect, e)
	userId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	deviceId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_DeviceID)
	platform := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Platform)
	clientIp := imcontext.GetContextAttrString(ctx, imcontext.StateKey_ClientIp)
	appKey := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey)
	clientSession := imcontext.GetContextAttrString(ctx, imcontext.StateKey_ClientSession)

	offlineMsg := &pbobjs.OnlineOfflineMsg{
		Type:          pbobjs.OnlineType_Offline,
		UserId:        userId,
		DeviceId:      deviceId,
		Platform:      platform,
		ClientIp:      clientIp,
		SessionId:     imcontext.GetConnSession(ctx),
		Timestamp:     time.Now().UnixMilli(),
		ClientSession: clientSession,
	}

	newCtx := listener.context(ctx)
	commonservices.SubOfflineEvent(newCtx, userId, offlineMsg)
	logmanager.WriteDisconnectionLog(newCtx, &pbobjs.DisconnectionLog{
		Timestamp:   logmanager.LogTimestamp(),
		ServiceName: "connect",
		Session:     imcontext.GetConnSession(ctx),
		Action:      string(imcontext.Action_Disconnect),
		AppKey:      appKey,
		UserId:      userId,
		Code:        0,
		Err:         e.Error(),
	})

	services.RemoveFromContextCache(ctx)
}

func (listener *ImListenerImpl) context(inboundCtx imcontext.WsHandleContext) context.Context {
	return imcontext.SendCtxFromNettyCtx(inboundCtx)
}

func (listener *ImListenerImpl) Connected(msg *codec.ConnectMsgBody, ctx imcontext.WsHandleContext) {
	if msg.ClientSession != "" {
		imcontext.SetContextAttr(ctx, imcontext.StateKey_ClientSession, msg.ClientSession)
	}
	clientIp := msg.ClientIp
	if clientIp == "" {
		clientIp = GetRemoteAddr(ctx)
	}
	//check something
	if code, ext := services.CheckLogin(ctx, msg); code > 0 {
		msgAck := codec.NewConnectAckMessage(&codec.ConnectAckMsgBody{
			Code:      code,
			Session:   imcontext.GetConnSession(ctx),
			Timestamp: time.Now().UnixMilli(),
			Ext:       ext,
		})
		ctx.Write(msgAck)
		go func() {
			time.Sleep(time.Millisecond * 50)
			ctx.Close(errors.New("Failed to Login"))
		}()
		return
	}
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Connected, "1")
	userId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	//success
	logs.Infof("session:%s\taction:%s\tappkey:%s\tuser_id:%s\tclient_ip:%s\tplatform:%s\tdevice_id:%s\tpush_token:%s\tclient_session:%s", imcontext.GetConnSession(ctx), imcontext.Action_Connect, msg.Appkey, userId, clientIp, msg.Platform, msg.DeviceId, msg.PushToken, msg.ClientSession)

	imcontext.SetContextAttr(ctx, imcontext.StateKey_Appkey, msg.Appkey)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_DeviceID, msg.DeviceId)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Platform, msg.Platform)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Version, msg.SdkVersion)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ClientIp, clientIp)
	services.PutInContextCache(ctx)

	//reg push token
	services.RegPushToken(ctx, msg.Appkey, userId, msg.DeviceId, msg.Platform, msg.PushChannel, msg.PackageName, msg.PushToken)

	onlineMsg := &pbobjs.OnlineOfflineMsg{
		Type:          pbobjs.OnlineType_Online,
		UserId:        userId,
		DeviceId:      msg.DeviceId,
		Platform:      msg.Platform,
		ClientIp:      clientIp,
		SessionId:     imcontext.GetConnSession(ctx),
		Timestamp:     time.Now().UnixMilli(),
		ConnectionExt: msg.Ext,
	}

	newCtx := listener.context(ctx)
	commonservices.SubOnlineEvent(newCtx, userId, onlineMsg)
	logmanager.WriteConnectionLog(newCtx, &pbobjs.ConnectionLog{
		Timestamp:   logmanager.LogTimestamp(),
		ServiceName: "connect",
		Session:     imcontext.GetConnSession(ctx),
		Action:      string(imcontext.Action_Connect),
		AppKey:      msg.Appkey,
		UserId:      userId,
		Platform:    msg.Platform,
		PushToken:   msg.PushToken,
		PushChannel: msg.PushChannel,
		ClientIp:    clientIp,
	})
	msgAck := codec.NewConnectAckMessage(&codec.ConnectAckMsgBody{
		Code:      int32(errs.IMErrorCode_SUCCESS),
		UserId:    userId,
		Session:   imcontext.GetConnSession(ctx),
		Timestamp: time.Now().UnixMilli(),
	})
	ctx.Write(msgAck)
}

func (listener *ImListenerImpl) Diconnected(msg *codec.DisconnectMsgBody, ctx imcontext.WsHandleContext) {
	logs.Infof("session:%s\taction:%s\tcode:%d", imcontext.GetConnSession(ctx), imcontext.Action_Disconnect, msg.Code)
	ctx.Close(fmt.Errorf("dissconnect"))

	userId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	deviceId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_DeviceID)
	platform := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Platform)
	clientIp := imcontext.GetContextAttrString(ctx, imcontext.StateKey_ClientIp)
	clientSession := imcontext.GetContextAttrString(ctx, imcontext.StateKey_ClientSession)

	commonservices.SubOfflineEvent(listener.context(ctx), userId, &pbobjs.OnlineOfflineMsg{
		Type:          pbobjs.OnlineType_Offline,
		UserId:        userId,
		DeviceId:      deviceId,
		Platform:      platform,
		ClientIp:      clientIp,
		SessionId:     imcontext.GetConnSession(ctx),
		Timestamp:     time.Now().UnixMilli(),
		ClientSession: clientSession,
	})

	services.RemoveFromContextCache(ctx)
}
func (*ImListenerImpl) PublishArrived(msg *codec.PublishMsgBody, qos int, ctx imcontext.WsHandleContext) {
	logs.Infof("session:%s\taction:%s\tseq_index:%d\ttopic:%s\ttarget_id:%s\tlen:%d", imcontext.GetConnSession(ctx), imcontext.Action_UserPub, msg.Index, msg.Topic, msg.TargetId, len(msg.Data))
	logmanager.WriteSdkRequestLog(imcontext.SendCtxFromNettyCtx(ctx), &pbobjs.SdkRequestLog{
		Timestamp:   logmanager.LogTimestamp(),
		ServiceName: "connect",
		Session:     imcontext.GetConnSession(ctx),
		Index:       uint32(msg.Index),
		Action:      string(imcontext.Action_UserPub),
		Method:      msg.Topic,
		TargetId:    msg.TargetId,
		Len:         uint32(len(msg.Data)),
		AppKey:      imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
	})
	//check params
	if msg.Topic == "" || msg.TargetId == "" {
		ack := codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
			Index: msg.Index,
			Code:  int32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
			MsgId: "",
		})
		ctx.Write(ack)
		logs.Infof("session:%s\taction:%s\tseq_index:%d\tcode:%d", imcontext.GetConnSession(ctx), imcontext.Action_UserPubAck, msg.Index, errs.IMErrorCode_CONNECT_PARAM_REQUIRED)
		logmanager.WriteSdkResponseLog(imcontext.SendCtxFromNettyCtx(ctx), &pbobjs.SdkResponseLog{
			Timestamp:   logmanager.LogTimestamp(),
			ServiceName: "connect",
			Session:     imcontext.GetConnSession(ctx),
			Index:       uint32(msg.Index),
			Action:      string(imcontext.Action_UserPubAck),
			Code:        uint32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
			Len:         0,
			AppKey:      imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
		})
		return
	}
	//check limiter
	limiter := imcontext.GetLimiter(ctx)
	if limiter != nil && !limiter.Allow() {
		ack := codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
			Index:     msg.Index,
			Code:      int32(errs.IMErrorCode_CONNECT_EXCEEDLIMITED),
			MsgId:     "",
			Timestamp: time.Now().UnixMilli(),
		})
		ctx.Write(ack)
		logs.Infof("session:%s\taction:%s\tseq_index:%d\tcode:%d", imcontext.GetConnSession(ctx), imcontext.Action_UserPubAck, msg.Index, errs.IMErrorCode_CONNECT_EXCEEDLIMITED)
		logmanager.WriteSdkResponseLog(imcontext.SendCtxFromNettyCtx(ctx), &pbobjs.SdkResponseLog{
			Timestamp:   logmanager.LogTimestamp(),
			ServiceName: "connect",
			Session:     imcontext.GetConnSession(ctx),
			Index:       uint32(msg.Index),
			Action:      string(imcontext.Action_UserPubAck),
			Code:        uint32(errs.IMErrorCode_CONNECT_EXCEEDLIMITED),
			Len:         0,
			AppKey:      imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
		})
		return
	}

	isSucc := bases.UnicastRoute(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_UserPub,
		AppKey:       imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
		Session:      imcontext.GetConnSession(ctx),
		Method:       "upstream",
		RequesterId:  msg.TargetId,
		ReqIndex:     msg.Index,
		Qos:          int32(qos),
		AppDataBytes: msg.Data,
		TargetId:     imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID),
		ExtParams: map[string]string{
			commonservices.RpcExtKey_RealMethod: msg.Topic,
		},
	}, "connect")

	if !isSucc {
		ack := codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
			Index:     msg.Index,
			Code:      int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC),
			MsgId:     "",
			Timestamp: time.Now().UnixMilli(),
		})
		ctx.Write(ack)
	}
}
func (*ImListenerImpl) PubAckArrived(msg *codec.PublishAckMsgBody, ctx imcontext.WsHandleContext) {
	index := msg.GetIndex()
	callback := imcontext.GetAndDeleteServerPubCallback(ctx, index)
	if callback != nil {
		callback()
	}
	logs.Infof("session:%s\taction:%s\tseq_index:%d", imcontext.GetConnSession(ctx), imcontext.Action_ServerPubAck, msg.Index)
	logmanager.WriteSdkResponseLog(imcontext.SendCtxFromNettyCtx(ctx), &pbobjs.SdkResponseLog{
		Timestamp:   logmanager.LogTimestamp(),
		ServiceName: "connect",
		Session:     imcontext.GetConnSession(ctx),
		Index:       uint32(msg.Index),
		Action:      string(imcontext.Action_ServerPubAck),
		Code:        0,
		Len:         0,
		AppKey:      imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
	})
}
func (listener *ImListenerImpl) QueryArrived(msg *codec.QueryMsgBody, ctx imcontext.WsHandleContext) {
	logs.Infof("session:%s\taction:%s\tseq_index:%d\ttopic:%s\ttarget_id:%s\tlen:%d", imcontext.GetConnSession(ctx), imcontext.Action_Query, msg.Index, msg.Topic, msg.TargetId, len(msg.Data))
	logmanager.WriteSdkRequestLog(listener.context(ctx), &pbobjs.SdkRequestLog{
		Timestamp:   logmanager.LogTimestamp(),
		ServiceName: "connect",
		Session:     imcontext.GetConnSession(ctx),
		Index:       uint32(msg.Index),
		Action:      string(imcontext.Action_Query),
		Method:      msg.Topic,
		TargetId:    msg.TargetId,
		Len:         uint32(len(msg.Data)),
		AppKey:      imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
	})

	if msg.Topic == "" || msg.TargetId == "" {
		ack := codec.NewQueryAckMessage(&codec.QueryAckMsgBody{
			Index:     msg.Index,
			Code:      int32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
			Timestamp: time.Now().UnixMilli(),
		}, codec.QoS_NoAck)
		ctx.Write(ack)
		logs.Infof("session:%s\taction:%s\tseq_index:%d\tcode:%d", imcontext.GetConnSession(ctx), imcontext.Action_QueryAck, msg.Index, errs.IMErrorCode_CONNECT_PARAM_REQUIRED)
		logmanager.WriteSdkResponseLog(listener.context(ctx), &pbobjs.SdkResponseLog{
			Timestamp:   logmanager.LogTimestamp(),
			ServiceName: "connect",
			Session:     imcontext.GetConnSession(ctx),
			Index:       uint32(msg.Index),
			Action:      string(imcontext.Action_QueryAck),
			Code:        uint32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
			Len:         uint32(len(msg.Data)),
			AppKey:      imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
		})
		return
	}
	//check limiter
	limiter := imcontext.GetLimiter(ctx)
	if limiter != nil && !limiter.Allow() {
		ack := codec.NewQueryAckMessage(&codec.QueryAckMsgBody{
			Index:     msg.Index,
			Code:      int32(errs.IMErrorCode_CONNECT_EXCEEDLIMITED),
			Timestamp: time.Now().UnixMilli(),
		}, codec.QoS_NoAck)
		ctx.Write(ack)
		logs.Infof("session:%s\taction:%s\tseq_index:%d\tcode:%d", imcontext.GetConnSession(ctx), imcontext.Action_QueryAck, msg.Index, errs.IMErrorCode_CONNECT_EXCEEDLIMITED)
		logmanager.WriteSdkResponseLog(listener.context(ctx), &pbobjs.SdkResponseLog{
			Timestamp:   logmanager.LogTimestamp(),
			ServiceName: "connect",
			Session:     imcontext.GetConnSession(ctx),
			Index:       uint32(msg.Index),
			Action:      string(imcontext.Action_QueryAck),
			Code:        uint32(errs.IMErrorCode_CONNECT_EXCEEDLIMITED),
			Len:         uint32(len(msg.Data)),
			AppKey:      imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
		})
		return
	}

	appkey := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey)
	userid := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	targetId := msg.TargetId

	if tId, ok := services.HisMsgRedirect(msg.Topic, msg.Data, userid, msg.TargetId); ok {
		targetId = tId
	}
	isSucc := bases.UnicastRoute(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_QueryMsg,
		AppKey:       appkey,
		Session:      imcontext.GetConnSession(ctx),
		Method:       msg.Topic,
		RequesterId:  userid,
		ReqIndex:     msg.Index,
		Qos:          int32(codec.QoS_NeedAck),
		AppDataBytes: msg.Data,
		TargetId:     targetId,
		TerminalNum:  services.GetConnectCountByUser(appkey, userid),
	}, "connect")
	if !isSucc {
		ack := codec.NewQueryAckMessage(&codec.QueryAckMsgBody{
			Index:     msg.Index,
			Code:      int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC),
			Timestamp: time.Now().UnixMilli(),
		}, codec.QoS_NoAck)
		ctx.Write(ack)
	}
}
func (*ImListenerImpl) QueryConfirmArrived(msg *codec.QueryConfirmMsgBody, ctx imcontext.WsHandleContext) {
	index := msg.GetIndex()
	callback := imcontext.GetAndDeleteQueryAckCallback(ctx, index)
	if callback != nil {
		callback()
	}
	logs.Infof("session:%s\taction:%s\tseq_index:%d", imcontext.GetConnSession(ctx), imcontext.Action_QueryConfirm, msg.Index)
	logmanager.WriteSdkResponseLog(imcontext.SendCtxFromNettyCtx(ctx), &pbobjs.SdkResponseLog{
		Timestamp:   logmanager.LogTimestamp(),
		ServiceName: "connect",
		Session:     imcontext.GetConnSession(ctx),
		Index:       uint32(msg.Index),
		Action:      string(imcontext.Action_QueryConfirm),
		Code:        0,
		Len:         0,
		AppKey:      imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
	})
}
func (*ImListenerImpl) PingArrived(ctx imcontext.WsHandleContext) {
	ctx.Write(codec.NewPongMessage())
}

func GetRemoteAddr(ctx imcontext.WsHandleContext) string {
	return ctx.RemoteAddr()
}
