package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type PushSwitchActor struct {
	bases.BaseActor
}

func (actor *PushSwitchActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.PushSwitch); ok {
		appkey := bases.GetAppKeyFromCtx(ctx)
		userId := bases.GetTargetIdFromCtx(ctx)
		services.SetPushSwitch(appkey, userId, req.Switch)
	}
}

func (actor *PushSwitchActor) CreateInputObj() proto.Message {
	return &pbobjs.PushSwitch{}
}
