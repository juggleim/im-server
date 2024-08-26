package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type DelMsgActor struct {
	bases.BaseActor
}

func (actor *DelMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if delMsg, ok := input.(*pbobjs.DelHisMsgsReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		msgIds := []string{}
		for _, msg := range delMsg.Msgs {
			msgIds = append(msgIds, msg.MsgId)
		}
		logs.WithContext(ctx).Infof("user_id:%s\ttargetId=%s\tchannel_type=%v\tmsg_id=%v", userId, delMsg.TargetId, delMsg.ChannelType, msgIds)
		code := services.DelHisMsg(ctx, delMsg)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result_code:%v", code)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *DelMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.DelHisMsgsReq{}
}
