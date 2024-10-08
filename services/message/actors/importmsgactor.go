package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type ImportPrivateHisMsgActor struct {
	bases.BaseActor
}

func (actor *ImportPrivateHisMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if msg, ok := input.(*pbobjs.UpMsg); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		exts := bases.GetExtsFromCtx(ctx)
		if receiverId, exist := exts[commonservices.RpcExtKey_RealTargetId]; exist {
			services.ImportPrivateHisMsg(ctx, userId, receiverId, msg)
		} else {
			logs.WithContext(ctx).Errorf("have no receiver")
		}
	} else {
		logs.WithContext(ctx).Errorf("input is illigal. val:%v", msg)
	}
}

func (actor *ImportPrivateHisMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.UpMsg{}
}
