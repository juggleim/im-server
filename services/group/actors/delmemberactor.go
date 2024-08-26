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

type DelMemberActor struct {
	bases.BaseActor
}

func (actor *DelMemberActor) OnReceive(ctx context.Context, input proto.Message) {
	if addMembersReq, ok := input.(*pbobjs.GroupMembersReq); ok {
		logs.WithContext(ctx).Infof("groupId:%s\tmembers:%v", addMembersReq.GroupId, addMembersReq.MemberIds)
		services.DelGroupMembers(ctx, addMembersReq.GroupId, addMembersReq.MemberIds)
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *DelMemberActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupMembersReq{}
}
