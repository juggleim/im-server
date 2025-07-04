package services

import (
	"im-server/commons/caches"
	"im-server/services/friends/storages"
	"strings"
	"time"
)

type FriendStatus struct {
	IsFriend          bool
	FriendDisplayName string
}

var friendStatusCache *caches.LruCache

func init() {
	friendStatusCache = caches.NewLruCacheWithAddReadTimeout("friendstatus_cache", 100000, func(key, value interface{}) {}, 10*time.Minute, 10*time.Minute)
}

func GetFriendStatus(appkey, userId, friendId string) *FriendStatus {
	key := getFriendStatusCacheKey(appkey, userId, friendId)
	if val, exist := friendStatusCache.Get(key); exist {
		return val.(*FriendStatus)
	} else {
		l := userLocks.GetLocks(key)
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
			}
			friendStatusCache.Add(key, status)
			return status
		}
	}
}

func getFriendStatusCacheKey(appkey, userId, friendId string) string {
	return strings.Join([]string{appkey, userId, friendId}, "_")
}
