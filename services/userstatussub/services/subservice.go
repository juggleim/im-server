package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/userstatussub/storages/dbs"
	"strings"
	"time"
)

var subRelCache *caches.LruCache
var userLocks *tools.SegmentatedLocks

func init() {
	subRelCache = caches.NewLruCacheWithAddReadTimeout("subrel_cache", 10000, nil, 10*time.Minute, 10*time.Minute)
	userLocks = tools.NewSegmentatedLocks(256)
}

func Subscribe(ctx context.Context, userId string, targetIds []string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	for _, targetId := range targetIds {
		//add for cache
		rels := GetSubRelationsFromCache(appkey, targetId)
		succ := rels.AddSubscriber(userId)

		if succ {
			//insert into db
			dao := dbs.SubRelationDao{}
			dao.Create(dbs.SubRelationDao{
				UserId:     targetId,
				Subscriber: userId,
				AppKey:     appkey,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func UnSubscribe(ctx context.Context, userId string, targetIds []string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	for _, targetId := range targetIds {
		//add for cache
		rels := GetSubRelationsFromCache(appkey, targetId)
		succ := rels.DelSubscriber(userId)

		if succ {
			//insert into db
			dao := dbs.SubRelationDao{}
			dao.Delete(appkey, targetId, userId)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

type SubRelations struct {
	AppKey      string
	UserId      string
	subscribers map[string]bool
}

func (rel *SubRelations) GetSubscriptions() []string {
	key := getKey(rel.AppKey, rel.UserId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	ids := []string{}
	for userId, _ := range rel.subscribers {
		ids = append(ids, userId)
	}
	return ids
}

func (rel *SubRelations) AddSubscriber(subscriber string) bool {
	key := getKey(rel.AppKey, rel.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if _, exist := rel.subscribers[subscriber]; exist {
		return false
	} else {
		rel.subscribers[subscriber] = true
		return true
	}
}

func (rel *SubRelations) DelSubscriber(subscriber string) bool {
	key := getKey(rel.AppKey, rel.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if _, exist := rel.subscribers[subscriber]; exist {
		delete(rel.subscribers, subscriber)
		return true
	}
	return false
}

func GetSubRelationsFromCache(appkey, userId string) *SubRelations {
	key := getKey(appkey, userId)
	if obj, exist := subRelCache.Get(key); exist {
		rel := obj.(*SubRelations)
		return rel
	} else {
		lock := userLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if obj, exist := subRelCache.Get(key); exist {
			return obj.(*SubRelations)
		} else {
			dao := dbs.SubRelationDao{}
			rels := &SubRelations{
				AppKey:      appkey,
				UserId:      userId,
				subscribers: make(map[string]bool),
			}
			var startId int64 = 0
			for {
				items, err := dao.QrySubscribers(appkey, userId, startId, 1000)
				if err != nil {
					break
				}
				for _, item := range items {
					rels.subscribers[item.Subscriber] = true
					if item.ID > startId {
						startId = item.ID
					}
				}
				if len(items) < 1000 {
					break
				}
			}
			subRelCache.Add(key, rels)
			return rels
		}
	}
}

func getKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}
