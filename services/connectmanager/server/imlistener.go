package server

import (
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
	ExceptionCaught(ctx imcontext.WsHandleContext, code errs.IMErrorCode, e error)

	Connected(msg *codec.ConnectMsgBody, ctx imcontext.WsHandleContext)
	Diconnected(msg *codec.DisconnectMsgBody, ctx imcontext.WsHandleContext)
	PublishArrived(msg *codec.PublishMsgBody, qos int, ctx imcontext.WsHandleContext)
	PubAckArrived(msg *codec.PublishAckMsgBody, ctx imcontext.WsHandleContext)
	QueryArrived(msg *codec.QueryMsgBody, ctx imcontext.WsHandleContext)
	QueryConfirmArrived(msg *codec.QueryConfirmMsgBody, ctx imcontext.WsHandleContext)
	PingArrived(ctx imcontext.WsHandleContext)
}

type ImListenerImpl struct{}

func (listener *ImListenerImpl) ExceptionCaught(ctx imcontext.WsHandleContext, code errs.IMErrorCode, e error) {
	logs.Infof("session:%s\taction:%s\tcode:%d\terr:%v", imcontext.GetConnSession(ctx), imcontext.Action_Disconnect, code, e)

	services.Offline(ctx, code)

	services.RemoveFromContextCache(ctx)
}

