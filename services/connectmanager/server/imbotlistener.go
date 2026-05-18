package server

import (
	"errors"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
	"im-server/services/connectmanager/services"
	"time"
)

type ImBotListenerImpl struct{}

func (listener *ImBotListenerImpl) ExceptionCaught(ctx imcontext.WsHandleContext, code errs.IMErrorCode, e error) {
}

func (listener *ImBotListenerImpl) Connected(msg *codec.ConnectMsgBody, ctx imcontext.WsHandleContext) {
	if msg == nil {
		logs.NewLogEntity().Error("connect body is nil")
		ctx.Close(errors.New("connect body is nil"))
		return
	}
	clientIp := imcontext.GetContextAttrString(ctx, imcontext.StateKey_ClientIp)
	clientHost := imcontext.GetContextAttrString(ctx, imcontext.StateKey_ClientHost)
	//check login
	if code := services.CheckBotLogin(ctx, msg); code != errs.IMErrorCode_SUCCESS {
		uId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
		if uId == "" {
			uId = msg.Token
		}
		msgAck := codec.NewConnectAckMessage(&codec.ConnectAckMsgBody{
			Code:      int32(code),
			Session:   imcontext.GetConnSession(ctx),
			Timestamp: time.Now().UnixMilli(),
		})
		ctx.Write(msgAck)
		go func() {
			time.Sleep(50 * time.Millisecond)
			ctx.Close(errors.New("failed to login"))
		}()
		logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
			WithField("session", imcontext.GetConnSession(ctx)).
			WithField("action", imcontext.Action_Connect).
			WithField("appkey", msg.Appkey).
			WithField("user_id", uId).
			WithField("client_ip", clientIp).
			WithField("client_host", clientHost).
			WithField("platform", msg.Platform).
			WithField("version", msg.SdkVersion).
			WithField("code", code).Info("")
		return
	}
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Connected, "1")
	userId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	//success
	logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
		WithField("session", imcontext.GetConnSession(ctx)).
		WithField("action", imcontext.Action_Connect).
		WithField("appkey", msg.Appkey).
		WithField("user_id", userId).
		WithField("client_ip", clientIp).
		WithField("client_host", clientHost).
		WithField("platform", msg.Platform).
		WithField("version", msg.SdkVersion).Info("")
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Appkey, msg.Appkey)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Version, msg.SdkVersion)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Platform, msg.Platform)
	services.PutInContextCache(ctx)
	msgAck := codec.NewConnectAckMessage(&codec.ConnectAckMsgBody{
		Code:      int32(errs.IMErrorCode_SUCCESS),
		UserId:    userId,
		Session:   imcontext.GetConnSession(ctx),
		Timestamp: time.Now().UnixMilli(),
	})
	ctx.Write(msgAck)
}

func (listener *ImBotListenerImpl) Diconnected(msg *codec.DisconnectMsgBody, ctx imcontext.WsHandleContext) {
	if msg == nil {
		logs.NewLogEntity().Error("disconnect body is nil")
		msg = &codec.DisconnectMsgBody{}
	}
	logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
		WithField("session", imcontext.GetConnSession(ctx)).
		WithField("action", imcontext.Action_Disconnect).
		WithField("code", msg.Code).Info("")

	services.RemoveFromContextCache(ctx)

	ctx.Close(fmt.Errorf("dissconnect"))
}

func (listener *ImBotListenerImpl) PublishArrived(msg *codec.PublishMsgBody, qos int, ctx imcontext.WsHandleContext) {
	if msg == nil {
		logs.NewLogEntity().Error("pub body is nil")
		return
	}
	logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
		WithField("session", imcontext.GetConnSession(ctx)).
		WithField("action", imcontext.Action_UserPub).
		WithField("seq_index", msg.Index).
		WithField("method", msg.Topic).
		WithField("target_id", msg.TargetId).
		WithField("len", len(msg.Data)).Info("")
	//check params
	if msg.Topic == "" || msg.TargetId == "" {
		ack := codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
			Index: msg.Index,
			Code:  int32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
			MsgId: "",
		})
		ctx.Write(ack)
		logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
			WithField("session", imcontext.GetConnSession(ctx)).
			WithField("action", imcontext.Action_UserPubAck).
			WithField("seq_index", msg.Index).
			WithField("code", errs.IMErrorCode_CONNECT_PARAM_REQUIRED).Info("")
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
		logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
			WithField("session", imcontext.GetConnSession(ctx)).
			WithField("action", imcontext.Action_UserPubAck).
			WithField("seq_index", msg.Index).
			WithField("code", errs.IMErrorCode_CONNECT_EXCEEDLIMITED).Info("")
		return
	}
	isFromApi := false
	userType := imcontext.GetUserType(ctx)
	if userType == pbobjs.UserType_Admin {
		isFromApi = true
	}
	isSucc := bases.UnicastRoute(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_UserPub,
		AppKey:       imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey),
		Session:      imcontext.GetConnSession(ctx),
		DeviceId:     imcontext.GetDeviceId(ctx),
		InstanceId:   imcontext.GetInstanceId(ctx),
		Platform:     imcontext.GetPlatform(ctx),
		Method:       "upstream",
		RequesterId:  msg.TargetId,
		ReqIndex:     msg.Index,
		Qos:          int32(qos),
		AppDataBytes: msg.Data,
		TargetId:     imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID),
		ExtParams: map[string]string{
			commonservices.RpcExtKey_RealMethod: msg.Topic,
		},
		IsFromApi: isFromApi,
	}, "connect")

	if !isSucc {
		ack := codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
			Index:     msg.Index,
			Code:      int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC),
			MsgId:     "",
			Timestamp: time.Now().UnixMilli(),
		})
		ctx.Write(ack)
		logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
			WithField("session", imcontext.GetConnSession(ctx)).
			WithField("action", imcontext.Action_UserPubAck).
			WithField("seq_index", msg.Index).
			WithField("code", errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC).Info("")
	}
}

