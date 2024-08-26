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

type CheckGroupMemberActor struct {
	bases.BaseActor
}

func (actor *CheckGroupMemberActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.CheckGroupMembersReq); ok {
		logs.WithContext(ctx).Infof("groupId:%s\tmembers:%v", req.GroupId, req.MemberIds)
		code, resp := services.CheckGroupMembers(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result_code:%d\tlen:%d", code, len(resp.MemberIdMap))
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, &pbobjs.GroupMembersResp{})
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *CheckGroupMemberActor) CreateInputObj() proto.Message {
	return &pbobjs.CheckGroupMembersReq{}
}
