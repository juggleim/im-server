package actors

import (
	"context"
	"time"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type ChatMsgActor struct {
	bases.BaseActor
}

func (actor *ChatMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if upMsg, ok := input.(*pbobjs.UpMsg); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tchat_id:%s\tmsg_type:%s", bases.GetRequesterIdFromCtx(ctx), bases.GetTargetIdFromCtx(ctx), upMsg.MsgType)
		code, msgId, msgTime, msgIndex := services.SendChatroomMsg(ctx, upMsg)
		userPubAck := bases.CreateUserPubAckWraper(ctx, code, msgId, msgTime, msgIndex)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%v\tmsg_id:%s\tmsg_time:%d\tmsg_index:%d", code, msgId, msgTime, msgIndex)
	} else {
		userPubAck := bases.CreateUserPubAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, "", time.Now().UnixMilli(), 0)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("upMsg is illigal. upMsg:%v", upMsg)
	}
}

func (actor *ChatMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.UpMsg{}
}
