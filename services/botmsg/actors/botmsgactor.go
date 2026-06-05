package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/botmsg/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type BotMsgActor struct {
	bases.BaseActor
}

func (actor *BotMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.DownMsg); ok {
		botId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("sender_id:%s\ttarget_id:%s\tbot_id:%s\tchannel_type:%d\tmsg_id:%s\tmsg_type:%s", req.SenderId, req.TargetId, botId, req.ChannelType, req.MsgId, req.MsgType)
		services.HandleBotMsg(ctx, req)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
	}
}

func (actor *BotMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.DownMsg{}
}
