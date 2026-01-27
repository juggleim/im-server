package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type QryFriendInfoActor struct {
	bases.BaseActor
}

func (actor *QryFriendInfoActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.FriendIdsReq); ok {
		code, resp := services.QryFriendInfos(ctx, req)
		queryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryFriendInfoActor) CreateInputObj() proto.Message {
	return &pbobjs.FriendIdsReq{}
}
