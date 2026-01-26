package friendcache

import (
	"im-server/commons/caches"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/friendmanager/storages"
	"strings"
	"time"
)

type FriendStatus struct {
	IsFriend          bool
	FriendDisplayName string
	UpdatedTime       int64
}

var friendStatusCache *caches.LruCache
var friendStatusLocks *tools.SegmentatedLocks

func init() {
	friendStatusCache = caches.NewLruCacheWithAddReadTimeout("friendstatus_cache", 100000, func(key, value interface{}) {}, 10*time.Minute, 10*time.Minute)
	friendStatusLocks = tools.NewSegmentatedLocks(512)
}

func GetFriendStatus(appkey, userId, friendId string) *FriendStatus {
	key := getFriendStatusCacheKey(appkey, userId, friendId)
	if val, exist := friendStatusCache.Get(key); exist {
		return val.(*FriendStatus)
	} else {
		l := friendStatusLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if val, exist := friendStatusCache.Get(key); exist {
			return val.(*FriendStatus)
		} else {
			status := &FriendStatus{}
			storage := storages.NewFriendRelStorage()
			rel, err := storage.GetFriendRel(appkey, userId, friendId)
			if err == nil && rel != nil {
				status.IsFriend = true
				status.FriendDisplayName = rel.DisplayName
				status.UpdatedTime = rel.UpdatedTime
			} else if err != nil {
				logs.NewLogEntity().Error(err.Error())
			}
			friendStatusCache.Add(key, status)
			return status
		}
	}
}

func BatchGetFriendStatus(appkey, userId string, friendIds []string) map[string]*FriendStatus {
	ret := map[string]*FriendStatus{}
	noCacheFriendIds := []string{}
	for _, friendId := range friendIds {
		key := getFriendStatusCacheKey(appkey, userId, friendId)
		if !friendStatusCache.Contains(key) {
			noCacheFriendIds = append(noCacheFriendIds, friendId)
		} else {
			ret[friendId] = GetFriendStatus(appkey, userId, friendId)
		}
	}
	if len(noCacheFriendIds) > 0 {
		storage := storages.NewFriendRelStorage()
		rels, err := storage.QueryFriendRelsByFriendIds(appkey, userId, noCacheFriendIds)
		if err == nil {
			for _, rel := range rels {
				status := &FriendStatus{
					IsFriend:          true,
					FriendDisplayName: rel.DisplayName,
					UpdatedTime:       rel.UpdatedTime,
				}
				ret[rel.FriendId] = status
				key := getFriendStatusCacheKey(appkey, userId, rel.FriendId)
				friendStatusCache.Add(key, status)
			}
		}
	}
	return ret
}

func RemoveFriendStatus(appkey, userId, friendId string) {
	key := getFriendStatusCacheKey(appkey, userId, friendId)
	// l := userLocks.GetLocks(key)
	// l.Lock()
	// defer l.Unlock()
	friendStatusCache.Remove(key)
}

func getFriendStatusCacheKey(appkey, userId, friendId string) string {
	return strings.Join([]string{appkey, userId, friendId}, "_")
}
