package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type AddAttActor struct {
	bases.BaseActor
}

func (actor *AddAttActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatAttReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tchatroom_id:%s\tkey:%s", bases.GetRequesterIdFromCtx(ctx), bases.GetTargetIdFromCtx(ctx), req.Key)
		code, resp := services.AddChatAtt(ctx, req)
		pubAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(pubAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *AddAttActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatAttReq{}
}

type BatchAddAttActor struct {
	bases.BaseActor
}

func (actor *BatchAddAttActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatAttBatchReq); ok {
		logs.WithContext(ctx).Infof("user_id:%schatroom_id:%s\tkvs:%v", bases.GetRequesterIdFromCtx(ctx), bases.GetTargetIdFromCtx(ctx), req.Atts)
		code, resp := services.BatchAddChatAtt(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *BatchAddAttActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatAttBatchReq{}
}
