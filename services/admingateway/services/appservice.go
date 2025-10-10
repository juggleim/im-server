package services

import (
	"context"
	"fmt"
	"im-server/commons/tools"
	"im-server/services/admingateway/ctxs"
	"im-server/services/admingateway/dbs"
	commonDbs "im-server/services/commonservices/dbs"
	"im-server/services/commonservices/logs"
	userStorage "im-server/services/usermanager/storages"
	"math"
	"time"
)

var appFieldsMap map[string]bool

func init() {
	appFieldsMap = make(map[string]bool)
	appFieldsMap["is_hide_msg_before_join_group"] = true
	appFieldsMap["file_config"] = true
	appFieldsMap["event_sub_config"] = true
	appFieldsMap["event_sub_switch"] = true
	appFieldsMap["his_msg_save_day"] = true
	appFieldsMap["kick_mode"] = true
}

func CreateApp(appInfo AppInfo) (AdminErrorCode, *AppInfo) {
	dao := commonDbs.AppInfoDao{}
	if appInfo.AppKey == "" {
		appInfo.AppKey = tools.RandStr(16)
	}
	dbAppInfo := dao.FindByAppkey(appInfo.AppKey)
	if dbAppInfo != nil && dbAppInfo.AppKey == appInfo.AppKey {
		return AdminErrorCode_AppHasExisted, &AppInfo{
			AppName:     dbAppInfo.AppName,
			AppKey:      dbAppInfo.AppKey,
			AppSecret:   dbAppInfo.AppSecret,
			CreatedTime: dbAppInfo.CreatedTime.UnixMilli(),
		}
	}
	if len(appInfo.AppSecret) != 16 {
		appInfo.AppSecret = tools.RandStr(16)
	}
	newApp := commonDbs.AppInfoDao{
		AppName:      appInfo.AppName,
		AppKey:       appInfo.AppKey,
		AppSecret:    appInfo.AppSecret,
		AppSecureKey: tools.RandStr(16),
		AppType:      appInfo.AppType,
		CreatedTime:  time.Now(),
		UpdatedTime:  time.Now(),
	}
	err := dao.Upsert(newApp)
	if err != nil {
		return AdminErrorCode_AddAppFail, nil
	}
	return AdminErrorCode_Success, &AppInfo{
		AppType:   newApp.AppType,
		AppName:   newApp.AppName,
		AppKey:    newApp.AppKey,
		AppSecret: newApp.AppSecret,
	}
}

