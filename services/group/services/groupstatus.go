package services

import (
	"im-server/commons/caches"
	"im-server/commons/tools"
	"strings"
	"time"
)

type GroupStatus struct {
	LatestMsgTime int64
}

var groupStatusCache *caches.LruCache
var groupStatusLocks *tools.SegmentatedLocks

func init() {
	groupStatusCache = caches.NewLruCacheWithReadTimeout("groupstatus_cache", 100000, func(key, value interface{}) {}, 10*time.Minute)
	groupStatusLocks = tools.NewSegmentatedLocks(512)
}

func GetGroupStatus(appkey, groupId string) *GroupStatus {
	key := strings.Join([]string{appkey, groupId}, "_")
	if val, exist := groupStatusCache.Get(key); exist {
		return val.(*GroupStatus)
	} else {
		l := groupStatusLocks.GetLocks(appkey, groupId)
		l.Lock()
		defer l.Unlock()
		if val, exist := groupStatusCache.Get(key); exist {
			return val.(*GroupStatus)
		} else {
			groupStatus := &GroupStatus{}
			groupStatusCache.Add(key, groupStatus)
			return groupStatus
		}
	}
}

func GetGroupSendTime(appkey, groupId string) int64 {
	group := GetGroupStatus(appkey, groupId)
	currentTime := time.Now().UnixMilli()

	l := groupStatusLocks.GetLocks(appkey, groupId)
	l.Lock()
	defer l.Unlock()

	ret := currentTime
	if currentTime > group.LatestMsgTime {
		group.LatestMsgTime = currentTime
	} else {
		ret = group.LatestMsgTime + 1
		group.LatestMsgTime = ret
	}
	return ret
}
