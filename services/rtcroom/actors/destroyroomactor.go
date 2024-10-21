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

type DestroyRoomActor struct {
	bases.BaseActor
}

func (actor *DestroyRoomActor) OnReceive(ctx context.Context, input proto.Message) {
	logs.WithContext(ctx).Infof("room_id:%s\tuser_id:%s", bases.GetTargetIdFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx))
	code := services.DestroyRtcRoom(ctx)
	qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
	actor.Sender.Tell(qryAck, actorsystem.NoSender)
}

func (actor *DestroyRoomActor) CreateInputObj() proto.Message {
	return &pbobjs.Nil{}
}
