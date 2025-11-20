package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/friendmanager/storages"
	"im-server/services/friendmanager/storages/models"
)

func AddFriends(ctx context.Context, req *pbobjs.FriendMembersReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	friendRels := []models.FriendRel{}
	friendIds := []string{}
	for _, friendMember := range req.FriendMembers {
		friendIds = append(friendIds, friendMember.FriendId)
		friendRels = append(friendRels, models.FriendRel{
			AppKey:      appkey,
			UserId:      userId,
			FriendId:    friendMember.FriendId,
			OrderTag:    friendMember.OrderTag,
			DisplayName: friendMember.DisplayName,
		})
	}
	err := storage.BatchUpsert(friendRels)
	if err != nil {
		logs.WithContext(ctx).Error(err.Error())
	}
	//sync to cache
	syncFriendRels(ctx, userId, friendIds)
	return errs.IMErrorCode_SUCCESS
}

func DelFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	err := storage.BatchDelete(appkey, userId, req.FriendIds)
	if err != nil {
		logs.WithContext(ctx).Error(err.Error())
	}
	//sync to cache
	syncFriendRels(ctx, userId, req.FriendIds)
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
	rels, err := storage.QueryFriendRels(appkey, userId, startId, req.Limit, req.Order > 0)
	if err == nil {
		for _, rel := range rels {
			ret.Offset, _ = tools.EncodeInt(rel.ID)
			ret.Items = append(ret.Items, &pbobjs.FriendMember{
				FriendId:    rel.FriendId,
				OrderTag:    rel.OrderTag,
				DisplayName: rel.DisplayName,
			})
		}
	} else {
		logs.WithContext(ctx).Error(err.Error())
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryFriendsWithPage(ctx context.Context, req *pbobjs.QryFriendsWithPageReq) (errs.IMErrorCode, *pbobjs.QryFriendsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	ret := &pbobjs.QryFriendsResp{
		Items: []*pbobjs.FriendMember{},
	}
	rels, err := storage.QueryFriendRelsWithPage(appkey, userId, req.OrderTag, req.Page, req.Size)
	if err == nil {
		for _, rel := range rels {
			ret.Items = append(ret.Items, &pbobjs.FriendMember{
				FriendId:    rel.FriendId,
				OrderTag:    rel.OrderTag,
				DisplayName: rel.DisplayName,
			})
		}
	} else {
		logs.WithContext(ctx).Error(err.Error())
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func CheckFriends(ctx context.Context, req *pbobjs.CheckFriendsReq) (errs.IMErrorCode, *pbobjs.CheckFriendsResp) {
	// appkey := bases.GetAppKeyFromCtx(ctx)
	// userId := bases.GetTargetIdFromCtx(ctx)
	ret := &pbobjs.CheckFriendsResp{
		CheckResults: make(map[string]bool),
	}
	// if len(req.FriendIds) <= 0 {
	// 	return errs.IMErrorCode_SUCCESS, ret
	// }
	// for _, friendId := range req.FriendIds {
	// 	status := GetFriendStatus(appkey, userId, friendId)
	// 	ret.CheckResults[friendId] = status.IsFriend
	// }
	return errs.IMErrorCode_SUCCESS, ret
}

func syncFriendRels(ctx context.Context, userId string, friendIds []string) {
	if len(friendIds) > 0 {
		bases.AsyncRpcCall(ctx, "sync_friend_rels", userId, &pbobjs.FriendIdsReq{
			FriendIds: friendIds,
		})
	}
}
