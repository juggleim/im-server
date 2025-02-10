package apis

import (
	"bytes"
	"encoding/base64"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/models"
	"im-server/services/appbusiness/services"
	"image/png"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func CreateGroup(ctx *httputils.HttpContext) {
	req := models.Group{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
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
	code, grpInfo := services.CreateGroup(ctx.ToRpcCtx(), &pbobjs.GroupMembersReq{
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
	code := services.UpdateGroup(ctx.ToRpcCtx(), &pbobjs.GroupInfo{
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

func DissolveGroup(ctx *httputils.HttpContext) {
	req := models.Group{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DissolveGroup(ctx.ToRpcCtx(), req.GroupId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func QuitGroup(ctx *httputils.HttpContext) {
	req := models.Group{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.QuitGroup(ctx.ToRpcCtx(), req.GroupId)
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
	code := services.AddGrpMembers(ctx.ToRpcCtx(), &pbobjs.GroupMembersReq{
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
	if err := ctx.BindJson(&req); err != nil || req.GroupId == "" || len(req.MemberIds) <= 0 {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelGrpMembers(ctx.ToRpcCtx(), &pbobjs.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: req.MemberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func QryGroupInfo(ctx *httputils.HttpContext) {
	groupId := ctx.Query("group_id")
	rpcCtx := ctx.ToRpcCtx()
	code, grpInfo := services.QryGroupInfo(rpcCtx, groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(grpInfo)
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
	code, grps := services.QueryMyGroups(ctx.ToRpcCtx(), &pbobjs.GroupInfoListReq{
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

func QryGrpMembers(ctx *httputils.HttpContext) {
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
	code, members := services.QueryGrpMembers(ctx.ToRpcCtx(), &pbobjs.QryGroupMembersReq{
		GroupId: groupId,
		Limit:   int64(limit),
		Offset:  offset,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ret := &models.GroupMembersResp{
		Items: []*models.GroupMember{},
	}
	for _, member := range members.Items {
		ret.Items = append(ret.Items, &models.GroupMember{
			UserObj: pbobjs.UserObj{
				UserId:   member.UserId,
				Nickname: member.Nickname,
				Avatar:   member.Avatar,
			},
		})
	}
	ret.Offset = members.Offset
	ctx.ResponseSucc(ret)
}

func CheckGroupMembers(ctx *httputils.HttpContext) {
	req := models.CheckGroupMembersReq{}
	if err := ctx.BindJson(&req); err != nil {
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

func SetGrpAnnouncement(ctx *httputils.HttpContext) {
	req := models.GroupAnnouncement{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGrpAnnouncement(ctx.ToRpcCtx(), &pbobjs.GrpAnnouncement{
		GroupId: req.GroupId,
		Content: req.Content,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(nil)
}

func GetGrpAnnouncement(ctx *httputils.HttpContext) {
	groupId := ctx.Query("group_id")
	code, grpAnnounce := services.GetGrpAnnouncement(ctx.ToRpcCtx(), groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(&models.GroupAnnouncement{
		GroupId: grpAnnounce.GroupId,
		Content: grpAnnounce.Content,
	})
}

func SetGrpDisplayName(ctx *httputils.HttpContext) {
	req := pbobjs.SetGroupDisplayNameReq{}
	if err := ctx.BindJson(&req); err != nil {
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

func QryGrpQrCode(ctx *httputils.HttpContext) {
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
	qrCode, _ := qr.Encode(tools.ToJson(m), qr.M, qr.Auto)
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
