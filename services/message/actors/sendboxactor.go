package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"
	"time"

	"google.golang.org/protobuf/proto"
)

type SendBoxActor struct {
	bases.BaseActor
}

func (actor *SendBoxActor) OnReceive(ctx context.Context, input proto.Message) {
	if downMsg, ok := input.(*pbobjs.DownMsg); ok {
		if downMsg.MsgTime == 0 || downMsg.MsgId == "" {
			appkey := bases.GetAppKeyFromCtx(ctx)
			msgTime := services.RegenateSendTime(appkey, downMsg.TargetId, time.Now().UnixMilli())
			msgId := tools.GenerateMsgId(msgTime, int32(downMsg.ChannelType), downMsg.TargetId)
			downMsg.MsgTime = msgTime
			downMsg.MsgId = msgId
		}
		if downMsg.SenderId == "" {
			downMsg.SenderId = bases.GetRequesterIdFromCtx(ctx)
		}
		//send to sender's other device
		services.MsgDirect(ctx, bases.GetRequesterIdFromCtx(ctx), downMsg)
		if !commonservices.IsStateMsg(downMsg.Flags) {
			//save msg to sendbox for sender
			services.SaveMsg2Sendbox(ctx, bases.GetAppKeyFromCtx(ctx), bases.GetRequesterIdFromCtx(ctx), downMsg)
			//record conversation for sender
			commonservices.SaveConversation(ctx, bases.GetRequesterIdFromCtx(ctx), downMsg)
		}
	} else {
		logs.WithContext(ctx).Error("input is illigal")
	}
}

func (actor *SendBoxActor) CreateInputObj() proto.Message {
	return &pbobjs.DownMsg{}
}
