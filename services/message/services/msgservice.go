package services

import (
	"context"
	"time"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/interceptors"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/logmanager/msglogs"

	"google.golang.org/protobuf/proto"
)

func SendPrivateMsg(ctx context.Context, senderId, receiverId string, upMsg *pbobjs.UpMsg) (errs.IMErrorCode, string, int64, int64, string, *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	converId := commonservices.GetConversationId(senderId, receiverId, pbobjs.ChannelType_Private)
	//statistic
	commonservices.ReportUpMsg(appkey, pbobjs.ChannelType_Private, 1)
	//check block user
	blockUser := GetBlockUserItem(appkey, receiverId, senderId)
	if blockUser.IsBlock {
		sendTime := time.Now().UnixMilli()
		msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Private), receiverId)
		return errs.IMErrorCode_MSG_BLOCK, msgId, sendTime, 0, upMsg.ClientUid, nil
	}
	//check msg interceptor
	var modifiedMsg *pbobjs.DownMsg = nil
	result, interceptorCode := commonservices.CheckMsgInterceptor(ctx, senderId, receiverId, pbobjs.ChannelType_Private, upMsg)
	if result == interceptors.InterceptorResult_Reject {
		sendTime := time.Now().UnixMilli()
		msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Private), receiverId)
		if interceptorCode == 0 {
			return errs.IMErrorCode_MSG_Hit_Sensitive, msgId, sendTime, 0, upMsg.ClientUid, nil
		} else {
			return errs.IMErrorCode(interceptorCode), msgId, sendTime, 0, upMsg.ClientUid, nil
		}
	} else if result == interceptors.InterceptorResult_Replace {
		modifiedMsg = &pbobjs.DownMsg{
			MsgType:    upMsg.MsgType,
			MsgContent: upMsg.MsgContent,
		}
	} else if result == interceptors.InterceptorResult_Silent {
		sendTime := time.Now().UnixMilli()
		msgId := tools.GenerateMsgId(sendTime, int32(pbobjs.ChannelType_Private), receiverId)
		return errs.IMErrorCode_SUCCESS, msgId, sendTime, 0, upMsg.ClientUid, nil
	}

	msgConverCache := commonservices.GetMsgConverCache(ctx, converId, pbobjs.ChannelType_Private)
	msgId, sendTime, msgSeq := msgConverCache.GenerateMsgId(converId, pbobjs.ChannelType_Private, time.Now().UnixMilli(), upMsg.Flags)
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
			return errs.IMErrorCode_SUCCESS, oldAck.MsgId, oldAck.MsgTime, oldAck.MsgSeq, upMsg.ClientUid, nil
		}
	} else {
		upMsg.ClientUid = tools.GenerateUUIDShort22()
	}

	downMsg4Sendbox := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       receiverId,
		ChannelType:    pbobjs.ChannelType_Private,
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
		TargetUserInfo: commonservices.GetTargetDisplayUserInfo(ctx, receiverId),
		MergedMsgs:     upMsg.MergedMsgs,
		PushData:       upMsg.PushData,
		SearchText:     upMsg.SearchText,
	}
	commonservices.Save2Sendbox(ctx, downMsg4Sendbox)
	msglogs.LogMsg(ctx, downMsg4Sendbox)

	if bases.GetOnlySendboxFromCtx(ctx) {
		return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq, upMsg.ClientUid, modifiedMsg
	}
	commonservices.SubPrivateMsg(ctx, msgId, downMsg4Sendbox)

	downMsg := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       senderId,
		ChannelType:    pbobjs.ChannelType_Private,
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
		MergedMsgs:     upMsg.MergedMsgs,
		PushData:       upMsg.PushData,
		SearchText:     upMsg.SearchText,
	}

	//check merged msg
	if msgdefines.IsMergedMsg(upMsg.Flags) && upMsg.MergedMsgs != nil && len(upMsg.MergedMsgs.Msgs) > 0 {
		bases.AsyncRpcCall(ctx, "merge_msgs", msgId, &pbobjs.MergeMsgReq{
			ParentMsgId: msgId,
			MergedMsgs:  upMsg.MergedMsgs,
		})
	}

	//save history msg
	if msgdefines.IsStoreMsg(upMsg.Flags) {
		commonservices.SaveHistoryMsg(ctx, senderId, receiverId, pbobjs.ChannelType_Private, downMsg, 0)
	}

	//dispatch to receiver
	if senderId != receiverId {
		dispatchMsg(ctx, receiverId, downMsg)
		commonservices.ReportDispatchMsg(appkey, pbobjs.ChannelType_Private, 1)
	}

	return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq, upMsg.ClientUid, modifiedMsg
}

