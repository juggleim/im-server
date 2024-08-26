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

type GroupMuteActor struct {
	bases.BaseActor
}

func (actor *GroupMuteActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.GroupMuteReq); ok {
		logs.WithContext(ctx).Infof("group_id:%s\tis_mute:%d", req.GroupId, req.IsMute)
		code := services.SetGroupMute(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("code:%d", code)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *GroupMuteActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupMuteReq{}
}
