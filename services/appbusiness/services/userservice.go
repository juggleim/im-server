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

func QryUserInfo(ctx context.Context, userId string) (errs.IMErrorCode, *pbobjs.UserObj) {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	userInfo := commonservices.GetTargetUserInfo(ctx, userId)
	ret := &pbobjs.UserObj{
		UserId:   userInfo.UserId,
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.UserPortrait,
	}
	if userId == requestId {
		ret.Settings = GetUserSettings(userInfo)
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func GetUserSettings(userInfo *commonservices.TargetUserInfo) *pbobjs.UserSettings {
	settings := &pbobjs.UserSettings{}
	for _, setting := range userInfo.SettingsFields {
		if setting.Key == models.UserExtKey_Language {
			settings.Language = setting.Value
		} else if setting.Key == models.UserExtKey_Undisturb {
			settings.Undisturb = setting.Value
		} else if setting.Key == models.UserExtKey_FriendVerifyType {
			verifyType := tools.ToInt(setting.Value)
			settings.FriendVerifyType = pbobjs.FriendVerifyType(verifyType)
		} else if setting.Key == models.UserExtKey_GrpVerifyType {
			verifyType := tools.ToInt(setting.Value)
			settings.GrpVerifyType = pbobjs.GrpVerifyType(verifyType)
		}
	}
	return settings
}

func SearchByPhone(ctx context.Context, phone string) (errs.IMErrorCode, *pbobjs.UserObjs) {
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
	users := &pbobjs.UserObjs{
		Items: []*pbobjs.UserObj{},
	}
	userInfo := respObj.(*pbobjs.UserInfo)
	users.Items = append(users.Items, &pbobjs.UserObj{
		UserId:   userInfo.UserId,
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.UserPortrait,
		IsFriend: checkFriend(ctx, requestId, targetUserId),
	})
	return errs.IMErrorCode_SUCCESS, users
}

func UpdateUser(ctx context.Context, req *pbobjs.UserObj) errs.IMErrorCode {
	requesterId := bases.GetRequesterIdFromCtx(ctx)
	AppSyncRpcCall(ctx, "upd_user_info", requesterId, req.UserId, &pbobjs.UserInfo{
		UserId:       req.UserId,
		Nickname:     req.Nickname,
		UserPortrait: req.Avatar,
	}, nil)
	return errs.IMErrorCode_SUCCESS
}

func UpdateUserSettings(ctx context.Context, req *pbobjs.UserSettings) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	settings := []*pbobjs.KvItem{
		{
			Key:   models.UserExtKey_Language,
			Value: req.Language,
		},
		{
			Key:   models.UserExtKey_Undisturb,
			Value: req.Undisturb,
		},
		{
			Key:   models.UserExtKey_FriendVerifyType,
			Value: tools.Int642String(int64(req.FriendVerifyType)),
		},
		{
			Key:   models.UserExtKey_GrpVerifyType,
			Value: tools.Int642String(int64(req.GrpVerifyType)),
		},
	}
	if len(settings) > 0 {
		AppSyncRpcCall(ctx, "upd_user_info", requestId, requestId, &pbobjs.UserInfo{
			UserId:   requestId,
			Settings: settings,
		}, nil)
	}
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
		Items: []*pbobjs.GrpInfo{},
	}
	for _, group := range groups {
		ret.Offset, _ = tools.EncodeInt(group.ID)
		grpInfo := commonservices.GetGroupInfoFromCache(ctx, group.GroupId)
		ret.Items = append(ret.Items, &pbobjs.GrpInfo{
			GroupId:       grpInfo.GroupId,
			GroupName:     grpInfo.GroupName,
			GroupPortrait: grpInfo.GroupPortrait,
			MemberCount:   grpInfo.MemberCount,
		})
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func GetUser(ctx context.Context, userId string) *pbobjs.UserObj {
	u := &pbobjs.UserObj{
		UserId: userId,
	}
	user := commonservices.GetTargetDisplayUserInfo(ctx, userId)
	if user != nil {
		u.Nickname = user.Nickname
		u.Avatar = user.UserPortrait
	}
	return u
}
