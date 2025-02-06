package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type CheckBlockUserActor struct {
	bases.BaseActor
}

func (actor *CheckBlockUserActor) OnReceive(ctx context.Context, input proto.Message) {
	targetUserId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	appkey := bases.GetAppKeyFromCtx(ctx)
	blockUser := services.GetBlockUserItem(appkey, targetUserId, userId)
	ret := &pbobjs.CheckBlockUserResp{
		IsBlock: blockUser.IsBlock,
	}
	ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, ret)
	actor.Sender.Tell(ack, actorsystem.NoSender)
}

func (actor *CheckBlockUserActor) CreateInputObj() proto.Message {
	return &pbobjs.Nil{}
}
