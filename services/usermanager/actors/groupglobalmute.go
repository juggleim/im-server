package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"

	"google.golang.org/protobuf/proto"
)

type GrpGlobalMuteActor struct {
	bases.BaseActor
}

func (actor *GrpGlobalMuteActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.BatchMuteUsersReq); ok {
		services.SetGroupGlobalMute(ctx, req)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *GrpGlobalMuteActor) CreateInputObj() proto.Message {
	return &pbobjs.BatchMuteUsersReq{}
}

type QryGrpGlobalMuteActor struct {
	bases.BaseActor
}

func (actor *QryGrpGlobalMuteActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryBlockUsersReq); ok {
		code, resp := services.QryGrpGlobalMuteUsers(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *QryGrpGlobalMuteActor) CreateInputObj() proto.Message {
	return &pbobjs.QryBlockUsersReq{}
}
