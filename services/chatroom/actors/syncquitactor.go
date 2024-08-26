package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/services"

	"google.golang.org/protobuf/proto"
)

type SyncQuitActor struct {
	bases.BaseActor
}

func (actor *SyncQuitActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatroomMember); ok {
		services.SyncQuitChatroom(ctx, req)
	}
}

func (actor *SyncQuitActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatroomMember{}
}
