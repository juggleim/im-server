package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/logmanager"
	"time"

	"im-server/commons/logs"
	"im-server/services/commonservices"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"

	"github.com/rfyiamcool/go-timewheel"
)

var callbackTimeoutTimer *timewheel.TimeWheel

func init() {
	t, err := timewheel.NewTimeWheel(1*time.Second, 360)
	if err != nil {
		logs.Error("can not init timeWheel for publish callback.")
	} else {
		callbackTimeoutTimer = t
		callbackTimeoutTimer.Start()
	}
}
func PublishServerPubMessage(appkey, userid, session string, msgId string, msgTime int64, serverPubMsg *codec.PublishMsgBody, publishType commonservices.PublishType, callback func(*pbobjs.OnlineStatus)) {
	userCtxMap := GetConnectCtxByUser(appkey, userid)
	onlineStatus := &pbobjs.OnlineStatus{
		Type:      pbobjs.OnlineType_Offline,
		Platforms: []string{},
	}
	if len(userCtxMap) > 0 { //target user is online
		onlineStatus.Type = pbobjs.OnlineType_Online
		isSetCallback := false
		for kSess, vCtx := range userCtxMap {
			if publishType == commonservices.PublishType_OnlineSelfSession && kSess != session { //publishType:1, 只给指定的session发送
				continue
			}
			if publishType == commonservices.PublishType_AllSessionExceptSelf && kSess == session { //publishType:2, 除了指定session以外，给该用户其他登录端发送
				continue
			}
			qos := codec.QoS_NoAck
			if callback != nil {
				qos = codec.QoS_NeedAck
			}
			index := imcontext.GetServerIndexAfterIncrease(vCtx)
			tmpPubMsg := codec.NewServerPublishMessage(&codec.PublishMsgBody{
				Index:     int32(index),
				Topic:     serverPubMsg.Topic,
				TargetId:  serverPubMsg.TargetId,
				Timestamp: time.Now().UnixMilli(),
				Data:      serverPubMsg.Data,
			}, qos)
			if callback != nil && !isSetCallback {
				isSetCallback = true
				if serverPubMsg.Topic == "msg" {
					task := callbackTimeoutTimer.Add(5*time.Second, func() {
						//do timeout
						imcontext.RemoveServerPubCallback(vCtx, tmpPubMsg.MsgBody.Index)
					})
					targetSession := imcontext.GetConnSession(vCtx)
					targetIndex := tmpPubMsg.MsgBody.Index
					imcontext.PutServerPubCallback(vCtx, tmpPubMsg.MsgBody.Index, func() {
						callbackTimeoutTimer.Remove(task)
						data, _ := tools.PbMarshal(&pbobjs.MsgAck{
							MsgId:   msgId,
							MsgTime: msgTime,
						})
						bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
							RpcMsgType:   pbobjs.RpcMsgType_UserPubAck,
							AppKey:       appkey,
							Session:      targetSession,
							ReqIndex:     targetIndex,
							Method:       "msg_ack",
							RequesterId:  userid,
							TargetId:     userid,
							AppDataBytes: data,
						})
					})
					onlineStatus.TargetSession = targetSession
					onlineStatus.TargetIndex = targetIndex
				}
			}
			vCtx.Write(tmpPubMsg)
			logs.Infof("session:%s\taction:%s\tseq_index:%d\ttopic:%s\tlen:%d", imcontext.GetConnSession(vCtx), imcontext.Action_ServerPub, tmpPubMsg.MsgBody.Index, tmpPubMsg.MsgBody.Topic, len(tmpPubMsg.MsgBody.Data))
			logmanager.WriteConnectionLog(context.TODO(), &pbobjs.ConnectionLog{
				AppKey:   appkey,
				Session:  imcontext.GetConnSession(vCtx),
				Index:    tmpPubMsg.MsgBody.Index,
				Action:   string(imcontext.Action_ServerPub),
				Method:   tmpPubMsg.MsgBody.Topic,
				TargetId: tmpPubMsg.MsgBody.TargetId,
				DataLen:  int32(len(tmpPubMsg.MsgBody.Data)),
			})
			onlineStatus.Platforms = append(onlineStatus.Platforms, imcontext.GetPlatform(vCtx))
		}
	}
	if callback != nil {
		callback(onlineStatus)
	}
}

func PublishQryAckMessage(session string, qryAckMsg *codec.QueryAckMsgBody, callback func(), notOnlineCallback func()) {
	ctx := GetConnectCtxBySession(session)
	if ctx != nil {
		qos := codec.QoS_NoAck
		if callback != nil {
			qos = codec.QoS_NeedAck
			task := callbackTimeoutTimer.Add(20*time.Second, func() {
				//do timeout
				imcontext.RemoveQueryAckCallback(ctx, qryAckMsg.Index)
			})
			imcontext.PutQueryAckCallback(ctx, qryAckMsg.Index, func() {
				callbackTimeoutTimer.Remove(task)
				callback()
			})
		}
		tmpQryAckMsg := codec.NewQueryAckMessage(qryAckMsg, qos)
		ctx.Write(tmpQryAckMsg)
		logs.Infof("session:%s\taction:%s\tseq_index:%d\tcode:%d\tlen:%d", imcontext.GetConnSession(ctx), imcontext.Action_QueryAck, qryAckMsg.Index, qryAckMsg.Code, len(qryAckMsg.Data))
		logmanager.WriteConnectionLog(context.TODO(), &pbobjs.ConnectionLog{
			AppKey:  imcontext.GetAppkey(ctx),
			Session: imcontext.GetConnSession(ctx),
			Index:   qryAckMsg.Index,
			Action:  string(imcontext.Action_QueryAck),
			DataLen: int32(len(qryAckMsg.Data)),
			Code:    qryAckMsg.Code,
		})
	} else {
		if notOnlineCallback != nil {
			notOnlineCallback()
		}
	}
}

func PublishUserPubAckMessage(appkey, userid, session string, pubAckMsg *codec.PublishAckMsgBody) {
	ctx := GetConnectCtxBySession(session)
	if ctx != nil {
		tmpPubAckMsg := codec.NewUserPublishAckMessage(pubAckMsg)
		ctx.Write(tmpPubAckMsg)
		logs.Infof("session:%s\taction:%s\tseq_index:%d\tcode:%d", imcontext.GetConnSession(ctx), imcontext.Action_UserPubAck, pubAckMsg.Index, pubAckMsg.Code)
		logmanager.WriteConnectionLog(context.TODO(), &pbobjs.ConnectionLog{
			AppKey:  imcontext.GetAppkey(ctx),
			Session: imcontext.GetConnSession(ctx),
			Index:   pubAckMsg.Index,
			Action:  string(imcontext.Action_UserPubAck),
			Code:    pubAckMsg.Code,
		})
	}
}
