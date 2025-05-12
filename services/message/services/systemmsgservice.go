package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/logmanager/msglogs"
	"time"
)

func SendSystemMsg(ctx context.Context, senderId, receiverId string, upMsg *pbobjs.UpMsg) (errs.IMErrorCode, string, int64, int64) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	converId := commonservices.GetConversationId(senderId, receiverId, pbobjs.ChannelType_Private)
	//statistic
	commonservices.ReportUpMsg(appkey, pbobjs.ChannelType_System, 1)
	commonservices.ReportDispatchMsg(appkey, pbobjs.ChannelType_System, 1)

	msgConverCache := commonservices.GetMsgConverCache(ctx, converId, pbobjs.ChannelType_System)
	msgId, sendTime, msgSeq := msgConverCache.GenerateMsgId(converId, pbobjs.ChannelType_System, time.Now().UnixMilli(), upMsg.Flags)
	preMsgId := bases.GetMsgIdFromCtx(ctx)
	if preMsgId != "" {
		msgId = preMsgId
	}

	downMsg4Sendbox := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       receiverId,
		ChannelType:    pbobjs.ChannelType_System,
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
	}
	msglogs.LogMsg(ctx, downMsg4Sendbox)
	//send to sender's other device
	// if !commonservices.IsStateMsg(upMsg.Flags) {
	// 	commonservices.Save2Sendbox(ctx, downMsg4Sendbox)
	// }
	if bases.GetOnlySendboxFromCtx(ctx) {
		return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq
	}
	//save msg to inbox for receiver
	downMsg := &pbobjs.DownMsg{
		SenderId:       senderId,
		TargetId:       senderId,
		ChannelType:    pbobjs.ChannelType_System,
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
	}
	//save history msg
	if msgdefines.IsStoreMsg(upMsg.Flags) {
		commonservices.SaveHistoryMsg(ctx, senderId, receiverId, pbobjs.ChannelType_System, downMsg, 0)
	}
	//dispatch to receiver
	if senderId != receiverId {
		dispatchMsg(ctx, receiverId, downMsg)
	}
	return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq
}
