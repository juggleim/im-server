package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	userStorage "im-server/services/usermanager/storages"
	userModels "im-server/services/usermanager/storages/models"
)

func SetPrivateGlobalMute(ctx context.Context, req *pbobjs.BatchMuteUsersReq) {
	if len(req.UserIds) <= 0 {
		return
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := userStorage.NewUserExtStorage()
	items := []userModels.UserExt{}
	for _, userId := range req.UserIds {
		user, exist := GetUserInfo(appkey, userId)
		if exist && user != nil {
			user.SetPriGlobalMute(req.IsDelete, 0)
		}
		if !req.IsDelete {
			items = append(items, userModels.UserExt{
				AppKey:    appkey,
				UserId:    userId,
				ItemKey:   string(commonservices.AttItemKey_PriGlobalMute),
				ItemValue: "0",
				ItemType:  int(commonservices.AttItemType_Setting),
			})
		}
	}
	if req.IsDelete {
		storage.BatchDelete(appkey, string(commonservices.AttItemKey_PriGlobalMute), req.UserIds)
	} else {
		storage.BatchUpsert(items)
	}
}

func QryPriGlobalMuteUsers(ctx context.Context, req *pbobjs.QryBlockUsersReq) (errs.IMErrorCode, *pbobjs.QryBlockUsersResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := userStorage.NewUserExtStorage()
	var startId int64 = 0
	if req.Offset != "" {
		offsetInt, err := tools.DecodeInt(req.Offset)
		if err == nil {
			startId = offsetInt
		}
	}
	ret := &pbobjs.QryBlockUsersResp{
		Items: []*pbobjs.BlockUser{},
	}
	var maxId int64 = 0
	exts, err := storage.QryExtsBaseItemKey(appkey, string(commonservices.AttItemKey_PriGlobalMute), startId, req.Limit)
	if err == nil {
		for _, ext := range exts {
			if ext.ID > maxId {
				maxId = ext.ID
			}
			ret.Items = append(ret.Items, &pbobjs.BlockUser{
				BlockUserId: ext.UserId,
				CreatedTime: ext.UpdatedTime.UnixMilli(),
			})
		}
		if maxId > 0 {
			ret.Offset, _ = tools.EncodeInt(maxId)
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func SetGroupGlobalMute(ctx context.Context, req *pbobjs.BatchMuteUsersReq) {
	if len(req.UserIds) <= 0 {
		return
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := userStorage.NewUserExtStorage()
	items := []userModels.UserExt{}
	for _, userId := range req.UserIds {
		user, exist := GetUserInfo(appkey, userId)
		if exist && user != nil {
			user.SetGrpGlobalMute(req.IsDelete, 0)
		}
		if !req.IsDelete {
			items = append(items, userModels.UserExt{
				AppKey:    appkey,
				UserId:    userId,
				ItemKey:   string(commonservices.AttItemKey_GrpGlobalMute),
				ItemValue: "0",
				ItemType:  int(commonservices.AttItemType_Setting),
			})
		}
	}
	if req.IsDelete {
		err := storage.BatchDelete(appkey, string(commonservices.AttItemKey_GrpGlobalMute), req.UserIds)
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
	} else {
		err := storage.BatchUpsert(items)
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
	}
}

func QryGrpGlobalMuteUsers(ctx context.Context, req *pbobjs.QryBlockUsersReq) (errs.IMErrorCode, *pbobjs.QryBlockUsersResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := userStorage.NewUserExtStorage()
	var startId int64 = 0
	if req.Offset != "" {
		offsetInt, err := tools.DecodeInt(req.Offset)
		if err == nil {
			startId = offsetInt
		}
	}
	ret := &pbobjs.QryBlockUsersResp{
		Items: []*pbobjs.BlockUser{},
	}
	var maxId int64 = 0
	exts, err := storage.QryExtsBaseItemKey(appkey, string(commonservices.AttItemKey_GrpGlobalMute), startId, req.Limit)
	if err == nil {
		for _, ext := range exts {
			if ext.ID > maxId {
				maxId = ext.ID
			}
			ret.Items = append(ret.Items, &pbobjs.BlockUser{
				BlockUserId: ext.UserId,
				CreatedTime: ext.UpdatedTime.UnixMilli(),
			})
		}
		if maxId > 0 {
			ret.Offset, _ = tools.EncodeInt(maxId)
		}
	} else {
		logs.WithContext(ctx).Error(err.Error())
	}
	return errs.IMErrorCode_SUCCESS, ret
}
