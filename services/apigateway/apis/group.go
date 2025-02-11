package apis

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func SetGroupSettings(ctx *gin.Context) {
	var req models.SetGroupSettingReq
	if err := ctx.BindJSON(&req); err != nil || req.GroupId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	kvMap := make(map[string]string)
	for k, v := range req.Settings {
		if commonservices.CheckGroupSettingKey(k) {
			kvMap[k] = fmt.Sprintf("%v", v)
		}
	}
	if len(kvMap) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, ""), "upd_group_info", req.GroupId, &pbobjs.GroupInfo{
		GroupId:  req.GroupId,
		Settings: commonservices.Map2KvItems(kvMap),
	})
	tools.SuccessHttpResp(ctx, nil)
}

func GetGroupSettings(ctx *gin.Context) {
	groupId := ctx.Query("group_id")

	groupReq := &pbobjs.GroupInfoReq{
		GroupId:    groupId,
		CareFields: []string{},
	}
	code, groupInfo, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "get_grp_setting", groupId, groupReq, func() proto.Message {
		return &pbobjs.GroupInfo{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	tools.SuccessHttpResp(ctx, groupInfo)
}

func GroupAddMembers(ctx *gin.Context) {
	var addMemberReq models.GroupMembersReq
	if err := ctx.BindJSON(&addMemberReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "g_add_members", addMemberReq.GroupId, &pbobjs.GroupMembersReq{
		GroupId:       addMemberReq.GroupId,
		GroupName:     addMemberReq.GroupName,
		GroupPortrait: addMemberReq.GroupPortrait,
		MemberIds:     addMemberReq.MemberIds,
		ExtFields:     commonservices.Map2KvItems(addMemberReq.ExtFields),
	}, func() proto.Message {
		return &pbobjs.GroupMembersResp{}
	})

	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	tools.SuccessHttpResp(ctx, nil)
}

func GroupDelMembers(ctx *gin.Context) {
	var delMemberReq models.GroupMembersReq
	if err := ctx.BindJSON(&delMemberReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "g_del_members", delMemberReq.GroupId, &pbobjs.GroupMembersReq{
		GroupId:   delMemberReq.GroupId,
		MemberIds: delMemberReq.MemberIds,
	}, func() proto.Message {
		return &pbobjs.GroupMembersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func GroupMemberUpdate(ctx *gin.Context) {
	var req models.GroupMemberUpdateReq
	if err := ctx.BindJSON(&req); err != nil || req.GroupId == "" || req.MemberId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	extFields := map[string]string{}
	if req.GrpDisplayName != "" {
		extFields[string(commonservices.AttItemKey_GrpDisplayName)] = req.GrpDisplayName
	}
	if len(req.ExtFields) > 0 {
		for k, v := range req.ExtFields {
			extFields[k] = v
		}
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "set_grp_member_setting", req.GroupId, &pbobjs.GroupMember{
		MemberId:  req.MemberId,
		ExtFields: commonservices.Map2KvItems(extFields),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func GroupDissolve(ctx *gin.Context) {
	var disReq models.GroupInfo
	if err := ctx.BindJSON(&disReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "g_dissolve", disReq.GroupId, &pbobjs.GroupMembersReq{
		GroupId: disReq.GroupId,
	}, func() proto.Message {
		return &pbobjs.GroupMembersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func GroupMute(ctx *gin.Context) {
	var muteReq models.GroupMuteReq
	if err := ctx.BindJSON(&muteReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "group_mute", muteReq.GrouopId, &pbobjs.GroupMuteReq{
		GroupId: muteReq.GrouopId,
		IsMute:  int32(muteReq.IsMute),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryGroupInfo(ctx *gin.Context) {
	groupId := ctx.Query("group_id")
	if groupId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	groupReq := &pbobjs.GroupInfoReq{
		GroupId:    groupId,
		CareFields: []string{},
	}
	code, groupInfo, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "qry_group_info", groupId, groupReq, func() proto.Message {
		return &pbobjs.GroupInfo{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	info := groupInfo.(*pbobjs.GroupInfo)

	tools.SuccessHttpResp(ctx, &models.GroupInfo{
		GroupId:       info.GroupId,
		GroupName:     info.GroupName,
		GroupPortrait: info.GroupPortrait,
		IsMute:        int(info.IsMute),
		UpdatedTime:   info.UpdatedTime,
		ExtFields:     commonservices.Kvitems2Map(info.ExtFields),
	})
}

func UpdateGroup(ctx *gin.Context) {
	var req models.GroupInfo
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "upd_group_info", req.GroupId, &pbobjs.GroupInfo{
		GroupId:       req.GroupId,
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
		ExtFields:     commonservices.Map2KvItems(req.ExtFields),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func GroupMemberMute(ctx *gin.Context) {
	var muteReq models.GroupMemberMuteReq
	if err := ctx.BindJSON(&muteReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	var muteEndAt int64 = 0
	if muteReq.IsMute > 0 && muteReq.MuteMinute > 0 {
		muteEndAt = time.Now().UnixMilli() + int64(muteReq.MuteMinute*60*1000)
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "group_member_mute", muteReq.GroupId, &pbobjs.GroupMemberMuteReq{
		GroupId:   muteReq.GroupId,
		MemberIds: muteReq.MemberIds,
		IsMute:    int32(muteReq.IsMute),
		MuteEndAt: muteEndAt,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func GroupMemberAllow(ctx *gin.Context) {
	var allowReq models.GroupMemberAllowReq
	if err := ctx.BindJSON(&allowReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "group_member_allow", allowReq.GroupId, &pbobjs.GroupMemberAllowReq{
		GroupId:   allowReq.GroupId,
		MemberIds: allowReq.MemberIds,
		IsAllow:   int32(allowReq.IsAllow),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func GroupMembers(ctx *gin.Context) {
	groupId := ctx.Query("group_id")
	offsetStr := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "g_qry_members", groupId, &pbobjs.QryGroupMembersReq{
		GroupId: groupId,
		Limit:   limit,
		Offset:  offsetStr,
	}, func() proto.Message {
		return &pbobjs.GroupMembersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	groupMembers := resp.(*pbobjs.GroupMembersResp)
	ret := &models.GroupMembersResp{
		Items:  []*models.GroupMember{},
		Offset: groupMembers.Offset,
	}
	for _, member := range groupMembers.Items {
		displayName := ""
		extMap := map[string]string{}
		if len(member.ExtFields) > 0 {
			for _, ext := range member.ExtFields {
				if ext.Key == string(commonservices.AttItemKey_GrpDisplayName) {
					displayName = ext.Value
				} else {
					extMap[ext.Key] = ext.Value
				}
			}
		}
		ret.Items = append(ret.Items, &models.GroupMember{
			MemberId:       member.MemberId,
			IsMute:         int(member.IsMute),
			IsAllow:        int(member.IsAllow),
			GrpDisplayName: displayName,
			ExtFields:      extMap,
		})
	}
	tools.SuccessHttpResp(ctx, ret)
}

func GroupMembersByIds(ctx *gin.Context) {
	var req models.GroupMembersReq
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "qry_group_members_by_ids", req.GroupId, &pbobjs.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: req.MemberIds,
	}, func() proto.Message {
		return &pbobjs.GroupMembersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	groupMembers := resp.(*pbobjs.GroupMembersResp)
	ret := &models.GroupMembersResp{
		Items:  []*models.GroupMember{},
		Offset: groupMembers.Offset,
	}
	for _, member := range groupMembers.Items {
		displayName := ""
		extMap := map[string]string{}
		if len(member.ExtFields) > 0 {
			for _, ext := range member.ExtFields {
				if ext.Key == string(commonservices.AttItemKey_GrpDisplayName) {
					displayName = ext.Value
				} else {
					extMap[ext.Key] = ext.Value
				}
			}
		}
		ret.Items = append(ret.Items, &models.GroupMember{
			MemberId:       member.MemberId,
			IsMute:         int(member.IsMute),
			IsAllow:        int(member.IsAllow),
			GrpDisplayName: displayName,
			ExtFields:      extMap,
		})
	}
	tools.SuccessHttpResp(ctx, ret)
}
