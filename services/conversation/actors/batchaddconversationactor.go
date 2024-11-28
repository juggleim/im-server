package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type BatchAddConversationActor struct {
	bases.BaseActor
}

func (actor *BatchAddConversationActor) OnReceive(ctx context.Context, input proto.Message) {
	if batchConvers, ok := input.(*pbobjs.BatchAddConvers); ok {
		for _, conver := range batchConvers.Convers {
			services.SaveConversationV2(bases.GetAppKeyFromCtx(ctx), conver.UserId, conver.Msg, true)
		}
	} else {
		logs.WithContext(ctx).Error("input is illigal.")
	}
}

func (actor *BatchAddConversationActor) CreateInputObj() proto.Message {
	return &pbobjs.BatchAddConvers{}
}
