package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type DestroyChrmActor struct {
	bases.BaseActor
}

func (actor *DestroyChrmActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatroomInfo); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tchatroom_id:%s", bases.GetRequesterIdFromCtx(ctx), req.ChatId)
		code := services.DestroyChatroom(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *DestroyChrmActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatroomInfo{}
}
