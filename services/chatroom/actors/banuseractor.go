package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type BanUserActor struct {
	bases.BaseActor
}

func (actor *BanUserActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.BatchBanUserReq); ok {
		logs.WithContext(ctx).Infof("chatroom_id:%s\tban_type:%d\tis_delete:%v\tmembers:%v", req.ChatId, req.BanType, req.IsDelete, req.MemberIds)
		code := services.HandleBanUsers(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *BanUserActor) CreateInputObj() proto.Message {
	return &pbobjs.BatchBanUserReq{}
}

type QryBanUsersActor struct {
	bases.BaseActor
}

func (actor *QryBanUsersActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryChrmBanUsersReq); ok {
		logs.WithContext(ctx).Infof("chatroom_id:%s\tban_type:%d\tlimit:%d\toffset:%s", req.ChatId, req.BanType, req.Limit, req.Offset)
		code, resp := services.QryBanUsers(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *QryBanUsersActor) CreateInputObj() proto.Message {
	return &pbobjs.QryChrmBanUsersReq{}
}
