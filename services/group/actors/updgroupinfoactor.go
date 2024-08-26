package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/group/services"

	"google.golang.org/protobuf/proto"
)

type UpdGroupInfoActor struct {
	bases.BaseActor
}

func (actor *UpdGroupInfoActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupInfo); ok {
		code := services.UpdGroupInfo(ctx, req)
		queryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *UpdGroupInfoActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupInfo{}
}