func dispatchMsg(ctx context.Context, receiverId string, msg *pbobjs.DownMsg) {
	data, _ := tools.PbMarshal(msg)
	bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_UserPub,
		AppKey:       bases.GetAppKeyFromCtx(ctx),
		Session:      bases.GetSessionFromCtx(ctx),
		Method:       "msg_dispatch",
		RequesterId:  bases.GetRequesterIdFromCtx(ctx),
		ReqIndex:     bases.GetSeqIndexFromCtx(ctx),
		Qos:          bases.GetQosFromCtx(ctx),
		AppDataBytes: data,
		TargetId:     receiverId,
	})
}

func MsgOrNtf(ctx context.Context, targetId string, downMsg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userStatus := GetUserStatus(appkey, targetId)
	if userStatus.IsOnline() {
		isNtf := GetUserStatus(appkey, targetId).CheckNtfWithSwitch()
		hasPush := false
		if userStatus.OpenPushSwitch() {
			hasPush = true
			SendPush(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, downMsg)
		}
		if isNtf { //发送通知
			logs.WithContext(ctx).Infof("ntf target_id:%s", targetId)
			rpcNtf := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "ntf", &pbobjs.Notify{
				Type:     pbobjs.NotifyType_Msg,
				SyncTime: downMsg.MsgTime,
			})
			rpcNtf.Qos = 0
			bases.UnicastRouteWithNoSender(rpcNtf)
			bases.UnicastRouteWithCallback(rpcNtf, &SendMsgAckActor{
				appkey:      appkey,
				senderId:    bases.GetRequesterIdFromCtx(ctx),
				targetId:    targetId,
				channelType: downMsg.ChannelType,
				Msg:         downMsg,
				ctx:         ctx,
				IsNotify:    isNtf,
				HasPush:     hasPush,
			}, 5*time.Second)
		} else {
			logs.WithContext(ctx).Infof("msg target_id:%s", targetId)
			rpcMsg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "msg", downMsg)
			rpcMsg.MsgId = downMsg.MsgId
			rpcMsg.MsgSendTime = downMsg.MsgTime
			bases.UnicastRouteWithCallback(rpcMsg, &SendMsgAckActor{
				appkey:      appkey,
				senderId:    bases.GetRequesterIdFromCtx(ctx),
				targetId:    targetId,
				channelType: downMsg.ChannelType,
				Msg:         downMsg,
				ctx:         ctx,
				IsNotify:    isNtf,
				HasPush:     hasPush,
			}, 5*time.Second)
		}
	} else { //for push
		SendPush(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, downMsg)
	}
}

func getTargetUserLanguage(ctx context.Context, userId string) string {
	appkey := bases.GetAppKeyFromCtx(ctx)
	language := "en_US"
	appinfo, exist := commonservices.GetAppInfo(appkey)
	if exist {
		language = appinfo.PushLanguage
	}
	uSetting := commonservices.GetTargetUserSettings(ctx, userId)
	if uSetting != nil && uSetting.Language != "" {
		language = uSetting.Language
	}
	return language
}

type SendMsgAckActor struct {
	actorsystem.UntypedActor
	appkey      string
	senderId    string
	targetId    string
	channelType pbobjs.ChannelType
	// pushData    *pbobjs.PushData
	Msg      *pbobjs.DownMsg
	ctx      context.Context
	IsNotify bool
	HasPush  bool
}

