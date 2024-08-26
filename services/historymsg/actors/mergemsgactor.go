package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type MergeMsgActor struct {
	bases.BaseActor
}

func (actor *MergeMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.MergeMsgReq); ok {
		mergedMsgIds := []string{}
		for _, msg := range req.MergedMsgs.Msgs {
			mergedMsgIds = append(mergedMsgIds, msg.MsgId)
		}
		logs.WithContext(ctx).Infof("parent_msgid:%s\tchannel_type:%s\tuser_id:%s\ttarget_id:%s\tmsg_len:%d\tmsgs:%v", req.ParentMsgId, req.MergedMsgs.ChannelType, req.MergedMsgs.UserId, req.MergedMsgs.TargetId, len(req.MergedMsgs.Msgs), mergedMsgIds)
		services.MergeMsg(ctx, req)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *MergeMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.MergeMsgReq{}
}
