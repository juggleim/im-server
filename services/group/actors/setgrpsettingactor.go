package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/group/services"

	"google.golang.org/protobuf/proto"
)

type SetGrpSettingActor struct {
	bases.BaseActor
}

func (actor *SetGrpSettingActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupInfo); ok {
		logs.WithContext(ctx).Infof("groupId:%s\tsettings:%s", req.GroupId, tools.ToJson(req.Settings))
		code := services.SetGroupSettings(ctx, req.GroupId, req.Settings)
		ack := bases.CreateQueryAckWraper(ctx, code, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *SetGrpSettingActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupInfo{}
}
