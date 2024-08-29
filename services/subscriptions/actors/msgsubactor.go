package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/subscriptions/services"

	"google.golang.org/protobuf/proto"
)

type MsgSubActor struct {
	bases.BaseActor
}

func (actor *MsgSubActor) OnReceive(ctx context.Context, input proto.Message) {
	if subMsgs, ok := input.(*pbobjs.SubMsgs); ok {
		services.MsgSubHandle(ctx, subMsgs)
	}
}

func (actor *MsgSubActor) CreateInputObj() proto.Message {
	return &pbobjs.SubMsgs{}
}
