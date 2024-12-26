package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type AddConversationActor struct {
	bases.BaseActor
}

func (actor *AddConversationActor) OnReceive(ctx context.Context, input proto.Message) {
	if conver, ok := input.(*pbobjs.Conversation); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		targetId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tconver_target_id:%s\tchannel_type:%v", userId, targetId, conver.TargetId, conver.ChannelType)
		if conver.TargetId == "" || conver.ChannelType == pbobjs.ChannelType_Unknown {
			logs.WithContext(ctx).Errorf("unknown conversation. user_id:%s\ttarget_id:%s", userId, targetId)
		} else if conver.Msg == nil {
			code, resp := services.SaveNilConversationV2(ctx, bases.GetAppKeyFromCtx(ctx), bases.GetTargetIdFromCtx(ctx), conver.TargetId, conver.ChannelType)
			qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
			actor.Sender.Tell(qryAck, actorsystem.NoSender)
		} else {
			uIds := bases.GetTargetIdsFromCtx(ctx)
			if len(uIds) > 0 {
				for _, uId := range uIds {
					services.SaveConversationV2(ctx, bases.GetAppKeyFromCtx(ctx), uId, conver.Msg, false)
				}
			} else {
				services.SaveConversationV2(ctx, bases.GetAppKeyFromCtx(ctx), bases.GetTargetIdFromCtx(ctx), conver.Msg, false)
			}
		}
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	}
}

func (actor *AddConversationActor) CreateInputObj() proto.Message {
	return &pbobjs.Conversation{}
}
