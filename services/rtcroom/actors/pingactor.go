package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/rtcroom/services"

	"google.golang.org/protobuf/proto"
)

type PingActor struct {
	bases.BaseActor
}

func (actor *PingActor) OnReceive(ctx context.Context, input proto.Message) {
	code := services.RtcPing(ctx)
	qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
	actor.Sender.Tell(qryAck, actorsystem.NoSender)
}

func (actor *PingActor) CreateInputObj() proto.Message {
	return &pbobjs.Nil{}
}
