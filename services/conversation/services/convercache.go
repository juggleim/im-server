package services

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/conversation/storages"
	"im-server/services/conversation/storages/models"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bytedance/gopkg/collection/zset"
)

var userConverCache *caches.LruCache
var userLocks *tools.SegmentatedLocks

var persistCache *caches.EphemeralCache

type ConverPersistIndex struct {
	Appkey      string
	UserId      string
	TargetId    string
	ChannelType pbobjs.ChannelType
}

func init() {
	userConverCache = caches.NewLruCacheWithReadTimeout("userconver_cache", 100000, nil, time.Hour)
	userLocks = tools.NewSegmentatedLocks(512)

	persistCache = caches.NewEphemeralCache(time.Second, 5*time.Second, func(key, value interface{}) {
		persistIndex, ok := value.(*ConverPersistIndex)
		if ok && persistIndex != nil {
			userConvers := getUserConvers(persistIndex.Appkey, persistIndex.UserId)
			item := userConvers.QryConver(persistIndex.TargetId, persistIndex.ChannelType)
			if item != nil {
				conversation := models.Conversation{
					AppKey:      item.AppKey,
					UserId:      item.UserId,
					TargetId:    item.TargetId,
					ChannelType: item.ChannelType,

					SortTime: item.SortTime,
					SyncTime: item.SyncTime,

					LatestMsgId: item.LatestMsgId,
					// LatestMsg:            item.LatestMsg,
					LatestUnreadMsgIndex: item.LatestUnreadMsgIndex,

					LatestReadMsgIndex: item.LatestReadMsgIndex,
					LatestReadMsgId:    item.LatestReadMsgId,
					LatestReadMsgTime:  item.LatestReadMsgTime,

					IsTop:          item.IsTop,
					TopUpdatedTime: item.TopUpdatedTime,
					UndisturbType:  item.UndisturbType,

					UnreadTag:  item.UnreadTag,
					ConverExts: item.ConverExts,
					IsDeleted:  item.IsDeleted,
				}
				storage := storages.NewConversationStorage()
				err := storage.Upsert(conversation)
				if err != nil {
					fmt.Println("save conver err:", err)
				}
			}
		}
	})
}

type UserConversations struct {
	Appkey   string
	UserId   string
	purgeing int32

	ConverItemMap map[string]*models.Conversation
	SyncTimeIndex *zset.Float64Set
	SortTimeIndex *zset.Float64Set

	TopIndexBaseTopTime  *zset.Float64Set
	TopIndexBaseSortTime *zset.Float64Set
}

func (uc *UserConversations) UpsertCovner(conver models.Conversation) {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	itemKey := getConverItemKey(conver.TargetId, conver.ChannelType)

	if cacheConver, exist := uc.ConverItemMap[itemKey]; exist {
		cacheConver.LatestMsgId = conver.LatestMsgId
		if conver.LatestUnreadMsgIndex > cacheConver.LatestUnreadMsgIndex {
			cacheConver.LatestUnreadMsgIndex = conver.LatestUnreadMsgIndex
		}
		if conver.SortTime > cacheConver.SortTime {
			cacheConver.SortTime = conver.SortTime
			uc.SortTimeIndex.Add(float64(conver.SortTime), itemKey)
		}
		if conver.SyncTime > cacheConver.SyncTime {
			cacheConver.SyncTime = conver.SyncTime
			uc.SyncTimeIndex.Add(float64(conver.SyncTime), itemKey)
		}
		cacheConver.IsDeleted = 0
	} else {
		item := &models.Conversation{
			AppKey:               uc.Appkey,
			UserId:               uc.UserId,
			TargetId:             conver.TargetId,
			ChannelType:          conver.ChannelType,
			LatestMsgId:          conver.LatestMsgId,
			LatestUnreadMsgIndex: conver.LatestUnreadMsgIndex,
			SortTime:             conver.SortTime,
			SyncTime:             conver.SyncTime,
		}
		uc.ConverItemMap[itemKey] = item
		uc.SyncTimeIndex.Add(float64(item.SyncTime), itemKey)
		uc.SortTimeIndex.Add(float64(item.SortTime), itemKey)
		evictCount := len(uc.ConverItemMap) - 10000
		if evictCount > 100 {
			go uc.purge()
		}
	}
}

