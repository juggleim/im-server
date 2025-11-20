package services

import (
	"im-server/commons/caches"
	"im-server/commons/tools"
	"time"
)

var friendRelCache *caches.LruCache
var friendRelLocks *tools.SegmentatedLocks

func init() {
	friendRelCache = caches.NewLruCacheWithAddReadTimeout("friend_cache", 10000, nil, 10*time.Minute, 10*time.Minute)
	friendRelLocks = tools.NewSegmentatedLocks(128)
}
