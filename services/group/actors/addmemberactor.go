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

type AddMemberActor struct {
	bases.BaseActor
}

func (actor *AddMemberActor) OnReceive(ctx context.Context, input proto.Message) {
	if addMembersReq, ok := input.(*pbobjs.GroupMembersReq); ok {
		logs.WithContext(ctx).Infof("groupId:%s\tmembers:%v", addMembersReq.GroupId, addMembersReq.MemberIds)
		code := services.AddGroupMembers(ctx, addMembersReq.GroupId, addMembersReq.GroupName, addMembersReq.GroupPortrait, addMembersReq.MemberIds, addMembersReq.ExtFields, addMembersReq.Settings)
		ack := bases.CreateQueryAckWraper(ctx, code, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *AddMemberActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupMembersReq{}
}
