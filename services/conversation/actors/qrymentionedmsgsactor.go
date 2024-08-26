package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type QryMentionMsgsActor struct {
	bases.BaseActor
}

func (actor *QryMentionMsgsActor) OnReceive(ctx context.Context, input proto.Message) {
	if qryMentionMsgReq, ok := input.(*pbobjs.QryMentionMsgsReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tchannel_type:%v\tstart_time:%d\tcount:%d\torder:%d\tread_index:%d", userId, qryMentionMsgReq.TargetId, qryMentionMsgReq.ChannelType, qryMentionMsgReq.StartTime, qryMentionMsgReq.Count, qryMentionMsgReq.Order, qryMentionMsgReq.LatestReadIndex)
		resp := services.QryMentionedMsgs(ctx, userId, qryMentionMsgReq)
		qryAck := bases.CreateQueryAckWraper(ctx, 0, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result_len:%d", len(resp.MentionMsgs))
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
	}
}

func (actor *QryMentionMsgsActor) CreateInputObj() proto.Message {
	return &pbobjs.QryMentionMsgsReq{}
}
