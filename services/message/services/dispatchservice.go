package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/botmsg/botclient"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"strings"
)

var MsgSinglePools *tools.SinglePools

func init() {
	MsgSinglePools = tools.NewSinglePools(8192)
}

func DispatchMsg(ctx context.Context, downMsg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	if downMsg.ChannelType == pbobjs.ChannelType_Private || downMsg.ChannelType == pbobjs.ChannelType_System {
		receiverId := bases.GetTargetIdFromCtx(ctx)
		MsgSinglePools.GetPool(strings.Join([]string{appkey, receiverId}, "_")).Submit(func() {
			doDispatch(ctx, receiverId, downMsg, false)
		})
	} else if downMsg.ChannelType == pbobjs.ChannelType_Group {
		memberIds := bases.GetTargetIdsFromCtx(ctx)
		bases.SetTargetIds2Ctx(ctx, []string{})
		threadhold := 1000
		appinfo, exist := commonservices.GetAppInfo(appkey)
		if exist && appinfo != nil {
			threadhold = appinfo.BigGrpThreshold
		}
		if downMsg.MemberCount > int32(threadhold) {
			preheat(ctx, appkey, memberIds)
			for _, receiverId := range memberIds {
				newDownMsg := copyDownMsg(downMsg)
				receId := receiverId
				userStatus := GetUserStatus(appkey, receId)
				closeOffline := (!userStatus.IsOnline())
				MsgSinglePools.GetPool(strings.Join([]string{appkey, receId}, "_")).Submit(func() {
					doDispatch(ctx, receId, newDownMsg, closeOffline)
				})
			}
		} else {
			for _, receiverId := range memberIds {
				newDownMsg := copyDownMsg(downMsg)
				receId := receiverId
				MsgSinglePools.GetPool(strings.Join([]string{appkey, receId}, "_")).Submit(func() {
					doDispatch(ctx, receId, newDownMsg, false)
				})
			}
		}
	}
}

func preheat(ctx context.Context, appkey string, memberIds []string) {
	//preheat user status
	noStatusCacheUids := []string{}
	for _, receiverId := range memberIds {
		if !UserStatusCacheContains(appkey, receiverId) {
			noStatusCacheUids = append(noStatusCacheUids, receiverId)
		}
	}
	if len(noStatusCacheUids) > 0 {
		BatchInitUserStatus(ctx, appkey, noStatusCacheUids)
	}
}

func doDispatch(ctx context.Context, receiverId string, msg *pbobjs.DownMsg, closeOffline bool) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//TODO save imediately when user online, other wise, use async queue.
	//TODO batch insert
	sendTime := RegenateSendTime(appkey, receiverId, msg.MsgTime)
	msg.MsgTime = sendTime
	//handle conversation check, such as undisturb, unread index
	HandleDownMsgByConver(ctx, receiverId, msg.TargetId, msg.SubChannel, msg.ChannelType, msg)
	targetUserInfo := commonservices.GetTargetUserInfo(ctx, receiverId)
	if targetUserInfo.UserType == pbobjs.UserType_Bot {
		if msg.ChannelType == pbobjs.ChannelType_Private || (msg.ChannelType == pbobjs.ChannelType_Group && commonservices.IsDirectMentionedMe(targetUserInfo.UserId, msg)) {
			botclient.SendMsg2Bot(ctx, receiverId, msg)
		}
	} else {
		if !msgdefines.IsStateMsg(msg.Flags) {
			if !closeOffline {
				commonservices.SaveConversation(ctx, receiverId, msg)
				SaveMsg2Inbox(appkey, receiverId, msg)
			} else {
				target := fmt.Sprintf("%s_%s_%d", appkey, msg.TargetId, msg.ChannelType)
				batchExecutorPool.GetBatchExecutor(target).Append(fmt.Sprintf("%s_%s", target, receiverId), &BatchConverItem{
					Appkey: appkey,
					UserId: receiverId,
					Msg:    msg,
				})
			}
			//send to client
			MsgOrNtf(ctx, receiverId, msg)
		} else {
			if !closeOffline {
				MsgDirect(ctx, receiverId, msg)
			}
		}
	}
}

func copyDownMsg(msg *pbobjs.DownMsg) *pbobjs.DownMsg {
	return &pbobjs.DownMsg{
		TargetId:          msg.TargetId,
		ChannelType:       msg.ChannelType,
		MsgType:           msg.MsgType,
		SenderId:          msg.SenderId,
		MsgId:             msg.MsgId,
		MsgSeqNo:          msg.MsgSeqNo,
		MsgContent:        msg.MsgContent,
		MsgTime:           msg.MsgTime,
		Flags:             msg.Flags,
		IsSend:            msg.IsSend,
		Platform:          msg.Platform,
		ClientUid:         msg.ClientUid,
		PushData:          msg.PushData,
		MentionInfo:       msg.MentionInfo,
		IsRead:            msg.IsRead,
		ReferMsg:          msg.ReferMsg,
		TargetUserInfo:    msg.TargetUserInfo,
		GroupInfo:         msg.GroupInfo,
		MergedMsgs:        msg.MergedMsgs,
		UndisturbType:     msg.UndisturbType,
		MemberCount:       msg.MemberCount,
		ReadCount:         msg.ReadCount,
		UnreadIndex:       msg.UnreadIndex,
		SearchText:        msg.SearchText,
		GrpMemberInfo:     msg.GrpMemberInfo,
		DestroyTime:       msg.DestroyTime,
		LifeTimeAfterRead: msg.LifeTimeAfterRead,
		SubChannel:        msg.SubChannel,
	}
}

func UserPush(ctx context.Context, msg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	targetId := bases.GetTargetIdFromCtx(ctx)
	if !UserStatusCacheContains(appkey, targetId) {
		BatchInitUserStatus(ctx, appkey, []string{targetId})
	}
	userStatus := GetUserStatus(appkey, targetId)
	if userStatus != nil {
		if msg.PushData != nil && msg.PushData.IsVoip {
			SendPush(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, msg)
		} else {
			if !userStatus.IsOnline() || userStatus.OpenPushSwitch() {
				SendPush(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, msg)
			}
		}
	}
}
