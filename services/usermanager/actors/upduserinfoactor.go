package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"

	"google.golang.org/protobuf/proto"
)

type UpdUserInfoActor struct {
	bases.BaseActor
}

func (actor *UpdUserInfoActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserInfo); ok {
		code := services.UpdUserInfo(ctx, req)
		queryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *UpdUserInfoActor) CreateInputObj() proto.Message {
	return &pbobjs.UserInfo{}
}
