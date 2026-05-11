package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/statussubscriptions/storages"
	"im-server/services/statussubscriptions/storages/models"
	"strings"
	"time"

	"github.com/bytedance/gopkg/collection/zset"
)

var subscriberUsersCache *caches.LruCache
var subscriberUsersLocks *tools.SegmentatedLocks

func init() {
	subscriberUsersCache = caches.NewLruCacheWithAddReadTimeout("subscriber_users_cache", 100000, nil, 10*time.Minute, 10*time.Minute)
	subscriberUsersLocks = tools.NewSegmentatedLocks(512)
}

func getSubscriberDeviceCacheKey(appley, subscriberId, deviceId string) string {
	return strings.Join([]string{appley, subscriberId, deviceId}, "_")
}

type SubUserItem struct {
	UserId      string
	RelId       int64
	CreatedTime int64
}

type SubscriberCacheItem struct {
	AppKey        string
	SubscriberId  string
	DeviceId      string
	SubUsers      map[string]*SubUserItem
	SubUsersIndex *zset.Float64Set
}

func (sub *SubscriberCacheItem) AddSubUsers(userIds []string, maxUserSubCount int) (added []*SubUserItem, evicted []*SubUserItem) {
	key := getSubscriberDeviceCacheKey(sub.AppKey, sub.SubscriberId, sub.DeviceId)
	lock := subscriberUsersLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	curr := time.Now().UnixMilli()
	var newlyAdded []*SubUserItem
	for _, raw := range userIds {
		userId := strings.TrimSpace(raw)
		if userId == "" {
			continue
		}
		if _, exist := sub.SubUsers[userId]; !exist {
			it := &SubUserItem{
				UserId:      userId,
				CreatedTime: curr,
			}
			sub.SubUsers[userId] = it
			sub.SubUsersIndex.Add(float64(curr), userId)
			newlyAdded = append(newlyAdded, it)
			curr = curr + 1
		}
	}
	if maxUserSubCount > 0 && len(sub.SubUsers) > maxUserSubCount {
		excess := len(sub.SubUsers) - maxUserSubCount
		nodes := sub.SubUsersIndex.RemoveRangeByRank(0, excess-1)
		evicted = make([]*SubUserItem, 0, len(nodes))
		for _, n := range nodes {
			if it, ok := sub.SubUsers[n.Value]; ok {
				evicted = append(evicted, it)
				delete(sub.SubUsers, n.Value)
			}
		}
	}
	for _, it := range newlyAdded {
		if it != nil && sub.SubUsers[it.UserId] == it {
			added = append(added, it)
		}
	}
	return added, evicted
}

func (sub *SubscriberCacheItem) RemoveSubUsers(userIds []string) []*SubUserItem {
	key := getSubscriberDeviceCacheKey(sub.AppKey, sub.SubscriberId, sub.DeviceId)
	lock := subscriberUsersLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	removed := make([]*SubUserItem, 0, len(userIds))
	for _, raw := range userIds {
		userId := strings.TrimSpace(raw)
		if userId == "" {
			continue
		}
		if it, ok := sub.SubUsers[userId]; ok {
			sub.SubUsersIndex.Remove(userId)
			delete(sub.SubUsers, userId)
			removed = append(removed, it)
		}
	}
	return removed
}

func GetSubscriberFromCache(appkey, subscriberId, deviceId string) *SubscriberCacheItem {
	key := getSubscriberDeviceCacheKey(appkey, subscriberId, deviceId)
	if obj, exist := subscriberUsersCache.Get(key); exist {
		return obj.(*SubscriberCacheItem)
	} else {
		lock := subscriberUsersLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if obj, exist := subscriberUsersCache.Get(key); exist {
			return obj.(*SubscriberCacheItem)
		} else {
			item := getSubscriberFromDb(appkey, subscriberId, deviceId)
			subscriberUsersCache.Add(key, item)
			return item
		}
	}
}

func getSubscriberFromDb(appkey, subscriberId, deviceId string) *SubscriberCacheItem {
	item := &SubscriberCacheItem{
		AppKey:        appkey,
		SubscriberId:  subscriberId,
		DeviceId:      deviceId,
		SubUsers:      make(map[string]*SubUserItem),
		SubUsersIndex: zset.NewFloat64(),
	}
	storage := storages.NewUserSubRelStorage()
	rels, err := storage.QryBySubscriber(appkey, subscriberId, deviceId, 10000)
	if err == nil {
		for _, rel := range rels {
			if _, exist := item.SubUsers[rel.UserId]; !exist {
				item.SubUsers[rel.UserId] = &SubUserItem{
					UserId: rel.UserId,
					RelId:  rel.ID,
				}
				item.SubUsersIndex.Add(float64(rel.CreatedTime), rel.UserId)
			}
		}
	}
	return item
}

