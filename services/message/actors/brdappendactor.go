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

type BrdAppendActor struct {
	bases.BaseActor
}

func (actor *BrdAppendActor) OnReceive(ctx context.Context, input proto.Message) {
	if msg, ok := input.(*pbobjs.DownMsg); ok {
		code := services.BrdAppendMsg(ctx, msg)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("input is illigal.")
	}
}

func (actor *BrdAppendActor) CreateInputObj() proto.Message {
	return &pbobjs.DownMsg{}
}
