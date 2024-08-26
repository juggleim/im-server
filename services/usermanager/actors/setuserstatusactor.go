package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"

	"google.golang.org/protobuf/proto"
)

type SetUserStatusActor struct {
	bases.BaseActor
}

func (actor *SetUserStatusActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserInfo); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		code := services.SetUserStatus(ctx, userId, req)
		queryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *SetUserStatusActor) CreateInputObj() proto.Message {
	return &pbobjs.UserInfo{}
}
