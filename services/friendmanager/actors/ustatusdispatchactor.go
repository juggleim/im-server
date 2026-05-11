package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/friendmanager/services"

	"google.golang.org/protobuf/proto"
)

type UserStatusDispatchActor struct {
	bases.BaseActor
}

func (actor *UserStatusDispatchActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserStatusFriDispatch); ok {
		services.DispatchUserStatus(ctx, req)
	}
}

func (actor *UserStatusDispatchActor) CreateInputObj() proto.Message {
	return &pbobjs.UserStatusFriDispatch{}
}