func (listener *ImBotListenerImpl) PubAckArrived(msg *codec.PublishAckMsgBody, ctx imcontext.WsHandleContext) {
	if msg == nil {
		logs.NewLogEntity().Error("pub_ack body is nil")
		return
	}
	index := msg.GetIndex()
	callback := imcontext.GetAndDeleteServerPubCallback(ctx, index)
	if callback != nil {
		callback()
	}
	logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
		WithField("session", imcontext.GetConnSession(ctx)).
		WithField("action", imcontext.Action_ServerPubAck).
		WithField("seq_index", msg.Index).Info("")
}

func (listener *ImBotListenerImpl) QueryArrived(msg *codec.QueryMsgBody, ctx imcontext.WsHandleContext) {
	if msg == nil {
		logs.NewLogEntity().Error("qry body is nil")
		return
	}
	logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
		WithField("session", imcontext.GetConnSession(ctx)).
		WithField("action", imcontext.Action_Query).
		WithField("seq_index", msg.Index).
		WithField("method", msg.Topic).
		WithField("target_id", msg.TargetId).
		WithField("len", len(msg.Data)).Info("")

	if msg.Topic == "" || msg.TargetId == "" {
		ack := codec.NewQueryAckMessage(&codec.QueryAckMsgBody{
			Index:     msg.Index,
			Code:      int32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
			Timestamp: time.Now().UnixMilli(),
		}, codec.QoS_NoAck)
		ctx.Write(ack)
		logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
			WithField("session", imcontext.GetConnSession(ctx)).
			WithField("action", imcontext.Action_QueryAck).
			WithField("seq_index", msg.Index).
			WithField("code", errs.IMErrorCode_CONNECT_PARAM_REQUIRED).Info("")
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
		logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
			WithField("session", imcontext.GetConnSession(ctx)).
			WithField("action", imcontext.Action_QueryAck).
			WithField("seq_index", msg.Index).
			WithField("code", errs.IMErrorCode_CONNECT_EXCEEDLIMITED).Info("")
		return
	}

	appkey := imcontext.GetContextAttrString(ctx, imcontext.StateKey_Appkey)
	userid := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)

	//preprocess rtc room creation
	services.PreProcessRtcCreate(msg)

	targetId := msg.TargetId
	if tId, ok := services.HisMsgRedirect(msg.Topic, msg.Data, userid, msg.TargetId); ok {
		targetId = tId
	}

	//chatroom
	if msg.Topic == "c_sync_msgs" || msg.Topic == "c_sync_atts" || msg.Topic == "c_sync_events" {
		targetId = userid
	}

	isSucc := bases.UnicastRoute(&pbobjs.RpcMessageWraper{
		RpcMsgType:    pbobjs.RpcMsgType_QueryMsg,
		AppKey:        appkey,
		Session:       imcontext.GetConnSession(ctx),
		ConnectedTime: imcontext.GetConnectCreateTime(ctx),
		DeviceId:      imcontext.GetDeviceId(ctx),
		InstanceId:    imcontext.GetInstanceId(ctx),
		Platform:      imcontext.GetPlatform(ctx),
		Method:        msg.Topic,
		RequesterId:   userid,
		ReqIndex:      msg.Index,
		Qos:           int32(codec.QoS_NeedAck),
		AppDataBytes:  msg.Data,
		TargetId:      targetId,
		TerminalNum:   services.GetConnectCountByUser(appkey, userid),
	}, "connect")
	if !isSucc {
		ack := codec.NewQueryAckMessage(&codec.QueryAckMsgBody{
			Index:     msg.Index,
			Code:      int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC),
			Timestamp: time.Now().UnixMilli(),
		}, codec.QoS_NoAck)
		ctx.Write(ack)
		logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
			WithField("session", imcontext.GetConnSession(ctx)).
			WithField("action", imcontext.Action_QueryAck).
			WithField("seq_index", msg.Index).
			WithField("code", errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC).Info("")
	}
}

func (listener *ImBotListenerImpl) QueryConfirmArrived(msg *codec.QueryConfirmMsgBody, ctx imcontext.WsHandleContext) {
	if msg == nil {
		logs.NewLogEntity().Error("qry_confirm body is nil")
		return
	}
	index := msg.GetIndex()
	callback := imcontext.GetAndDeleteQueryAckCallback(ctx, index)
	if callback != nil {
		callback()
	}
	logs.NewLogEntity().WithField("service_name", imcontext.ServiceName).
		WithField("session", imcontext.GetConnSession(ctx)).
		WithField("action", imcontext.Action_QueryConfirm).
		WithField("seq_index", msg.Index).Info("")
}

func (listener *ImBotListenerImpl) PingArrived(ctx imcontext.WsHandleContext) {
	ctx.Write(codec.NewPongMessage())
}
