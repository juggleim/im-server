package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/services"
)

func QryMyGrpApplications(ctx *httputils.HttpContext) {
	startTimeStr := ctx.Query("start")
	start, err := tools.String2Int64(startTimeStr)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := tools.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := tools.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryMyGrpApplications(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.QryGrpApplicationsReq{
		StartTime: start,
		Count:     int32(count),
		Order:     int32(order),
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func QryMyPendingGrpInvitations(ctx *httputils.HttpContext) {
	startTimeStr := ctx.Query("start")
	start, err := tools.String2Int64(startTimeStr)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := tools.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := tools.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryMyPendingGrpInvitations(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.QryGrpApplicationsReq{
		StartTime: start,
		Count:     int32(count),
		Order:     int32(order),
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func QryGrpInvitations(ctx *httputils.HttpContext) {
	groupId := ctx.Query("group_id")
	startTimeStr := ctx.Query("start")
	start, err := tools.String2Int64(startTimeStr)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := tools.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := tools.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryGrpInvitations(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.QryGrpApplicationsReq{
		StartTime: start,
		Count:     int32(count),
		Order:     int32(order),
		GroupId:   groupId,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func QryGrpPendingApplications(ctx *httputils.HttpContext) {
	groupId := ctx.Query("group_id")
	startTimeStr := ctx.Query("start")
	start, err := tools.String2Int64(startTimeStr)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := tools.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := tools.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryGrpPendingApplications(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.QryGrpApplicationsReq{
		StartTime: start,
		Count:     int32(count),
		Order:     int32(order),
		GroupId:   groupId,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}
