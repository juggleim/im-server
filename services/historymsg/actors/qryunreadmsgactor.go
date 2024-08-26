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

type QryFirstUnreadMsgActor struct {
	bases.BaseActor
}

func (actor *QryFirstUnreadMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryFirstUnreadMsgReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tchannel_type:%v", bases.GetRequesterIdFromCtx(ctx), req.TargetId, req.ChannelType)
		code, resp := services.QryFirstUnreadMsg(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%d", code)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *QryFirstUnreadMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.QryFirstUnreadMsgReq{}
}
