package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/models"
	"im-server/services/commonservices"
	"im-server/services/group/dbs"

	"google.golang.org/protobuf/proto"
)

func QryUserInfo(ctx context.Context, userId string) (errs.IMErrorCode, *pbobjs.UserInfo) {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, respObj, err := AppSyncRpcCall(ctx, "qry_user_info", requestId, userId, &pbobjs.UserIdReq{
		UserId: userId,
	}, func() proto.Message {
		return &pbobjs.UserInfo{}
	})
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	return errs.IMErrorCode_SUCCESS, respObj.(*pbobjs.UserInfo)
}

func SearchByPhone(ctx context.Context, phone string) (errs.IMErrorCode, *pbobjs.UserInfos) {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	targetUserId := tools.ShortMd5(phone)
	code, respObj, err := AppSyncRpcCall(ctx, "qry_user_info", requestId, targetUserId, &pbobjs.UserIdReq{
		UserId:   targetUserId,
		AttTypes: []int32{int32(commonservices.AttItemType_Att)},
	}, func() proto.Message {
		return &pbobjs.UserInfo{}
	})
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	users := &pbobjs.UserInfos{
		UserInfos: []*pbobjs.UserInfo{},
	}
	users.UserInfos = append(users.UserInfos, respObj.(*pbobjs.UserInfo))
	return errs.IMErrorCode_SUCCESS, users
}

func UpdateUser(ctx context.Context, req *pbobjs.UserInfo) errs.IMErrorCode {
	requesterId := bases.GetRequesterIdFromCtx(ctx)
	AppSyncRpcCall(ctx, "upd_user_info", requesterId, req.UserId, req, nil)
	return errs.IMErrorCode_SUCCESS
}

func QueryMyGroups(ctx context.Context, req *pbobjs.GroupInfoListReq) (errs.IMErrorCode, *pbobjs.GroupInfoListResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	memberId := bases.GetRequesterIdFromCtx(ctx)
	dao := dbs.GroupMemberDao{}
	var startId int64
	if req.Offset != "" {
		startId, _ = tools.DecodeInt(req.Offset)
	}
	groups, err := dao.QueryGroupsByMemberId(appkey, memberId, startId, int64(req.Limit))
	if err != nil {
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	ret := &pbobjs.GroupInfoListResp{
		Items: []*pbobjs.GroupInfo{},
	}
	for _, group := range groups {
		ret.Offset, _ = tools.EncodeInt(group.ID)
		ret.Items = append(ret.Items, &pbobjs.GroupInfo{
			GroupId: group.GroupId,
		})
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func GetUser(ctx context.Context, userId string) *models.User {
	u := &models.User{
		UserId: userId,
	}
	user := commonservices.GetTargetDisplayUserInfo(ctx, userId)
	if user != nil {
		u.Nickname = user.Nickname
		u.Avatar = user.UserPortrait
	}
	return u
}
