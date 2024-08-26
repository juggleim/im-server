package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/pushmanager/services"

	"google.golang.org/protobuf/proto"
)

type RegPushTokenActor struct {
	bases.BaseActor
}

func (actor *RegPushTokenActor) OnReceive(ctx context.Context, input proto.Message) {
	method := bases.GetMethodFromCtx(ctx)
	if req, ok := input.(*pbobjs.RegPushTokenReq); ok {
		targetId := bases.GetTargetIdFromCtx(ctx)
		services.RegPushToken(ctx, targetId, req)
		logs.WithContext(ctx).Infof("target_id:%s\treq:%v", targetId, req)
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).WithField("method", method).Infof("input is illegal")
	}
}

func (actor *RegPushTokenActor) CreateInputObj() proto.Message {
	return &pbobjs.RegPushTokenReq{}
}
