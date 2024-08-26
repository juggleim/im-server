package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/group/services"

	"google.golang.org/protobuf/proto"
)

type UpdGrpConverActor struct {
	bases.BaseActor
}

func (actor *UpdGrpConverActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UpdLatestMsgReq); ok {
		services.Dispatch2Conversation(ctx, bases.GetTargetIdFromCtx(ctx), req)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *UpdGrpConverActor) CreateInputObj() proto.Message {
	return &pbobjs.UpdLatestMsgReq{}
}
