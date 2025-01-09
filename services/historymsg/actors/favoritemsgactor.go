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

type AddFavoriteMsgActor struct {
	bases.BaseActor
}

func (actor *AddFavoriteMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.AddFavoriteMsgReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tsender_id:%s\treceiver_id:%s\tchannel_type:%v\tmsg_id:%s", userId, req.SenderId, req.ReceiverId, req.ChannelType, req.MsgId)
		code := services.AddFavoriteMsg(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *AddFavoriteMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.AddFavoriteMsgReq{}
}

type QryFavoriteMsgsActor struct {
	bases.BaseActor
}

func (actor *QryFavoriteMsgsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryFavoriteMsgsReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tlimit:%d\toffset:%s", userId, req.Limit, req.Offset)
		code, resp := services.QryFavoriteMsgs(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *QryFavoriteMsgsActor) CreateInputObj() proto.Message {
	return &pbobjs.QryFavoriteMsgsReq{}
}
