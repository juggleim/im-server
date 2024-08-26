package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/pushmanager/services"

	"google.golang.org/protobuf/proto"
)

type PushActor struct {
	bases.BaseActor
}

func (actor *PushActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.PushData); ok {
		services.SendPush(ctx, bases.GetTargetIdFromCtx(ctx), req)
		logs.WithContext(ctx).Infof("user_id:%s\t%v", bases.GetTargetIdFromCtx(ctx), req)
	} else {
		logs.WithContext(ctx).Infof("input is illegal")
	}
}

func (actor *PushActor) CreateInputObj() proto.Message {
	return &pbobjs.PushData{}
}
