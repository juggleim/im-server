package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type MarkGrpMsgReadActor struct {
	bases.BaseActor
}

func (actor *MarkGrpMsgReadActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.MarkGrpMsgReadReq); ok {
		services.MarkGrpMsgRead(ctx, req)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *MarkGrpMsgReadActor) CreateInputObj() proto.Message {
	return &pbobjs.MarkGrpMsgReadReq{}
}
