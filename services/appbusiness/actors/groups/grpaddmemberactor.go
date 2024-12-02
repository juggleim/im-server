package groups

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/services"

	"google.golang.org/protobuf/proto"
)

type GrpAddMemberActor struct {
	bases.BaseActor
}

func (actor *GrpAddMemberActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupMembersReq); ok {
		code := services.AddGrpMembers(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *GrpAddMemberActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupMembersReq{}
}
