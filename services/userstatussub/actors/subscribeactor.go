package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/userstatussub/services"

	"google.golang.org/protobuf/proto"
)

type SubscribeActor struct {
	bases.BaseActor
}

func (actor *SubscribeActor) OnReceive(ctx context.Context, input proto.Message) {
	if uIds, ok := input.(*pbobjs.UserIdsReq); ok {
		if len(uIds.UserIds) > 0 {
			bases.GroupRpcCall(ctx, "inner_sub_users", uIds.UserIds, &pbobjs.UserIdsReq{})
		}
		queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, nil)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)

	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *SubscribeActor) CreateInputObj() proto.Message {
	return &pbobjs.UserIdsReq{}
}

type InnerSubscribeActor struct {
	bases.BaseActor
}

func (actor *InnerSubscribeActor) OnReceive(ctx context.Context, input proto.Message) {
	if _, ok := input.(*pbobjs.UserIdsReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		targetIds := bases.GetTargetIdsFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_ids:%v", userId, targetIds)
		code := services.Subscribe(ctx, userId, targetIds)
		queryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}
func (actor *InnerSubscribeActor) CreateInputObj() proto.Message {
	return &pbobjs.UserIdsReq{}
}
