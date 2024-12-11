package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/friends/storages"
	"im-server/services/friends/storages/models"
)

func AddFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	friendRels := []models.FriendRel{}
	for _, friendId := range req.FriendIds {
		friendRels = append(friendRels, models.FriendRel{
			AppKey:   appkey,
			UserId:   userId,
			FriendId: friendId,
		})
	}
	storage.BatchUpsert(friendRels)
	return errs.IMErrorCode_SUCCESS
}

func DelFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	storage.BatchDelete(appkey, userId, req.FriendIds)
	return errs.IMErrorCode_SUCCESS
}

func QryFriends(ctx context.Context, req *pbobjs.QryFriendsReq) (errs.IMErrorCode, *pbobjs.QryFriendsResp) {
	return errs.IMErrorCode_SUCCESS, nil
}
