package services

import (
	"im-server/commons/caches"
	"im-server/services/commonservices/dbs"
	"time"
)

var userActivityCache *caches.LruCache
var userActivityNotExist *UserActivity

func init() {
	userActivityCache = caches.NewLruCacheWithAddReadTimeout("useractivity_cache", 200000, func(key, value interface{}) {}, time.Hour, time.Hour)
	userActivityNotExist = &UserActivity{}
}

type UserActivity struct {
	appkey string
	userId string

	LatestActivityTime int64
}

func GetLatestUserActivity(appkey, userId string) *UserActivity {
	key := getKey(appkey, userId)
	if val, exist := userActivityCache.Get(key); exist {
		activity := val.(*UserActivity)
		if activity == userActivityNotExist {
			return nil
		} else {
			return activity
		}
	} else {
		l := userLocks.GetLocks(appkey, userId)
		l.Lock()
		defer l.Unlock()
		if val, exist := userActivityCache.Get(key); exist {
			activity := val.(*UserActivity)
			if activity == userActivityNotExist {
				return nil
			} else {
				return activity
			}
		} else {
			activity := initUserActivity(appkey, userId)
			userActivityCache.Add(key, activity)
			return activity
		}
	}
}

func initUserActivity(appkey, userId string) *UserActivity {
	dao := dbs.UserActivityDao{}
	dbActivity, err := dao.QryLatestUserActivity(appkey, userId)
	if err != nil || dbActivity == nil {
		return userActivityNotExist
	} else {
		return &UserActivity{
			appkey: appkey,
			userId: userId,

			LatestActivityTime: dbActivity.TimeMark,
		}
	}
}
