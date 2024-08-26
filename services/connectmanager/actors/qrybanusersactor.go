package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/connectmanager/services"

	"google.golang.org/protobuf/proto"
)

type QryBanUsersActor struct {
	bases.BaseActor
}

func (actor *QryBanUsersActor) OnReceive(ctx context.Context, input proto.Message) {
	ret := &pbobjs.QryBanUsersResp{
		Items: []*pbobjs.BanUser{},
	}
	code := errs.IMErrorCode_SUCCESS
	if banReq, ok := input.(*pbobjs.QryBanUsersReq); ok {
		logs.WithContext(ctx).WithField("method", "qry_ban_users").Infof("limit:%d\toffset:%s", banReq.Limit, banReq.Offset)
		code, ret.Items, ret.Offset = services.QryBanUsers(ctx, banReq.Limit, banReq.Offset)
	} else {
		logs.WithContext(ctx).Errorf("qry_ban_users, input is illigal.")
	}
	ack := bases.CreateQueryAckWraper(ctx, code, ret)
	actor.Sender.Tell(ack, actorsystem.NoSender)
	logs.WithContext(ctx).WithField("method", "qry_ban_users").Infof("result_len:%d\toffset:%s", len(ret.Items), ret.Offset)
}

func (actor *QryBanUsersActor) CreateInputObj() proto.Message {
	return &pbobjs.QryBanUsersReq{}
}
