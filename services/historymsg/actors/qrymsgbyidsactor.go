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

type QryHistoryMsgByIdsActor struct {
	bases.BaseActor
}

func (actor *QryHistoryMsgByIdsActor) OnReceive(ctx context.Context, input proto.Message) {
	if qryHisMsgsReq, ok := input.(*pbobjs.QryHisMsgByIdsReq); ok {
		logs.WithContext(ctx).Infof("target_id:%s\tchannel_type:%v\tmsg_ids:%v", qryHisMsgsReq.TargetId, qryHisMsgsReq.ChannelType, qryHisMsgsReq.MsgIds)
		resp := services.QryHisMsgByIds(ctx, qryHisMsgsReq)
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result:%d", len(resp.Msgs))
	} else {
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Error("failed to decode input")
	}
}

func (actor *QryHistoryMsgByIdsActor) CreateInputObj() proto.Message {
	return &pbobjs.QryHisMsgByIdsReq{}
}
