package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type AddHisMsgActor struct {
	bases.BaseActor
}

func (actor *AddHisMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if addHisMsg, ok := input.(*pbobjs.AddHisMsgReq); ok {
		if addHisMsg != nil {
			converId := commonservices.GetConversationId(addHisMsg.SenderId, addHisMsg.TargetId, addHisMsg.ChannelType)
			logs.WithContext(ctx).Infof("channel_type:%s\tsender:%s\treceiver:%s", addHisMsg.ChannelType, addHisMsg.SenderId, addHisMsg.TargetId)
			latestMsg := services.GetLatestMsg(ctx, converId, addHisMsg.ChannelType)
			latestMsg.Update(addHisMsg.Msg)
			if addHisMsg.ChannelType == pbobjs.ChannelType_Private {
				services.SavePrivateHisMsg(ctx, converId, addHisMsg.SenderId, addHisMsg.TargetId, addHisMsg.Msg)
			} else if addHisMsg.ChannelType == pbobjs.ChannelType_Group {
				services.SaveGroupHisMsg(ctx, converId, addHisMsg.Msg, int(addHisMsg.GroupMemberCount))
			} else if addHisMsg.ChannelType == pbobjs.ChannelType_System {
				services.SaveSystemHisMsg(ctx, converId, addHisMsg.SenderId, addHisMsg.TargetId, addHisMsg.Msg)
			} else if addHisMsg.ChannelType == pbobjs.ChannelType_GroupCast {
				services.SaveGroupCastHisMsg(ctx, converId, addHisMsg.SenderId, addHisMsg.TargetId, addHisMsg.Msg)
			} else if addHisMsg.ChannelType == pbobjs.ChannelType_BroadCast {
				services.SaveBroadCastHisMsg(ctx, converId, addHisMsg.SenderId, addHisMsg.Msg)
			}
		}
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *AddHisMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.AddHisMsgReq{}
}
