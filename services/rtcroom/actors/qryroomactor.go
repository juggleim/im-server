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

type QryRtcRoomActor struct {
	bases.BaseActor
}

func (actor *QryRtcRoomActor) OnReceive(ctx context.Context, input proto.Message) {
	logs.WithContext(ctx).Infof("room_id:%s", bases.GetTargetIdFromCtx(ctx))
	code, resp := services.QryRtcRoom(ctx, bases.GetTargetIdFromCtx(ctx))
	qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
	actor.Sender.Tell(qryAck, actorsystem.NoSender)
}

func (actor *QryRtcRoomActor) CreateInputObj() proto.Message {
	return &pbobjs.Nil{}
}
