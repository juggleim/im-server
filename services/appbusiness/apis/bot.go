package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/services"
	"strconv"
)

func QryBots(ctx *httputils.HttpContext) {
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
	code, bots := services.QryAiBots(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.QryAiBotsReq{
		Limit:  int64(count),
		Offset: offset,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(bots)
}
