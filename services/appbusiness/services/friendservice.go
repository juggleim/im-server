package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/storages"
	"im-server/services/appbusiness/storages/models"
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
		ret.Items = append(ret.Items, &pbobjs.UserInfo{
			UserId: rel.FriendId,
		})
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func AddFriends(ctx context.Context, req *pbobjs.FriendsAddReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	for _, friendId := range req.FriendIds {
		storage.Upsert(models.FriendRel{
			AppKey:   appkey,
			UserId:   userId,
			FriendId: friendId,
		})
	}
	return errs.IMErrorCode_SUCCESS
}
