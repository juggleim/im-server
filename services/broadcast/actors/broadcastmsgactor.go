package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/broadcast/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type BroadcastMsgActor struct {
	bases.BaseActor
}

func (actor *BroadcastMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if upMsg, ok := input.(*pbobjs.UpMsg); ok {
		logs.WithContext(ctx).Infof("sender_id:%s\tmsg_type:%s", bases.GetRequesterIdFromCtx(ctx), upMsg.MsgType)
		code, msgId, sendTime, msgSeq := services.BroadcastMsg(ctx, upMsg)
		userPubAck := bases.CreateUserPubAckWraper(ctx, code, msgId, sendTime, msgSeq, "")
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *BroadcastMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.UpMsg{}
}
