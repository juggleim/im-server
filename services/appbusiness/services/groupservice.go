package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/models"

	"google.golang.org/protobuf/proto"
)

func QryGroupInfo(ctx context.Context, groupId string) (errs.IMErrorCode, *pbobjs.GroupInfo) {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, respObj, err := AppSyncRpcCall(ctx, "qry_group_info", requestId, groupId, &pbobjs.GroupInfoReq{
		GroupId: groupId,
	}, func() proto.Message {
		return &pbobjs.GroupInfo{}
	})
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	return errs.IMErrorCode_SUCCESS, respObj.(*pbobjs.GroupInfo)
}

func CreateGroup(ctx context.Context, req *pbobjs.GroupMembersReq) (errs.IMErrorCode, *pbobjs.GroupInfo) {
	grpId := tools.GenerateUUIDShort11()
	requestId := bases.GetRequesterIdFromCtx(ctx)
	memberIds := []string{}
	containCurrUserId := false
	for _, memberId := range req.MemberIds {
		memberIds = append(memberIds, memberId)
		if memberId == requestId {
			containCurrUserId = true
		}
	}
	if !containCurrUserId {
		memberIds = append(memberIds, requestId)
	}
	code, _, err := AppSyncRpcCall(ctx, "g_add_members", requestId, grpId, &pbobjs.GroupMembersReq{
		GroupId:       grpId,
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
		MemberIds:     memberIds,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	//send notify msg
	targetUsers := []*models.User{}
	for _, memberId := range req.MemberIds {
		targetUsers = append(targetUsers, &models.User{
			UserId: memberId,
		})
	}
	notify := &models.GroupNotify{
		Operator: &models.User{
			UserId: requestId,
		},
		Members: targetUsers,
		Type:    models.GroupNotifyType_AddMember,
	}
	SendGrpNotify(ctx, grpId, notify)
	return errs.IMErrorCode_SUCCESS, &pbobjs.GroupInfo{
		GroupId:       grpId,
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
	}
}

func UpdateGroup(ctx context.Context, req *pbobjs.GroupInfo) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, _, err := AppSyncRpcCall(ctx, "upd_group_info", requestId, req.GroupId, req, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	SendGrpNotify(ctx, req.GroupId, &models.GroupNotify{
		Operator: &models.User{
			UserId: requestId,
		},
		Name: req.GroupName,
		Type: models.GroupNotifyType_Rename,
	})
	return errs.IMErrorCode_SUCCESS
}

func AddGrpMembers(ctx context.Context, grpMembers *pbobjs.GroupMembersReq) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, _, err := AppSyncRpcCall(ctx, "g_add_members", requestId, grpMembers.GroupId, &pbobjs.GroupMembersReq{
		GroupId:   grpMembers.GroupId,
		MemberIds: grpMembers.MemberIds,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	//send notify msg
	targetUsers := []*models.User{}
	for _, memberId := range grpMembers.MemberIds {
		targetUsers = append(targetUsers, &models.User{
			UserId: memberId,
		})
	}
	notify := &models.GroupNotify{
		Operator: &models.User{
			UserId: requestId,
		},
		Members: targetUsers,
		Type:    models.GroupNotifyType_AddMember,
	}
	//send notify msg
	SendGrpNotify(ctx, grpMembers.GroupId, notify)
	return errs.IMErrorCode_SUCCESS
}

func DelGrpMembers(ctx context.Context, req *pbobjs.GroupMembersReq) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, _, err := AppSyncRpcCall(ctx, "g_del_members", requestId, req.GroupId, &pbobjs.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: req.MemberIds,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	//send notify msg
	targetUsers := []*models.User{}
	for _, memberId := range req.MemberIds {
		targetUsers = append(targetUsers, &models.User{
			UserId: memberId,
		})
	}
	SendGrpNotify(ctx, req.GroupId, &models.GroupNotify{
		Operator: &models.User{
			UserId: requestId,
		},
		Members: targetUsers,
		Type:    models.GroupNotifyType_RemoveMember,
	})
	return errs.IMErrorCode_SUCCESS
}
