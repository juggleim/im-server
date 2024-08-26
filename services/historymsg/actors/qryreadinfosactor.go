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

type QryReadInfosActor struct {
	bases.BaseActor
}

func (actor *QryReadInfosActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryReadInfosReq); ok {
		logs.WithContext(ctx).Infof("target_id:%s\tchannel_type:%v\tmsg_ids:%v", req.TargetId, req.ChannelType, req.MsgIds)
		code, resp := services.QryReadInfos(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result:%d", len(resp.Items))
	} else {
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Error("failed to decode input")
	}
}

func (actor *QryReadInfosActor) CreateInputObj() proto.Message {
	return &pbobjs.QryReadInfosReq{}
}