func (actor *SendMsgAckActor) OnReceive(ctx context.Context, input proto.Message) {
	if rpcMsg, ok := input.(*pbobjs.RpcMessageWraper); ok {
		data := rpcMsg.AppDataBytes
		onlineStatus := &pbobjs.OnlineStatus{}
		err := tools.PbUnMarshal(data, onlineStatus)
		if err == nil {
			logs.WithContext(actor.ctx).Infof("target_id:%s\tonline_type:%d", actor.targetId, onlineStatus.Type)
			if onlineStatus.Type == pbobjs.OnlineType_Offline { //receiver is offline
				RecordUserOnlineStatus(actor.appkey, actor.targetId, false, 0)
				if !actor.HasPush {
					SendPush(actor.ctx, actor.senderId, actor.targetId, actor.Msg)
				}
			}
		}
	}
}

func GetPushData(ctx context.Context, msg *pbobjs.DownMsg, pushLanguage string) *pbobjs.PushData {
	if msg == nil {
		return nil
	}
	var (
		title  string
		prefix string
	)
	nickName := msg.TargetUserInfo.GetNickname()
	if msg.ChannelType == pbobjs.ChannelType_Group {
		title = msg.GroupInfo.GroupName
		if nickName != "" {
			prefix = nickName + ": "
		}
	} else {
		title = nickName
	}
	retPushData := &pbobjs.PushData{}
	if msg.PushData != nil {
		retPushData.Title = msg.PushData.Title
		retPushData.PushId = msg.PushData.PushId
		retPushData.PushText = msg.PushData.PushText
		retPushData.PushExtraData = msg.PushData.PushExtraData

		retPushData.IsVoip = msg.PushData.IsVoip
		retPushData.RtcRoomId = msg.PushData.RtcRoomId
		retPushData.RtcInviterId = msg.PushData.RtcInviterId
		retPushData.RtcRoomType = msg.PushData.RtcRoomType
		retPushData.RtcMediaType = msg.PushData.RtcMediaType
	}
	if retPushData.Title == "" {
		retPushData.Title = title
	}
	if retPushData.PushText != "" {
		pushText := retPushData.PushText
		//handle template
		pushText = TemplateI18nAssign(ctx, pushText, pushLanguage)
		retPushData.PushText = prefix + pushText
	} else {
		if msg.MsgType == msgdefines.InnerMsgType_Text {
			txtMsg := &msgdefines.TextMsg{}
			err := tools.JsonUnMarshal(msg.MsgContent, txtMsg)
			pushText := txtMsg.Content
			charArr := []rune(pushText)
			if len(charArr) > 20 {
				pushText = string(charArr[:20]) + "..."
			}
			if err == nil {
				retPushData.PushText = prefix + pushText
			} else {
				retPushData.PushText = prefix + commonservices.GetInnerI18nStr(pushLanguage, commonservices.PlaceholderKey_Text, "[Text]")
			}
		} else if msg.MsgType == msgdefines.InnerMsgType_Img {
			retPushData.PushText = prefix + commonservices.GetInnerI18nStr(pushLanguage, commonservices.PlaceholderKey_Image, "[Image]")
		} else if msg.MsgType == msgdefines.InnerMsgType_Voice {
			retPushData.PushText = prefix + commonservices.GetInnerI18nStr(pushLanguage, commonservices.PlaceholderKey_Voice, "[Voice]")
		} else if msg.MsgType == msgdefines.InnerMsgType_File {
			retPushData.PushText = prefix + commonservices.GetInnerI18nStr(pushLanguage, commonservices.PlaceholderKey_File, "[File]")
		} else if msg.MsgType == msgdefines.InnerMsgType_Video {
			retPushData.PushText = prefix + commonservices.GetInnerI18nStr(pushLanguage, commonservices.PlaceholderKey_Video, "[Video]")
		} else if msg.MsgType == msgdefines.InnerMsgType_Merge {
			retPushData.PushText = prefix + commonservices.GetInnerI18nStr(pushLanguage, commonservices.PlaceholderKey_Merge, "[Merge]")
		} else if msg.MsgType == msgdefines.InnerMsgType_VoiceCall {
			retPushData.PushText = prefix + commonservices.GetInnerI18nStr(pushLanguage, commonservices.PlaceholderKey_RtcCall, "invites you to a voice call")
		} else {
			return nil
		}
	}

	//add internal fields
	retPushData.MsgId = msg.MsgId
	retPushData.SenderId = msg.SenderId
	retPushData.ConverId = msg.TargetId
	retPushData.ChannelType = msg.ChannelType
	return retPushData
}

