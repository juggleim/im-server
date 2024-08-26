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

type QryUploadTokenActor struct {
	bases.BaseActor
}

func (actor QryUploadTokenActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryFileCredReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tfile_type:%v", bases.GetRequesterIdFromCtx(ctx), req.FileType)
		code, resp := services.GetFileCred(ctx, req)
		qrAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qrAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor QryUploadTokenActor) CreateInputObj() proto.Message {
	return &pbobjs.QryFileCredReq{}
}
