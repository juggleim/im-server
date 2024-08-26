package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
)

func UpdLatestMsg(ctx context.Context, userId string, req *pbobjs.UpdLatestMsgReq) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	saveConversationByCache(&ConversationCacheItem{
		Appkey:      appkey,
		UserId:      userId,
		TargetId:    req.TargetId,
		ChannelType: req.ChannelType,
		LatestMsgId: req.LatestMsgId,
		LatestMsg:   req.Msg,
		OnlyUpdMsg:  true,
	})
}
