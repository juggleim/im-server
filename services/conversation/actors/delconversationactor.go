package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type DelConversationsActor struct {
	bases.BaseActor
}

func (actor *DelConversationsActor) OnReceive(ctx context.Context, input proto.Message) {
	if delConverReq, ok := input.(*pbobjs.ConversationsReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)

		logs.WithContext(ctx).WithField("method", "del_convers").Infof("user_id:%s\tconversations=%v", userId, delConverReq.Conversations)

		code := services.DelConversationV2(ctx, userId, delConverReq.Conversations)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).WithField("method", "del_convers").Infof("user_id:%s\tinput is illegal", bases.GetRequesterIdFromCtx(ctx))
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_MSG_DEFAULT, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *DelConversationsActor) CreateInputObj() proto.Message {
	return &pbobjs.ConversationsReq{}
}
