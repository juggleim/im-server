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

type GetGrpSettingActor struct {
	bases.BaseActor
}

func (actor *GetGrpSettingActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupInfo); ok {
		logs.WithContext(ctx).Infof("groupId:%s\t", req.GroupId)
		code, grpSetting := services.GetGroupSettings(ctx, req.GroupId)
		ack := bases.CreateQueryAckWraper(ctx, code, grpSetting)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *GetGrpSettingActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupInfo{}
}
