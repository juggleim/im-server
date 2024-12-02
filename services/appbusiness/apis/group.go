package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/models"
	"im-server/services/appbusiness/services"
	"strconv"
)

func CreateGroup(ctx *httputils.HttpContext) {
	req := models.Group{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	memberIds := []string{}
	for _, user := range req.GrpMembers {
		memberIds = append(memberIds, user.UserId)
	}
	code, grpInfo := services.CreateGroup(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.GroupMembersReq{
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
		MemberIds:     memberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(&models.Group{
		GroupId:       grpInfo.GroupId,
		GroupName:     grpInfo.GroupName,
		GroupPortrait: grpInfo.GroupPortrait,
	})
}

func UpdateGroup(ctx *httputils.HttpContext) {
	req := models.Group{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.UpdateGroup(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.GroupInfo{
		GroupId:       req.GroupId,
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func AddGrpMembers(ctx *httputils.HttpContext) {
	req := models.Group{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	memberIds := []string{}
	for _, user := range req.GrpMembers {
		memberIds = append(memberIds, user.UserId)
	}
	code := services.AddGrpMembers(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: memberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func DelGrpMembers(ctx *httputils.HttpContext) {
	req := models.Group{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	memberIds := []string{}
	for _, user := range req.GrpMembers {
		memberIds = append(memberIds, user.UserId)
	}
	code := services.DelGrpMembers(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: memberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func QryGroupInfo(ctx *httputils.HttpContext) {
	groupId := ctx.Query("group_id")
	rpcCtx := ctx.ToRpcCtx(ctx.CurrentUserId)
	code, grpInfo := services.QryGroupInfo(rpcCtx, groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(&models.Group{
		GroupId:       groupId,
		GroupName:     grpInfo.GroupName,
		GroupPortrait: grpInfo.GroupPortrait,
	})
}

func QryMyGroups(ctx *httputils.HttpContext) {
	offset := ctx.Query("offset")
	count := 20
	countStr := ctx.Query("count")
	var err error
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			count = 20
		}
	}
	code, grps := services.QueryMyGroups(ctx.ToRpcCtx(ctx.CurrentUserId), &pbobjs.GroupInfoListReq{
		Limit:  int64(count),
		Offset: offset,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ret := &models.Groups{
		Offset: grps.Offset,
		Items:  []*models.Group{},
	}
	for _, grp := range grps.Items {
		ret.Items = append(ret.Items, &models.Group{
			GroupId:       grp.GroupId,
			GroupName:     grp.GroupName,
			GroupPortrait: grp.GroupPortrait,
		})
	}
	ctx.ResponseSucc(ret)
}
