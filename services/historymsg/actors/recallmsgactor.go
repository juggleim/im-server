package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type RecallMsgActor struct {
	bases.BaseActor
}

func (actor *RecallMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if recallMsg, ok := input.(*pbobjs.RecallMsgReq); ok {
		logs.WithContext(ctx).WithField("method", "recall_msg").Infof("from_id:%s\ttarget_id:%schannelType:%v\tmsg_id:%s\tmsg_time:%d", bases.GetRequesterIdFromCtx(ctx), recallMsg.TargetId, recallMsg.ChannelType, recallMsg.MsgId, recallMsg.MsgTime)
		code := services.RecallMsg(ctx, recallMsg)
		userPubAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		logs.WithContext(ctx).WithField("method", "recall_msg").Infof("code:%d", code)
	} else {
		logs.WithContext(ctx).WithField("method", "recall_msg").Errorf("input is illigal. input:%v", input)
		userPubAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
	}
}

func (actor *RecallMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.RecallMsgReq{}
}