func (uc *UserConversations) AppendMention(targetId string, channelType pbobjs.ChannelType, mentionMsg *models.MentionMsg) {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	itemKey := getConverItemKey(targetId, channelType)
	if cacheConver, exist := uc.ConverItemMap[itemKey]; exist {
		if cacheConver.MentionInfo != nil && mentionMsg == nil {
			return
		}
		if cacheConver.MentionInfo == nil {
			mentionInfo := &models.ConverMentionInfo{
				MentionMsgs: []*models.MentionMsg{},
				SenderIds:   []string{},
			}
			//read from db
			storage := storages.NewMentionMsgStorage()
			mentionMsgs, err := storage.QryMentionSenderIdsBaseIndex(uc.Appkey, uc.UserId, targetId, channelType, cacheConver.LatestReadMsgIndex, 100)
			if err == nil {
				mentionInfo.MentionMsgs = append(mentionInfo.MentionMsgs, mentionMsgs...)
			}
			cacheConver.MentionInfo = mentionInfo
		}
		if mentionMsg != nil {
			//add new mention msg
			cacheConver.MentionInfo.MentionMsgs = append(cacheConver.MentionInfo.MentionMsgs, &models.MentionMsg{
				SenderId:    mentionMsg.SenderId,
				MsgId:       mentionMsg.MsgId,
				MsgTime:     mentionMsg.MsgTime,
				MsgIndex:    mentionMsg.MsgIndex,
				MentionType: mentionMsg.MentionType,
			})
			length := len(cacheConver.MentionInfo.MentionMsgs)
			if length > 100 {
				cacheConver.MentionInfo.MentionMsgs = cacheConver.MentionInfo.MentionMsgs[length-100:]
			}
		}
		//calculate mention count
		cacheConver.MentionInfo.MentionMsgCount = len(cacheConver.MentionInfo.MentionMsgs)
		cacheConver.MentionInfo.IsMentioned = cacheConver.MentionInfo.MentionMsgCount > 0
		//generate sender info
		cacheConver.MentionInfo.SenderIds = []string{}
		tmpMap := map[string]int{}
		for _, mMsg := range cacheConver.MentionInfo.MentionMsgs {
			if _, exist := tmpMap[mMsg.SenderId]; !exist {
				tmpMap[mMsg.SenderId] = 1
				cacheConver.MentionInfo.SenderIds = append(cacheConver.MentionInfo.SenderIds, mMsg.SenderId)
			}
		}
	}
}

func (uc *UserConversations) GetMentionInfo(targetId string, channelType pbobjs.ChannelType) *models.ConverMentionInfo {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()

	itemKey := getConverItemKey(targetId, channelType)
	if cacheConver, exist := uc.ConverItemMap[itemKey]; exist {
		if cacheConver.MentionInfo != nil {
			retMentionInfo := &models.ConverMentionInfo{
				IsMentioned:     cacheConver.MentionInfo.IsMentioned,
				MentionMsgCount: cacheConver.MentionInfo.MentionMsgCount,
				SenderIds:       []string{},
				MentionMsgs:     []*models.MentionMsg{},
			}
			retMentionInfo.SenderIds = append(retMentionInfo.SenderIds, cacheConver.MentionInfo.SenderIds...)
			retMentionInfo.MentionMsgs = append(retMentionInfo.MentionMsgs, cacheConver.MentionInfo.MentionMsgs...)
			return retMentionInfo
		}
		return nil
	}
	return &models.ConverMentionInfo{
		SenderIds:   []string{},
		MentionMsgs: []*models.MentionMsg{},
	}
}