func SubUsers(ctx context.Context, req *pbobjs.SubUsersReq) (errs.IMErrorCode, *pbobjs.UserStatusList) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	subscriberId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := strings.TrimSpace(bases.GetDeviceIdFromCtx(ctx))
	maxUserSubCount := 1000
	if appInfo, ok := commonservices.GetAppInfo(appkey); ok && appInfo != nil && appInfo.MaxUserSubscriptionCount > 0 {
		maxUserSubCount = appInfo.MaxUserSubscriptionCount
	}
	subscriberItem := GetSubscriberFromCache(appkey, subscriberId, deviceId)
	var uids []string
	if req != nil {
		uids = req.GetUserIds()
	}
	added, evicted := subscriberItem.AddSubUsers(uids, maxUserSubCount)
	storage := storages.NewUserSubRelStorage()
	if len(added) > 0 {
		rels := make([]*models.UserSubRel, 0, len(added))
		for _, it := range added {
			if it == nil {
				continue
			}
			rels = append(rels, &models.UserSubRel{
				AppKey:             appkey,
				UserId:             it.UserId,
				SubscriberId:       subscriberId,
				SubscriberDeviceId: deviceId,
			})
		}
		if len(rels) > 0 {
			targetIds := make([]string, 0, len(rels))
			if err := storage.BatchCreate(rels); err != nil {
				logs.WithContext(ctx).Error(err.Error())
			} else {
				j := 0
				for _, it := range added {
					if it == nil {
						continue
					}
					it.RelId = rels[j].ID
					if rels[j].CreatedTime > 0 {
						it.CreatedTime = rels[j].CreatedTime
					}
					j++
					targetIds = append(targetIds, it.UserId)
				}
			}
			if len(targetIds) > 0 {
				syncSubRelChg(ctx, pbobjs.StatusSubBusType_UserStatus, true, targetIds)
			}
		}
	}
	if len(evicted) > 0 {
		relIds := make([]int64, 0, len(evicted))
		for _, e := range evicted {
			if e == nil {
				continue
			}
			if e.RelId > 0 {
				relIds = append(relIds, e.RelId)
				continue
			}
			if err := storage.Delete(appkey, e.UserId, subscriberId, deviceId); err != nil {
				logs.WithContext(ctx).Error(err.Error())
			}
		}
		if len(relIds) > 0 {
			if err := storage.DeleteByRelIDs(relIds); err != nil {
				logs.WithContext(ctx).Error(err.Error())
			} else {
				targetIds := make([]string, 0, len(relIds))
				for _, it := range evicted {
					if it == nil {
						continue
					}
					targetIds = append(targetIds, it.UserId)
				}
				if len(targetIds) > 0 {
					syncSubRelChg(ctx, pbobjs.StatusSubBusType_UserStatus, false, targetIds)
				}
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, QryUserStatusList(ctx, req.UserIds)
}

func UnSubUsers(ctx context.Context, req *pbobjs.SubUsersReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	subscriberId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := strings.TrimSpace(bases.GetDeviceIdFromCtx(ctx))
	var uids []string
	if req != nil {
		uids = req.GetUserIds()
	}
	subscriberItem := GetSubscriberFromCache(appkey, subscriberId, deviceId)
	removed := subscriberItem.RemoveSubUsers(uids)
	if len(removed) == 0 {
		return errs.IMErrorCode_SUCCESS
	}
	storage := storages.NewUserSubRelStorage()
	relIds := make([]int64, 0, len(removed))
	for _, it := range removed {
		if it == nil {
			continue
		}
		if it.RelId > 0 {
			relIds = append(relIds, it.RelId)
			continue
		}
		if err := storage.Delete(appkey, it.UserId, subscriberId, deviceId); err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
	}
	if len(relIds) > 0 {
		if err := storage.DeleteByRelIDs(relIds); err != nil {
			logs.WithContext(ctx).Error(err.Error())
		} else {
			targetIds := make([]string, 0, len(relIds))
			for _, it := range removed {
				if it == nil {
					continue
				}
				targetIds = append(targetIds, it.UserId)
			}
			if len(targetIds) > 0 {
				syncSubRelChg(ctx, pbobjs.StatusSubBusType_UserStatus, false, targetIds)
			}
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func syncSubRelChg(ctx context.Context, busType pbobjs.StatusSubBusType, isAdd bool, targetIds []string) {
	req := &pbobjs.SubRelChangeReq{
		BusType: busType,
		IsAdd:   isAdd,
	}
	data, _ := tools.PbMarshal(req)
	groups := bases.GroupTargets("sync_sub_change", targetIds)
	for _, ids := range groups {
		bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
			RpcMsgType:   pbobjs.RpcMsgType_UserPub,
			AppKey:       bases.GetAppKeyFromCtx(ctx),
			Session:      bases.GetSessionFromCtx(ctx),
			Method:       "sync_sub_change",
			RequesterId:  bases.GetRequesterIdFromCtx(ctx),
			DeviceId:     bases.GetDeviceIdFromCtx(ctx),
			ReqIndex:     bases.GetSeqIndexFromCtx(ctx),
			Qos:          bases.GetQosFromCtx(ctx),
			AppDataBytes: data,
			TargetId:     ids[0],
			TargetIds:    ids,
		})
	}
}
