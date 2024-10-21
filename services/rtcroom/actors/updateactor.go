package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/rtcroom/services"

	"google.golang.org/protobuf/proto"
)

type UpdateActor struct {
	bases.BaseActor
}

func (actor *UpdateActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.RtcMember); ok {
		logs.WithContext(ctx).Infof("room_id:%s\tuser_id:%s\tmember:[%v]", bases.GetTargetIdFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx), req)
		code := services.UpdRtcMemberState(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
	}
}

func (actor *UpdateActor) CreateInputObj() proto.Message {
	return &pbobjs.RtcMember{}
}
