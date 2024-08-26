package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type UpdStreamActor struct {
	bases.BaseActor
}

func (actor *UpdStreamActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.StreamDownMsg); ok {
		services.UpdStreamMsg(ctx, req)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *UpdStreamActor) CreateInputObj() proto.Message {
	return &pbobjs.StreamDownMsg{}
}
