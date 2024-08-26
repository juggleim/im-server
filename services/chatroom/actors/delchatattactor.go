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

type DelAttActor struct {
	bases.BaseActor
}

func (actor *DelAttActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatAttReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tchatroom_id:%s\tkey:%s", bases.GetRequesterIdFromCtx(ctx), bases.GetTargetIdFromCtx(ctx), req.Key)
		code, resp := services.DelChatAtt(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *DelAttActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatAttReq{}
}

type BatchDelAttActor struct {
	bases.BaseActor
}

func (actor *BatchDelAttActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatAttBatchReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tchatroom_id:%s\tkvs:%v", bases.GetRequesterIdFromCtx(ctx), bases.GetTargetIdFromCtx(ctx), req.Atts)
		code, resp := services.BatchDelChatAtt(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *BatchDelAttActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatAttBatchReq{}
}
