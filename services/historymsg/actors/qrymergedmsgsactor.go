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

type QryMergedMsgsActor struct {
	bases.BaseActor
}

func (actor *QryMergedMsgsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryMergedMsgsReq); ok {
		parentMsgId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("parent_msgid:%s\tstart:%d\tcount:%d\torder:%d", parentMsgId, req.StartTime, req.Count, req.Order)
		code, msgs := services.QryMergedMsgs(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, msgs)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result:%d", len(msgs.Msgs))
	} else {
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *QryMergedMsgsActor) CreateInputObj() proto.Message {
	return &pbobjs.QryMergedMsgsReq{}
}
