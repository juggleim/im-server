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

type JoinChatroomActor struct {
	bases.BaseActor
}

func (actor *JoinChatroomActor) OnReceive(ctx context.Context, input proto.Message) {
	if chatroom, ok := input.(*pbobjs.ChatroomInfo); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tchatroom_id:%s", bases.GetRequesterIdFromCtx(ctx), bases.GetTargetIdFromCtx(ctx))
		code := services.JoinChatroom(ctx, chatroom)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%v", code)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *JoinChatroomActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatroomInfo{}
}
