package services

import (
	"context"
	"strings"
	"time"

	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/convercache"
	"im-server/services/commonservices/logs"
	"im-server/services/statussubscriptions/storages"
)

var userSubRelCache *caches.LruCache
var userSubRelLocks *tools.SegmentatedLocks

func init() {
	userSubRelCache = caches.NewLruCacheWithAddReadTimeout("user_sub_rel_cache", 100000, nil, 10*time.Minute, 10*time.Minute)
	userSubRelLocks = tools.NewSegmentatedLocks(512)
}

type UserSubscribers struct {
	AppKey      string
	UserId      string
	Subscribers map[string]*Subscriber
}

type Subscriber struct {
	SubscriberId string
	DeviceIdMap  map[string]bool
}

func (rel *UserSubscribers) AddSubscriber(subscriberId, deviceId string) {
	key := getUserSubscribersCacheKey(rel.AppKey, rel.UserId)
	lock := userSubRelLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	subscriberId = strings.TrimSpace(subscriberId)
	deviceId = strings.TrimSpace(deviceId)
	if subscriberId == "" || deviceId == "" {
		return
	}
	if rel.Subscribers == nil {
		rel.Subscribers = make(map[string]*Subscriber)
	}
	sub, ok := rel.Subscribers[subscriberId]
	if !ok {
		rel.Subscribers[subscriberId] = &Subscriber{
			SubscriberId: subscriberId,
			DeviceIdMap:  map[string]bool{deviceId: true},
		}
		return
	}
	if sub.DeviceIdMap == nil {
		sub.DeviceIdMap = make(map[string]bool)
	}
	if !sub.DeviceIdMap[deviceId] {
		sub.DeviceIdMap[deviceId] = true
	}
}

func (rel *UserSubscribers) RemoveSubscriber(subscriberId, deviceId string) {
	key := getUserSubscribersCacheKey(rel.AppKey, rel.UserId)
	lock := userSubRelLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	subscriberId = strings.TrimSpace(subscriberId)
	deviceId = strings.TrimSpace(deviceId)
	if subscriberId == "" {
		return
	}
	if rel.Subscribers == nil {
		return
	}
	sub, ok := rel.Subscribers[subscriberId]
	if !ok {
		return
	}
	if deviceId == "" {
		delete(rel.Subscribers, subscriberId)
		return
	}
	if sub.DeviceIdMap == nil {
		delete(rel.Subscribers, subscriberId)
		return
	}
	delete(sub.DeviceIdMap, deviceId)
	if len(sub.DeviceIdMap) == 0 {
		delete(rel.Subscribers, subscriberId)
	}
}

func getUserSubscribersCacheKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}

func GetUserSubscribers(appkey, userId string) *UserSubscribers {
	key := getUserSubscribersCacheKey(appkey, userId)
	if userSub, exist := userSubRelCache.Get(key); exist {
		return userSub.(*UserSubscribers)
	} else {
		lock := userSubRelLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if userSub, exist := userSubRelCache.Get(key); exist {
			return userSub.(*UserSubscribers)
		} else {
			userSub := getUserSubscribersFromDb(appkey, userId)
			userSubRelCache.Add(key, userSub)
			return userSub
		}
	}
}

func getUserSubscribersFromDb(appkey, userId string) *UserSubscribers {
	var targetUserSubPageSize int = 10000
	out := &UserSubscribers{
		AppKey:      appkey,
		UserId:      userId,
		Subscribers: make(map[string]*Subscriber),
	}
	appkey = strings.TrimSpace(appkey)
	userId = strings.TrimSpace(userId)
	if appkey == "" || userId == "" {
		return out
	}
	storage := storages.NewUserSubRelStorage()
	var afterID int64
	for {
		rels, err := storage.QryByUserID(appkey, userId, afterID, targetUserSubPageSize)
		if err != nil {
			logs.WithContext(context.Background()).Errorf("QryByUserID appkey=%s userId=%s afterID=%d err=%v", appkey, userId, afterID, err)
			break
		}
		if len(rels) == 0 {
			break
		}
		for _, rel := range rels {
			if rel == nil {
				continue
			}
			sid := strings.TrimSpace(rel.SubscriberId)
			dev := strings.TrimSpace(rel.SubscriberDeviceId)
			if sid == "" || dev == "" {
				continue
			}
			sub, ok := out.Subscribers[sid]
			if !ok {
				sub = &Subscriber{
					SubscriberId: sid,
					DeviceIdMap:  make(map[string]bool),
				}
				out.Subscribers[sid] = sub
			}
			sub.DeviceIdMap[dev] = true
		}
		last := rels[len(rels)-1]
		afterID = last.ID
		if len(rels) < targetUserSubPageSize {
			break
		}
	}
	return out
}

func PublishUserStatus(ctx context.Context, upMsg *pbobjs.UpMsg) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	msgConverCache := convercache.GetMsgConverCache(ctx, userId, "", pbobjs.ChannelType_SubStatus)
	msgId, sendTime, _ := msgConverCache.GenerateMsgId(userId, pbobjs.ChannelType_SubStatus, time.Now().UnixMilli(), upMsg.Flags)
	downMsg := &pbobjs.DownMsg{
		SenderId:    userId,
		TargetId:    userId,
		ChannelType: pbobjs.ChannelType_SubStatus,
		MsgType:     upMsg.MsgType,
		MsgContent:  upMsg.MsgContent,
		MsgId:       msgId,
		MsgTime:     sendTime,
		Flags:       upMsg.Flags,
	}
	userSubscribers := GetUserSubscribers(appkey, userId)
	var memberIds []string
	if len(userSubscribers.Subscribers) > 0 {
		memberIds := make([]string, 0, len(userSubscribers.Subscribers))
		for sid := range userSubscribers.Subscribers {
			if sid != "" {
				memberIds = append(memberIds, sid)
			}
		}
		if len(memberIds) > 0 {
			Dispatch2Message(ctx, downMsg, memberIds)
		}
	}
	if appinfo, exist := commonservices.GetAppInfo(appkey); exist && appinfo != nil && appinfo.OpenFriendStatusSub {
		bases.AsyncRpcCall(ctx, "dispatch_user_status", userId, &pbobjs.UserStatusFriDispatch{
			Msg:             downMsg,
			ExcludedUserIds: memberIds,
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func Dispatch2Message(ctx context.Context, downMsg *pbobjs.DownMsg, memberIds []string) {
	bases.GroupRpcCall(ctx, "msg_dispatch", memberIds, downMsg)
}
