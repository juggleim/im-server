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
	code, user := services.QryUserInfo(rpcCtx, userId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(user)
}

func UpdateUser(ctx *httputils.HttpContext) {
	req := &pbobjs.UserObj{}
	if err := ctx.BindJson(req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	rpcCtx := ctx.ToRpcCtx(req.UserId)
	services.UpdateUser(rpcCtx, req)
	ctx.ResponseSucc(nil)
}

func UpdateUserSettings(ctx *httputils.HttpContext) {
	req := &pbobjs.UserSettings{}
	if err := ctx.BindJson(req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.UpdateUserSettings(ctx.ToRpcCtx(ctx.CurrentUserId), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func SearchByPhone(ctx *httputils.HttpContext) {
	req := &pbobjs.UserObj{}
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
		Items: []*pbobjs.UserObj{},
	}
	for _, user := range users.UserInfos {
		ret.Items = append(ret.Items, &pbobjs.UserObj{
			UserId:   user.UserId,
			Nickname: user.Nickname,
			Avatar:   user.UserPortrait,
			IsFriend: false,
		})
	}
	ctx.ResponseSucc(ret)
}
