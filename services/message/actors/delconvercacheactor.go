package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type DelConverCacheActor struct {
	bases.BaseActor
}

func (actor *DelConverCacheActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ConversationsReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		for _, conver := range req.Conversations {
			services.ClearConversation(ctx, userId, conver.TargetId, conver.ChannelType)
		}
	}
}

func (actor *DelConverCacheActor) CreateInputObj() proto.Message {
	return &pbobjs.ConversationsReq{}
}
