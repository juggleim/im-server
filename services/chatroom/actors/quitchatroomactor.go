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

type QuitChatroomActor struct {
	bases.BaseActor
}

func (actor *QuitChatroomActor) OnReceive(ctx context.Context, input proto.Message) {
	if chatroom, ok := input.(*pbobjs.ChatroomInfo); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tchatroom_id:%s", userId, chatroom.ChatId)
		code := services.QuitChatroom(ctx, userId, chatroom)
		qrAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qrAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *QuitChatroomActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatroomInfo{}
}
