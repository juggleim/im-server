package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/friendmanager/services"

	"google.golang.org/protobuf/proto"
)

type AddFriendActor struct {
	bases.BaseActor
}

func (actor *AddFriendActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.FriendMembersReq); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tfriends:%v", userId, req.FriendMembers)
		code := services.AddFriends(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *AddFriendActor) CreateInputObj() proto.Message {
	return &pbobjs.FriendMembersReq{}
}
