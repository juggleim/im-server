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

type QryChatroomActor struct {
	bases.BaseActor
}

func (actor *QryChatroomActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatroomReq); ok {
		chatId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("chat_id:%s\tcount:%d\torder:%d", chatId, req.Count, req.Order)
		code, resp := services.QryChatroomInfo(ctx, chatId, int(req.Count), int(req.Order))
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *QryChatroomActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatroomReq{}
}
