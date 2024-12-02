package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/models"
	"im-server/services/appbusiness/services"
)

func QryUserInfo(ctx *httputils.HttpContext) {
	userId := ctx.Query("user_id")
	rpcCtx := ctx.ToRpcCtx(ctx.CurrentUserId)
	code, userInfo := services.QryUserInfo(rpcCtx, userId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(&models.User{
		UserId:   userInfo.UserId,
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.UserPortrait,
	})
}

func UpdateUser(ctx *httputils.HttpContext) {
	req := &models.User{}
	if err := ctx.BindJson(req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	rpcCtx := ctx.ToRpcCtx(req.UserId)
	services.UpdateUser(rpcCtx, &pbobjs.UserInfo{
		UserId:       req.UserId,
		Nickname:     req.Nickname,
		UserPortrait: req.Avatar,
	})
	ctx.ResponseSucc(nil)
}

func SearchByPhone(ctx *httputils.HttpContext) {
	req := &models.User{}
	if err := ctx.BindJson(req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	rpcCtx := ctx.ToRpcCtx(ctx.CurrentUserId)
	code, users := services.SearchByPhone(rpcCtx, req.Phone)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ret := &models.Users{
		Items: []*models.User{},
	}
	for _, user := range users.UserInfos {
		ret.Items = append(ret.Items, &models.User{
			UserId:   user.UserId,
			Nickname: user.Nickname,
			Avatar:   user.UserPortrait,
			IsFriend: false,
		})
	}
	ctx.ResponseSucc(ret)
}
