package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"time"
)

func SendGroupMsg(ctx context.Context, upMsg *pbobjs.UpMsg) (errs.IMErrorCode, string, int64, int64, string, int32) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	groupId := bases.GetTargetIdFromCtx(ctx)
	senderId := bases.GetRequesterIdFromCtx(ctx)

	//statistic
	commonservices.ReportUpMsg(appkey, pbobjs.ChannelType_Group, 1)

	//check user is member of group
	isFromApi := bases.GetIsFromApiFromCtx(ctx)
	if !isFromApi {
		if !checkIsMember(ctx, groupId, bases.GetRequesterIdFromCtx(ctx)) {
			sendTime := time.Now().UnixMilli()
			msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Group), groupId)
			return errs.IMErrorCode_GROUP_NOTGROUPMEMBER, msgId, sendTime, 0, upMsg.ClientUid, 0
		}
		//check group member mute
		if checkGroupMemberIsMute(ctx, groupId, senderId) {
			sendTime := time.Now().UnixMilli()
			msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Group), groupId)
			return errs.IMErrorCode_GROUP_GROUPMEMBERMUTE, msgId, sendTime, 0, upMsg.ClientUid, 0
		}
		//check group mute
		if checkGroupIsMute(ctx, groupId) {
			//check group member allow
			if !checkGroupMemberIsAllow(ctx, groupId, senderId) {
				sendTime := time.Now().UnixMilli()
				msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Group), groupId)
				return errs.IMErrorCode_GROUP_GROUPMUTE, msgId, sendTime, 0, upMsg.ClientUid, 0
			}
		}
	}

	//check msg interceptor
	if code := commonservices.CheckMsgInterceptor(ctx, senderId, groupId, pbobjs.ChannelType_Group, upMsg); code != errs.IMErrorCode_SUCCESS {
		sendTime := time.Now().UnixMilli()
		msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Group), groupId)
		return code, msgId, sendTime, 0, upMsg.ClientUid, 0
	}
	msgConverCache := commonservices.GetMsgConverCache(ctx, groupId, pbobjs.ChannelType_Group)
	msgId, sendTime, msgSeq := msgConverCache.GenerateMsgId(groupId, pbobjs.ChannelType_Group, time.Now().UnixMilli(), upMsg.Flags)
	preMsgId := bases.GetMsgIdFromCtx(ctx)
	if preMsgId != "" {
		msgId = preMsgId
	}
	if upMsg.ClientUid != "" {
		if oldAck, filter := commonservices.FilterDuplicateMsg(upMsg.ClientUid, commonservices.MsgAck{
			MsgId:   msgId,
			MsgTime: sendTime,
			MsgSeq:  msgSeq,
		}); filter {
			return errs.IMErrorCode_SUCCESS, oldAck.MsgId, oldAck.MsgTime, oldAck.MsgSeq, upMsg.ClientUid, 0
		}
	} else {
		upMsg.ClientUid = tools.GenerateUUIDShort22()
	}

	groupInfo := GetGroupInfo4Msg(ctx, groupId)

	//update mentioned user's info
	UpdateMentionedUserInfo(ctx, upMsg)

	var memberIds []string
	//oriented msgs
	if len(upMsg.ToUserIds) > 0 {
		newMemberIds := []string{}
		for _, id := range upMsg.ToUserIds {
			if id != senderId && checkIsMember(ctx, groupId, id) {
				newMemberIds = append(newMemberIds, id)
			}
		}
		memberIds = newMemberIds
	} else {
		memberIds = getMembersExceptMe(ctx, groupId)
	}
	memberCount := len(memberIds)

	downMsg4Sendbox := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       groupId,
		ChannelType:    pbobjs.ChannelType_Group,
		MsgType:        upMsg.MsgType,
		MsgId:          msgId,
		MsgSeqNo:       msgSeq,
		MsgContent:     upMsg.MsgContent,
		MsgTime:        sendTime,
		Flags:          upMsg.Flags,
		ClientUid:      upMsg.ClientUid,
		IsSend:         true,
		MentionInfo:    upMsg.MentionInfo,
		ReferMsg:       commonservices.FillReferMsg(ctx, upMsg),
		TargetUserInfo: commonservices.GetSenderUserInfo(ctx),
		GroupInfo:      groupInfo,
		MergedMsgs:     upMsg.MergedMsgs,
		MemberCount:    int32(memberCount),
		PushData:       upMsg.PushData,
	}
	commonservices.Save2Sendbox(ctx, downMsg4Sendbox)

	if bases.GetOnlySendboxFromCtx(ctx) {
		return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq, upMsg.ClientUid, int32(memberCount)
	}

	downMsg := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       groupId,
		ChannelType:    pbobjs.ChannelType_Group,
		MsgType:        upMsg.MsgType,
		MsgId:          msgId,
		MsgSeqNo:       msgSeq,
		MsgContent:     upMsg.MsgContent,
		MsgTime:        sendTime,
		Flags:          upMsg.Flags,
		ClientUid:      upMsg.ClientUid,
		MentionInfo:    upMsg.MentionInfo,
		ReferMsg:       commonservices.FillReferMsg(ctx, upMsg),
		TargetUserInfo: commonservices.GetSenderUserInfo(ctx),
		GroupInfo:      groupInfo,
		MergedMsgs:     upMsg.MergedMsgs,
		MemberCount:    int32(memberCount),
		PushData:       upMsg.PushData,
	}

	commonservices.SubGroupMsg(ctx, msgId, downMsg4Sendbox)

	//check merged msg
	if commonservices.IsMergedMsg(upMsg.Flags) && upMsg.MergedMsgs != nil && len(upMsg.MergedMsgs.Msgs) > 0 {
		bases.AsyncRpcCall(ctx, "merge_msgs", msgId, &pbobjs.MergeMsgReq{
			ParentMsgId: msgId,
			MergedMsgs:  upMsg.MergedMsgs,
		})
	}

	if !commonservices.IsStateMsg(upMsg.Flags) {
		//save history msg
		commonservices.SaveHistoryMsg(ctx, bases.GetRequesterIdFromCtx(ctx), groupId, pbobjs.ChannelType_Group, downMsg, memberCount)
	}

	if len(memberIds) > 0 {
		//statistic
		commonservices.ReportDispatchMsg(appkey, pbobjs.ChannelType_Group, int64(len(memberIds)))
		Dispatch2Message(ctx, groupId, memberIds, downMsg)
	}

	return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq, upMsg.ClientUid, int32(memberCount)
}

