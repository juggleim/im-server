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

type TopConversActor struct {
	bases.BaseActor
}

func (actor *TopConversActor) OnReceive(ctx context.Context, input proto.Message) {
	if conversReq, ok := input.(*pbobjs.ConversationsReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tconversations=%v", userId, conversReq.Conversations)
		code, cmdMsgTime := services.SetTopConvers(ctx, conversReq)
		qryAck := bases.CreateQueryAckWraperWithTime(ctx, code, cmdMsgTime, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Infof("user_id:%s\tinput is illegal", bases.GetRequesterIdFromCtx(ctx))
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_MSG_DEFAULT, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *TopConversActor) CreateInputObj() proto.Message {
	return &pbobjs.ConversationsReq{}
}
