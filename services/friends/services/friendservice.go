package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
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
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	var startId int64 = 0
	if req.Offset != "" {
		startId, _ = tools.DecodeInt(req.Offset)
	}
	ret := &pbobjs.QryFriendsResp{
		Items: []*pbobjs.FriendMember{},
	}
	rels, err := storage.QueryFriendRels(appkey, userId, startId, req.Limit)
	if err == nil {
		for _, rel := range rels {
			ret.Offset, _ = tools.EncodeInt(rel.ID)
			ret.Items = append(ret.Items, &pbobjs.FriendMember{
				FriendId: rel.FriendId,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func CheckFriends(ctx context.Context, req *pbobjs.CheckFriendsReq) (errs.IMErrorCode, *pbobjs.CheckFriendsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	ret := &pbobjs.CheckFriendsResp{
		CheckResults: make(map[string]bool),
	}
	if len(req.FriendIds) <= 0 {
		return errs.IMErrorCode_SUCCESS, ret
	}
	for _, friendId := range req.FriendIds {
		ret.CheckResults[friendId] = false
	}
	storage := storages.NewFriendRelStorage()
	rels, err := storage.QueryFriendRelsByFriendIds(appkey, userId, req.FriendIds)
	if err == nil {
		for _, rel := range rels {
			ret.CheckResults[rel.FriendId] = true
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}
