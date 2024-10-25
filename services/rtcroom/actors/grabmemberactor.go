package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/rtcroom/services"

	"google.golang.org/protobuf/proto"
)

type GrabMemberActor struct {
	bases.BaseActor
}

func (actor *GrabMemberActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.MemberState); ok {
		code := services.GrabMemberState(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *GrabMemberActor) CreateInputObj() proto.Message {
	return &pbobjs.MemberState{}
}
