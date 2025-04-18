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
	LatestMsgTime *int64 // latest msg time
	// LatestSendMsgTime *int64
	TerminalNum  int
	OnlineStatus bool //online state

	isNtf bool //is ntf

	PushSwitch int32
	PushBadge  int32
	CanPush    int32
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
	user.OnlineStatus = onlineStatus
	user.TerminalNum = terminalNum
}

func (user *UserStatus) IsOnline() bool {
	return user.OnlineStatus
}

func (user *UserStatus) SetPushStatus(canPush int32) {
	atomic.StoreInt32(&user.CanPush, canPush)
}

func (user *UserStatus) SetPushSwitch(pushSwitch int32) {
	atomic.StoreInt32(&user.PushSwitch, pushSwitch)
}

func (user *UserStatus) OpenPushSwitch() bool {
	return user.PushSwitch > 0
}

func (user *UserStatus) SetBadge(badge int32) {
	atomic.StoreInt32(&user.PushBadge, badge)
}

func (user *UserStatus) BadgeIncr() int32 {
	atomic.AddInt32(&user.PushBadge, 1)
	return user.PushBadge
}

func (user *UserStatus) CheckNtfWithSwitch() bool {
	if !user.OnlineStatus || user.TerminalNum > 1 {
		return true
	}
	if user.isNtf {
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
	if user.LatestMsgTime != nil && *user.LatestMsgTime == ackTime {
		user.isNtf = false
	}
}

func (user *UserStatus) SetLatestMsgTime(time int64) {
	key := getKey(user.appkey, user.userId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if user.LatestMsgTime == nil || *user.LatestMsgTime < time {
		user.LatestMsgTime = &time
	}
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
							OnlineStatus: item.IsOnline,
							CanPush:      1,
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
	if user.LatestMsgTime == nil || currentTime > *user.LatestMsgTime {
		user.LatestMsgTime = &currentTime
	} else {
		ret = *user.LatestMsgTime + 1
		user.LatestMsgTime = &ret
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
		OnlineStatus: true,
		CanPush:      1,
	}
}