func GetGroupInfo4Msg(ctx context.Context, groupId string) *pbobjs.GroupInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	groupInfo, exist := GetGroupInfoFromCache(ctx, appkey, groupId)
	if exist && groupInfo != nil {
		retGrpInfo := &pbobjs.GroupInfo{
			GroupId:       groupId,
			GroupName:     groupInfo.GroupName,
			GroupPortrait: groupInfo.GroupPortrait,
			IsMute:        groupInfo.IsMute,
			UpdatedTime:   groupInfo.UpdatedTime,
			ExtFields:     []*pbobjs.KvItem{},
		}
		for k, v := range groupInfo.ExtFields {
			retGrpInfo.ExtFields = append(retGrpInfo.ExtFields, &pbobjs.KvItem{
				Key:   k,
				Value: v,
			})
		}
		return retGrpInfo
	}
	return &pbobjs.GroupInfo{
		GroupId: groupId,
	}
}

func UpdateMentionedUserInfo(ctx context.Context, upMsg *pbobjs.UpMsg) {
	if upMsg != nil && upMsg.MentionInfo != nil {
		if upMsg.MentionInfo.MentionType == pbobjs.MentionType_AllAndSomeone || upMsg.MentionInfo.MentionType == pbobjs.MentionType_Someone {
			for _, userInfo := range upMsg.MentionInfo.TargetUsers {
				uinfo := commonservices.GetTargetDisplayUserInfo(ctx, userInfo.UserId)
				if uinfo != nil {
					userInfo.Nickname = uinfo.Nickname
					userInfo.UpdatedTime = uinfo.UpdatedTime
					userInfo.UserPortrait = uinfo.UserPortrait
					userInfo.ExtFields = uinfo.ExtFields
				}
			}
		}
	}
}

