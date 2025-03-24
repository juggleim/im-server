package apis

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"

	"github.com/gin-gonic/gin"
)

func AddFriends(ctx *gin.Context) {
	var req models.FriendIds
	if err := ctx.BindJSON(&req); err != nil || req.UserId == "" || len(req.FriendIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	members := []*pbobjs.FriendMember{}
	for _, userId := range req.FriendIds {
		members = append(members, &pbobjs.FriendMember{
			FriendId: userId,
		})
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, req.UserId), "add_friends", req.UserId, &pbobjs.FriendMembersReq{
		FriendMembers: members,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func DelFriends(ctx *gin.Context) {
	var req models.FriendIds
	if err := ctx.BindJSON(&req); err != nil || req.UserId == "" || len(req.FriendIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, req.UserId), "del_friends", req.UserId, &pbobjs.FriendIdsReq{
		FriendIds: req.FriendIds,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}
