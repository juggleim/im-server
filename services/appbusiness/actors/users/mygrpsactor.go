package users

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/services"

	"google.golang.org/protobuf/proto"
)

type MyGroupsActor struct {
	bases.BaseActor
}

func (actor *MyGroupsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupInfoListReq); ok {
		code, resp := services.QueryMyGroups(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *MyGroupsActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupInfoListReq{}
}
