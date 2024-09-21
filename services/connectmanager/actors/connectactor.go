package actors

import (
	"context"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/services"
	"time"

	"google.golang.org/protobuf/proto"
)

type ConnectActor struct {
	actorsystem.UntypedActor
}

func (actor *ConnectActor) OnReceive(ctx context.Context, input proto.Message) {
	if rpcMsg, ok := input.(*pbobjs.RpcMessageWraper); ok {
		if rpcMsg.RpcMsgType == pbobjs.RpcMsgType_UserPubAck {
			services.PublishUserPubAckMessage(rpcMsg.AppKey, rpcMsg.RequesterId, rpcMsg.Session, &codec.PublishAckMsgBody{
				Index:       rpcMsg.ReqIndex,
				Code:        rpcMsg.ResultCode,
				MsgId:       rpcMsg.MsgId,
				Timestamp:   rpcMsg.MsgSendTime,
				MsgSeqNo:    rpcMsg.MsgSeqNo,
				MemberCount: rpcMsg.MemberCount,
			})
		} else if rpcMsg.RpcMsgType == pbobjs.RpcMsgType_QueryAck {
			var callback func()
			var ontOnlineCallback func()
			if int(rpcMsg.Qos) == codec.QoS_NeedAck || actor.Sender != actorsystem.NoSender {
				callback = func() {}
				ontOnlineCallback = func() {}
			}
			timestamp := rpcMsg.MsgSendTime
			if timestamp <= 0 {
				timestamp = time.Now().UnixMilli()
			}
			services.PublishQryAckMessage(rpcMsg.Session, &codec.QueryAckMsgBody{
				Index:     rpcMsg.ReqIndex,
				Code:      rpcMsg.ResultCode,
				Timestamp: timestamp,
				Data:      rpcMsg.AppDataBytes,
			}, callback, ontOnlineCallback)
		} else if rpcMsg.RpcMsgType == pbobjs.RpcMsgType_ServerPub {
			var onlineCallback func(*pbobjs.OnlineStatus)
			if int(rpcMsg.Qos) == codec.QoS_NeedAck || actor.Sender != actorsystem.NoSender {
				pubAckCallback := &PublishAckCallback{
					sender: actor.Sender,
				}
				onlineCallback = pubAckCallback.OnlineCallback
			}
			services.PublishServerPubMessage(rpcMsg.AppKey, rpcMsg.TargetId, rpcMsg.Session, rpcMsg.MsgId, rpcMsg.MsgSendTime, &codec.PublishMsgBody{
				Topic:     rpcMsg.Method,
				TargetId:  rpcMsg.TargetId,
				Timestamp: rpcMsg.MsgSendTime,
				Data:      rpcMsg.AppDataBytes,
			}, commonservices.PublishType(rpcMsg.PublishType), onlineCallback)
		}
	}
}

func (actor *ConnectActor) CreateInputObj() proto.Message {
	return &pbobjs.RpcMessageWraper{}
}

type PublishAckCallback struct {
	sender actorsystem.ActorRef
}

func (callback *PublishAckCallback) OnlineCallback(onlineStatus *pbobjs.OnlineStatus) {
	if onlineStatus != nil {
		data, _ := tools.PbMarshal(onlineStatus)
		callback.sender.TellAndNoSender(&pbobjs.RpcMessageWraper{
			RpcMsgType:   pbobjs.RpcMsgType_ServerPubAck,
			AppDataBytes: data,
		})
	}
}
