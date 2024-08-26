package bases

import (
	"context"
	"im-server/commons/caches"
	"im-server/commons/tools"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

var (
	limiterCache *caches.LruCache
	limiterLocks *tools.SegmentatedLocks
)

func init() {
	limiterCache = caches.NewLruCacheWithReadTimeout(1000000, nil, time.Minute)
	limiterLocks = tools.NewSegmentatedLocks(512)
}

func CheckLimiter(ctx context.Context) bool {
	limiter := GetLimiter(ctx)
	return limiter.Allow()
}

func GetLimiter(ctx context.Context) *rate.Limiter {
	appkey := GetAppKeyFromCtx(ctx)
	targetId := GetTargetIdFromCtx(ctx)
	scene := GetMethodFromCtx(ctx)
	key := strings.Join([]string{scene, appkey, targetId}, "_")
	if limiterObj, exist := limiterCache.Get(key); exist {
		return limiterObj.(*rate.Limiter)
	} else {
		lock := limiterLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if limiterObj, exist := limiterCache.Get(key); exist {
			return limiterObj.(*rate.Limiter)
		} else {
			limiter := rate.NewLimiter(100, 10)
			limiterCache.Add(key, limiter)
			return limiter
		}
	}
}
