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
	pushTokenCache = caches.NewLruCacheWithAddReadTimeout("pushtoken_cache", 100000, nil, 5*time.Minute, 5*time.Minute)
	pushTokenLocks = tools.NewSegmentatedLocks(512)
	notExistPushToken = &UserPushToken{DeviceId: ""}
}

type UserPushToken struct {
	AppKey        string
	UserId        string
	DeviceId      string
	Platform      pbobjs.Platform
	PushChannel   pbobjs.PushChannel
	PushToken     string
	VoipPushToken string
	PackageName   string
}

func (pushToken *UserPushToken) IsSame(pushToken2 *UserPushToken) bool {
	if pushToken.DeviceId != pushToken2.DeviceId ||
		pushToken.Platform != pushToken2.Platform ||
		pushToken.PushChannel != pushToken2.PushChannel ||
		pushToken.PackageName != pushToken2.PackageName ||
		(pushToken2.PushToken != "" && pushToken.PushToken != pushToken2.PushToken) ||
		(pushToken2.VoipPushToken != "" && pushToken.VoipPushToken != pushToken2.VoipPushToken) {
		return false
	}
	return true
}

func getUserPushTokenKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}

func AddPushToken(appkey, userId string, userPushToken *UserPushToken) {
	cachePushToken := GetPushToken(appkey, userId)
	if !cachePushToken.IsSame(userPushToken) {
		cachePushToken.DeviceId = userPushToken.DeviceId
		cachePushToken.Platform = userPushToken.Platform
		cachePushToken.PushChannel = userPushToken.PushChannel
		cachePushToken.PackageName = userPushToken.PackageName
		if userPushToken.PushToken != "" {
			cachePushToken.PushToken = userPushToken.PushToken
		}
		if userPushToken.VoipPushToken != "" {
			cachePushToken.VoipPushToken = userPushToken.VoipPushToken
		}
		//save to db
		dao := dbs.PushTokenDao{}
		dao.Upsert(dbs.PushTokenDao{
			AppKey:      appkey,
			UserId:      userId,
			DeviceId:    cachePushToken.DeviceId,
			Platform:    commonservices.Platform2Str(cachePushToken.Platform),
			PushChannel: commonservices.PushChannel2Str(cachePushToken.PushChannel),
			Package:     cachePushToken.PackageName,
			PushToken:   cachePushToken.PushToken,
			VoipToken:   cachePushToken.VoipPushToken,
		})
		//清除之前的用户缓存（同一设备先登录A，退出A登录B， A不应该收到推送消息）
		items, err := dao.QueryByDeviceId(appkey, userPushToken.DeviceId)
		if err == nil {
			for _, item := range items {
				key := getUserPushTokenKey(appkey, item.UserId)
				pushTokenCache.Remove(key)
			}
		}
		//clearn other user for this device
		dao.DeleteByDeviceId(appkey, userPushToken.DeviceId, userId)
	}
}

func RemovePushToken(appkey, userId string) {
	key := getUserPushTokenKey(appkey, userId)
	pushToken := GetPushToken(appkey, userId)
	if pushToken != nil && (pushToken.PushToken != "" || pushToken.VoipPushToken != "") {
		dao := dbs.PushTokenDao{}
		dao.DeleteByUserId(appkey, userId)
		pushTokenCache.Remove(key)
	}
}

func GetPushToken(appkey, userId string) *UserPushToken {
	key := getUserPushTokenKey(appkey, userId)
	if obj, exist := pushTokenCache.Get(key); exist {
		pushToken := obj.(*UserPushToken)
		return pushToken
	} else {
		lock := pushTokenLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		userPushToken := &UserPushToken{
			AppKey: appkey,
			UserId: userId,
		}
		//read from db
		dao := dbs.PushTokenDao{}
		dbPushToken, err := dao.FindByUserId(appkey, userId)
		if err == nil && dbPushToken != nil {
			userPushToken.DeviceId = dbPushToken.DeviceId
			userPushToken.Platform = commonservices.Str2Platform(dbPushToken.Platform)
			userPushToken.PushChannel = commonservices.Str2PushChannel(dbPushToken.PushChannel)
			userPushToken.PackageName = dbPushToken.Package
			userPushToken.PushToken = dbPushToken.PushToken
			userPushToken.VoipPushToken = dbPushToken.VoipToken
		}
		pushTokenCache.Add(key, userPushToken)
		return userPushToken
	}
}

func RegPushToken(ctx context.Context, userId string, req *pbobjs.RegPushTokenReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	AddPushToken(appkey, userId, &UserPushToken{
		DeviceId:      req.DeviceId,
		Platform:      req.Platform,
		PushChannel:   req.PushChannel,
		PushToken:     req.PushToken,
		VoipPushToken: req.VoipToken,
		PackageName:   req.PackageName,
	})
	if req.PushToken != "" {
		//ntf open push
		bases.AsyncRpcCall(ctx, "upd_push_status", userId, &pbobjs.UserPushStatus{
			CanPush: true,
		})
	}
	return errs.IMErrorCode_SUCCESS
}
