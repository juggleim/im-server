package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type BrdcastInboxActor struct {
	bases.BaseActor
}

func (actor *BrdcastInboxActor) OnReceive(ctx context.Context, input proto.Message) {
	if msg, ok := input.(*pbobjs.DownMsg); ok {
		logs.WithContext(ctx).Infof("msg_type:%s\tmsg_id:%s\tmsg_time:%d", msg.MsgType, msg.MsgId, msg.MsgTime)
		code := services.SaveBroadcastMsg(ctx, msg)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("input is illigal.")
	}
}

func (actor *BrdcastInboxActor) CreateInputObj() proto.Message {
	return &pbobjs.DownMsg{}
}
