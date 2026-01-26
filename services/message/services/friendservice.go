package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/friendcache"
)

func SyncFriendRels(ctx context.Context, req *pbobjs.FriendIdsReq) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	for _, friendId := range req.FriendIds {
		friendcache.RemoveFriendStatus(appkey, userId, friendId)
	}
}
