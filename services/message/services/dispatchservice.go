package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/botmsg/botclient"
	"im-server/services/commonservices"
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
			doDispatch(ctx, receiverId, downMsg)
		})
	} else if downMsg.ChannelType == pbobjs.ChannelType_Group {
		memberIds := bases.GetTargetIdsFromCtx(ctx)
		bases.SetTargetIds2Ctx(ctx, []string{})
		threadhold := 1000
		appinfo, exist := commonservices.GetAppInfo(appkey)
		if exist && appinfo != nil {
			threadhold = appinfo.BigGrpThreshold
		}
		if downMsg.MemberCount < int32(threadhold) {
			//save msg to inbox and record conversation for each member
			for _, receiverId := range memberIds {
				newDownMsg := copyDownMsg(downMsg)
				recvId := receiverId
				MsgSinglePools.GetPool(strings.Join([]string{appkey, recvId}, "_")).Submit(func() {
					doDispatch(ctx, recvId, newDownMsg)
				})
			}
		} else {
			offlineMemberIds := []string{}
			for _, receiverId := range memberIds {
				userStatus := GetUserStatus(appkey, receiverId)
				if userStatus.IsOnline() {
					newDownMsg := copyDownMsg(downMsg)
					recvId := receiverId
					MsgSinglePools.GetPool(strings.Join([]string{appkey, recvId}, "_")).Submit(func() {
						doDispatch(ctx, recvId, newDownMsg)
					})
				} else {
					offlineMemberIds = append(offlineMemberIds, receiverId)
				}
			}
			if len(offlineMemberIds) > 0 {
				if !commonservices.IsStateMsg(downMsg.Flags) {
					commonservices.BatchSaveConversations(ctx, offlineMemberIds, downMsg)
				}
			}
		}
	}
}

func doDispatch(ctx context.Context, receiverId string, msg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//TODO save imediately when user online, other wise, user async queue.
	//TODO batch insert & regenate msg time.
	sendTime := RegenateSendTime(appkey, receiverId, msg.MsgTime)
	msg.MsgTime = sendTime
	//handle conversation check, such as undisturb, unread index
	HandleDownMsgByConver(ctx, receiverId, msg.TargetId, msg.ChannelType, msg)
	targetUserInfo := commonservices.GetTargetUserInfo(ctx, receiverId)
	if targetUserInfo.UserType == pbobjs.UserType_Bot {
		botclient.SendMsg2Bot(ctx, receiverId, msg)
	} else {
		if !commonservices.IsStateMsg(msg.Flags) {
			//record conversation
			commonservices.SaveConversation(ctx, receiverId, msg)
			SaveMsg2Inbox(appkey, receiverId, msg)
			//send to client
			MsgOrNtf(ctx, receiverId, msg)
		} else {
			MsgDirect(ctx, receiverId, msg)
		}
	}
}

func copyDownMsg(msg *pbobjs.DownMsg) *pbobjs.DownMsg {
	return &pbobjs.DownMsg{
		TargetId:       msg.TargetId,
		ChannelType:    msg.ChannelType,
		MsgType:        msg.MsgType,
		SenderId:       msg.SenderId,
		MsgId:          msg.MsgId,
		MsgSeqNo:       msg.MsgSeqNo,
		MsgContent:     msg.MsgContent,
		MsgTime:        msg.MsgTime,
		Flags:          msg.Flags,
		IsSend:         msg.IsSend,
		Platform:       msg.Platform,
		ClientUid:      msg.ClientUid,
		PushData:       msg.PushData,
		MentionInfo:    msg.MentionInfo,
		IsRead:         msg.IsRead,
		ReferMsg:       msg.ReferMsg,
		TargetUserInfo: msg.TargetUserInfo,
		GroupInfo:      msg.GroupInfo,
		MergedMsgs:     msg.MergedMsgs,
		UndisturbType:  msg.UndisturbType,
		MemberCount:    msg.MemberCount,
		ReadCount:      msg.ReadCount,
		UnreadIndex:    msg.UnreadIndex,
	}
}
