package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"
	"time"

	"google.golang.org/protobuf/proto"
)

type SetMsgExtActor struct {
	bases.BaseActor
}

func (actor *SetMsgExtActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.MsgExt); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tmsg_id:%s\ttarget_id:%s\tchannel_type:%v\tkey_value:%v", bases.GetRequesterIdFromCtx(ctx), req.MsgId, req.TargetId, req.ChannelType, req.Ext)
		code := services.SetMsgExt(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal. input:%v", input)
		userPubAck := bases.CreateUserPubAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, "", time.Now().UnixMilli(), 0)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
	}
}

func (actor *SetMsgExtActor) CreateInputObj() proto.Message {
	return &pbobjs.MsgExt{}
}

type DelMsgExtActor struct {
	bases.BaseActor
}

func (actor *DelMsgExtActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.MsgExt); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tmsg_id:%s\ttarget_id:%s\tchannel_type:%v\tkey_value:%v", bases.GetRequesterIdFromCtx(ctx), req.MsgId, req.TargetId, req.ChannelType, req.Ext)
		code := services.DelMsgExt(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal. input:%v", input)
		userPubAck := bases.CreateUserPubAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, "", time.Now().UnixMilli(), 0)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
	}
}

func (actor *DelMsgExtActor) CreateInputObj() proto.Message {
	return &pbobjs.MsgExt{}
}
