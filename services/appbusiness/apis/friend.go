package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
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
		Items:  []*pbobjs.UserObj{},
		Offset: friends.Offset,
	}
	for _, friend := range friends.Items {
		ret.Items = append(ret.Items, &pbobjs.UserObj{
			UserId:   friend.UserId,
			Nickname: friend.Nickname,
			Avatar:   friend.UserPortrait,
			IsFriend: true,
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
	code := services.AddFriends(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.FriendIdsReq{
		FriendIds: []string{req.FriendId},
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func DelFriend(ctx *httputils.HttpContext) {
	req := models.FriendIds{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelFriends(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.FriendIdsReq{
		FriendIds: req.FriendIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func ApplyFriend(ctx *httputils.HttpContext) {
	req := models.ApplyFriend{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.ApplyFriend(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.ApplyFriend{
		FriendId: req.FriendId,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func MyFriendApplications(ctx *httputils.HttpContext) {
	startTimeStr := ctx.Query("start")
	start, err := tools.String2Int64(startTimeStr)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := tools.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := tools.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryMyFriendApplications(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.QryFriendApplicationsReq{
		StartTime: start,
		Count:     int32(count),
		Order:     int32(order),
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func MyPendingFriendApplications(ctx *httputils.HttpContext) {
	startTimeStr := ctx.Query("start")
	start, err := tools.String2Int64(startTimeStr)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := tools.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := tools.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryMyPendingFriendApplications(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.QryFriendApplicationsReq{
		StartTime: start,
		Count:     int32(count),
		Order:     int32(order),
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}
