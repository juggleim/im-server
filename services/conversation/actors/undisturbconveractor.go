package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type UndisturbConversActor struct {
	bases.BaseActor
}

func (actor *UndisturbConversActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UndisturbConversReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tconvers=%v", bases.GetTargetIdFromCtx(ctx), req.Items)
		code := services.DoUndisturbConvers(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Info("input is illegal")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_MSG_DEFAULT, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *UndisturbConversActor) CreateInputObj() proto.Message {
	return &pbobjs.UndisturbConversReq{}
}
