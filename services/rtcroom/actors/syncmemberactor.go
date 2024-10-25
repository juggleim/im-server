package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/rtcroom/services"

	"google.golang.org/protobuf/proto"
)

type SyncMemberActor struct {
	bases.BaseActor
}

func (actor *SyncMemberActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.SyncMemberStateReq); ok {
		code := services.SyncMemberState(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *SyncMemberActor) CreateInputObj() proto.Message {
	return &pbobjs.SyncMemberStateReq{}
}
