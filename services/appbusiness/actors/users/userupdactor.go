package users

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/services"

	"google.golang.org/protobuf/proto"
)

type UserUpdateActor struct {
	bases.BaseActor
}

func (actor *UserUpdateActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserInfo); ok {
		requesterId := bases.GetRequesterIdFromCtx(ctx)
		services.AppAsyncRpcCallWithSender(ctx, "upd_user_info", requesterId, req.UserId, req, actor.Sender)
	}
}

func (actor *UserUpdateActor) CreateInputObj() proto.Message {
	return &pbobjs.UserInfo{}
}
