package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
)

func HandleChrmDispatch(ctx context.Context, req *pbobjs.ChrmDispatchReq) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := req.ChatId
	if req.DispatchType == pbobjs.ChrmDispatchType_CreateChatroom {
		initChatroomContainer(appkey, chatId)
	} else if req.DispatchType == pbobjs.ChrmDispatchType_DestroyChatroom {
		key := getChrmKey(appkey, chatId)
		if obj, exist := chrmCache.Get(key); exist {
			container := obj.(*ChatroomContainer)
			container.Destroy()
		}
		chrmCache.Remove(key)
	}
}

func HandleChatMembersDispatch(ctx context.Context, dispatch *pbobjs.ChatMembersDispatchReq) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := dispatch.ChatId
	memberIds := dispatch.MemberIds
	if dispatch.DispatchType == pbobjs.ChatMembersDispatchType_JoinChatroom {
		for _, memberId := range memberIds {
			AddPartialChatroomMember2Cache(ctx, appkey, chatId, memberId)
		}
	} else if dispatch.DispatchType == pbobjs.ChatMembersDispatchType_QuitChatroom {
		for _, memberId := range memberIds {
			DelPartialChatroomMemberFromCache(ctx, appkey, chatId, memberId)
		}
	}
}

func DelPartialChatroomMemberFromCache(ctx context.Context, appkey, chatId, memberId string) bool {
	container, exist := GetChrmContainer(ctx, appkey, chatId)
	if !exist {
		logs.WithContext(ctx).Errorf("chatroom not exist. chat_id:%s", chatId)
		return false
	}
	return container.DelMember(memberId)
}

func AddPartialChatroomMember2Cache(ctx context.Context, appkey, chatId, memberId string) bool {
	container, exist := GetChrmContainer(ctx, appkey, chatId)
	if !exist {
		logs.WithContext(ctx).Errorf("chatroom not exist. chat_id:%s", chatId)
		return false
	}
	return container.AddMember(memberId)
}
