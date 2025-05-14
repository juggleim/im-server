package apis

import (
	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"

	"strconv"
)

func AssistantAnswer(ctx *HttpContext) {
	req := apimodels.AssistantAnswerReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.AutoAnswer(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func PromptAdd(ctx *HttpContext) {
	req := apimodels.Prompt{}
	if err := ctx.BindJSON(&req); err != nil || req.Prompts == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.PromptAdd(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	ctx.ResponseSucc(resp)
}

func PromptUpdate(ctx *HttpContext) {
	req := apimodels.Prompt{}
	if err := ctx.BindJSON(&req); err != nil || req.Id == "" || req.Prompts == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptUpdate(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func PromptDel(ctx *HttpContext) {
	req := apimodels.Prompt{}
	if err := ctx.BindJSON(&req); err != nil || req.Id == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptDel(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func PromptBatchDel(ctx *HttpContext) {
	req := apimodels.PromptIds{}
	if err := ctx.BindJSON(&req); err != nil || len(req.Ids) <= 0 {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptBatchDel(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func QryPrompts(ctx *HttpContext) {
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
	code, prompts := services.QryPrompts(ctx.ToRpcCtx(), int64(count), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, prompts)
}
