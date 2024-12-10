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

type DissolveGroupActor struct {
	bases.BaseActor
}

func (actor *DissolveGroupActor) OnReceive(ctx context.Context, input proto.Message) {
	if addMembersReq, ok := input.(*pbobjs.GroupMembersReq); ok {
		logs.WithContext(ctx).Infof("group_id:%s", addMembersReq.GroupId)
		services.DissolveGroup(ctx, addMembersReq.GroupId)
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *DissolveGroupActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupMembersReq{}
}
