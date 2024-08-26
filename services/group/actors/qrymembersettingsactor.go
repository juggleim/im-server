package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/group/services"

	"google.golang.org/protobuf/proto"
)

type QryMemberSettingsActor struct {
	bases.BaseActor
}

func (actor *QryMemberSettingsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryGrpMemberSettingsReq); ok {
		groupId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("group_id:%s\tmeber_id:%s", groupId, req.MemberId)
		code, info := services.QryMemberSettings(ctx, groupId, req.MemberId)
		ack := bases.CreateQueryAckWraper(ctx, code, info)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%d", code)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryMemberSettingsActor) CreateInputObj() proto.Message {
	return &pbobjs.QryGrpMemberSettingsReq{}
}
