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

type QryRtcMemberRoomsActor struct {
	bases.BaseActor
}

func (actor *QryRtcMemberRoomsActor) OnReceive(ctx context.Context, input proto.Message) {
	logs.WithContext(ctx).Infof("user_id:%s", bases.GetTargetIdFromCtx(ctx))
	code, resp := services.QryRtcMemberRooms(ctx)
	qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
	actor.Sender.Tell(qryAck, actorsystem.NoSender)
}

func (actor *QryRtcMemberRoomsActor) CreateInputObj() proto.Message {
	return &pbobjs.Nil{}
}
