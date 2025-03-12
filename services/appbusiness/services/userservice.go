package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/apimodels"
	"im-server/services/appbusiness/services/imsdk"
	"im-server/services/commonservices"
	"im-server/services/friends/storages"
	"im-server/services/group/dbs"
	userStorage "im-server/services/usermanager/storages"

	juggleimsdk "github.com/juggleim/imserver-sdk-go"
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
	} else {
		ret.IsFriend = checkFriend(ctx, requestId, userId)
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func GetUserSettings(userInfo *commonservices.TargetUserInfo) *pbobjs.UserSettings {
	settings := &pbobjs.UserSettings{}
	for _, setting := range userInfo.SettingsFields {
		if setting.Key == apimodels.UserExtKey_Language {
			settings.Language = setting.Value
		} else if setting.Key == apimodels.UserExtKey_Undisturb {
			settings.Undisturb = setting.Value
		} else if setting.Key == apimodels.UserExtKey_FriendVerifyType {
			verifyType := tools.ToInt(setting.Value)
			settings.FriendVerifyType = pbobjs.FriendVerifyType(verifyType)
		} else if setting.Key == apimodels.UserExtKey_GrpVerifyType {
			verifyType := tools.ToInt(setting.Value)
			settings.GrpVerifyType = pbobjs.GrpVerifyType(verifyType)
		}
	}
	return settings
}

func SearchByPhone(ctx context.Context, phone string) (errs.IMErrorCode, *pbobjs.UserObjs) {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	appkey := bases.GetAppKeyFromCtx(ctx)
	targetUserId := tools.ShortMd5(phone)
	storage := userStorage.NewUserStorage()
	user, err := storage.FindByPhone(appkey, phone)
	if err == nil && user != nil {
		targetUserId = user.UserId
	}
	code, respObj, err := bases.SyncRpcCall(ctx, "qry_user_info", targetUserId, &pbobjs.UserIdReq{
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
	appkey := bases.GetAppKeyFromCtx(ctx)
	bases.SyncRpcCall(ctx, "upd_user_info", req.UserId, &pbobjs.UserInfo{
		UserId:       req.UserId,
		Nickname:     req.Nickname,
		UserPortrait: req.Avatar,
	}, nil)
	if req.Nickname != "" {
		//update assistant
		sdk := imsdk.GetImSdk(appkey)
		if sdk != nil {
			sdk.AddBot(juggleimsdk.BotInfo{
				BotId:    GetAssistantId(req.UserId),
				Nickname: GetAssistantNickname(req.Nickname),
				Portrait: req.Avatar,
				BotType:  int(commonservices.BotType_Custom),
			})
		}
		// update order tag for friends
		storage := storages.NewFriendRelStorage()
		appkey := bases.GetAppKeyFromCtx(ctx)
		storage.UpdateOrderTag(appkey, req.UserId, tools.GetFirstLetter(req.Nickname))
	}
	return errs.IMErrorCode_SUCCESS
}

func UpdateUserSettings(ctx context.Context, req *pbobjs.UserSettings) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	settings := []*pbobjs.KvItem{
		{
			Key:   apimodels.UserExtKey_Language,
			Value: req.Language,
		},
		{
			Key:   apimodels.UserExtKey_Undisturb,
			Value: req.Undisturb,
		},
		{
			Key:   apimodels.UserExtKey_FriendVerifyType,
			Value: tools.Int642String(int64(req.FriendVerifyType)),
		},
		{
			Key:   apimodels.UserExtKey_GrpVerifyType,
			Value: tools.Int642String(int64(req.GrpVerifyType)),
		},
	}
	if len(settings) > 0 {
		bases.SyncRpcCall(ctx, "upd_user_info", requestId, &pbobjs.UserInfo{
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
