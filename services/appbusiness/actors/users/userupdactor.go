package users

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"

	"google.golang.org/protobuf/proto"
)

type UserUpdateActor struct {
	bases.BaseActor
}

func (actor *UserUpdateActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserInfo); ok {
		bases.AsyncRpcCallWithSender(ctx, "upd_user_info", req.UserId, req, actor.Sender)
	}
}

func (actor *UserUpdateActor) CreateInputObj() proto.Message {
	return &pbobjs.UserInfo{}
}
