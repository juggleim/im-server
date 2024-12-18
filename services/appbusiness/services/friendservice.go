package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	apiModels "im-server/services/appbusiness/models"
	"im-server/services/appbusiness/storages"
	"im-server/services/appbusiness/storages/models"
	"im-server/services/commonservices"
	"time"

	"google.golang.org/protobuf/proto"
)

func QryFriends(ctx context.Context, req *pbobjs.FriendListReq) (errs.IMErrorCode, *pbobjs.UserObjs) {
	userId := bases.GetRequesterIdFromCtx(ctx)

	code, respObj, err := AppSyncRpcCall(ctx, "qry_friends", userId, userId, &pbobjs.QryFriendsReq{
		Limit:  req.Limit,
		Offset: req.Offset,
	}, func() proto.Message {
		return &pbobjs.QryFriendsResp{}
	})
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	resp := respObj.(*pbobjs.QryFriendsResp)
	ret := &pbobjs.UserObjs{
		Items:  []*pbobjs.UserObj{},
		Offset: resp.Offset,
	}
	for _, rel := range resp.Items {
		friend := commonservices.GetTargetDisplayUserInfo(ctx, rel.FriendId)

		ret.Items = append(ret.Items, &pbobjs.UserObj{
			UserId:   friend.UserId,
			Nickname: friend.Nickname,
			Avatar:   friend.UserPortrait,
		})
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func AddFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	userId := bases.GetRequesterIdFromCtx(ctx)
	for _, friendId := range req.FriendIds {
		AppSyncRpcCall(ctx, "add_friends", userId, userId, &pbobjs.FriendIdsReq{
			FriendIds: []string{friendId},
		}, nil)
		AppSyncRpcCall(ctx, "add_friends", userId, friendId, &pbobjs.FriendIdsReq{
			FriendIds: []string{userId},
		}, nil)
		//send notify msg
		SendFriendNotify(ctx, friendId, &apiModels.FriendNotify{
			Type: 0,
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func DelFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	userId := bases.GetRequesterIdFromCtx(ctx)
	AppSyncRpcCall(ctx, "del_friends", userId, userId, &pbobjs.FriendIdsReq{
		FriendIds: req.FriendIds,
	}, nil)
	return errs.IMErrorCode_SUCCESS
}

func ApplyFriend(ctx context.Context, req *pbobjs.ApplyFriend) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	//check friend relation
	if checkFriend(ctx, req.FriendId, userId) {
		return errs.IMErrorCode_APP_FRIEND_APPLY_REPEATED
	}
	friendUserInfo := commonservices.GetTargetUserInfo(ctx, req.FriendId)
	friendSettings := GetUserSettings(friendUserInfo)
	if friendSettings.FriendVerifyType == pbobjs.FriendVerifyType_DeclineFriend {
		return errs.IMErrorCode_APP_FRIEND_APPLY_DECLINE
	} else if friendSettings.FriendVerifyType == pbobjs.FriendVerifyType_NeedFriendVerify {
		storage := storages.NewFriendApplicationStorage()
		storage.Upsert(models.FriendApplication{
			RecipientId: req.FriendId,
			SponsorId:   userId,
			ApplyTime:   time.Now().UnixMilli(),
			Status:      models.FriendApplicationStatus(models.FriendApplicationStatus_Apply),
			AppKey:      appkey,
		})
		//send notify msg
		SendFriendApplyNotify(ctx, req.FriendId, &apiModels.FriendApplyNotify{
			SponsorId:   userId,
			RecipientId: req.FriendId,
		})
	} else if friendSettings.FriendVerifyType == pbobjs.FriendVerifyType_NoNeedFriendVerify {
		AppSyncRpcCall(ctx, "add_friends", userId, userId, &pbobjs.FriendIdsReq{
			FriendIds: []string{req.FriendId},
		}, nil)
		AppSyncRpcCall(ctx, "add_friends", userId, req.FriendId, &pbobjs.FriendIdsReq{
			FriendIds: []string{userId},
		}, nil)
		//send notify msg
		SendFriendNotify(ctx, req.FriendId, &apiModels.FriendNotify{
			Type: 0,
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func ConfirmFriend(ctx context.Context, req *pbobjs.ConfirmFriend) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewFriendApplicationStorage()
	if req.IsAgree {
		//add friend
		AppSyncRpcCall(ctx, "add_friends", userId, userId, &pbobjs.FriendIdsReq{
			FriendIds: []string{req.SponsorId},
		}, nil)
		AppSyncRpcCall(ctx, "add_friends", userId, req.SponsorId, &pbobjs.FriendIdsReq{
			FriendIds: []string{userId},
		}, nil)
		//send notify msg
		SendFriendNotify(ctx, req.SponsorId, &apiModels.FriendNotify{
			Type: 1,
		})
		storage.UpdateStatus(appkey, req.SponsorId, userId, models.FriendApplicationStatus_Agree)
	} else {
		storage.UpdateStatus(appkey, req.SponsorId, userId, models.FriendApplicationStatus_Decline)
	}
	return errs.IMErrorCode_SUCCESS
}

func checkFriend(ctx context.Context, userId, friendId string) bool {
	results := CheckFriends(ctx, userId, []string{friendId})
	if isFriend, exist := results[friendId]; exist {
		return isFriend
	}
	return false
}

func QryMyFriendApplications(ctx context.Context, req *pbobjs.QryFriendApplicationsReq) (errs.IMErrorCode, *pbobjs.QryFriendApplicationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewFriendApplicationStorage()
	ret := &pbobjs.QryFriendApplicationsResp{
		Items: []*pbobjs.FriendApplicationItem{},
	}
	applications, err := storage.QueryMyApplications(appkey, userId, req.StartTime, int64(req.Count), req.Order > 0)
	if err == nil {
		for _, application := range applications {
			ret.Items = append(ret.Items, &pbobjs.FriendApplicationItem{
				Recipient: &pbobjs.UserObj{
					UserId: application.RecipientId,
				},
				Status:    int32(application.Status),
				ApplyTime: application.ApplyTime,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryMyPendingFriendApplications(ctx context.Context, req *pbobjs.QryFriendApplicationsReq) (errs.IMErrorCode, *pbobjs.QryFriendApplicationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewFriendApplicationStorage()
	ret := &pbobjs.QryFriendApplicationsResp{
		Items: []*pbobjs.FriendApplicationItem{},
	}
	applications, err := storage.QueryPendingApplications(appkey, userId, req.StartTime, int64(req.Count), req.Order > 0)
	if err == nil {
		for _, application := range applications {
			ret.Items = append(ret.Items, &pbobjs.FriendApplicationItem{
				Sponsor: &pbobjs.UserObj{
					UserId: application.SponsorId,
				},
				Status:    int32(application.Status),
				ApplyTime: application.ApplyTime,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryFriendApplications(ctx context.Context, req *pbobjs.QryFriendApplicationsReq) (errs.IMErrorCode, *pbobjs.QryFriendApplicationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewFriendApplicationStorage()
	ret := &pbobjs.QryFriendApplicationsResp{
		Items: []*pbobjs.FriendApplicationItem{},
	}
	applications, err := storage.QueryApplications(appkey, userId, req.StartTime, int64(req.Count), req.Order > 0)
	if err == nil {
		for _, application := range applications {
			item := &pbobjs.FriendApplicationItem{
				Status:    int32(application.Status),
				ApplyTime: application.ApplyTime,
			}
			if userId == application.SponsorId {
				item.IsSponsor = true
				item.TargetUser = GetUser(ctx, application.RecipientId)
			} else {
				item.TargetUser = GetUser(ctx, application.SponsorId)
			}
			ret.Items = append(ret.Items, item)
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func CheckFriends(ctx context.Context, userId string, friendIds []string) map[string]bool {
	ret := make(map[string]bool)
	if len(friendIds) <= 0 {
		return ret
	}
	for _, friend := range friendIds {
		ret[friend] = false
	}
	code, respObj, err := bases.SyncRpcCall(ctx, "check_friends", userId, &pbobjs.CheckFriendsReq{
		FriendIds: friendIds,
	}, func() proto.Message {
		return &pbobjs.CheckFriendsResp{}
	})
	if err == nil && code == errs.IMErrorCode_SUCCESS {
		checkFriendResp := respObj.(*pbobjs.CheckFriendsResp)
		if checkFriendResp != nil {
			for friendId, isFriend := range checkFriendResp.CheckResults {
				ret[friendId] = isFriend
			}
		}
	}
	return ret
}
