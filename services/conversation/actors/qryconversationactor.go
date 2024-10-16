package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type QryConversationActor struct {
	bases.BaseActor
}

func (actor *QryConversationActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryConverReq); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		if !req.IsInner {
			logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tchannel_type:%v", userId, req.TargetId, req.ChannelType)
		}
		var code errs.IMErrorCode
		var resp proto.Message
		if len(req.UserIds) > 0 {
			code, resp = services.BatchQryConvers(ctx, req)
		} else {
			code, resp = services.QryConver(ctx, userId, req)
		}
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Info("input is illegal")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_MSG_DEFAULT, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *QryConversationActor) CreateInputObj() proto.Message {
	return &pbobjs.QryConverReq{}
}
