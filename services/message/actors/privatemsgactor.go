package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type PrivateMsgActor struct {
	bases.BaseActor
}

func (actor *PrivateMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if upMsg, ok := input.(*pbobjs.UpMsg); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		exts := bases.GetExtsFromCtx(ctx)
		if receiverId, exist := exts[commonservices.RpcExtKey_RealTargetId]; exist {
			logs.WithContext(ctx).Infof("sender:%s\treceiver:%s\tflag:%d", userId, receiverId, upMsg.Flags)
			code, msgId, sendTime, msgSeq, clientMsgId, modifiedMsg := services.SendPrivateMsg(ctx, userId, receiverId, upMsg)
			userPubAck := bases.CreateUserPubAckWraper(ctx, code, msgId, sendTime, msgSeq, clientMsgId, modifiedMsg)
			actor.Sender.Tell(userPubAck, actorsystem.NoSender)
			logs.WithContext(ctx).Infof("code:%d", code)
		} else {
			logs.WithContext(ctx).Errorf("have no receiver")
		}
	} else {
		logs.WithContext(ctx).Errorf("upMsg is illigal. upMsg:%v", upMsg)
	}
}

func (actor *PrivateMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.UpMsg{}
}
