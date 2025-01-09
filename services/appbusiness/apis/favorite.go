package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/services"
	"strconv"
)

func AddFavoriteMsg(ctx *httputils.HttpContext) {
	req := &pbobjs.FavoriteMsg{}
	if err := ctx.BindJson(req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.AddFavoriteMsg(ctx.ToRpcCtx(), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func QryFavoriteMsgs(ctx *httputils.HttpContext) {
	offset := ctx.Query("offset")
	limit := 20
	var err error
	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			limit = 20
		}
	}
	code, msgs := services.QryFavoriteMsgs(ctx.ToRpcCtx(), int64(limit), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(msgs)
}
