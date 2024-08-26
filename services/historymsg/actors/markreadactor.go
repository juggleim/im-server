package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type MarkReadActor struct {
	bases.BaseActor
}

func (actor *MarkReadActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.MarkReadReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		targetId := bases.GetTargetIdFromCtx(ctx)

		logs.WithContext(ctx).Infof("user_id:%s\ttargetId=%s\tchannel_type=%v\tmsg_len=%d\tscope:%v", userId, targetId, req.ChannelType, len(req.Msgs), req.IndexScopes)

		code := services.MarkRead(ctx, userId, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result_code:%v", code)
	} else {
		logs.WithContext(ctx).Infof("user_id:%s\tinput is illegal", bases.GetRequesterIdFromCtx(ctx))
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_MSG_DEFAULT, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *MarkReadActor) CreateInputObj() proto.Message {
	return &pbobjs.MarkReadReq{}
}
