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

func SendChatroomMsg(ctx context.Context, upMsg *pbobjs.UpMsg) (errs.IMErrorCode, string, int64, int64) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	chatId := bases.GetTargetIdFromCtx(ctx)

	container, exist := getChatroomContainer(appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST, "", time.Now().UnixMilli(), 0
	}
	isFromApi := bases.GetIsFromApiFromCtx(ctx)
	if !isFromApi && !container.CheckMemberExist(userId) {
		return errs.IMErrorCode_CHATROOM_NOTMEMBER, "", time.Now().UnixMilli(), 0
	}
	//check mute
	if container.IsMute {
		if !container.CheckMemberAllow(userId) {
			return errs.IMErrorCode_CHATROOM_MUTE, "", time.Now().UnixMilli(), 0
		}
	} else {
		if container.CheckMemberMute(userId) {
			return errs.IMErrorCode_CHATROOM_MUTE, "", time.Now().UnixMilli(), 0
		}
	}

	msgTime, msgSeq := container.GetMsgTimeSeq(time.Now().UnixMilli(), upMsg.Flags)
	msgId := tools.GenerateMsgId(msgTime, int32(pbobjs.ChannelType_Chatroom), chatId)

	downMsg := &pbobjs.DownMsg{
		SenderId:    userId,
		TargetId:    chatId,
		ChannelType: pbobjs.ChannelType_Chatroom,
		MsgType:     upMsg.MsgType,
		MsgId:       msgId,
		MsgSeqNo:    msgSeq,
		MsgContent:  upMsg.MsgContent,
		MsgTime:     msgTime,
		Flags:       upMsg.Flags,
		ClientUid:   upMsg.ClientUid,
		MentionInfo: upMsg.MentionInfo,
		ReferMsg:    commonservices.FillReferMsg(ctx, upMsg),
	}
	bases.Broadcast(ctx, "c_msgs_dispatch", &pbobjs.DownMsgSet{
		Msgs: []*pbobjs.DownMsg{
			downMsg,
		},
	})
	return errs.IMErrorCode_SUCCESS, msgId, msgTime, msgSeq
}
