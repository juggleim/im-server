package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/friends/services"

	"google.golang.org/protobuf/proto"
)

type QryFriendsActor struct {
	bases.BaseActor
}

func (actor *QryFriendsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryFriendsReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\treq:%v", bases.GetTargetIdFromCtx(ctx), req)
		code, members := services.QryFriends(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, members)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *QryFriendsActor) CreateInputObj() proto.Message {
	return &pbobjs.QryFriendsReq{}
}
