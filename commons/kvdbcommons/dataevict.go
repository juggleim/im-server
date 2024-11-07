package kvdbcommons

import (
	"fmt"
	"im-server/commons/tools"
	"sync"
	"time"
)

var evictTicker *time.Ticker
var evictLock *sync.RWMutex
var evictQueuePrefix []byte = []byte("kvdb_evict_queue")
var expirePrefix []byte = []byte("kvdb_expire_prefix")

func init() {
	evictLock = &sync.RWMutex{}
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
				timestamp = timestamp >> 10
				fmt.Println("clean:", timestamp)
				if timestamp < timeLine {
					dataKey := item.Val

					Delete(item.Key)
					//delete expire record
					expireKey := []byte{}
					expireKey = append(expireKey, expirePrefix...)
					expireKey = append(expireKey, dataKey...)
					val, err := Get(expireKey)
					if err == nil {
						expireAt := tools.BytesToInt64(val)
						if expireAt < timeLine {
							Delete(expireKey)
							//delete data
							Delete(dataKey)
						}
					}
				} else {
					break
				}
			}
		}
	}
}

func ExpireAt(key []byte, expiredAt int64) {
	evictLock.Lock()
	defer evictLock.Unlock()

	retryCount := 3
	var evictKey []byte
	for retryCount > 0 {
		randNum := tools.RandInt(1024)
		evictKey = []byte{}
		evictKey = append(evictKey, evictQueuePrefix...)
		evictKey = append(evictKey, tools.Int64ToBytes(expiredAt<<10+int64(randNum))...)
		exist, err := Exist(evictKey)
		if err != nil {
			return
		}
		if !exist {
			break
		}
		retryCount--
	}
	Set(evictKey, key)
	expireKey := []byte{}
	expireKey = append(expireKey, expirePrefix...)
	expireKey = append(expireKey, key...)
	Set(expireKey, tools.Int64ToBytes(expiredAt))
}

func Expire(key []byte, duration time.Duration) {
	curr := time.Now().UnixMilli()
	curr = curr + int64(duration/1000/1000)
	ExpireAt(key, curr)
}

func Ttl(key []byte) int64 {
	expireKey := []byte{}
	expireKey = append(expireKey, expirePrefix...)
	expireKey = append(expireKey, key...)
	val, err := Get(expireKey)
	if err != nil {
		return 0
	}
	fmt.Println("expiredAt:", tools.BytesToInt64(val))
	result := tools.BytesToInt64(val) - time.Now().UnixMilli()
	return result
}
