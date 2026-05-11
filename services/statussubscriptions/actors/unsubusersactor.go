package actors

import (
	"context"

	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/statussubscriptions/services"

	"google.golang.org/protobuf/proto"
)

type UnSubUsersActor struct {
	bases.BaseActor
}

func (actor *UnSubUsersActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.SubUsersReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		deviceId := bases.GetDeviceIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tdevice:%s\tuids:%v", userId, deviceId, req.UserIds)
		code := services.UnSubUsers(ctx, req)
		queryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *UnSubUsersActor) CreateInputObj() proto.Message {
	return &pbobjs.SubUsersReq{}
}
