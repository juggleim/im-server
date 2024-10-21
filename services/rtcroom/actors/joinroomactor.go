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

type JoinRoomActor struct {
	bases.BaseActor
}

func (actor *JoinRoomActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.RtcRoomReq); ok {
		logs.WithContext(ctx).Infof("room_id:%s\tuser_id:%s", bases.GetTargetIdFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx))
		code, resp := services.JoinRtcRoom(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
	}
}

func (actor *JoinRoomActor) CreateInputObj() proto.Message {
	return &pbobjs.RtcRoomReq{}
}
