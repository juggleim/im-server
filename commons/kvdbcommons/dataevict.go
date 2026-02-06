package kvdbcommons

import (
	"im-server/commons/kvdbcommons/kvobjs"
	"im-server/commons/logs"
	"im-server/commons/tools"
	"time"
)

var evictTicker *time.Ticker
var evictQueuePrefix []byte = []byte("kvdb_evict_queue")
var metaPrefix []byte = []byte("kvdb_meta_")
var evictTaskPools *tools.SinglePools

func init() {
	evictTaskPools = tools.NewSinglePools(32, true)
}

func startEvictTask() {
	if evictTicker != nil {
		evictTicker.Stop()
	}
	evictTicker = time.NewTicker(10 * time.Second)
	go func() {
		for task := range evictTicker.C {
			current := time.Now().UnixMilli()
			if current-task.UnixMilli() > 500 {
				continue
			}
			cleanTimeoutData(current)
		}
	}()
}

func cleanTimeoutData(timeLine int64) {
	items, err := Scan(evictQueuePrefix, []byte{}, 1000)
	if err == nil {
		for _, item := range items {
			if len(item.Key) > 8 {
				bs := item.Key[len(item.Key)-8:]
				timestamp := tools.BytesToInt64(bs)
				if timestamp < timeLine {
					evictItem := &kvobjs.KvEvictItem{}
					err = tools.PbUnMarshal(item.Val, evictItem)
					if err == nil {
						for _, key := range evictItem.Keys {
							//check expired & delete
							metaKey := []byte{}
							metaKey = append(metaKey, metaPrefix...)
							metaKey = append(metaKey, key...)
							metaVal, err := Get(metaKey)
							if err == nil {
								meta := &kvobjs.KvMeta{}
								err = tools.PbUnMarshal(metaVal, meta)
								if err == nil {
									if meta.ExpiredAt < timeLine {
										Delete(key)
										Delete(metaKey)
									}
								}
							}
						}
					}
					Delete(item.Key)
				} else {
					break
				}
			}
		}
	}
}

func appendEvictQueue(key []byte, expiredAt int64) {
	evictKey := []byte{}
	evictKey = append(evictKey, evictQueuePrefix...)
	evictKey = append(evictKey, tools.Int64ToBytes(expiredAt)...)
	evictTaskPools.GetPool(Bytes2SafeKey(evictKey)).Submit(func() {
		isExist, err := Exist(evictKey)
		if err == nil {
			if isExist {
				val, err := Get(evictKey)
				if err == nil {
					var evictItem kvobjs.KvEvictItem
					err = tools.PbUnMarshal(val, &evictItem)
					if err == nil {
						evictItem.Keys = append(evictItem.Keys, key)
						newVal, _ := tools.PbMarshal(&evictItem)
						err = Set(evictKey, newVal)
						if err != nil {
							logs.Errorf("[kvdb]failed store evict item err:%v", err)
						}
					}
				}
			} else {
				evictItem := &kvobjs.KvEvictItem{
					Keys: [][]byte{},
				}
				evictItem.Keys = append(evictItem.Keys, key)
				newVal, _ := tools.PbMarshal(evictItem)
				err = Set(evictKey, newVal)
				if err != nil {
					logs.Errorf("[kvdb]failed store evict item err:%v", err)
				}
			}
		}
	})
}

func ExpireAt(key []byte, expiredAt int64) {
	//record meta
	metaKey := []byte{}
	metaKey = append(metaKey, metaPrefix...)
	metaKey = append(metaKey, key...)
	metaBs, _ := tools.PbMarshal(&kvobjs.KvMeta{
		ExpiredAt: expiredAt,
	})
	Set(metaKey, metaBs)
	//record to evict queue
	appendEvictQueue(key, expiredAt)
}

func Expire(key []byte, duration time.Duration) {
	expiredAt := time.Now().UnixMilli() + int64(duration/1000/1000)
	ExpireAt(key, expiredAt)
}

func Ttl(key []byte) int64 {
	metaKey := []byte{}
	metaKey = append(metaKey, metaPrefix...)
	metaKey = append(metaKey, key...)
	val, err := Get(metaKey)
	if err == nil {
		var kvMeta kvobjs.KvMeta
		err = tools.PbUnMarshal(val, &kvMeta)
		if err == nil {
			result := kvMeta.ExpiredAt - time.Now().UnixMilli()
			if result > 0 {
				return result
			}
		}
	}
	return 0
}
