package services

import (
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	"time"
)

func CreateApp(appInfo AppInfo) (AdminErrorCode, *AppInfo) {
	dao := dbs.AppInfoDao{}
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
	newApp := dbs.AppInfoDao{
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
