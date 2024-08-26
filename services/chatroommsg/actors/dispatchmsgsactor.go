package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroommsg/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type DispatchMsgsActor struct {
	bases.BaseActor
}

func (actor DispatchMsgsActor) OnReceive(ctx context.Context, input proto.Message) {
	if msgSet, ok := input.(*pbobjs.DownMsgSet); ok {
		logs.WithContext(ctx).Infof("dispatch chat msgs.")
		services.HandleChatMsgsDispatch(ctx, msgSet)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
	}
}

func (actor DispatchMsgsActor) CreateInputObj() proto.Message {
	return &pbobjs.DownMsgSet{}
}
