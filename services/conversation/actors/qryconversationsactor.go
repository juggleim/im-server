package actors

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type QryConversationsActor struct {
	bases.BaseActor
}

func (actor *QryConversationsActor) OnReceive(ctx context.Context, input proto.Message) {
	if qryConverReq, ok := input.(*pbobjs.QryConversationsReq); ok {
		requesterId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tstart:%d\tcount:%d\ttarget_id:%s\tchannel_type:%d", requesterId, qryConverReq.StartTime, qryConverReq.Count, qryConverReq.TargetId, qryConverReq.ChannelType)

		resp := services.QryConversations(ctx, qryConverReq)
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%d\tlen:%d", errs.IMErrorCode_SUCCESS, len(resp.Conversations))
	} else {
		fmt.Println("p_msg, upMsg is illigal.")
	}
}

func (actor *QryConversationsActor) CreateInputObj() proto.Message {
	return &pbobjs.QryConversationsReq{}
}