func checkGroupExist(ctx context.Context, groupId string) bool {
	appkey := bases.GetAppKeyFromCtx(ctx)
	_, exist := GetGroupInfoFromCache(ctx, appkey, groupId)
	return exist
}

func checkIsMember(ctx context.Context, groupId, userId string) bool {
	appkey := bases.GetAppKeyFromCtx(ctx)
	memberContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	if exist && memberContainer != nil && memberContainer.Members != nil {
		memberMap := memberContainer.CheckGroupMembers([]string{userId})
		if _, exist := memberMap[userId]; exist {
			return true
		}
	}
	return false
}

func checkGroupIsMute(ctx context.Context, groupId string) bool {
	appkey := bases.GetAppKeyFromCtx(ctx)
	groupInfo, exist := GetGroupInfoFromCache(ctx, appkey, groupId)
	if exist && groupInfo.IsMute > 0 {
		return true
	}
	return false
}

func checkGroupMemberIsMute(ctx context.Context, groupId, memberId string) bool {
	appkey := bases.GetAppKeyFromCtx(ctx)
	groupContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	if exist {
		member := groupContainer.GetMember(memberId)
		if member != nil && member.IsMute > 0 {
			if member.MuteEndAt <= 0 {
				return true
			} else if member.MuteEndAt > time.Now().UnixMilli() {
				return true
			}
		}
	}
	return false
}

func checkGroupMemberIsAllow(ctx context.Context, groupId, memberId string) bool {
	appkey := bases.GetAppKeyFromCtx(ctx)
	groupContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	if exist {
		member := groupContainer.GetMember(memberId)
		if member != nil && member.IsAllow > 0 {
			return true
		}
	}
	return false
}

func getMembersExceptMe(ctx context.Context, groupId string) []string {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	groupContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	memberIds := []string{}
	if exist {
		memberMap := groupContainer.GetMemberMap()
		for memberId := range memberMap {
			if memberId != userId {
				memberIds = append(memberIds, memberId)
			}
		}
	}
	return memberIds
}

func ImportGroupHisMsg(ctx context.Context, msg *pbobjs.UpMsg) {
	groupId := bases.GetTargetIdFromCtx(ctx)
	senderId := bases.GetRequesterIdFromCtx(ctx)
	groupInfo := GetGroupInfo4Msg(ctx, groupId)
	memberIds := getMembersExceptMe(ctx, groupId)
	memberCount := len(memberIds)
	msgId := tools.GenerateMsgId(msg.MsgTime, int32(pbobjs.ChannelType_Group), groupId)
	// downMsg4Sendbox := &pbobjs.DownMsg{
	// 	SenderId:       senderId,
	// 	TargetId:       groupId,
	// 	ChannelType:    pbobjs.ChannelType_Group,
	// 	MsgType:        msg.MsgType,
	// 	MsgContent:     msg.MsgContent,
	// 	MsgId:          msgId,
	// 	MsgSeqNo:       -1,
	// 	MsgTime:        msg.MsgTime,
	// 	Flags:          msg.Flags,
	// 	IsSend:         true,
	// 	TargetUserInfo: commonservices.GetSenderUserInfo(ctx),
	// 	GroupInfo:      groupInfo,
	// 	MemberCount:    int32(memberCount),
	// }
	// if commonservices.IsStoreMsg(msg.Flags) {
	// 	//add conver for sender
	// 	commonservices.BatchSaveConversations(ctx, []string{senderId}, downMsg4Sendbox)
	// }
	downMsg := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       groupId,
		ChannelType:    pbobjs.ChannelType_Group,
		MsgType:        msg.MsgType,
		MsgContent:     msg.MsgContent,
		MsgId:          msgId,
		MsgSeqNo:       -1,
		MsgTime:        msg.MsgTime,
		Flags:          msg.Flags,
		TargetUserInfo: commonservices.GetSenderUserInfo(ctx),
		GroupInfo:      groupInfo,
		MemberCount:    int32(memberCount),
	}
	if commonservices.IsStoreMsg(msg.Flags) {
		//add hismsg
		commonservices.SaveHistoryMsg(ctx, senderId, groupId, pbobjs.ChannelType_Group, downMsg, memberCount)
		//add conver for receivers
		// commonservices.BatchSaveConversations(ctx, memberIds, downMsg)
	}

}
