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

func AddFriends(ctx *gin.Context) {
	var req models.FriendIds
	if err := ctx.BindJSON(&req); err != nil || req.UserId == "" || (len(req.Friends) <= 0 && len(req.FriendIds) <= 0) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	members := []*pbobjs.FriendMember{}
	if len(req.Friends) > 0 {
		for _, friend := range req.Friends {
			members = append(members, &pbobjs.FriendMember{
				FriendId:    friend.FriendId,
				DisplayName: friend.DisplayName,
			})
		}
	} else if len(req.FriendIds) > 0 {
		for _, friendId := range req.FriendIds {
			members = append(members, &pbobjs.FriendMember{
				FriendId: friendId,
			})
		}
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

func QryFriends(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	if userId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	offsetStr := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	orderStr := ctx.Query("order")
	var order int32 = 0
	if orderStr != "" {
		intVal, err := tools.String2Int64(orderStr)
		if err == nil && intVal > 0 {
			order = int32(intVal)
		}
	}
	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, userId), "qry_friends", userId, &pbobjs.QryFriendsReq{
		Limit:  limit,
		Offset: offsetStr,
		Order:  order,
	}, func() proto.Message {
		return &pbobjs.QryFriendsResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	friends := resp.(*pbobjs.QryFriendsResp)
	ret := &models.FriendsResp{
		Items: []*models.FriendItem{},
	}
	if friends != nil && len(friends.Items) > 0 {
		ret.Offset = friends.Offset
		for _, friend := range friends.Items {
			ret.Items = append(ret.Items, &models.FriendItem{
				FriendId:    friend.FriendId,
				DisplayName: friend.DisplayName,
			})
		}
	}
	tools.SuccessHttpResp(ctx, ret)
}

func SetFriendDisplayName(ctx *gin.Context) {
	var req models.FriendItem
	if err := ctx.BindJSON(&req); err != nil || req.UserId == "" || req.FriendId == "" || req.DisplayName == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	members := []*pbobjs.FriendMember{
		{
			FriendId:    req.FriendId,
			DisplayName: req.DisplayName,
		},
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
