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

type GroupMemberMuteActor struct {
	bases.BaseActor
}

func (actor *GroupMemberMuteActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupMemberMuteReq); ok {
		logs.WithContext(ctx).Infof("group_id:%s\tis_mute:%d\tmember_ids:%v", req.GroupId, req.IsMute, req.MemberIds)
		code := services.SetGroupMemberMute(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%d", code)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *GroupMemberMuteActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupMemberMuteReq{}
}
