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

type InviteActor struct {
	bases.BaseActor
}

func (actor *InviteActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.RtcInviteReq); ok {
		logs.WithContext(ctx).Infof("room_id:%s\troom_type:%v\tuser_id:%s\ttargets:%v", req.RoomId, req.RoomType, bases.GetRequesterIdFromCtx(ctx), req.TargetIds)
		code := services.RtcInvite(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
	}
}

func (actor *InviteActor) CreateInputObj() proto.Message {
	return &pbobjs.RtcInviteReq{}
}
