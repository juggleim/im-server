package services

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/services/conversation/storages"
	"im-server/services/conversation/storages/models"
	"time"
)

var offlineUserConverCache *caches.EphemeralCache

func init() {
	offlineUserConverCache = caches.NewEphemeralCache(time.Second, 5*time.Second, func(key, value interface{}) {
		conver, ok := value.(*ConversationCacheItem)
		if ok && conver != nil {
			converStorage := storages.NewConversationStorage()
			converStorage.UpsertConversation(models.Conversation{
				UserId:               conver.UserId,
				TargetId:             conver.TargetId,
				ChannelType:          conver.ChannelType,
				SortTime:             conver.SortTime,
				LatestMsgId:          conver.LatestMsgId,
				LatestUnreadMsgIndex: conver.UnReadIndex,
				SyncTime:             conver.SyncTime,
				AppKey:               conver.Appkey,
			})
		}
	})
}

func UpsertOfflineConversation(item *ConversationCacheItem) {
	key := fmt.Sprintf("%s_%s_%s_%d", item.Appkey, item.UserId, item.TargetId, item.ChannelType)
	offlineUserConverCache.Upsert(key, func(oldVal interface{}) interface{} {
		var converItem *ConversationCacheItem
		if oldVal != nil {
			converItem = oldVal.(*ConversationCacheItem)
			converItem.LatestMsgId = item.LatestMsgId
			if item.SortTime > converItem.SortTime {
				converItem.SortTime = item.SortTime
			}
			if item.UnReadIndex > converItem.UnReadIndex {
				converItem.UnReadIndex = item.UnReadIndex
			}
			if item.SyncTime > converItem.SyncTime {
				converItem.SyncTime = item.SyncTime
			}
		} else {
			converItem = item
		}
		return converItem
	})
}
