package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/tools"
	"strings"
	"time"
)

var memberCache *caches.LruCache
var memberLocks *tools.SegmentatedLocks

func init() {
	memberCache = caches.NewLruCacheWithReadTimeout(100000, nil, 10*time.Minute)
	memberLocks = tools.NewSegmentatedLocks(256)
}

type MemberStatus struct {
	AppKey   string
	MemberId string

	StopDispatch bool
}

func getMemberKey(appkey, memberId string) string {
	return strings.Join([]string{appkey, memberId}, "_")
}
func GetMemberStatus(ctx context.Context, memberId string) *MemberStatus {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getMemberKey(appkey, memberId)
	if obj, exist := memberCache.Get(key); exist {
		return obj.(*MemberStatus)
	} else {
		lock := memberLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if obj, exist := memberCache.Get(key); exist {
			return obj.(*MemberStatus)
		} else {
			member := &MemberStatus{
				AppKey:   appkey,
				MemberId: memberId,
			}
			memberCache.Add(key, member)
			return member
		}
	}
}
