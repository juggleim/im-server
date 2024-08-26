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

type ClearTotalUnreadActor struct {
	bases.BaseActor
}

func (actor *ClearTotalUnreadActor) OnReceive(ctx context.Context, input proto.Message) {
	if _, ok := input.(*pbobjs.QryTotalUnreadCountReq); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s", userId)
		code := services.ClearTotalUnread(ctx, userId)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result_code:%v", code)
	} else {
		logs.WithContext(ctx).Infof("input is illegal")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *ClearTotalUnreadActor) CreateInputObj() proto.Message {
	return &pbobjs.QryTotalUnreadCountReq{}
}
