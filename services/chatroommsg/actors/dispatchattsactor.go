package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroommsg/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type DispatchAttsActor struct {
	bases.BaseActor
}

func (actor DispatchAttsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatAtts); ok {
		logs.WithContext(ctx).Infof("dispatch chat atts. chat_id:%s", req.ChatId)
		services.HandleChatAttsDispatch(ctx, req)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
	}
}

func (actor DispatchAttsActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatAtts{}
}
