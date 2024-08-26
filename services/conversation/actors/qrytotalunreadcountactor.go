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

type QryTotalUnreadCountActor struct {
	bases.BaseActor
}

func (actor *QryTotalUnreadCountActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryTotalUnreadCountReq); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s", userId)

		resp := services.QryTotalUnreadCount(ctx, userId, req)
		qryAck := bases.CreateQueryAckWraper(ctx, 0, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)

		logs.WithContext(ctx).Infof("total_count:%d", resp.TotalCount)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
	}
}

func (actor *QryTotalUnreadCountActor) CreateInputObj() proto.Message {
	return &pbobjs.QryTotalUnreadCountReq{}
}
