package actors

import (
	"context"
	"im-server/commons/bases"
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
		services.ReportClientLogState(ctx, req)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
	}
}

func (actor ReportClientLogStateActor) CreateInputObj() proto.Message {
	return &pbobjs.UploadLogStatusReq{}
}
