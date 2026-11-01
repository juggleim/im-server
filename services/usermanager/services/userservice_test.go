package services

import (
	"reflect"
	"testing"
	"time"

	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/usermanager/storages/models"
)

type batchUserStorageStub struct {
	models.IUserStorage
	calls   int
	userIds []string
}

func (stub *batchUserStorageStub) FindByUserIds(_ string, userIds []string) ([]*models.User, error) {
	stub.calls++
	stub.userIds = append([]string(nil), userIds...)
	return []*models.User{{
		AppKey:      "batch-test-app",
		UserId:      "user-1",
		Nickname:    "nickname",
		UpdatedTime: time.UnixMilli(1000),
	}}, nil
}

type batchUserExtStorageStub struct {
	models.IUserExtStorage
	calls   int
	userIds []string
}

func (stub *batchUserExtStorageStub) QryExtFieldsByUserIds(_ string, userIds []string) ([]*models.UserExt, error) {
	stub.calls++
	stub.userIds = append([]string(nil), userIds...)
	return []*models.UserExt{{
		UserId:      "user-1",
		ItemKey:     "setting",
		ItemValue:   "setting-value",
		ItemType:    int(commonservices.AttItemType_Setting),
		UpdatedTime: time.UnixMilli(2000),
	}}, nil
}

func TestBatchGetUserInfosLoadsCacheMissesInBatches(t *testing.T) {
	appkey := "batch-test-app"
	t.Cleanup(func() {
		userCache.Remove(appkey + "_user-1")
		userCache.Remove(appkey + "_missing")
	})
	userStorage := &batchUserStorageStub{}
	extStorage := &batchUserExtStorageStub{}
	userIds := []string{"user-1", "missing", "user-1", ""}

	got := batchGetUserInfos(appkey, userIds, userStorage, extStorage)
	if userStorage.calls != 1 || !reflect.DeepEqual(userStorage.userIds, []string{"user-1", "missing"}) {
		t.Fatalf("user batch calls=%d ids=%v", userStorage.calls, userStorage.userIds)
	}
	if extStorage.calls != 1 || !reflect.DeepEqual(extStorage.userIds, []string{"user-1"}) {
		t.Fatalf("user ext batch calls=%d ids=%v", extStorage.calls, extStorage.userIds)
	}
	if got["user-1"] == nil || got["user-1"].SettingFields["setting"] != "setting-value" {
		t.Fatalf("user info=%+v", got["user-1"])
	}
	if _, exists := got["missing"]; exists {
		t.Fatal("missing user must not be returned")
	}

	got = batchGetUserInfos(appkey, userIds, userStorage, extStorage)
	if userStorage.calls != 1 || extStorage.calls != 1 {
		t.Fatalf("cache hit queried storage again: user=%d ext=%d", userStorage.calls, extStorage.calls)
	}
	if got["user-1"] == nil {
		t.Fatal("cached user info is missing")
	}
}

func TestUserModelsToUserInfosAndApplyUserExts(t *testing.T) {
	userUpdatedAt := time.UnixMilli(1000)
	attUpdatedAt := time.UnixMilli(2000)
	settingUpdatedAt := time.UnixMilli(3000)
	statusUpdatedAt := time.UnixMilli(4000)
	userInfoMap := userModelsToUserInfos([]*models.User{
		{
			AppKey:       "app",
			UserId:       "user-1",
			UserType:     pbobjs.UserType_Bot,
			Nickname:     "nickname",
			UserPortrait: "portrait",
			UpdatedTime:  userUpdatedAt,
		},
	})

	applyUserExts(userInfoMap, []*models.UserExt{
		{UserId: "user-1", ItemKey: "profile", ItemValue: "profile-value", ItemType: int(commonservices.AttItemType_Att), UpdatedTime: attUpdatedAt},
		{UserId: "user-1", ItemKey: "setting", ItemValue: "setting-value", ItemType: int(commonservices.AttItemType_Setting), UpdatedTime: settingUpdatedAt},
		{UserId: "user-1", ItemKey: "status", ItemValue: "status-value", ItemType: int(commonservices.AttItemType_Status), UpdatedTime: statusUpdatedAt},
		{UserId: "missing", ItemKey: "ignored", ItemValue: "ignored", ItemType: int(commonservices.AttItemType_Att), UpdatedTime: statusUpdatedAt},
	})

	userInfo := userInfoMap["user-1"]
	if userInfo == nil {
		t.Fatal("user info is missing")
	}
	if userInfo.AppKey != "app" || userInfo.UserType != int(pbobjs.UserType_Bot) {
		t.Fatalf("unexpected base user info: %+v", userInfo)
	}
	if got := userInfo.ExtFields["profile"]; got != "profile-value" {
		t.Fatalf("profile=%q", got)
	}
	if got := userInfo.SettingFields["setting"]; got != "setting-value" {
		t.Fatalf("setting=%q", got)
	}
	if got := userInfo.Statuses["status"]; got == nil || got.ItemValue != "status-value" || got.UpdatedTime != statusUpdatedAt.UnixMilli() {
		t.Fatalf("status=%+v", got)
	}
	if userInfo.UpdatedTime != attUpdatedAt.UnixMilli() {
		t.Fatalf("updated time=%d, want=%d", userInfo.UpdatedTime, attUpdatedAt.UnixMilli())
	}
}
