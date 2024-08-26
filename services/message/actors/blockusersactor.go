package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type BlockUsersActor struct {
	bases.BaseActor
}

func (actor *BlockUsersActor) OnReceive(ctx context.Context, input proto.Message) {
	code := errs.IMErrorCode_SUCCESS
	if blockReq, ok := input.(*pbobjs.BlockUsersReq); ok {
		logs.WithContext(ctx).Infof("is_add:%v\tuser_ids:%v", blockReq.IsAdd, blockReq.UserIds)
		if len(blockReq.UserIds) > 0 {
			if blockReq.IsAdd {
				services.AddBlockUsers(ctx, blockReq.UserIds)
			} else {
				services.RemoveBlockUsers(ctx, blockReq.UserIds)
			}
		}
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
	}
	ack := bases.CreateQueryAckWraper(ctx, code, nil)
	actor.Sender.Tell(ack, actorsystem.NoSender)
}

func (actor *BlockUsersActor) CreateInputObj() proto.Message {
	return &pbobjs.BlockUsersReq{}
}
