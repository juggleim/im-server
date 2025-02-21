package apis

import (
	"im-server/commons/errs"
	"im-server/services/appbusiness/apimodels"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/services"
	"strconv"
)

func AssistantAnswer(ctx *httputils.HttpContext) {
	req := apimodels.AssistantAnswerReq{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.AssistantAnswer(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func PromptAdd(ctx *httputils.HttpContext) {
	req := apimodels.Prompt{}
	if err := ctx.BindJson(&req); err != nil || req.Prompts == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptAdd(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func PromptUpdate(ctx *httputils.HttpContext) {
	req := apimodels.Prompt{}
	if err := ctx.BindJson(&req); err != nil || req.Id == "" || req.Prompts == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptUpdate(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func PromptDel(ctx *httputils.HttpContext) {
	req := apimodels.Prompt{}
	if err := ctx.BindJson(&req); err != nil || req.Id == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptDel(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func PromptBatchDel(ctx *httputils.HttpContext) {
	req := apimodels.PromptIds{}
	if err := ctx.BindJson(&req); err != nil || len(req.Ids) <= 0 {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptBatchDel(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func QryPrompts(ctx *httputils.HttpContext) {
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
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(prompts)
}
