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

type SetTopMsgActor struct {
	bases.BaseActor
}

func (actor *SetTopMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.TopMsgReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tchannel_type:%v\tmsg_id:%s", userId, req.TargetId, req.ChannelType, req.MsgId)
		code := services.SetTopMsg(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *SetTopMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.TopMsgReq{}
}

type DelTopMsgActor struct {
	bases.BaseActor
}

func (actor *DelTopMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.TopMsgReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tchannel_type:%v\tmsg_id:%s", userId, req.TargetId, req.ChannelType, req.MsgId)
		code := services.DelTopMsg(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *DelTopMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.TopMsgReq{}
}

type GetTopMsgActor struct {
	bases.BaseActor
}

func (actor *GetTopMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GetTopMsgReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tchannel_type:%v", bases.GetRequesterIdFromCtx(ctx), req.TargetId, req.ChannelType)
		code, resp := services.GetTopMsg(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *GetTopMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.GetTopMsgReq{}
}
