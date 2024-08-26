package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"sync/atomic"
)

type MsgBucket struct {
}

func getMaxMsgCount(appkey string) int {
	count := 50
	if appinfo, exist := commonservices.GetAppInfo(appkey); exist {
		count = appinfo.ChrmMsgCacheMaxCount
	}
	return count
}

func AppendMsg(ctx context.Context, appkey, chatId string, msg *pbobjs.DownMsg) {
	container, exist := GetChrmContainer(ctx, appkey, chatId)
	if !exist {
		logs.WithContext(ctx).Errorf("chatroom not exist. chat_id:%s", chatId)
		return
	}
	container.AppendMsg(ctx, msg)
}

func SyncChatroomMsgs(ctx context.Context, sync *pbobjs.SyncChatroomReq) (errs.IMErrorCode, *pbobjs.SyncChatroomMsgResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := sync.ChatroomId
	userId := bases.GetRequesterIdFromCtx(ctx)
	container, exist := GetChrmContainer(ctx, appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST, nil
	}
	container.CleanUnread(userId)
	msgs, code := container.GetMsgsBaseTime(ctx, userId, sync.SyncTime)
	if code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	return errs.IMErrorCode_SUCCESS, &pbobjs.SyncChatroomMsgResp{
		Msgs: msgs,
	}
}

func HandleChatMsgsDispatch(ctx context.Context, msg *pbobjs.DownMsgSet) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	senderId := bases.GetRequesterIdFromCtx(ctx)
	for _, msg := range msg.Msgs {
		chatId := msg.TargetId
		//add in cache
		AppendMsg(ctx, appkey, chatId, msg)
		//notify membners
		ntfChatMembers(ctx, chatId, senderId, msg.MsgTime, pbobjs.NotifyType_ChatroomMsg)
	}
}

// TODO: combine the ntf
func ntfChatMembers(ctx context.Context, chatId, senderId string, msgTime int64, ntfType pbobjs.NotifyType) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//notify online chatroom members
	container, exist := GetChrmContainer(ctx, appkey, chatId)
	if exist {
		needCleanMembers := []string{}
		container.ForeachMembers(func(member *ChatroomMember) bool {
			memberId := member.MemberId
			unreadCount := atomic.AddInt32(&member.UnReadCount, 1)
			if unreadCount > 30 {
				needCleanMembers = append(needCleanMembers, memberId)
				return true
			}
			rpcNtf := bases.CreateServerPubWraper(ctx, senderId, memberId, "ntf", &pbobjs.Notify{
				Type:       ntfType,
				SyncTime:   msgTime,
				ChatroomId: chatId,
			})
			rpcNtf.Qos = 0
			if memberId == senderId {
				rpcNtf.PublishType = int32(commonservices.PublishType_AllSessionExceptSelf)
			}
			bases.UnicastRouteWithNoSender(rpcNtf)
			return true
		})
		//remove timeout members
		for _, memberId := range needCleanMembers {
			//remove from cache
			container.DelMember(memberId)
			//notify chatroom
			bases.AsyncRpcCall(ctx, "c_sync_quit", chatId, &pbobjs.ChatroomMember{
				MemberId: memberId,
			})
		}
	}
}
