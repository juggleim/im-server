package apis

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strconv"

	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/utils"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func CreateGroup(ctx *HttpContext) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	memberIds := req.MemberIds
	if len(memberIds) <= 0 && len(req.GrpMembers) > 0 {
		ids := []string{}
		for _, member := range req.GrpMembers {
			ids = append(ids, member.UserId)
		}
		memberIds = ids
	}
	code, grpInfo := services.CreateGroup(ctx.ToRpcCtx(), &apimodels.GroupMembersReq{
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
		MemberIds:     memberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(&apimodels.Group{
		GroupId:       grpInfo.GroupId,
		GroupName:     grpInfo.GroupName,
		GroupPortrait: grpInfo.GroupPortrait,
	})
}

func UpdateGroup(ctx *HttpContext) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.UpdateGroup(ctx.ToRpcCtx(), &apimodels.GroupInfo{
		GroupId:       req.GroupId,
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func DissolveGroup(ctx *HttpContext) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DissolveGroup(ctx.ToRpcCtx(), req.GroupId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func QuitGroup(ctx *HttpContext) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.QuitGroup(ctx.ToRpcCtx(), req.GroupId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func AddGrpMembers(ctx *HttpContext) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	memberIds := []string{}
	for _, user := range req.GrpMembers {
		memberIds = append(memberIds, user.UserId)
	}
	code := services.AddGrpMembers(ctx.ToRpcCtx(), &apimodels.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: memberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func DelGrpMembers(ctx *HttpContext) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil || req.GroupId == "" || len(req.MemberIds) <= 0 {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelGrpMembers(ctx.ToRpcCtx(), &apimodels.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: req.MemberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func QryGroupInfo(ctx *HttpContext) {
	groupId := ctx.Query("group_id")
	rpcCtx := ctx.ToRpcCtx()
	code, grpInfo := services.QryGroupInfo(rpcCtx, groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(grpInfo)
}

func QryMyGroups(ctx *HttpContext) {
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
	code, grps := services.QueryMyGroups(ctx.ToRpcCtx(), int64(count), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ret := &apimodels.Groups{
		Offset: grps.Offset,
		Items:  []*apimodels.Group{},
	}
	for _, grp := range grps.Items {
		ret.Items = append(ret.Items, &apimodels.Group{
			GroupId:       grp.GroupId,
			GroupName:     grp.GroupName,
			GroupPortrait: grp.GroupPortrait,
		})
	}
	ctx.ResponseSucc(ret)
}

func QryGrpMembers(ctx *HttpContext) {
	groupId := ctx.Query("group_id")
	offset := ctx.Query("offset")
	limit := 100
	limitStr := ctx.Query("limit")
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			limit = 100
		}
	}
	code, members := services.QueryGrpMembers(ctx.ToRpcCtx(), groupId, int64(limit), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ret := &apimodels.GroupMembersResp{
		Items: []*apimodels.GroupMember{},
	}
	for _, member := range members.Items {
		ret.Items = append(ret.Items, &apimodels.GroupMember{
			UserObj: apimodels.UserObj{
				UserId:   member.UserId,
				Nickname: member.Nickname,
				Avatar:   member.Avatar,
			},
		})
	}
	ret.Offset = members.Offset
	ctx.ResponseSucc(ret)
}

func CheckGroupMembers(ctx *HttpContext) {
	req := apimodels.CheckGroupMembersReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.CheckGroupMembers(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(resp)
}

func SetGrpAnnouncement(ctx *HttpContext) {
	req := apimodels.GroupAnnouncement{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGrpAnnouncement(ctx.ToRpcCtx(), &apimodels.GroupAnnouncement{
		GroupId: req.GroupId,
		Content: req.Content,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func GetGrpAnnouncement(ctx *HttpContext) {
	groupId := ctx.Query("group_id")
	code, grpAnnounce := services.GetGrpAnnouncement(ctx.ToRpcCtx(), groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(&apimodels.GroupAnnouncement{
		GroupId: grpAnnounce.GroupId,
		Content: grpAnnounce.Content,
	})
}

func SetGrpDisplayName(ctx *HttpContext) {
	req := apimodels.SetGroupDisplayNameReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGrpDisplayName(ctx.ToRpcCtx(), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func QryGrpQrCode(ctx *HttpContext) {
	grpId := ctx.Query("group_id")
	if grpId == "" {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	userId := ctx.CurrentUserId

	m := map[string]interface{}{
		"action":   "join_group",
		"group_id": grpId,
		"user_id":  userId,
	}
	buf := bytes.NewBuffer([]byte{})
	qrCode, _ := qr.Encode(utils.ToJson(m), qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 400, 400)
	err := png.Encode(buf, qrCode)
	if err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_DEFAULT)
		return
	}
	ctx.ResponseSucc(map[string]string{
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})
}
