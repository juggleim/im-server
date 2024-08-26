package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
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

	bases.SetPreProcess(func(ctx context.Context, sender actorsystem.ActorRef) bool {
		limiter := GetLimiter(ctx)
		if limiter != nil {
			if limiter.Allow() {
				return true
			} else {
				rpcType := bases.GetRpcTypeFromCtx(ctx)
				if rpcType == pbobjs.RpcMsgType_UserPub {
					sender.Tell(bases.CreateUserPubAckWraper(ctx, errs.IMErrorCode_CONNECT_EXCEEDLIMITED, "", 0, 0), actorsystem.NoSender)
				} else if rpcType == pbobjs.RpcMsgType_QueryMsg {
					sender.Tell(bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_CONNECT_EXCEEDLIMITED, nil), actorsystem.NoSender)
				}
				return false
			}
		} else {
			return true
		}
	})
}

func CheckLimiter(ctx context.Context) bool {
	limiter := GetLimiter(ctx)
	return limiter.Allow()
}

func GetLimiter(ctx context.Context) *rate.Limiter {
	appkey := bases.GetAppKeyFromCtx(ctx)
	targetId := bases.GetTargetIdFromCtx(ctx)
	scene := bases.GetMethodFromCtx(ctx)
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
