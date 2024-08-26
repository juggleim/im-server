package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
)

func HandleChatAttsDispatch(ctx context.Context, req *pbobjs.ChatAtts) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	var attTime int64 = 0
	for _, att := range req.Atts {
		AppendChatAttChg(ctx, appkey, req.ChatId, att)
		if att.AttTime > attTime {
			attTime = att.AttTime
		}
	}
	if attTime > 0 {
		//ntf
		ntfChatMembers(ctx, req.ChatId, bases.GetRequesterIdFromCtx(ctx), attTime, pbobjs.NotifyType_ChatroomAtt)
	}
}

func AppendChatAttChg(ctx context.Context, appkey, chatId string, att *pbobjs.ChatAttItem) {
	container, exist := GetChrmContainer(ctx, appkey, chatId)
	if !exist {
		logs.WithContext(ctx).Errorf("chatroom not exist. chat_id:%s", chatId)
		return
	}
	container.AppendAtt(ctx, att)
}

func SyncChatroomAtts(ctx context.Context, sync *pbobjs.SyncChatroomReq) (errs.IMErrorCode, *pbobjs.SyncChatroomAttResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	chatId := sync.ChatroomId
	userId := bases.GetRequesterIdFromCtx(ctx)
	container, exist := GetChrmContainer(ctx, appkey, chatId)
	if !exist {
		return errs.IMErrorCode_CHATROOM_NOTEXIST, nil
	}
	container.CleanUnread(userId)
	atts, code := container.GetAttsBaseTime(ctx, sync.SyncTime)
	if code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	return errs.IMErrorCode_SUCCESS, &pbobjs.SyncChatroomAttResp{
		Atts: atts,
	}
}
