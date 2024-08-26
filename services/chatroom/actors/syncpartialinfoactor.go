package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/services"

	"google.golang.org/protobuf/proto"
)

type SyncPartialInfoActor struct {
	bases.BaseActor
}

func (actor *SyncPartialInfoActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatMsgNode); ok {
		chatId := bases.GetTargetIdFromCtx(ctx)
		resp, code := services.GetPartialMembers(ctx, chatId, req.NodeName, req.Method)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, &pbobjs.ChatroomInfo{})
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *SyncPartialInfoActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatMsgNode{}
}
