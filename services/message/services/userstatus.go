package services

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"

	"google.golang.org/protobuf/proto"
)

type UserStatus struct {
	appkey string
	userId string
	// LastSyncTime        *int64
	// LastSendBoxSyncTime *int64
	latestMsgTime *int64 // latest msg time
	// LatestSendMsgTime *int64
	terminalNum  int
	onlineStatus bool //online state

	isNtf bool //is ntf

	pushSwitch int32
	pushBadge  int32
	canPush    int32
}

var userOnlineStatusCache *caches.LruCache
var userLocks *tools.SegmentatedLocks

func init() {
	userOnlineStatusCache = caches.NewLruCacheWithReadTimeout("useronlinestatus_cache", 100000, func(key, value interface{}) {}, time.Hour)
	userLocks = tools.NewSegmentatedLocks(512)
}

/*
record user's  status when sync msg
*/
func RecordUserOnlineStatus(appKey, userId string, onlineStatus bool, terminalNum int) {
	user := GetUserStatus(appKey, userId)
	key := getKey(appKey, userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	user.onlineStatus = onlineStatus
	user.terminalNum = terminalNum
}

func (user *UserStatus) IsOnline() bool {
	key := getKey(user.appkey, user.userId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	return user.onlineStatus
}

func (user *UserStatus) SetPushStatus(canPush int32) {
	atomic.StoreInt32(&user.canPush, canPush)
}

func (user *UserStatus) CanPush() bool {
	return atomic.LoadInt32(&user.canPush) > 0
}

func (user *UserStatus) SetPushSwitch(pushSwitch int32) {
	atomic.StoreInt32(&user.pushSwitch, pushSwitch)
}

func (user *UserStatus) OpenPushSwitch() bool {
	return atomic.LoadInt32(&user.pushSwitch) > 0
}

func (user *UserStatus) SetBadge(badge int32) {
	atomic.StoreInt32(&user.pushBadge, badge)
}

func (user *UserStatus) BadgeIncr() int32 {
	atomic.AddInt32(&user.pushBadge, 1)
	return atomic.LoadInt32(&user.pushBadge)
}

func (user *UserStatus) mustNtf() bool {
	key := getKey(user.appkey, user.userId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	if !user.onlineStatus || user.terminalNum > 1 {
		return true
	}
	if user.isNtf {
		return true
	}
	return false
}

func (user *UserStatus) CheckNtfWithSwitch() bool {
	if user.mustNtf() {
		return true
	} else {
		key := getKey(user.appkey, user.userId)
		lock := userLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if user.isNtf {
			return true
		} else {
			ret := user.isNtf
			user.isNtf = true
			return ret
		}
	}
}

func (user *UserStatus) SetNtfStatus(isNtf bool) {
	key := getKey(user.appkey, user.userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	user.isNtf = isNtf
}

func (user *UserStatus) CloseNtf(ackTime int64) {
	key := getKey(user.appkey, user.userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if user.latestMsgTime != nil && *user.latestMsgTime == ackTime {
		user.isNtf = false
	}
}

func (user *UserStatus) SetLatestMsgTime(time int64) {
	key := getKey(user.appkey, user.userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if user.latestMsgTime == nil || *user.latestMsgTime < time {
		user.latestMsgTime = &time
	}
}

func (user *UserStatus) GetLatestMsgTime() (int64, bool) {
	key := getKey(user.appkey, user.userId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()

	if user.latestMsgTime == nil {
		return 0, false
	}
	return *user.latestMsgTime, true
}

func UserStatusCacheContains(appkey, userId string) bool {
	key := getKey(appkey, userId)
	return userOnlineStatusCache.Contains(key)
}

func GetUserStatus(appKey, userId string) *UserStatus {
	key := getKey(appKey, userId)
	if val, exist := userOnlineStatusCache.Get(key); exist {
		return val.(*UserStatus)
	} else {
		l := userLocks.GetLocks(appKey, userId)
		l.Lock()
		defer l.Unlock()
		if val, exist := userOnlineStatusCache.Get(key); exist {
			return val.(*UserStatus)
		} else {
			userInfo := initUserStatus(appKey, userId)
			userOnlineStatusCache.Add(key, userInfo)
			return userInfo
		}
	}
}

func CacheUserStatus(appkey, userId string, status *UserStatus) {
	key := getKey(appkey, userId)
	l := userLocks.GetLocks(key)
	l.Lock()
	defer l.Unlock()
	if !UserStatusCacheContains(appkey, userId) {
		userOnlineStatusCache.Add(key, status)
	}
}

func BatchInitUserStatus(ctx context.Context, appkey string, userIds []string) {
	//check status from connect manager
	groups := bases.GroupTargets("qry_online_status", userIds)
	wg := sync.WaitGroup{}
	for _, ids := range groups {
		wg.Add(1)
		uIds := ids
		go func() {
			defer wg.Done()
			_, resp, err := bases.SyncRpcCall(ctx, "qry_online_status", uIds[0], &pbobjs.UserOnlineStatusReq{
				UserIds: uIds,
			}, func() proto.Message {
				return &pbobjs.UserOnlineStatusResp{}
			})
			if err == nil {
				onlineResp, ok := resp.(*pbobjs.UserOnlineStatusResp)
				if ok && len(onlineResp.Items) > 0 {
					for _, item := range onlineResp.Items {
						CacheUserStatus(appkey, item.UserId, &UserStatus{
							appkey:       appkey,
							userId:       item.UserId,
							onlineStatus: item.IsOnline,
							canPush:      1,
						})
					}
				}
			}
		}()
	}
	wg.Wait()
}

func RegenateSendTime(appkey, userId string, currentTime int64) int64 {
	user := GetUserStatus(appkey, userId)

	key := getKey(appkey, userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	ret := currentTime
	if user.latestMsgTime == nil || currentTime > *user.latestMsgTime {
		user.latestMsgTime = &currentTime
	} else {
		ret = *user.latestMsgTime + 1
		user.latestMsgTime = &ret
	}
	return ret
}

func getKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}

func initUserStatus(appkey, userId string) *UserStatus {
	return &UserStatus{
		appkey:       appkey,
		userId:       userId,
		onlineStatus: true,
		canPush:      1,
	}
}
