package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/services"
)

func ChgGroupOwner(ctx *httputils.HttpContext) {
	req := &pbobjs.GroupOwnerChgReq{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.ChgGroupOwner(ctx.ToRpcCtx(ctx.CurrentUserId), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func AddGrpAdministrator(ctx *httputils.HttpContext) {
	req := &pbobjs.GroupAdministratorsReq{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.AddGroupAdministrators(ctx.ToRpcCtx(ctx.CurrentUserId), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func DelGrpAdministrator(ctx *httputils.HttpContext) {
	req := &pbobjs.GroupAdministratorsReq{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelGroupAdministrators(ctx.ToRpcCtx(ctx.CurrentUserId), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func QryGrpAdministrators(ctx *httputils.HttpContext) {
	groupId := ctx.Query("group_id")
	code, resp := services.QryGroupAdministrators(ctx.ToRpcCtx(ctx.CurrentUserId), groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func SetGroupMute(ctx *httputils.HttpContext) {
	req := &pbobjs.SetGroupMuteReq{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGroupMute(ctx.ToRpcCtx(ctx.CurrentUserId), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func SetGrpVerifyType(ctx *httputils.HttpContext) {
	req := &pbobjs.SetGroupVerifyTypeReq{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGroupVerifyType(ctx.ToRpcCtx(ctx.CurrentUserId), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func SetGrpHisMsgVisible(ctx *httputils.HttpContext) {
	req := &pbobjs.SetGroupHisMsgVisibleReq{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGroupHisMsgVisible(ctx.ToRpcCtx(ctx.CurrentUserId), req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}