func (actor *SendMsgAckActor) CreateInputObj() proto.Message {
	return &pbobjs.RpcMessageWraper{}
}
func (actor *SendMsgAckActor) OnTimeout() {

}

func SendPush(ctx context.Context, senderId, receiverId string, msg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	appInfo, exist := commonservices.GetAppInfo(appkey)
	if exist && appInfo != nil && appInfo.IsOpenPush {
		//check close push threshold
		if msg.ChannelType == pbobjs.ChannelType_Group && appInfo.ClosePushGrpThreshold > 0 && msg.MemberCount > int32(appInfo.ClosePushGrpThreshold) {
			if msg.PushData == nil || msg.PushData.PushLevel == pbobjs.PushLevel_DefaultPuhsLevel {
				return
			}
		}
		//undisturb
		if msgdefines.IsUndisturbMsg(msg.Flags) {
			if msg.PushData == nil || msg.PushData.PushLevel < pbobjs.PushLevel_IgnoreUndisturb {
				return
			}
		}
		pushData := GetPushData(ctx, msg, getTargetUserLanguage(ctx, receiverId))
		if pushData != nil {
			//badge
			userStatus := GetUserStatus(appkey, receiverId)
			pushData.Badge = userStatus.BadgeIncr()
			if userStatus.CanPush > 0 {
				pushRpc := bases.CreateServerPubWraper(ctx, senderId, receiverId, "push", pushData)
				bases.UnicastRouteWithNoSender(pushRpc)
			}
		}
	}
}

func ImportPrivateHisMsg(ctx context.Context, senderId, targetId string, msg *pbobjs.UpMsg) {
	msgId := tools.GenerateMsgId(msg.MsgTime, int32(pbobjs.ChannelType_Private), targetId)
	/*
		downMsg4Sendbox := &pbobjs.DownMsg{
			SenderId:    senderId,
			TargetId:    targetId,
			ChannelType: pbobjs.ChannelType_Private,
			MsgType:     msg.MsgType,
			MsgContent:  msg.MsgContent,
			MsgId:       msgId,
			MsgSeqNo:    -1,
			MsgTime:     msg.MsgTime,
			Flags:       msg.Flags,
			IsSend:      true,
			//TargetUserInfo: commonservices.GetTargetDisplayUserInfo(ctx, targetId),
		}*/
	// add conver for sender
	// if commonservices.IsStoreMsg(msg.Flags) {
	// 	commonservices.BatchSaveConversations(ctx, []string{senderId}, downMsg4Sendbox)
	// }

	downMsg := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       senderId,
		ChannelType:    pbobjs.ChannelType_Private,
		MsgType:        msg.MsgType,
		MsgContent:     msg.MsgContent,
		MsgId:          msgId,
		MsgSeqNo:       -1,
		MsgTime:        msg.MsgTime,
		Flags:          msg.Flags,
		TargetUserInfo: commonservices.GetSenderUserInfo(ctx),
	}
	//add hismsg
	if msgdefines.IsStoreMsg(msg.Flags) {
		commonservices.SaveHistoryMsg(ctx, senderId, targetId, pbobjs.ChannelType_Private, downMsg, 0)

		//add conver for receiver
		//commonservices.BatchSaveConversations(ctx, []string{targetId}, downMsg)
	}
}
