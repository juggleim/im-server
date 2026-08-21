package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type MarkMentionReadActor struct {
	bases.BaseActor
}

func (actor *MarkMentionReadActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.MarkReadReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		services.MarkMentionRead(ctx, userId, req)
	}
}

func (actor *MarkMentionReadActor) CreateInputObj() proto.Message {
	return &pbobjs.MarkReadReq{}
}
