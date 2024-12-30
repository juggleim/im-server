package services

import (
	"fmt"
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	userDbs "im-server/services/usermanager/dbs"
	"math"
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

func QryApps(limit int64, offset string) *Apps {
	apps := &Apps{
		Items:   []*SimpleApp{},
		HasMore: false,
		Offset:  "",
	}
	dao := dbs.AppInfoDao{}
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
			userDao := userDbs.UserDao{}
			app.CurUserCount = userDao.Count(dbApp.AppKey)
			apps.Items = append(apps.Items, app)
			if dbApp.ID < id {
				id = dbApp.ID
			}
		}
		if id > 0 {
			offset, _ := tools.EncodeInt(id)
			apps.Offset = offset
		}
	}
	return apps
}

func QryApp(appkey string) *AppInfo {
	dao := dbs.AppInfoDao{}
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
	userDao := userDbs.UserDao{}
	appInfo.CurUserCount = userDao.Count(dbApp.AppKey)
	//appext
	extDao := dbs.AppExtDao{}
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
	dao := dbs.AppExtDao{}
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
	dao := dbs.AppExtDao{}
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
