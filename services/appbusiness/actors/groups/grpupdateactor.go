package groups

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/services"

	"google.golang.org/protobuf/proto"
)

type GrpUpdateActor struct {
	bases.BaseActor
}

func (actor *GrpUpdateActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupInfo); ok {
		requesterId := bases.GetRequesterIdFromCtx(ctx)
		services.AppAsyncRpcCallWithSender(ctx, "upd_group_info", requesterId, req.GroupId, req, actor.Sender)
	}
}

func (actor *GrpUpdateActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupInfo{}
}
