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

type ModifyMsgActor struct {
	bases.BaseActor
}

func (actor *ModifyMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if modifyReq, ok := input.(*pbobjs.ModifyMsgReq); ok {
		logs.WithContext(ctx).Infof("from_id:%s\ttarget_id:%schannelType:%v\tmsg_id:%s\tmsg_time:%d", bases.GetRequesterIdFromCtx(ctx), modifyReq.TargetId, modifyReq.ChannelType, modifyReq.MsgId, modifyReq.MsgTime)
		code := services.ModifyMsg(ctx, modifyReq)
		userPubAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%d", code)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal. input:%v", input)
		userPubAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
	}
}

func (actor *ModifyMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.ModifyMsgReq{}
}
