package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type UpdLatestMsgActor struct {
	bases.BaseActor
}

func (actor *UpdLatestMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UpdLatestMsgReq); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_id:%s\tchannel_type:%v\taction:%s", userId, req.TargetId, req.ChannelType, req.Action)
		uIds := bases.GetTargetIdsFromCtx(ctx)
		if len(uIds) > 0 {
			for _, uId := range uIds {
				services.UpdLatestMsg(ctx, uId, req)
			}
		} else {
			services.UpdLatestMsg(ctx, userId, req)
		}
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
	}
}

func (actor *UpdLatestMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.UpdLatestMsgReq{}
}
