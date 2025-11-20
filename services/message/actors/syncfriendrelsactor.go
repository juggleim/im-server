package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type SyncFriendRelsActor struct {
	bases.BaseActor
}

func (actor *SyncFriendRelsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.FriendIdsReq); ok {
		services.SyncFriendRels(ctx, req)
	}
}

func (actor *SyncFriendRelsActor) CreateInputObj() proto.Message {
	return &pbobjs.FriendIdsReq{}
}