func QryApps(ctx context.Context, account string, limit int64, offset string) (AdminErrorCode, *Apps) {
	curAccount, exist := GetAccountInfo(ctxs.GetAccountFromCtx(ctx))
	if !exist || curAccount == nil {
		return AdminErrorCode_AccountNotExist, nil
	}
	apps := &Apps{
		Items:   []*SimpleApp{},
		HasMore: false,
		Offset:  "",
	}
	if curAccount.RoleType == RoleType_SuperAdmin && account == "" {
		dao := commonDbs.AppInfoDao{}
		offsetInt, err := tools.DecodeInt(offset)
		if err != nil {
			offsetInt = math.MaxInt64
		}
		dbApps, err := dao.QryApps(limit+1, offsetInt)
		if err == nil {
			if len(dbApps) > int(limit) {
				dbApps = dbApps[:len(dbApps)-1]
				apps.HasMore = true
			}
			var id int64 = math.MaxInt64
			for _, dbApp := range dbApps {
				app := &SimpleApp{
					AppKey:       dbApp.AppKey,
					AppName:      dbApp.AppName,
					CreatedTime:  dbApp.CreatedTime.UnixMilli(),
					AppType:      dbApp.AppType,
					MaxUserCount: 100,
				}
				storage := userStorage.NewUserStorage()
				app.CurUserCount = storage.Count(dbApp.AppKey)
				apps.Items = append(apps.Items, app)
				if id > 0 {
					offset, _ := tools.EncodeInt(id)
					apps.Offset = offset
				}
			}
		} else {
			logs.NewLogEntity().Error(err.Error())
		}
	} else {
		acc := curAccount.Account
		if curAccount.RoleType == RoleType_SuperAdmin && account != "" {
			acc = account
		}
		dao := dbs.AccountAppRelDao{}
		offsetInt, err := tools.DecodeInt(offset)
		if err != nil {
			offsetInt = math.MaxInt64
		}
		dbApps, err := dao.QryApps(acc, limit+1, offsetInt)
		if err == nil {
			if len(dbApps) > int(limit) {
				dbApps = dbApps[:len(dbApps)-1]
				apps.HasMore = true
			}
			var id int64 = math.MaxInt64
			for _, dbApp := range dbApps {
				app := &SimpleApp{
					AppKey:       dbApp.AppKey,
					AppName:      dbApp.AppName,
					CreatedTime:  dbApp.CreatedTime.UnixMilli(),
					AppType:      dbApp.AppType,
					MaxUserCount: 100,
				}
				storage := userStorage.NewUserStorage()
				app.CurUserCount = storage.Count(dbApp.AppKey)
				apps.Items = append(apps.Items, app)
				if dbApp.ID < id {
					id = dbApp.ID
				}
			}
			if id > 0 {
				offset, _ := tools.EncodeInt(id)
				apps.Offset = offset
			}
		} else {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	return AdminErrorCode_Success, apps
}

func QryApp(appkey string) *AppInfo {
	dao := commonDbs.AppInfoDao{}
	dbApp := dao.FindByAppkey(appkey)
	if dbApp == nil {
		return nil
	}
	appInfo := &AppInfo{
		AppType:      dbApp.AppType,
		AppName:      dbApp.AppName,
		AppKey:       dbApp.AppKey,
		AppSecret:    dbApp.AppSecret,
		CreatedTime:  dbApp.CreatedTime.UnixMilli(),
		UpdateTime:   dbApp.UpdatedTime.UnixMilli(),
		AppStatus:    dbApp.AppStatus,
		ConfigFields: make(map[string]string),
		MaxUserCount: 100,
	}
	storage := userStorage.NewUserStorage()
	appInfo.CurUserCount = storage.Count(dbApp.AppKey)
	//appext
	extDao := commonDbs.AppExtDao{}
	dbExts := extDao.FindListByAppkey(appkey)
	for _, dbExt := range dbExts {
		appInfo.ConfigFields[dbExt.AppItemKey] = dbExt.AppItemValue
	}

	return appInfo
}

func UpdateAppConfigs(appkey string, configFields map[string]interface{}) AdminErrorCode {
	//check fields
	// for fieldKey, _ := range configFields {
	// 	if _, exist := appFieldsMap[fieldKey]; !exist {
	// 		return AdminErrorCode_NotSupportField
	// 	}
	// }
	dao := commonDbs.AppExtDao{}
	for fieldKey, fieldValue := range configFields {
		dao.CreateOrUpdate(appkey, fieldKey, fmt.Sprintf("%s", fieldValue))
	}
	return AdminErrorCode_Success
}

func QryAppConfigs(appkey string, keys []string) (AdminErrorCode, *AppConfigs) {
	ret := &AppConfigs{
		AppKey:  appkey,
		Configs: map[string]interface{}{},
	}
	dao := commonDbs.AppExtDao{}
	extList, err := dao.FindByItemKeys(appkey, keys)
	extMap := map[string]string{}
	if err == nil {
		for _, ext := range extList {
			extMap[ext.AppItemKey] = ext.AppItemValue
		}
	}
	for _, key := range keys {
		if val, ok := extMap[key]; ok {
			ret.Configs[key] = val
		} else {
			ret.Configs[key] = ""
		}
	}
	return AdminErrorCode_Success, ret
}

type AppConfigs struct {
	AppKey  string                 `json:"app_key"`
	Configs map[string]interface{} `json:"configs"`
}