func (listener *ImListenerImpl) Connected(msg *codec.ConnectMsgBody, ctx imcontext.WsHandleContext) {
	if msg.InstanceId != "" {
		imcontext.SetContextAttr(ctx, imcontext.StateKey_InstanceId, msg.InstanceId)
	}
	clientIp := msg.ClientIp
	if clientIp == "" {
		clientIp = GetRemoteAddr(ctx)
	}
	ucLog := &pbobjs.UserConnectLog{
		AppKey:   msg.Appkey,
		UserId:   msg.Token,
		Session:  imcontext.GetConnSession(ctx),
		Platform: msg.Platform,
		ClientIp: clientIp,
		Version:  msg.SdkVersion,
	}
	//check something
	if code, ext := services.CheckLogin(ctx, msg); code > 0 {
		//connect log
		uId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
		if uId != "" {
			ucLog.UserId = uId
		}
		ucLog.Code = code
		logmanager.WriteUserConnectLog(imcontext.GetRpcContext(ctx), ucLog)
		msgAck := codec.NewConnectAckMessage(&codec.ConnectAckMsgBody{
			Code:      code,
			Session:   imcontext.GetConnSession(ctx),
			Timestamp: time.Now().UnixMilli(),
			Ext:       ext,
		})
		ctx.Write(msgAck)
		go func() {
			time.Sleep(50 * time.Millisecond)
			ctx.Close(errors.New("failed to login"))
		}()
		logs.Infof("session:%s\taction:%s\tappkey:%s\tuser_id:%s\tclient_ip:%s\tplatform:%s\tdevice_id:%s\tpush_token:%s\tinstance_id:%s\tcode:%d", imcontext.GetConnSession(ctx), imcontext.Action_Connect, msg.Appkey, ucLog.UserId, clientIp, msg.Platform, msg.DeviceId, msg.PushToken, msg.InstanceId, ucLog.Code)
		return
	}
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Connected, "1")
	userId := imcontext.GetContextAttrString(ctx, imcontext.StateKey_UserID)
	//success
	logs.Infof("session:%s\taction:%s\tappkey:%s\tuser_id:%s\tclient_ip:%s\tplatform:%s\tdevice_id:%s\tpush_token:%s\tinstance_id:%s", imcontext.GetConnSession(ctx), imcontext.Action_Connect, msg.Appkey, userId, clientIp, msg.Platform, msg.DeviceId, msg.PushToken, msg.InstanceId)
	commonservices.ReportUserLogin(msg.Appkey, userId)

	imcontext.SetContextAttr(ctx, imcontext.StateKey_Appkey, msg.Appkey)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_DeviceID, msg.DeviceId)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Platform, msg.Platform)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Version, msg.SdkVersion)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ClientIp, clientIp)
	services.PutInContextCache(ctx)

	//reg push token
	services.RegPushToken(ctx, msg.Appkey, userId, msg.DeviceId, msg.Platform, msg.PushChannel, msg.PackageName, msg.PushToken)

	services.Online(ctx, msg.Ext)

	//connect log
	ucLog.UserId = userId
	ucLog.Code = int32(errs.IMErrorCode_SUCCESS)
	logmanager.WriteUserConnectLog(imcontext.GetRpcContext(ctx), ucLog)
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

	services.Offline(ctx, errs.IMErrorCode(msg.Code))

	ctx.Close(fmt.Errorf("dissconnect"))
	services.RemoveFromContextCache(ctx)
}
func (*ImListenerImpl) PublishArrived(msg *codec.PublishMsgBody, qos int, ctx imcontext.WsHandleContext) {
	logs.Infof("session:%s\taction:%s\tseq_index:%d\ttopic:%s\ttarget_id:%s\tlen:%d", imcontext.GetConnSession(ctx), imcontext.Action_UserPub, msg.Index, msg.Topic, msg.TargetId, len(msg.Data))
	logmanager.WriteConnectionLog(imcontext.GetRpcContext(ctx), &pbobjs.ConnectionLog{
		AppKey:   imcontext.GetAppkey(ctx),
		Session:  imcontext.GetConnSession(ctx),
		Index:    msg.Index,
		Action:   string(imcontext.Action_UserPub),
		Method:   msg.Topic,
		TargetId: msg.TargetId,
		DataLen:  int32(len(msg.Data)),
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
		logmanager.WriteConnectionLog(imcontext.GetRpcContext(ctx), &pbobjs.ConnectionLog{
			AppKey:   imcontext.GetAppkey(ctx),
			Session:  imcontext.GetConnSession(ctx),
			Index:    msg.Index,
			Action:   string(imcontext.Action_UserPubAck),
			Method:   msg.Topic,
			TargetId: msg.TargetId,
			Code:     int32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
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
		logmanager.WriteConnectionLog(imcontext.GetRpcContext(ctx), &pbobjs.ConnectionLog{
			AppKey:   imcontext.GetAppkey(ctx),
			Session:  imcontext.GetConnSession(ctx),
			Index:    msg.Index,
			Action:   string(imcontext.Action_UserPubAck),
			Method:   msg.Topic,
			TargetId: msg.TargetId,
			Code:     int32(errs.IMErrorCode_CONNECT_EXCEEDLIMITED),
		})
		return
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
	logmanager.WriteConnectionLog(imcontext.GetRpcContext(ctx), &pbobjs.ConnectionLog{
		AppKey:  imcontext.GetAppkey(ctx),
		Session: imcontext.GetConnSession(ctx),
		Index:   msg.Index,
		Action:  string(imcontext.Action_ServerPubAck),
	})
}
func (listener *ImListenerImpl) QueryArrived(msg *codec.QueryMsgBody, ctx imcontext.WsHandleContext) {
	logs.Infof("session:%s\taction:%s\tseq_index:%d\ttopic:%s\ttarget_id:%s\tlen:%d", imcontext.GetConnSession(ctx), imcontext.Action_Query, msg.Index, msg.Topic, msg.TargetId, len(msg.Data))
	logmanager.WriteConnectionLog(imcontext.GetRpcContext(ctx), &pbobjs.ConnectionLog{
		AppKey:   imcontext.GetAppkey(ctx),
		Session:  imcontext.GetConnSession(ctx),
		Index:    msg.Index,
		Action:   string(imcontext.Action_Query),
		Method:   msg.Topic,
		TargetId: msg.TargetId,
		DataLen:  int32(len(msg.Data)),
	})

	if msg.Topic == "" || msg.TargetId == "" {
		ack := codec.NewQueryAckMessage(&codec.QueryAckMsgBody{
			Index:     msg.Index,
			Code:      int32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
			Timestamp: time.Now().UnixMilli(),
		}, codec.QoS_NoAck)
		ctx.Write(ack)
		logs.Infof("session:%s\taction:%s\tseq_index:%d\tcode:%d", imcontext.GetConnSession(ctx), imcontext.Action_QueryAck, msg.Index, errs.IMErrorCode_CONNECT_PARAM_REQUIRED)
		logmanager.WriteConnectionLog(imcontext.GetRpcContext(ctx), &pbobjs.ConnectionLog{
			AppKey:  imcontext.GetAppkey(ctx),
			Session: imcontext.GetConnSession(ctx),
			Index:   msg.Index,
			Action:  string(imcontext.Action_QueryAck),
			Code:    int32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
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
		logmanager.WriteConnectionLog(imcontext.GetRpcContext(ctx), &pbobjs.ConnectionLog{
			AppKey:  imcontext.GetAppkey(ctx),
			Session: imcontext.GetConnSession(ctx),
			Index:   msg.Index,
			Action:  string(imcontext.Action_QueryAck),
			Code:    int32(errs.IMErrorCode_CONNECT_EXCEEDLIMITED),
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
		DeviceId:     imcontext.GetDeviceId(ctx),
		InstanceId:   imcontext.GetInstanceId(ctx),
		Platform:     imcontext.GetPlatform(ctx),
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
	logmanager.WriteConnectionLog(imcontext.GetRpcContext(ctx), &pbobjs.ConnectionLog{
		AppKey:  imcontext.GetAppkey(ctx),
		Session: imcontext.GetConnSession(ctx),
		Index:   msg.Index,
		Action:  string(imcontext.Action_QueryConfirm),
	})
}
func (*ImListenerImpl) PingArrived(ctx imcontext.WsHandleContext) {
	ctx.Write(codec.NewPongMessage())
}

func GetRemoteAddr(ctx imcontext.WsHandleContext) string {
	return ctx.RemoteAddr()
}
