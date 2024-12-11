package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	apiModels "im-server/services/appbusiness/models"
	"im-server/services/commonservices"
	"im-server/services/friends/storages"
)

func QryFriends(ctx context.Context, req *pbobjs.FriendListReq) (errs.IMErrorCode, *pbobjs.FriendListResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	var startId int64 = 0
	if req.Offset != "" {
		startId, _ = tools.DecodeInt(req.Offset)
	}
	storage := storages.NewFriendRelStorage()
	rels, err := storage.QueryFriendRels(appkey, userId, startId, req.Limit)
	if err != nil {
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	ret := &pbobjs.FriendListResp{
		Items: []*pbobjs.UserInfo{},
	}
	for _, rel := range rels {
		ret.Offset, _ = tools.EncodeInt(rel.ID)
		ret.Items = append(ret.Items, commonservices.GetTargetDisplayUserInfo(ctx, rel.FriendId))
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func AddFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	userId := bases.GetRequesterIdFromCtx(ctx)
	for _, friendId := range req.FriendIds {
		AppSyncRpcCall(ctx, "add_friend", userId, userId, &pbobjs.FriendIdsReq{
			FriendIds: []string{friendId},
		}, nil)
		AppSyncRpcCall(ctx, "add_friend", userId, friendId, &pbobjs.FriendIdsReq{
			FriendIds: []string{userId},
		}, nil)
		//send notify msg
		SendFriendNotify(ctx, friendId, &apiModels.FriendNotify{
			Type: 0,
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func ApplyFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	return errs.IMErrorCode_SUCCESS
}
