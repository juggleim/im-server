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

type QryGroupMembersByIdsActor struct {
	bases.BaseActor
}

func (actor *QryGroupMembersByIdsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupMembersReq); ok {
		logs.WithContext(ctx).Infof("group_id:%s\tmember_ids:%v", req.GroupId, req.MemberIds)
		code, members := services.QryGroupMembersByIds(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, members)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%d", code)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryGroupMembersByIdsActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupMembersReq{}
}
