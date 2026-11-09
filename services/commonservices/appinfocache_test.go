package commonservices

import (
	"im-server/services/commonservices/dbs"
	"testing"
)

func TestAppExtensionsProtectCoreFieldsAndPreserveWhitelist(t *testing.T) {
	app := &AppInfo{AppKey: "app", AppSecret: "primary", AppSecureKey: "secure", AppStatus: 1}
	extMap := buildAppExtMap([]*dbs.AppExtDao{
		nil,
		{AppItemKey: "APP_KEY", AppItemValue: "override"},
		{AppItemKey: "app_secret", AppItemValue: "override"},
		{AppItemKey: " APP_SECRET ", AppItemValue: "override"},
		{AppItemKey: "app_secure_key", AppItemValue: "override"},
		{AppItemKey: "app_status", AppItemValue: "0"},
		{AppItemKey: "temporary_secret", AppItemValue: "temporary"},
		{AppItemKey: "is_sub_api_msg", AppItemValue: "true"},
		{AppItemKey: "api_msg_pre_callback_user_whitelist", AppItemValue: `["allowed"]`},
	})
	FillObjField(app, extMap)
	app.loadApiMsgPreCallbackUserWhitelist()
	if app.AppKey != "app" || app.AppSecret != "primary" || app.AppSecureKey != "secure" || app.AppStatus != 1 {
		t.Fatalf("core fields overridden: %+v", app)
	}
	if app.TemporarySecret != "temporary" || !app.IsSubApiMsg {
		t.Fatal("extension settings lost")
	}
	if !app.IsApiMsgPreCallbackUserAllowed("allowed") || app.IsApiMsgPreCallbackUserAllowed("other") {
		t.Fatal("whitelist lost")
	}
}
