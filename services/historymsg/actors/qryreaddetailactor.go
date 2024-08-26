package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type QryReadDetailActor struct {
	bases.BaseActor
}

func (actor *QryReadDetailActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryReadDetailReq); ok {
		logs.WithContext(ctx).Infof("target_id:%s\tchannel_type:%v\tmsg_id:%v", req.TargetId, req.ChannelType, req.MsgId)
		code, resp := services.QryReadDetail(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Error("failed to decode input")
	}
}

func (actor *QryReadDetailActor) CreateInputObj() proto.Message {
	return &pbobjs.QryReadDetailReq{}
}
