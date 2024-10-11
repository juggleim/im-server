package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/fileplugin/services"

	"google.golang.org/protobuf/proto"
)

type ReportClientLogStateActor struct {
	bases.BaseActor
}

func (actor ReportClientLogStateActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UploadLogStatusReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tmsg_id:%s\tstate:%d", bases.GetRequesterIdFromCtx(ctx), req.MsgId, req.State)
		code := services.ReportClientLogState(ctx, req)
		qrAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qrAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor ReportClientLogStateActor) CreateInputObj() proto.Message {
	return &pbobjs.UploadLogStatusReq{}
}
