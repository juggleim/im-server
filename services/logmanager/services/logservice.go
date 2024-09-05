package services

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"strings"
	"time"
)

var logCache *caches.LruCache
var logLocks *tools.SegmentatedLocks

type LogCacheItem struct {
	cacheKey   string
	LatestTime int64
}

func (log *LogCacheItem) GetLogTime(curr int64) int64 {
	l := logLocks.GetLocks(log.cacheKey)
	l.Lock()
	defer l.Unlock()

	if curr > log.LatestTime {
		log.LatestTime = curr
	} else {
		log.LatestTime = log.LatestTime + 1
	}
	return log.LatestTime
}

func init() {
	logCache = caches.NewLruCacheWithReadTimeout(100000, nil, 30*time.Second)
	logLocks = tools.NewSegmentatedLocks(256)
}

func getLogCache(key string) *LogCacheItem {
	if obj, exist := logCache.Get(key); exist {
		return obj.(*LogCacheItem)
	} else {
		l := logLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if obj, exist := logCache.Get(key); exist {
			return obj.(*LogCacheItem)
		} else {
			item := &LogCacheItem{
				cacheKey:   key,
				LatestTime: 0,
			}
			logCache.Add(key, item)
			return item
		}
	}
}

func WriteUserConnectLog(data *pbobjs.UserConnectLog) error {
	data.RealTime = data.Timestamp
	cacheKey := strings.Join([]string{data.AppKey, data.UserId}, "_")
	logCache := getLogCache(cacheKey)
	data.Timestamp = logCache.GetLogTime(data.Timestamp)

	key := strings.Join([]string{data.AppKey, data.UserId, tools.Int642String(data.Timestamp)}, "_")
	return writeLog(string(ServerLogType_UserConnect), key, tools.ToJson(data))
}

func QryUserConnectLogs(appkey, userId string, start, count int64) ([]LogEntity, error) {
	prefix := fmt.Sprintf("%s_%s_", appkey, userId)
	startKey := ""
	if start > 0 {
		startKey = fmt.Sprintf("%s_%s_%d", appkey, userId, start)
	}
	return qryLogs(string(ServerLogType_UserConnect), prefix, startKey, int(count))
}

func WriteConnectLog(data *pbobjs.ConnectionLog) error {
	data.RealTime = data.Timestamp
	cacheKey := strings.Join([]string{data.AppKey, data.Session}, "_")
	logCache := getLogCache(cacheKey)
	data.Timestamp = logCache.GetLogTime(data.Timestamp)

	key := strings.Join([]string{data.AppKey, data.Session, tools.Int642String(data.Timestamp)}, "_")
	return writeLog(string(ServerLogType_Connect), key, tools.ToJson(data))
}

func QryConnectLogs(appkey, session string, start, count int64) ([]LogEntity, error) {
	prefix := fmt.Sprintf("%s_%s_", appkey, session)
	startKey := ""
	if start > 0 {
		startKey = fmt.Sprintf("%s%d", prefix, start)
	}
	return qryLogs(string(ServerLogType_Connect), prefix, startKey, int(count))
}