func (uc *UserConversations) innerClearMentionMsgs(item *models.Conversation, msgIndex int64) {
	if item.MentionInfo != nil {
		mentionMsgs := []*models.MentionMsg{}
		for _, mentionMsg := range item.MentionInfo.MentionMsgs {
			if mentionMsg.MsgIndex > msgIndex {
				mentionMsgs = append(mentionMsgs, mentionMsg)
			}
		}
		newMentionCount := len(mentionMsgs)
		if newMentionCount < item.MentionInfo.MentionMsgCount {
			item.MentionInfo.MentionMsgCount = newMentionCount
			item.MentionInfo.IsMentioned = newMentionCount > 0
			item.MentionInfo.MentionMsgs = mentionMsgs
			item.MentionInfo.SenderIds = []string{}
			//generate sender info
			tmpMap := map[string]int{}
			for _, mMsg := range item.MentionInfo.MentionMsgs {
				if _, exist := tmpMap[mMsg.SenderId]; !exist {
					tmpMap[mMsg.SenderId] = 1
					item.MentionInfo.SenderIds = append(item.MentionInfo.SenderIds, mMsg.SenderId)
				}
			}
			go func() {
				//clear readed mention msgs
				storage := storages.NewMentionMsgStorage()
				storage.CleanMentionMsgsBaseIndex(uc.Appkey, uc.UserId, item.TargetId, item.ChannelType, msgIndex)
			}()
		}
	}
}

func (uc *UserConversations) purge() {
	if atomic.CompareAndSwapInt32(&uc.purgeing, 0, 1) {
		key := getUserConverCacheKey(uc.Appkey, uc.UserId)
		lock := userLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		evictCount := len(uc.ConverItemMap) - 10000
		if evictCount > 100 {
			list := uc.SortTimeIndex.RangeByScore(0, float64(time.Now().UnixMilli()))
			for i, item := range list {
				uc.SortTimeIndex.Remove(item.Value)
				uc.SyncTimeIndex.Remove(item.Value)
				delete(uc.ConverItemMap, item.Value)
				if i >= evictCount {
					break
				}
			}
		}
		atomic.CompareAndSwapInt32(&uc.purgeing, 1, 0)
	}
}

func (uc *UserConversations) PersistConver(targetId string, channelType pbobjs.ChannelType) {
	key := fmt.Sprintf("%s_%s_%s_%d", uc.Appkey, uc.UserId, targetId, channelType)
	persistCache.Add(key, &ConverPersistIndex{
		Appkey:      uc.Appkey,
		UserId:      uc.UserId,
		TargetId:    targetId,
		ChannelType: channelType,
	})
}

func (uc *UserConversations) SyncConvers(startTime int64, count int32) []*models.Conversation {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()

	nodes := uc.SyncTimeIndex.RangeByScore(float64(startTime+1), math.MaxInt64)
	resp := []*models.Conversation{}
	var index int32 = 0
	for _, node := range nodes {
		itemKey := node.Value
		if conver, exist := uc.ConverItemMap[itemKey]; exist {
			resp = append(resp, &models.Conversation{
				UserId:      conver.UserId,
				TargetId:    conver.TargetId,
				ChannelType: conver.ChannelType,
				SortTime:    conver.SortTime,
				SyncTime:    conver.SyncTime,

				LatestMsgId:          conver.LatestMsgId,
				LatestUnreadMsgIndex: conver.LatestUnreadMsgIndex,

				LatestReadMsgIndex: conver.LatestReadMsgIndex,
				LatestReadMsgId:    conver.LatestReadMsgId,
				LatestReadMsgTime:  conver.LatestReadMsgTime,

				IsTop:          conver.IsTop,
				TopUpdatedTime: conver.TopUpdatedTime,
				UndisturbType:  conver.UndisturbType,
				IsDeleted:      conver.IsDeleted,
				UnreadTag:      conver.UnreadTag,
				ConverExts:     conver.ConverExts,
			})
			index++
		}
		if index >= count {
			break
		}
	}
	return resp
}

func (uc *UserConversations) QryConver(targetId string, channelType pbobjs.ChannelType) *models.Conversation {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	itemKey := getConverItemKey(targetId, channelType)
	if conver, exist := uc.ConverItemMap[itemKey]; exist {
		ret := &models.Conversation{
			AppKey:      conver.AppKey,
			UserId:      conver.UserId,
			TargetId:    targetId,
			ChannelType: channelType,
			SortTime:    conver.SortTime,
			SyncTime:    conver.SyncTime,

			LatestMsgId:          conver.LatestMsgId,
			LatestUnreadMsgIndex: conver.LatestUnreadMsgIndex,

			LatestReadMsgIndex: conver.LatestReadMsgIndex,
			LatestReadMsgId:    conver.LatestReadMsgId,
			LatestReadMsgTime:  conver.LatestReadMsgTime,

			IsTop:          conver.IsTop,
			TopUpdatedTime: conver.TopUpdatedTime,
			UndisturbType:  conver.UndisturbType,
			IsDeleted:      conver.IsDeleted,
			UnreadTag:      conver.UnreadTag,
			ConverExts:     conver.ConverExts,
		}
		return ret
	}
	return nil
}

