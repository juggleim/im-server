package groups

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"

	"google.golang.org/protobuf/proto"
)

type GrpUpdateActor struct {
	bases.BaseActor
}

func (actor *GrpUpdateActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupInfo); ok {
		bases.AsyncRpcCallWithSender(ctx, "upd_group_info", req.GroupId, req, actor.Sender)
	}
}

func (actor *GrpUpdateActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupInfo{}
}
