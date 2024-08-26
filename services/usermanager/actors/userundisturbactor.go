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

type SetUserUndisturbActor struct {
	bases.BaseActor
}

func (actor *SetUserUndisturbActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserUndisturb); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		code := services.SetUserUndisturb(ctx, userId, req)
		queryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *SetUserUndisturbActor) CreateInputObj() proto.Message {
	return &pbobjs.UserUndisturb{}
}

type GetUserUndisturbActor struct {
	bases.BaseActor
}

func (actor *GetUserUndisturbActor) OnReceive(ctx context.Context, input proto.Message) {
	userId := bases.GetTargetIdFromCtx(ctx)
	resp, code := services.GetUserUndisturb(ctx, userId)
	queryAck := bases.CreateQueryAckWraper(ctx, code, resp)
	actor.Sender.Tell(queryAck, actorsystem.NoSender)
}

func (actor *GetUserUndisturbActor) CreateInputObj() proto.Message {
	return &pbobjs.Nil{}
}
