package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroommsg/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type SyncChatroomMsgsActor struct {
	bases.BaseActor
}

func (actor *SyncChatroomMsgsActor) OnReceive(ctx context.Context, input proto.Message) {
	if sync, ok := input.(*pbobjs.SyncChatroomReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tchatroom_id:%s\tsync_time:%d", bases.GetRequesterIdFromCtx(ctx), sync.ChatroomId, sync.SyncTime)
		code, downMsgs := services.SyncChatroomMsgs(ctx, sync)
		qryAck := bases.CreateQueryAckWraper(ctx, code, downMsgs)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("chatroom_id:%s\tresult_code:%v\tmsg_len:%d", sync.ChatroomId, code, len(downMsgs.Msgs))
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *SyncChatroomMsgsActor) CreateInputObj() proto.Message {
	return &pbobjs.SyncChatroomReq{}
}
