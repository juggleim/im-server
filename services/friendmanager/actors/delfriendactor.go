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

type DelFriendActor struct {
	bases.BaseActor
}

func (actor *DelFriendActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.FriendIdsReq); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tfriend_ids:%v", userId, req.FriendIds)
		code := services.DelFriends(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *DelFriendActor) CreateInputObj() proto.Message {
	return &pbobjs.FriendIdsReq{}
}
