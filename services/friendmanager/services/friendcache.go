package services

import (
	"context"
	"im-server/commons/caches"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/friendmanager/storages"
	"strings"
	"time"
)

var friendRelCache *caches.LruCache
var friendRelLocks *tools.SegmentatedLocks

func init() {
	friendRelCache = caches.NewLruCacheWithAddReadTimeout("friend_cache", 10000, nil, 10*time.Minute, 10*time.Minute)
	friendRelLocks = tools.NewSegmentatedLocks(128)
}

type FriendContainer struct {
	AppKey    string
	UserId    string
	FriendMap map[string]int64
}

func (container *FriendContainer) AddFriends(friendIds []string) {
	key := getFriendKey(container.AppKey, container.UserId)
	l := friendRelLocks.GetLocks(key)
	l.Lock()
	defer l.Unlock()
	for _, friendId := range friendIds {
		container.FriendMap[friendId] = 0
	}
}

func (container *FriendContainer) DelFriends(friendIds []string) {
	key := getFriendKey(container.AppKey, container.UserId)
	l := friendRelLocks.GetLocks(key)
	l.Lock()
	defer l.Unlock()
	for _, friendId := range friendIds {
		delete(container.FriendMap, friendId)
	}
}

func (container *FriendContainer) ForeachFriends(f func(friendId string)) {
	key := getFriendKey(container.AppKey, container.UserId)
	l := friendRelLocks.GetLocks(key)
	l.RLock()
	defer l.RUnlock()
	for k := range container.FriendMap {
		f(k)
	}
}

func (container *FriendContainer) CheckFriend(friendId string) bool {
	key := getFriendKey(container.AppKey, container.UserId)
	l := friendRelLocks.GetLocks(key)
	l.RLock()
	defer l.RUnlock()
	if _, exist := container.FriendMap[friendId]; exist {
		return true
	}
	return false
}

func (container *FriendContainer) BatchCheckFriends(friendIds []string) map[string]bool {
	key := getFriendKey(container.AppKey, container.UserId)
	l := friendRelLocks.GetLocks(key)
	l.RLock()
	defer l.RUnlock()
	ret := make(map[string]bool)
	for _, friendId := range friendIds {
		if _, exist := container.FriendMap[friendId]; exist {
			ret[friendId] = true
		} else {
			ret[friendId] = false
		}
	}
	return ret
}

func getFriendKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}

func GetFriendContainer(ctx context.Context, appkey, userId string) *FriendContainer {
	key := getFriendKey(appkey, userId)
	if container, exist := friendRelCache.Get(key); exist {
		return container.(*FriendContainer)
	} else {
		l := friendRelLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if container, exist := friendRelCache.Get(key); exist {
			return container.(*FriendContainer)
		} else {
			container := &FriendContainer{
				AppKey:    appkey,
				UserId:    userId,
				FriendMap: map[string]int64{},
			}
			storage := storages.NewFriendRelStorage()
			var startId int64 = 0
			var limit int64 = 10000
			for {
				friendRels, err := storage.QueryFriendRels(appkey, userId, startId, limit, true)
				if err == nil {
					for _, rel := range friendRels {
						container.FriendMap[rel.FriendId] = 0
						if rel.ID > startId {
							startId = rel.ID
						}
					}
					if len(friendRels) < int(limit) {
						break
					}
				} else {
					logs.WithContext(ctx).Error(err.Error())
					break
				}
			}
			friendRelCache.Add(key, container)
			return container
		}
	}
}
