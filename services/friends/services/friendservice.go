package services

import (
	"context"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
)

func AddFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	return errs.IMErrorCode_SUCCESS
}

func DelFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	return errs.IMErrorCode_SUCCESS
}

func QryFriends(ctx context.Context, req *pbobjs.QryFriendsReq) (errs.IMErrorCode, *pbobjs.QryFriendsResp) {
	return errs.IMErrorCode_SUCCESS, nil
}
