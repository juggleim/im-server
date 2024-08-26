package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/connectmanager/services"

	"google.golang.org/protobuf/proto"
)

type KickUserActor struct {
	bases.BaseActor
}

func (actor *KickUserActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.KickUserReq); ok {
		logs.WithContext(ctx).Infof("\tuser_id:%s\tplatform:%v", req.UserId, req.Platforms)
		services.KickUser(ctx, req)
	} else {
		logs.WithContext(ctx).Errorf("ban_users, input is illigal.")
	}
	ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, nil)
	actor.Sender.Tell(ack, actorsystem.NoSender)
}

func (actor *KickUserActor) CreateInputObj() proto.Message {
	return &pbobjs.KickUserReq{}
}
