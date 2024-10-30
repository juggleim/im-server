package apis

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func AddPrivateGlobalMuteMembers(ctx *gin.Context) {
	var req models.UserIds
	if err := ctx.BindJSON(&req); err != nil || len(req.UserIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	groups := bases.GroupTargets("pri_global_mute", req.UserIds)
	for _, ids := range groups {
		uIds := ids
		go func() {
			services.AsyncApiCall(ctx, "pri_global_mute", "", uIds[0], &pbobjs.BatchMuteUsersReq{
				UserIds:  uIds,
				IsDelete: false,
			})
		}()
	}
	tools.SuccessHttpResp(ctx, nil)
}

func DelPrivateGlobalMuteMembers(ctx *gin.Context) {
	var req models.UserIds
	if err := ctx.BindJSON(&req); err != nil || len(req.UserIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	groups := bases.GroupTargets("pri_global_mute", req.UserIds)
	for _, ids := range groups {
		uIds := ids
		go func() {
			services.AsyncApiCall(ctx, "pri_global_mute", "", uIds[0], &pbobjs.BatchMuteUsersReq{
				UserIds:  uIds,
				IsDelete: true,
			})
		}()
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryPrivateGlobalMuteMembers(ctx *gin.Context) {
	offset := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 {
			limit = intVal
		}
	}
	if limit > 1000 {
		limit = 1000
	}

	code, resp, err := services.SyncApiCall(ctx, "qry_pri_global_mute", "", tools.RandStr(10), &pbobjs.QryBlockUsersReq{
		Limit:  limit,
		Offset: offset,
	}, func() proto.Message {
		return &pbobjs.QryBlockUsersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	ret := &models.QryMuteUsersResp{
		Items: []*models.MuteUser{},
	}
	if muteUsers, ok := resp.(*pbobjs.QryBlockUsersResp); ok {
		ret.Offset = muteUsers.Offset
		for _, muteUser := range muteUsers.Items {
			ret.Items = append(ret.Items, &models.MuteUser{
				UserId:      muteUser.BlockUserId,
				CreatedTime: muteUser.CreatedTime,
			})
		}
	}
	tools.SuccessHttpResp(ctx, ret)
}
