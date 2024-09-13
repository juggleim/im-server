package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/pushmanager/storages/dbs"
	"strings"
	"time"
)

var (
	pushTokenCache *caches.LruCache
	pushTokenLocks *tools.SegmentatedLocks

	notExistPushToken *UserPushToken
)

func init() {
	pushTokenCache = caches.NewLruCacheWithAddReadTimeout(100000, nil, 5*time.Minute, 5*time.Minute)
	pushTokenLocks = tools.NewSegmentatedLocks(512)
	notExistPushToken = &UserPushToken{DeviceId: ""}
}

type UserPushToken struct {
	DeviceId    string
	Platform    pbobjs.Platform
	PushChannel pbobjs.PushChannel
	PushToken   string
	PackageName string
}

func (pushToken *UserPushToken) IsSame(pushToken2 *UserPushToken) bool {
	if pushToken.DeviceId == pushToken2.DeviceId && pushToken.Platform == pushToken2.Platform && pushToken.PushChannel == pushToken2.PushChannel && pushToken.PushToken == pushToken2.PushToken && pushToken.PackageName == pushToken2.PackageName {
		return true
	}
	return false
}

func getUserPushTokenKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}

func AddPushToken(appkey, userId string, userPushToken *UserPushToken) {
	key := getUserPushTokenKey(appkey, userId)
	cachePushToken, exist := GetPushToken(appkey, userId)
	if !exist || !cachePushToken.IsSame(userPushToken) {
		//add to cache
		pushTokenCache.Add(key, userPushToken)
		//save to db
		dao := dbs.PushTokenDao{}
		dao.UpsertPushToken(dbs.PushTokenDao{
			AppKey:      appkey,
			UserId:      userId,
			DeviceId:    userPushToken.DeviceId,
			Platform:    commonservices.Platform2Str(userPushToken.Platform),
			PushChannel: commonservices.PushChannel2Str(userPushToken.PushChannel),
			Package:     userPushToken.PackageName,
			PushToken:   userPushToken.PushToken,
		})
		//清除之前的用户缓存（同一设备先登录A，退出A登录B， A不应该收到推送消息）
		prevItem, _ := dao.GetUserWithToken(appkey, userPushToken.PushToken)
		if prevItem != nil {
			pushConfCache.Remove(getUserPushTokenKey(appkey, prevItem.UserId))
		}
		//clearn other user for this device
		dao.DeleteByDeviceId(appkey, userPushToken.DeviceId, userId)
	}
}

func RemovePushToken(appkey, userId string) {
	key := getUserPushTokenKey(appkey, userId)
	_, exist := GetPushToken(appkey, userId)
	if exist {
		dao := dbs.PushTokenDao{}
		dao.DeleteByUserId(appkey, userId)
		pushConfCache.Remove(key)
	}
}

func GetPushToken(appkey, userId string) (*UserPushToken, bool) {
	key := getUserPushTokenKey(appkey, userId)
	if obj, exist := pushTokenCache.Get(key); exist {
		pushToken := obj.(*UserPushToken)
		if pushToken != notExistPushToken {
			return pushToken, true
		} else {
			return nil, false
		}
	} else {
		lock := pushTokenLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		//read from db
		dao := dbs.PushTokenDao{}
		dbPushToken, err := dao.FindByUserId(appkey, userId)
		if err == nil && dbPushToken != nil {
			userPushToken := &UserPushToken{
				DeviceId:    dbPushToken.DeviceId,
				Platform:    commonservices.Str2Platform(dbPushToken.Platform),
				PushChannel: commonservices.Str2PushChannel(dbPushToken.PushChannel),
				PushToken:   dbPushToken.PushToken,
				PackageName: dbPushToken.Package,
			}
			pushTokenCache.Add(key, userPushToken)
			return userPushToken, true
		} else {
			pushTokenCache.Add(key, notExistPushToken)
			return nil, false
		}
	}
}

func RegPushToken(ctx context.Context, userId string, req *pbobjs.RegPushTokenReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	AddPushToken(appkey, userId, &UserPushToken{
		DeviceId:    req.DeviceId,
		Platform:    req.Platform,
		PushChannel: req.PushChannel,
		PushToken:   req.PushToken,
		PackageName: req.PackageName,
	})
	return errs.IMErrorCode_SUCCESS
}
