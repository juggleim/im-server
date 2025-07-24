package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"time"
)

func SendGroupCastMsg(ctx context.Context, upMsg *pbobjs.UpMsg) (errs.IMErrorCode, string, int64, int64) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	grpCastId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)

	//statistic
	commonservices.ReportUpMsg(appkey, pbobjs.ChannelType_GroupCast, 1)

	converId := commonservices.GetConversationId(userId, grpCastId, pbobjs.ChannelType_GroupCast)
	msgConverCache := commonservices.GetMsgConverCache(ctx, converId, "", pbobjs.ChannelType_GroupCast)
	msgId, sendTime, msgSeq := msgConverCache.GenerateMsgId(converId, pbobjs.ChannelType_GroupCast, time.Now().UnixMilli(), upMsg.Flags)

	downMsg4Sendbox := &pbobjs.DownMsg{
		SenderId:       userId,
		TargetId:       grpCastId,
		ChannelType:    pbobjs.ChannelType_GroupCast,
		MsgType:        upMsg.MsgType,
		MsgId:          msgId,
		MsgSeqNo:       msgSeq,
		MsgContent:     upMsg.MsgContent,
		MsgTime:        sendTime,
		Flags:          upMsg.Flags,
		ClientUid:      upMsg.ClientUid,
		IsSend:         true,
		TargetUserInfo: bases.GetSenderInfoFromCtx(ctx),
	}
	if !msgdefines.IsStateMsg(upMsg.Flags) {
		//save msg to sendbox for sender
		//record conversation for sender
		commonservices.Save2Sendbox(ctx, downMsg4Sendbox)
		//save history msg
		commonservices.SaveHistoryMsg(ctx, userId, grpCastId, pbobjs.ChannelType_GroupCast, downMsg4Sendbox, 0)
	}
	return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq
}
