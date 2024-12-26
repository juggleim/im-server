package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type UserPushActor struct {
	bases.BaseActor
}

func (actor *UserPushActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.DownMsg); ok {
		services.UserPush(ctx, req)
	}
}

func (actor *UserPushActor) CreateInputObj() proto.Message {
	return &pbobjs.DownMsg{}
}
