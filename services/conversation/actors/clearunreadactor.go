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

type ClearUnReadActor struct {
	bases.BaseActor
}

func (actor *ClearUnReadActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ClearUnreadReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		targetId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tconversations:%v", userId, targetId, req.Conversations)
		code := services.ClearUnread(ctx, targetId, req.Conversations, req.NoCmdMsg)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result_code:%v", code)
	} else {
		logs.WithContext(ctx).Infof("user_id:%s\tinput is illegal", bases.GetRequesterIdFromCtx(ctx))
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *ClearUnReadActor) CreateInputObj() proto.Message {
	return &pbobjs.ClearUnreadReq{}
}