func (uc *UserConversations) QryConvers(startTime int64, count int32, isPositive bool, targetId string, channelType pbobjs.ChannelType, tag string) []*models.Conversation {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	var nodes []zset.Float64Node
	if isPositive {
		nodes = uc.SortTimeIndex.RangeByScore(float64(startTime+1), math.MaxFloat64)
	} else {
		nodes = uc.SortTimeIndex.RevRangeByScore(float64(startTime-1), 0)
	}
	resp := []*models.Conversation{}
	var index int32 = 0
	for _, node := range nodes {
		itemKey := node.Value
		if conver, exist := uc.ConverItemMap[itemKey]; exist {
			if channelType != pbobjs.ChannelType_Unknown && conver.ChannelType != channelType {
				continue
			}
			if tag != "" && (conver.ConverExts == nil || len(conver.ConverExts.ConverTags) <= 0 || !conver.ConverExts.ConverTags[tag]) {
				continue
			}

			resp = append(resp, &models.Conversation{
				AppKey:      conver.AppKey,
				UserId:      conver.UserId,
				TargetId:    conver.TargetId,
				ChannelType: conver.ChannelType,
				SortTime:    conver.SortTime,
				SyncTime:    conver.SyncTime,

				LatestMsgId:          conver.LatestMsgId,
				LatestUnreadMsgIndex: conver.LatestUnreadMsgIndex,

				LatestReadMsgIndex: conver.LatestReadMsgIndex,
				LatestReadMsgId:    conver.LatestReadMsgId,
				LatestReadMsgTime:  conver.LatestReadMsgTime,

				IsTop:          conver.IsTop,
				TopUpdatedTime: conver.TopUpdatedTime,
				UndisturbType:  conver.UndisturbType,
				IsDeleted:      conver.IsDeleted,
				UnreadTag:      conver.UnreadTag,

				ConverExts: conver.ConverExts,
			})
			index++
		}
		if index >= count {
			break
		}
	}
	return resp
}

func (uc *UserConversations) QryTopConvers(startTime int64, count int32, sortType pbobjs.TopConverSortType, isPositive bool) []*models.Conversation {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	var nodes []zset.Float64Node
	var topIndex *zset.Float64Set
	if sortType == pbobjs.TopConverSortType_BySortTime {
		topIndex = uc.TopIndexBaseSortTime
	} else {
		topIndex = uc.TopIndexBaseTopTime
	}
	if isPositive {
		nodes = topIndex.RangeByScore(float64(startTime+1), math.MaxFloat64)
	} else {
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
		nodes = topIndex.RevRangeByScore(float64(startTime-1), 0)
	}
	resp := []*models.Conversation{}
	var index int32 = 0
	for _, node := range nodes {
		itemKey := node.Value
		if conver, exist := uc.ConverItemMap[itemKey]; exist {
			resp = append(resp, &models.Conversation{
				UserId:      conver.UserId,
				TargetId:    conver.TargetId,
				ChannelType: conver.ChannelType,
				SortTime:    conver.SortTime,
				SyncTime:    conver.SyncTime,

				LatestMsgId:          conver.LatestMsgId,
				LatestUnreadMsgIndex: conver.LatestUnreadMsgIndex,

				LatestReadMsgIndex: conver.LatestReadMsgIndex,
				LatestReadMsgId:    conver.LatestReadMsgId,
				LatestReadMsgTime:  conver.LatestReadMsgTime,

				IsTop:          conver.IsTop,
				TopUpdatedTime: conver.TopUpdatedTime,
				UndisturbType:  conver.UndisturbType,
				IsDeleted:      conver.IsDeleted,
				UnreadTag:      conver.UnreadTag,
			})
			index++
		}
		if index >= count {
			break
		}
	}
	return resp
}

