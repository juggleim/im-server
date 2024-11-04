package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type MarkUnreadActor struct {
	bases.BaseActor
}

func (actor *MarkUnreadActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ConversationsReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tconvers:%v", userId, req.Conversations)
		code := services.MarkUnreadV2(ctx, userId, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *MarkUnreadActor) CreateInputObj() proto.Message {
	return &pbobjs.ConversationsReq{}
}
