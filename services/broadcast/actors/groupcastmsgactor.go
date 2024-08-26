package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/broadcast/services"
	"im-server/services/commonservices/logs"
	"time"

	"google.golang.org/protobuf/proto"
)

type GroupCastMsgActor struct {
	bases.BaseActor
}

func (actor *GroupCastMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if upMsg, ok := input.(*pbobjs.UpMsg); ok {
		logs.WithContext(ctx).Infof("target_id:%s\tmsg_type:%s\tflag:%d", bases.GetTargetIdFromCtx(ctx), upMsg.MsgType, upMsg.Flags)
		code, msgId, sendTime, msgSeq := services.SendGroupCastMsg(ctx, upMsg)
		userPubAck := bases.CreateUserPubAckWraper(ctx, code, msgId, sendTime, msgSeq)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%d\tmsg_id:%s", code, msgId)
	} else {
		userPubAck := bases.CreateUserPubAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, "", time.Now().UnixMilli(), 0)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("upMsg is illigal. upMsg:%v", upMsg)
	}
}

func (actor *GroupCastMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.UpMsg{}
}
