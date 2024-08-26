package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroommsg/services"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type DispatchMembersActor struct {
	bases.BaseActor
}

func (actor DispatchMembersActor) OnReceive(ctx context.Context, input proto.Message) {
	if chatroom, ok := input.(*pbobjs.ChatMembersDispatchReq); ok {
		logs.WithContext(ctx).Infof("chatroom_id:%s\tmembers:%v\tdispatch_type:%v", chatroom.ChatId, chatroom.MemberIds, chatroom.DispatchType)
		services.HandleChatMembersDispatch(ctx, chatroom)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal.")
	}
}

func (actor DispatchMembersActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatMembersDispatchReq{}
}
