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

type CheckFriendActor struct {
	bases.BaseActor
}

func (actor *CheckFriendActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.CheckFriendsReq); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tfriend_ids:%v", userId, req.FriendIds)
		code, resp := services.CheckFriends(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *CheckFriendActor) CreateInputObj() proto.Message {
	return &pbobjs.CheckFriendsReq{}
}
