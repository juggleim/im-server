package friends

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/services"

	"google.golang.org/protobuf/proto"
)

type FriendAddActor struct {
	bases.BaseActor
}

func (actor *FriendAddActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.FriendsAddReq); ok {
		code := services.AddFriends(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *FriendAddActor) CreateInputObj() proto.Message {
	return &pbobjs.FriendsAddReq{}
}
