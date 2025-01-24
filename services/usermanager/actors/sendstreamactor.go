package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"

	"google.golang.org/protobuf/proto"
)

type SendStreamActor struct {
	bases.BaseActor
}

func (actor *SendStreamActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.StreamDownMsg); ok {
		sender := bases.GetRequesterIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("sender:%s\treceiver:%s\tchannel_type:%v\tmsg_type:%s\tmsg_len:%d", sender, req.TargetId, req.ChannelType, req.MsgType, len(req.MsgItems))
		code := services.HandleStreamMsg(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal")
	}
}

func (acctor *SendStreamActor) CreateInputObj() proto.Message {
	return &pbobjs.StreamDownMsg{}
}
