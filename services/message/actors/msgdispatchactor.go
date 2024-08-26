package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type MsgDispatchActor struct {
	bases.BaseActor
}

func (actor *MsgDispatchActor) OnReceive(ctx context.Context, input proto.Message) {
	if downMsg, ok := input.(*pbobjs.DownMsg); ok {
		logs.WithContext(ctx).Infof("channel_type:%v\treceiver:%s\tgroup_id:%s\ttargetids_len:%d", downMsg.ChannelType, bases.GetTargetIdFromCtx(ctx), bases.GetGroupIdFromCtx(ctx), len(bases.GetTargetIdsFromCtx(ctx)))
		services.DispatchMsg(ctx, downMsg)
		logs.WithContext(ctx).Info("finish")
	} else {
		logs.WithContext(ctx).Error("input is illigal")
	}
}

func (actor *MsgDispatchActor) CreateInputObj() proto.Message {
	return &pbobjs.DownMsg{}
}
