package apis

import (
	"im-server/commons/errs"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/models"
	"im-server/services/appbusiness/services"
)

func AssistantAnswer(ctx *httputils.HttpContext) {
	req := models.AssistantAnswerReq{}
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
