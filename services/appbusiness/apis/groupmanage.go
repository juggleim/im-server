package apis

import (
	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
)

func ChgGroupOwner(ctx *HttpContext) {
	req := &apimodels.GroupOwnerChgReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.ChgGroupOwner(ctx.ToRpcCtx(), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func AddGrpAdministrator(ctx *HttpContext) {
	req := &apimodels.GroupAdministratorsReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.AddGroupAdministrators(ctx.ToRpcCtx(), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func DelGrpAdministrator(ctx *HttpContext) {
	req := &apimodels.GroupAdministratorsReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelGroupAdministrators(ctx.ToRpcCtx(), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func QryGrpAdministrators(ctx *HttpContext) {
	groupId := ctx.Query("group_id")
	code, resp := services.QryGroupAdministrators(ctx.ToRpcCtx(), groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func SetGroupMute(ctx *HttpContext) {
	req := &apimodels.SetGroupMuteReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGroupMute(ctx.ToRpcCtx(), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func SetGrpVerifyType(ctx *HttpContext) {
	req := &apimodels.SetGroupVerifyTypeReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGroupVerifyType(ctx.ToRpcCtx(), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func SetGrpHisMsgVisible(ctx *HttpContext) {
	req := &apimodels.SetGroupHisMsgVisibleReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGroupHisMsgVisible(ctx.ToRpcCtx(), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}
