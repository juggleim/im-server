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

type HangupActor struct {
	bases.BaseActor
}

func (actor *HangupActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.RtcAnswerReq); ok {
		logs.WithContext(ctx).Infof("room_id:%s\tuser_id:%s\target_id:%s", req.RoomId, bases.GetRequesterIdFromCtx(ctx), req.TargetId)
		code := services.RtcHangup(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *HangupActor) CreateInputObj() proto.Message {
	return &pbobjs.RtcAnswerReq{}
}
