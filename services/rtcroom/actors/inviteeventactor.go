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

type InviteEventActor struct {
	bases.BaseActor
}

func (actor *InviteEventActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.RtcInviteReq); ok {
		logs.WithContext(ctx).Infof("room_id:%s\tfrom_user_id:%s\tinvite_type:%v\ttargets:%v", bases.GetTargetIdFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx), req.InviteType, req.TargetIds)
		code := services.HandleInviteEvent(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
	}
}

func (actor *InviteEventActor) CreateInputObj() proto.Message {
	return &pbobjs.RtcInviteReq{}
}
