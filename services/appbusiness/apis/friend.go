package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/models"
	"im-server/services/appbusiness/services"
	"strconv"
)

func QryFriends(ctx *httputils.HttpContext) {
	offset := ctx.Query("offset")
	count := 20
	var err error
	countStr := ctx.Query("count")
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			count = 20
		}
	}
	code, friends := services.QryFriends(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.FriendListReq{
		Limit:  int64(count),
		Offset: offset,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ret := &models.Friends{
		Items:  []*models.User{},
		Offset: friends.Offset,
	}
	for _, friend := range friends.Items {
		ret.Items = append(ret.Items, &models.User{
			UserId:   friend.UserId,
			Nickname: friend.Nickname,
			Avatar:   friend.UserPortrait,
		})
	}
	ctx.ResponseSucc(ret)
}

func AddFriend(ctx *httputils.HttpContext) {
	req := models.Friend{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.AddFriends(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.FriendsAddReq{
		FriendIds: []string{req.FriendId},
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	//send notify msg
	services.SendFriendNotify(ctx.ToRpcCtx(ctx.CurrentUserId), req.FriendId, &models.FriendNotify{
		Type: 0,
	})
	ctx.ResponseSucc(nil)
}
