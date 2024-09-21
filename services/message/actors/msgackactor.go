package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type MsgAckActor struct {
	bases.BaseActor
}

func (actor *MsgAckActor) OnReceive(ctx context.Context, input proto.Message) {
	if ack, ok := input.(*pbobjs.MsgAck); ok {
		appkey := bases.GetAppKeyFromCtx(ctx)
		userId := bases.GetTargetIdFromCtx(ctx)
		userStatus := services.GetUserStatus(appkey, userId)
		if userStatus != nil {
			userStatus.CloseNtf(ack.MsgTime)
			//statistic
			commonservices.ReportDownMsg(appkey, tools.ParseChannelTypeFromMsgId(ack.MsgId), 1)
		}
	}
}

func (actor *MsgAckActor) CreateInputObj() proto.Message {
	return &pbobjs.MsgAck{}
}
