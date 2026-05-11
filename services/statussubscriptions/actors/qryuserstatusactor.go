package actors

import (
	"context"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/statussubscriptions/services"

	"google.golang.org/protobuf/proto"
)

type QryUserStatusActor struct {
	bases.BaseActor
}

func (actor *QryUserStatusActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserOnlineStatusReq); ok {
		logs.WithContext(ctx).Infof("user_ids:%v", req.UserIds)
		code, resp := services.QryUserStatus(ctx, req)
		queryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illegal")
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *QryUserStatusActor) CreateInputObj() proto.Message {
	return &pbobjs.UserOnlineStatusReq{}
}
