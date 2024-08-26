package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroommsg/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type DispatchChrmActor struct {
	bases.BaseActor
}

func (actor *DispatchChrmActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChrmDispatchReq); ok {
		logs.WithContext(ctx).Infof("chatroom_id:%s", req.ChatId)
		services.HandleChrmDispatch(ctx, req)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *DispatchChrmActor) CreateInputObj() proto.Message {
	return &pbobjs.ChrmDispatchReq{}
}
