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

type QryBlockUsersActor struct {
	bases.BaseActor
}

func (actor *QryBlockUsersActor) OnReceive(ctx context.Context, input proto.Message) {
	ret := &pbobjs.QryBlockUsersResp{
		Items: []*pbobjs.BlockUser{},
	}
	code := errs.IMErrorCode_SUCCESS
	if qryBlockUsersReq, ok := input.(*pbobjs.QryBlockUsersReq); ok {
		logs.WithContext(ctx).Infof("limit:%d\toffset:%s", qryBlockUsersReq.Limit, qryBlockUsersReq.Offset)
		code, ret.Items, ret.Offset = services.QryBlockUsers(ctx, qryBlockUsersReq.UserId, qryBlockUsersReq.Limit, qryBlockUsersReq.Offset)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
	}
	ack := bases.CreateQueryAckWraper(ctx, code, ret)
	actor.Sender.Tell(ack, actorsystem.NoSender)
	logs.WithContext(ctx).Infof("result_len:%d\toffset:%s", len(ret.Items), ret.Offset)
}

func (actor *QryBlockUsersActor) CreateInputObj() proto.Message {
	return &pbobjs.QryBlockUsersReq{}
}
