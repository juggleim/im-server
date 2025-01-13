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

type AddFavoriteMsgsActor struct {
	bases.BaseActor
}

func (actor *AddFavoriteMsgsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.FavoriteMsgIds); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tmsgs:%v", userId, req)
		code := services.AddFavoriteMsgs(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *AddFavoriteMsgsActor) CreateInputObj() proto.Message {
	return &pbobjs.FavoriteMsgIds{}
}

type DelFavoriteMsgsActor struct {
	bases.BaseActor
}

func (actor *DelFavoriteMsgsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.FavoriteMsgIds); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\tmsgs:%v", userId, req)
		code := services.DelFavoriteMsgs(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *DelFavoriteMsgsActor) CreateInputObj() proto.Message {
	return &pbobjs.FavoriteMsgIds{}
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
