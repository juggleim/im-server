package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
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

func QryFriendInfos(ctx context.Context, req *pbobjs.FriendIdsReq) (errs.IMErrorCode, *pbobjs.FriendInfos) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	ret := &pbobjs.FriendInfos{
		Items: []*pbobjs.FriendInfo{},
	}
	for _, friendId := range req.FriendIds {
		info := &pbobjs.FriendInfo{
			FriendId: friendId,
			IsFriend: false,
		}
		friStatus := friendcache.GetFriendStatus(appkey, userId, friendId)
		if friStatus != nil && friStatus.IsFriend {
			info.IsFriend = true
			info.FriendDisplayName = friStatus.FriendDisplayName
			info.UpdatedTime = friStatus.UpdatedTime
		}
		ret.Items = append(ret.Items, info)
	}
	return errs.IMErrorCode_SUCCESS, ret
}
