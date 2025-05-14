package kvdbcommons

import (
	"encoding/base64"
	"im-server/commons/caches"
	"im-server/commons/tools"
	"time"
)

var tsdbCache *caches.LruCache

type TsdbCacheItem struct {
	cacheKey   string
	LatestTime int64
}

func (item *TsdbCacheItem) GetTimestamp(curr int64) int64 {
	l := kvdbLocks.GetLocks(item.cacheKey)
	l.Lock()
	defer l.Unlock()

	if curr > item.LatestTime {
		item.LatestTime = curr
	} else {
		item.LatestTime = item.LatestTime + 1
	}
	return item.LatestTime
}

func init() {
	tsdbCache = caches.NewLruCacheWithReadTimeout("tsdb_cache", 100000, nil, 30*time.Second)
}

func getTsdbCache(keyBs []byte) *TsdbCacheItem {
	key := base64.URLEncoding.EncodeToString(keyBs)
	if obj, exist := tsdbCache.Get(key); exist {
		return obj.(*TsdbCacheItem)
	} else {
		l := kvdbLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if obj, exist := tsdbCache.Get(key); exist {
			return obj.(*TsdbCacheItem)
		} else {
			item := &TsdbCacheItem{
				cacheKey:   key,
				LatestTime: 0,
			}
			tsdbCache.Add(key, item)
			return item
		}
	}
}

func TsAppend(key []byte, value []byte) (int64, error) {
	cacheItem := getTsdbCache(key)
	timestamp := cacheItem.GetTimestamp(time.Now().UnixMilli())
	keyBs := []byte{}
	keyBs = append(keyBs, key...)
	keyBs = append(keyBs, tools.Int64ToBytes(timestamp)...)
	return timestamp, Set(keyBs, value)
}

func TsAppendWithTime(key []byte, value []byte, timestamp int64) error {
	keyBs := []byte{}
	keyBs = append(keyBs, key...)
	keyBs = append(keyBs, tools.Int64ToBytes(timestamp)...)
	return Set(keyBs, value)
}

type TsItem struct {
	Key       []byte
	Timestamp int64
	Value     []byte
}

func TsScan(key []byte, startTime int64, count int) ([]TsItem, error) {
	start := []byte{}
	if startTime > 0 {
		start = tools.Int64ToBytes(startTime)
	}
	items, err := Scan(key, start, count)
	if err != nil {
		return []TsItem{}, err
	}
	ret := []TsItem{}
	for _, item := range items {
		var timestamp int64 = 0
		if len(item.Key) > 8 {
			bs := item.Key[len(item.Key)-8:]
			timestamp = tools.BytesToInt64(bs)
			item.Key = item.Key[:len(item.Key)-8]
		}
		ret = append(ret, TsItem{
			Key:       item.Key,
			Value:     item.Val,
			Timestamp: timestamp,
		})
	}
	return ret, nil
}
