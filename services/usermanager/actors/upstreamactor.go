package actors

import (
	"context"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/friendcache"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"

	"google.golang.org/protobuf/proto"
)

type UpstreamActor struct {
	actorsystem.UntypedActor
}

func (actor *UpstreamActor) OnReceive(ctx context.Context, input proto.Message) {
	if rpcMsg, ok := input.(*pbobjs.RpcMessageWraper); ok {
		appkey := rpcMsg.AppKey
		userId := rpcMsg.TargetId
		targetId := rpcMsg.RequesterId
		realMethod := rpcMsg.ExtParams[commonservices.RpcExtKey_RealMethod]

		userInfo := &pbobjs.UserInfo{
			UserId: userId,
		}
		user, exist := services.GetUserInfo(appkey, userId)
		if exist && user != nil {
			userInfo.Nickname = user.Nickname
			userInfo.UserPortrait = user.UserPortrait
			userInfo.UserType = pbobjs.UserType(user.UserType)
			userInfo.UpdatedTime = user.UpdatedTime
			userInfo.ExtFields = commonservices.Map2KvItems(user.ExtFields)
			//check private global mute
			if realMethod == "p_msg" && user.CheckPrivateGlobalMute() {
				userPubAck := &pbobjs.RpcMessageWraper{
					RpcMsgType:   pbobjs.RpcMsgType_UserPubAck,
					ResultCode:   int32(errs.IMErrorCode_MSG_BLOCK),
					MsgId:        "",
					MsgSendTime:  0,
					MsgSeqNo:     0,
					ReqIndex:     rpcMsg.ReqIndex,
					AppKey:       appkey,
					Qos:          rpcMsg.Qos,
					Session:      rpcMsg.Session,
					Method:       rpcMsg.Method,
					SourceMethod: rpcMsg.SourceMethod,
					RequesterId:  rpcMsg.RequesterId,
					TargetId:     rpcMsg.TargetId,
					PublishType:  rpcMsg.PublishType,
				}
				actor.Sender.Tell(userPubAck, actorsystem.NoSender)
				return
			}
		}

		//friend info
		senderFriendInfo := &pbobjs.FriendInfo{}
		if realMethod == "p_msg" {
			if appinfo, exist := commonservices.GetAppInfo(appkey); exist && appinfo != nil && appinfo.OpenRemark {
				friStatus := friendcache.GetFriendStatus(appkey, userId, targetId)
				if friStatus != nil && friStatus.IsFriend {
					senderFriendInfo.IsFriend = true
					senderFriendInfo.FriendDisplayName = friStatus.FriendDisplayName
					senderFriendInfo.UpdatedTime = friStatus.UpdatedTime
				}
			}
		}

		exts := map[string]string{}
		if realMethod == "p_msg" || realMethod == "imp_pri_msg" {
			exts[commonservices.RpcExtKey_RealTargetId] = targetId
			targetId = commonservices.GetConversationId(userId, targetId, pbobjs.ChannelType_Private)
		} else if realMethod == "s_msg" {
			exts[commonservices.RpcExtKey_RealTargetId] = targetId
			targetId = commonservices.GetConversationId(userId, targetId, pbobjs.ChannelType_System)
		}
		isSucc := bases.UnicastRouteWithSenderActor(&pbobjs.RpcMessageWraper{
			RpcMsgType:       pbobjs.RpcMsgType_UserPub,
			AppKey:           appkey,
			Session:          rpcMsg.Session,
			DeviceId:         rpcMsg.DeviceId,
			InstanceId:       rpcMsg.InstanceId,
			Platform:         rpcMsg.Platform,
			Method:           realMethod,
			RequesterId:      userId,
			ReqIndex:         rpcMsg.ReqIndex,
			Qos:              rpcMsg.Qos,
			AppDataBytes:     rpcMsg.AppDataBytes,
			TargetId:         targetId,
			SenderInfo:       userInfo,
			SenderFriendInfo: senderFriendInfo,
			ExtParams:        exts,
			IsFromApi:        rpcMsg.IsFromApi,
			NoSendbox:        rpcMsg.NoSendbox,
			OnlySendbox:      rpcMsg.OnlySendbox,
			MsgId:            rpcMsg.MsgId,
		}, actor.Sender)
		if !isSucc {
			userPubAck := &pbobjs.RpcMessageWraper{
				RpcMsgType:   pbobjs.RpcMsgType_UserPubAck,
				ResultCode:   int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC),
				MsgId:        "",
				MsgSendTime:  0,
				MsgSeqNo:     0,
				ReqIndex:     rpcMsg.ReqIndex,
				AppKey:       appkey,
				Qos:          rpcMsg.Qos,
				Session:      rpcMsg.Session,
				Method:       rpcMsg.Method,
				SourceMethod: rpcMsg.SourceMethod,
				RequesterId:  rpcMsg.RequesterId,
				TargetId:     rpcMsg.TargetId,
				PublishType:  rpcMsg.PublishType,
			}
			actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		}
	} else {
		logs.WithContext(ctx).Error("input illegal")
	}
}

func (actor *UpstreamActor) CreateInputObj() proto.Message {
	return &pbobjs.RpcMessageWraper{}
}
