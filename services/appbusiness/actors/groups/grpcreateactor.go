package groups

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/services"

	"google.golang.org/protobuf/proto"
)

type GrpCreateActor struct {
	bases.BaseActor
}

func (actor *GrpCreateActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupMembersReq); ok {
		code, grpInfo := services.CreateGroup(ctx, req)

		ack := bases.CreateQueryAckWraper(ctx, code, grpInfo)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *GrpCreateActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupMembersReq{}
}
