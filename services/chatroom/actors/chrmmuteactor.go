package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type ChrmMuteActor struct {
	bases.BaseActor
}

func (actor *ChrmMuteActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatroomInfo); ok {
		logs.WithContext(ctx).Infof("chatroom_id:%s\tis_mute:%d", req.ChatId, req.IsMute)
		code := services.SetChrmMute(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *ChrmMuteActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatroomInfo{}
}
