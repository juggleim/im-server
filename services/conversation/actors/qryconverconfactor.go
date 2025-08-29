package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type QryUserConverConfActor struct {
	bases.BaseActor
}

func (actor *QryUserConverConfActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ConverIndex); ok {
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tchannel_type:%v\tsub_channel:%s", bases.GetRequesterIdFromCtx(ctx), req.TargetId, req.ChannelType, req.SubChannel)
		code, resp := services.QryConverConf(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *QryUserConverConfActor) CreateInputObj() proto.Message {
	return &pbobjs.ConverIndex{}
}

type QryGlobalConverConfActor struct {
	bases.BaseActor
}

func (actor *QryGlobalConverConfActor) OnReceive(ctx context.Context, input proto.Message) {

}

func (actor *QryGlobalConverConfActor) CreateInputObj() proto.Message {
	return &pbobjs.ConverIndex{}
}
