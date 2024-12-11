package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	apiModels "im-server/services/appbusiness/models"
	"im-server/services/commonservices"
	"im-server/services/friends/storages"
	"im-server/services/friends/storages/models"
	"time"

	"google.golang.org/protobuf/proto"
)

func QryFriends(ctx context.Context, req *pbobjs.FriendListReq) (errs.IMErrorCode, *pbobjs.FriendListResp) {
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
	ret := &pbobjs.FriendListResp{
		Items:  []*pbobjs.UserInfo{},
		Offset: resp.Offset,
	}
	for _, rel := range resp.Items {
		ret.Items = append(ret.Items, commonservices.GetTargetDisplayUserInfo(ctx, rel.FriendId))
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

func ApplyFriend(ctx context.Context, req *pbobjs.ApplyFriend) (errs.IMErrorCode, *pbobjs.ApplyFriendResp) {
	resp := &pbobjs.ApplyFriendResp{}
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	//check friend relation
	if checkFriend(ctx, req.FriendId, userId) {
		resp.Reason = pbobjs.ApplyFriendResultReason_ApplyRepeated
		return errs.IMErrorCode_SUCCESS, resp
	}
	friendUserInfo := commonservices.GetTargetUserInfo(ctx, req.FriendId)
	friendSettings := GetUserSettings(friendUserInfo)
	if friendSettings.FriendVerifyType == pbobjs.FriendVerifyType_DeclineFriend {
		resp.Reason = pbobjs.ApplyFriendResultReason_ApplyDecline
	} else if friendSettings.FriendVerifyType == pbobjs.FriendVerifyType_NeedFriendVerify {
		storage := storages.NewFriendApplicationStorage()
		storage.Upsert(models.FriendApplication{
			RecipientId: req.FriendId,
			SponsorId:   userId,
			ApplyTime:   time.Now().UnixMilli(),
			Status:      models.FriendApplicationStatus(models.FriendApplicationStatus_Apply),
			AppKey:      appkey,
		})
		resp.Reason = pbobjs.ApplyFriendResultReason_ApplySendOut
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
		resp.Reason = pbobjs.ApplyFriendResultReason_ApplySucc
	}
	return errs.IMErrorCode_SUCCESS, resp
}

func checkFriend(ctx context.Context, userId, friendId string) bool {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, respObj, err := AppSyncRpcCall(ctx, "check_friends", requestId, userId, &pbobjs.CheckFriendsReq{
		FriendIds: []string{friendId},
	}, func() proto.Message {
		return &pbobjs.CheckFriendsResp{}
	})
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return false
	}
	resp := respObj.(*pbobjs.CheckFriendsResp)
	if isFriend, exist := resp.CheckResults[friendId]; exist {
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
				RecipientId: application.RecipientId,
				SponsorId:   application.SponsorId,
				Status:      int32(application.Status),
				ApplyTime:   application.ApplyTime,
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
				RecipientId: application.RecipientId,
				SponsorId:   application.SponsorId,
				Status:      int32(application.Status),
				ApplyTime:   application.ApplyTime,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}
