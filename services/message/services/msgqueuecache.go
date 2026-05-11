package services

import (
	"context"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"strings"
	"time"

	"github.com/bytedance/gopkg/collection/skipmap"
)

var (
	userMsgQueueCache *caches.LruCache
)

func init() {
	userMsgQueueCache = caches.NewLruCacheWithReadTimeout("user_offline_msg_queue_cache", 100000, nil, time.Hour)
}

type MessageQueueContainer struct {
	AppKey string
	UserId string
	MsgSet *skipmap.Int64Map
}

func (container *MessageQueueContainer) AppendMsg(ctx context.Context, msg *pbobjs.DownMsg) {
	key := getMessageQueueCacheKey(container.AppKey, container.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	container.MsgSet.Store(int64(msg.MsgTime), msg)
	msgLen := container.MsgSet.Len()
	msgQueueMaxCount := getMsgQueueMaxCount(container.AppKey)
	if msgLen > msgQueueMaxCount {
		evictCount := msgLen - msgQueueMaxCount
		delTimes := make([]int64, 0, evictCount)
		container.MsgSet.Range(func(key int64, value interface{}) bool {
			delTimes = append(delTimes, key)
			evictCount--
			return evictCount > 0
		})
		for _, t := range delTimes {
			container.MsgSet.Delete(t)
		}
	}
}

func (container *MessageQueueContainer) GetMsgsBaseTime(start int64, count int) []*pbobjs.DownMsg {
	if count <= 0 {
		return []*pbobjs.DownMsg{}
	}
	retMsgs := make([]*pbobjs.DownMsg, 0, count)
	container.MsgSet.Range(func(key int64, value interface{}) bool {
		if key <= start {
			return true
		}
		retMsgs = append(retMsgs, value.(*pbobjs.DownMsg))
		return len(retMsgs) < count
	})
	return retMsgs
}

func getMsgQueueMaxCount(appkey string) int {
	count := 100
	if appinfo, exist := commonservices.GetAppInfo(appkey); exist && appinfo != nil {
		count = appinfo.MsgQueueMaxCount
	}
	return count
}

func getMessageQueueCacheKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}

func GetMessageQueueContainer(appkey, userId string) *MessageQueueContainer {
	key := getMessageQueueCacheKey(appkey, userId)
	if obj, exist := userMsgQueueCache.Get(key); exist {
		return obj.(*MessageQueueContainer)
	} else {
		lock := userLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if obj, exist := userMsgQueueCache.Get(key); exist {
			return obj.(*MessageQueueContainer)
		} else {
			container := &MessageQueueContainer{
				AppKey: appkey,
				UserId: userId,
				MsgSet: skipmap.NewInt64(),
			}
			userMsgQueueCache.Add(key, container)
			return container
		}
	}
}