func (uc *UserConversations) ClearTotalUnread() {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	for _, conver := range uc.ConverItemMap {
		conver.LatestReadMsgIndex = conver.LatestUnreadMsgIndex
		conver.UnreadTag = 0
		conver.MentionInfo = &models.ConverMentionInfo{
			SenderIds:   []string{},
			MentionMsgs: []*models.MentionMsg{},
		}
	}
}

func (uc *UserConversations) TotalUnreadCount() int64 {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	var count int64 = 0
	for _, conver := range uc.ConverItemMap {
		if conver.IsDeleted == 0 {
			c := conver.LatestUnreadMsgIndex - conver.LatestReadMsgIndex
			if c == 0 && conver.UnreadTag == 1 {
				count = count + 1
			} else {
				count = count + c
			}
		}
	}
	return count
}

func (uc *UserConversations) ClearUnread(targetId string, channelType pbobjs.ChannelType, readMsgIndex int64, readMsgId string, readMsgTime int64) bool {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	itemKey := getConverItemKey(targetId, channelType)
	if item, exist := uc.ConverItemMap[itemKey]; exist {
		if readMsgIndex > item.LatestReadMsgIndex && readMsgIndex <= item.LatestUnreadMsgIndex {
			item.LatestReadMsgIndex = readMsgIndex
			item.LatestReadMsgId = readMsgId
			item.LatestReadMsgTime = readMsgTime
			item.UnreadTag = 0
			uc.innerClearMentionMsgs(item, readMsgIndex)
			return true
		} else if item.UnreadTag > 0 {
			item.UnreadTag = 0
			return true
		}
	}
	return false
}

func (uc *UserConversations) DefaultClearUnread(targetId string, channelType pbobjs.ChannelType) bool {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	itemKey := getConverItemKey(targetId, channelType)
	if item, exist := uc.ConverItemMap[itemKey]; exist {
		if item.LatestReadMsgIndex < item.LatestUnreadMsgIndex {
			item.LatestReadMsgIndex = item.LatestUnreadMsgIndex
			item.LatestReadMsgId = ""
			item.LatestReadMsgTime = time.Now().UnixMilli()
			item.UnreadTag = 0
			uc.innerClearMentionMsgs(item, item.LatestUnreadMsgIndex)
			return true
		} else if item.UnreadTag > 0 {
			item.UnreadTag = 0
			return true
		}
	}
	return false
}

func (uc *UserConversations) DelConversation(targetId string, channelType pbobjs.ChannelType) bool {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	itemKey := getConverItemKey(targetId, channelType)
	if item, exist := uc.ConverItemMap[itemKey]; exist {
		if item.IsDeleted == 0 {
			item.IsDeleted = 1
			uc.SortTimeIndex.Remove(itemKey)
			uc.innerClearMentionMsgs(item, item.LatestUnreadMsgIndex)
			return true
		}
	}
	return false
}

func (uc *UserConversations) UpdateUnreadTag(targetId string, channelType pbobjs.ChannelType, unreadTag int) bool {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	itemKey := getConverItemKey(targetId, channelType)
	if item, exist := uc.ConverItemMap[itemKey]; exist {
		if item.UnreadTag != unreadTag {
			item.UnreadTag = unreadTag
			return true
		}
	}
	return false
}

func (uc *UserConversations) UpdTopState(targetId string, channelType pbobjs.ChannelType, isTop int, topTime int64) bool {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	itemKey := getConverItemKey(targetId, channelType)
	if item, exist := uc.ConverItemMap[itemKey]; exist {
		if item.IsTop != isTop {
			item.IsTop = isTop
			item.TopUpdatedTime = topTime
			if isTop == 0 {
				uc.TopIndexBaseSortTime.Remove(itemKey)
				uc.TopIndexBaseTopTime.Remove(itemKey)
			} else {
				uc.TopIndexBaseSortTime.Add(float64(item.SortTime), itemKey)
				uc.TopIndexBaseTopTime.Add(float64(topTime), itemKey)
			}
			return true
		}
	}
	return false
}

