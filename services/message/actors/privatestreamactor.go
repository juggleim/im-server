package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type PrivateStreamActor struct {
	bases.BaseActor
}

func (actor *PrivateStreamActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.StreamDownMsg); ok {
		services.HandleStreamMsg(ctx, req)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
	}
}

func (actor *PrivateStreamActor) CreateInputObj() proto.Message {
	return &pbobjs.StreamDownMsg{}
}
