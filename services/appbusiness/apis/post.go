package apis

import (
	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/utils"
)

func PostAdd(ctx *HttpContext) {
	req := apimodels.Post{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.PostAdd(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func QryPosts(ctx *HttpContext) {
	var limit int64 = 20
	limitStr := ctx.Query("limit")
	var err error
	if limitStr != "" {
		limit, err = utils.String2Int64(limitStr)
		if err != nil {
			limit = 20
		}
	}
	var start int64
	startTimeStr := ctx.Query("start")
	start, err = utils.String2Int64(startTimeStr)
	if err != nil {
		start = 0
	}
	var isPositive bool = false
	orderStr := ctx.Query("order")
	order, err := utils.String2Int64(orderStr)
	if err == nil {
		if order == 1 {
			isPositive = true
		}
	}
	code, resp := services.QryPosts(ctx.ToRpcCtx(), start, limit, isPositive)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func PostInfo(ctx *HttpContext) {
	postId := ctx.Query("post_id")
	if postId == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.QryPostInfo(ctx.ToRpcCtx(), postId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func QryPostComments(ctx *HttpContext) {
	postId := ctx.Query("post_id")
	if postId == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	var limit int64 = 20
	limitStr := ctx.Query("limit")
	var err error
	if limitStr != "" {
		limit, err = utils.String2Int64(limitStr)
		if err != nil {
			limit = 20
		}
	}
	var start int64
	startTimeStr := ctx.Query("start")
	start, err = utils.String2Int64(startTimeStr)
	if err != nil {
		start = 0
	}
	var isPositive bool = false
	orderStr := ctx.Query("order")
	order, err := utils.String2Int64(orderStr)
	if err == nil {
		if order == 1 {
			isPositive = true
		}
	}
	code, resp := services.QryPostComments(ctx.ToRpcCtx(), postId, start, limit, isPositive)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func PostCommentAdd(ctx *HttpContext) {
	req := apimodels.PostComment{}
	if err := ctx.BindJSON(&req); err != nil || req.PostId == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.PostCommentAdd(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}