func (uc *UserConversations) UpdateUndisturbType(targetId string, channelType pbobjs.ChannelType, undisturbType int32) bool {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	itemKey := getConverItemKey(targetId, channelType)
	if item, exist := uc.ConverItemMap[itemKey]; exist {
		if item.UndisturbType != undisturbType {
			item.UndisturbType = undisturbType
			return true
		}
	}
	return false
}

func (uc *UserConversations) TagAddConvers(tag string, convers []*pbobjs.SimpleConversation) bool {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	ret := false
	for _, conver := range convers {
		itemKey := getConverItemKey(conver.TargetId, conver.ChannelType)
		if item, exist := uc.ConverItemMap[itemKey]; exist {
			if item.ConverExts == nil {
				item.ConverExts = &pbobjs.ConverExts{
					ConverTags: make(map[string]bool),
				}
			}
			if _, exist := item.ConverExts.ConverTags[tag]; !exist {
				item.ConverExts.ConverTags[tag] = true
				ret = ret || true
			}
		}
	}
	return ret
}

func (uc *UserConversations) TagDelConvers(tag string, convers []*pbobjs.SimpleConversation) bool {
	key := getUserConverCacheKey(uc.Appkey, uc.UserId)
	lock := userLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	ret := false
	for _, conver := range convers {
		itemKey := getConverItemKey(conver.TargetId, conver.ChannelType)
		if item, exist := uc.ConverItemMap[itemKey]; exist && item.ConverExts != nil {
			if _, exist := item.ConverExts.ConverTags[tag]; exist {
				delete(item.ConverExts.ConverTags, tag)
				ret = ret || true
			}
		}
	}
	return ret
}

func getUserConverCacheKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}

func getConverItemKey(targetId string, channelType pbobjs.ChannelType) string {
	return fmt.Sprintf("%s_%d", targetId, channelType)
}

func getUserConvers(appkey, userId string) *UserConversations {
	key := getUserConverCacheKey(appkey, userId)
	if obj, exist := userConverCache.Get(key); exist {
		return obj.(*UserConversations)
	} else {
		l := userLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()

		if obj, exist := userConverCache.Get(key); exist {
			return obj.(*UserConversations)
		} else {
			//get from db
			storage := storages.NewConversationStorage()
			var startTime int64 = time.Now().UnixMilli()
			userConvers := &UserConversations{
				Appkey:               appkey,
				UserId:               userId,
				ConverItemMap:        make(map[string]*models.Conversation),
				SyncTimeIndex:        zset.NewFloat64(),
				SortTimeIndex:        zset.NewFloat64(),
				TopIndexBaseTopTime:  zset.NewFloat64(),
				TopIndexBaseSortTime: zset.NewFloat64(),
			}
			var count int32 = 1000
			index := 0
			for {
				dbConvers, err := storage.QryConvers(appkey, userId, startTime, count)
				if err != nil {
					break
				}
				converLen := len(dbConvers)
				if converLen > 0 {
					for _, dbConver := range dbConvers {
						if dbConver.SyncTime < startTime {
							startTime = dbConver.SyncTime
						}
						itemKey := getConverItemKey(dbConver.TargetId, dbConver.ChannelType)
						userConvers.ConverItemMap[itemKey] = dbConver
						// build index
						userConvers.SyncTimeIndex.Add(float64(dbConver.SyncTime), itemKey)
						if dbConver.IsDeleted == 0 {
							userConvers.SortTimeIndex.Add(float64(dbConver.SortTime), itemKey)
						}
						if dbConver.IsTop > 0 {
							userConvers.TopIndexBaseSortTime.Add(float64(dbConver.SortTime), itemKey)
							userConvers.TopIndexBaseTopTime.Add(float64(dbConver.TopUpdatedTime), itemKey)
						}
					}
				}
				if converLen < int(count) {
					break
				}
				index++
				if index >= 10 {
					break
				}
			}
			userConverCache.Add(key, userConvers)
			return userConvers
		}
	}
}

func UserConversContains(appkey, userId string) bool {
	key := getUserConverCacheKey(appkey, userId)
	return userConverCache.Contains(key)
}
