package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/group/services"

	"google.golang.org/protobuf/proto"
)

type SetGrpMemberSettingActor struct {
	bases.BaseActor
}

func (actor *SetGrpMemberSettingActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupMember); ok {
		grpId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("group_id:%s\tmember_id:%s\tsettings%s", grpId, req.MemberId, tools.ToJson(req.Settings))
		code := services.SetGroupMemberSettings(ctx, grpId, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *SetGrpMemberSettingActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupMember{}
}
