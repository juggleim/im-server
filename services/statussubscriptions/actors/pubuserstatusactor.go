package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/statussubscriptions/services"

	"google.golang.org/protobuf/proto"
)

type PubUserStatusActor struct {
	bases.BaseActor
}

func (actor *PubUserStatusActor) OnReceive(ctx context.Context, input proto.Message) {
	if upMsg, ok := input.(*pbobjs.UpMsg); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tmsg_type:%s", userId, upMsg.MsgType)
		services.PublishUserStatus(ctx, upMsg)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *PubUserStatusActor) CreateInputObj() proto.Message {
	return &pbobjs.UpMsg{}
}
