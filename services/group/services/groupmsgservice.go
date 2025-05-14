package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/interceptors"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/logmanager/msglogs"
	"time"
)

func SendGroupMsg(ctx context.Context, upMsg *pbobjs.UpMsg) (errs.IMErrorCode, string, int64, int64, string, int32, *pbobjs.DownMsg) {
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
			return errs.IMErrorCode_GROUP_NOTGROUPMEMBER, msgId, sendTime, 0, upMsg.ClientUid, 0, nil
		}
		//check group member mute
		if checkGroupMemberIsMute(ctx, groupId, senderId) {
			sendTime := time.Now().UnixMilli()
			msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Group), groupId)
			return errs.IMErrorCode_GROUP_GROUPMEMBERMUTE, msgId, sendTime, 0, upMsg.ClientUid, 0, nil
		}
		//check group mute
		if checkGroupIsMute(ctx, groupId) {
			//check group member allow
			if !checkGroupMemberIsAllow(ctx, groupId, senderId) {
				sendTime := time.Now().UnixMilli()
				msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Group), groupId)
				return errs.IMErrorCode_GROUP_GROUPMUTE, msgId, sendTime, 0, upMsg.ClientUid, 0, nil
			}
		}
	}

	//check msg interceptor
	var modifiedMsg *pbobjs.DownMsg
	result, interceptorCode := commonservices.CheckMsgInterceptor(ctx, senderId, groupId, pbobjs.ChannelType_Group, upMsg)
	if result == interceptors.InterceptorResult_Reject {
		sendTime := time.Now().UnixMilli()
		msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Group), groupId)
		if interceptorCode == 0 {
			return errs.IMErrorCode_MSG_Hit_Sensitive, msgId, sendTime, 0, upMsg.ClientUid, 0, nil
		} else {
			return errs.IMErrorCode(interceptorCode), msgId, sendTime, 0, upMsg.ClientUid, 0, nil
		}
	} else if result == interceptors.InterceptorResult_Replace {
		modifiedMsg = &pbobjs.DownMsg{
			MsgType:    upMsg.MsgType,
			MsgContent: upMsg.MsgContent,
		}
	} else if result == interceptors.InterceptorResult_Silent {
		sendTime := time.Now().UnixMilli()
		msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Group), groupId)
		return errs.IMErrorCode_SUCCESS, msgId, sendTime, 0, upMsg.ClientUid, 0, nil
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
			return errs.IMErrorCode_SUCCESS, oldAck.MsgId, oldAck.MsgTime, oldAck.MsgSeq, upMsg.ClientUid, 0, nil
		}
	} else {
		upMsg.ClientUid = tools.GenerateUUIDShort22()
	}

	groupInfo := GetGroupInfo4Msg(ctx, groupId)

	//update mentioned user's info
	UpdateMentionedUserInfo(ctx, upMsg)
	//get group member info
	grpMemberInfo := getGroupMemberInfo(ctx, groupId, senderId)

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
		SearchText:     upMsg.SearchText,
		GrpMemberInfo:  grpMemberInfo,
	}
	commonservices.Save2Sendbox(ctx, downMsg4Sendbox)
	msglogs.LogMsg(ctx, downMsg4Sendbox)

	if bases.GetOnlySendboxFromCtx(ctx) {
		return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq, upMsg.ClientUid, int32(memberCount), modifiedMsg
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
		SearchText:     upMsg.SearchText,
		GrpMemberInfo:  grpMemberInfo,
	}

	commonservices.SubGroupMsg(ctx, msgId, downMsg4Sendbox)

	//check merged msg
	if msgdefines.IsMergedMsg(upMsg.Flags) && upMsg.MergedMsgs != nil && len(upMsg.MergedMsgs.Msgs) > 0 {
		bases.AsyncRpcCall(ctx, "merge_msgs", msgId, &pbobjs.MergeMsgReq{
			ParentMsgId: msgId,
			MergedMsgs:  upMsg.MergedMsgs,
		})
	}

	if !msgdefines.IsStateMsg(upMsg.Flags) {
		//save history msg
		commonservices.SaveHistoryMsg(ctx, bases.GetRequesterIdFromCtx(ctx), groupId, pbobjs.ChannelType_Group, downMsg, memberCount)
	}

	if len(memberIds) > 0 {
		//statistic
		commonservices.ReportDispatchMsg(appkey, pbobjs.ChannelType_Group, int64(len(memberIds)))
		Dispatch2Message(ctx, groupId, memberIds, downMsg)
	}

	return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq, upMsg.ClientUid, int32(memberCount), modifiedMsg
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

func getGroupMemberInfo(ctx context.Context, groupId, memberId string) *pbobjs.GrpMemberInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	memberAtts := GetGrpMemberAttsFromCache(ctx, appkey, groupId, memberId)
	if memberAtts != nil && len(memberAtts.ExtFields) > 0 {
		ret := &pbobjs.GrpMemberInfo{
			ExtFields:   []*pbobjs.KvItem{},
			UpdatedTime: memberAtts.UpdatedTime,
		}
		for key, val := range memberAtts.ExtFields {
			if key == string(commonservices.AttItemKey_GrpDisplayName) {
				ret.GrpDisplayName = val
			} else {
				ret.ExtFields = append(ret.ExtFields, &pbobjs.KvItem{
					Key:   key,
					Value: val,
				})
			}
		}
		return ret
	}
	return nil
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
	if msgdefines.IsStoreMsg(msg.Flags) {
		//add hismsg
		commonservices.SaveHistoryMsg(ctx, senderId, groupId, pbobjs.ChannelType_Group, downMsg, memberCount)
		//add conver for receivers
		// commonservices.BatchSaveConversations(ctx, memberIds, downMsg)
	}

}
