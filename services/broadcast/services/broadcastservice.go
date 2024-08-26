package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"time"
)

func BroadcastMsg(ctx context.Context, msg *pbobjs.UpMsg) (errs.IMErrorCode, string, int64, int64) {
	senderId := bases.GetRequesterIdFromCtx(ctx)

	converId := commonservices.GetConversationId(senderId, senderId, pbobjs.ChannelType_BroadCast)
	msgConverCache := commonservices.GetMsgConverCache(ctx, converId, pbobjs.ChannelType_BroadCast)
	msgId, sendTime, msgSeq := msgConverCache.GenerateMsgId(converId, pbobjs.ChannelType_BroadCast, time.Now().UnixMilli(), msg.Flags)

	downMsg := &pbobjs.DownMsg{
		SenderId:    senderId,
		TargetId:    senderId,
		ChannelType: pbobjs.ChannelType_BroadCast,
		MsgType:     msg.MsgType,
		MsgId:       msgId,
		MsgSeqNo:    msgSeq,
		UnreadIndex: msgSeq,
		MsgTime:     sendTime,
		MsgContent:  msg.MsgContent,
		Flags:       msg.Flags,
	}
	//save to history msg
	if !commonservices.IsStateMsg(msg.Flags) {
		commonservices.SaveHistoryMsg(ctx, senderId, "", pbobjs.ChannelType_BroadCast, downMsg, 0)
	}

	//save to brd inbox
	bases.AsyncRpcCall(ctx, "brd_inbox", senderId, downMsg)

	//broadcast to all nodes of message.
	bases.Broadcast(ctx, "brd_append", downMsg)

	return errs.IMErrorCode_SUCCESS, msgId, sendTime, msgSeq
}
