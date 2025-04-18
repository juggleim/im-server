package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

var (
	batchExecutorPool *tools.BatchExecutorPool
	converCache       *caches.LruCache

	UndisturbType_None   int32 = 0
	UndisturbType_Normal int32 = 1
)

func init() {
	converCache = caches.NewLruCacheWithAddReadTimeout("msg_conver_cache", 100000, nil, 10*time.Minute, 10*time.Minute)
	batchExecutorPool = tools.NewBatchExecutorPool(128, 100, 5*time.Second, batchSaveConver)
}

type UserConversationItem struct {
	key               string
	UndisturbType     int32
	UnreadIndex       int64
	LatestReadMsgTime int64
	ConverTags        []*pbobjs.ConverTag
}

func (conver *UserConversationItem) GetUnreadIndex() int64 {
	lock := userLocks.GetLocks(conver.key)
	lock.Lock()
	defer lock.Unlock()
	conver.UnreadIndex++
	return conver.UnreadIndex
}

func UserConverCacheContains(appkey, userId, targetId string, channelType pbobjs.ChannelType) bool {
	key := fmt.Sprintf("%s_%s_%s_%d", appkey, userId, targetId, channelType)
	return converCache.Contains(key)
}

func GetConversation(ctx context.Context, userId, targetId string, channelType pbobjs.ChannelType) *UserConversationItem {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := fmt.Sprintf("%s_%s_%s_%d", appkey, userId, targetId, channelType)
	if converObj, exist := converCache.Get(key); exist {
		converItem := converObj.(*UserConversationItem)
		return converItem
	} else {
		lock := userLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if undisturbObj, exist := converCache.Get(key); exist {
			undisturbItem := undisturbObj.(*UserConversationItem)
			return undisturbItem
		} else {
			undisturbItem := QryConversation(ctx, userId, targetId, channelType)
			undisturbItem.key = key
			converCache.Add(key, undisturbItem)
			return undisturbItem
		}
	}
}

func CacheUserConver(appkey, userId, targetId string, channelType pbobjs.ChannelType, conver *UserConversationItem) {
	key := fmt.Sprintf("%s_%s_%s_%d", appkey, userId, targetId, channelType)
	l := userLocks.GetLocks(key)
	l.Lock()
	defer l.Unlock()
	if !UserConverCacheContains(appkey, userId, targetId, channelType) {
		converCache.Add(key, conver)
	}
}

func BatchInitUserConvers(ctx context.Context, targetId string, channelType pbobjs.ChannelType, userIds []string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	groups := bases.GroupTargets("qry_conver", userIds)
	wg := sync.WaitGroup{}
	for _, ids := range groups {
		wg.Add(1)
		uIds := ids
		go func() {
			defer wg.Done()
			_, resp, err := bases.SyncRpcCall(ctx, "qry_conver", uIds[0], &pbobjs.QryConverReq{
				TargetId:    targetId,
				ChannelType: channelType,
				IsInner:     true,
				UserIds:     uIds,
			}, func() proto.Message {
				return &pbobjs.QryConversationsResp{}
			})
			if err == nil {
				convers, ok := resp.(*pbobjs.QryConversationsResp)
				if ok && convers != nil {
					for _, conver := range convers.Conversations {
						key := fmt.Sprintf("%s_%s_%s_%d", appkey, conver.UserId, targetId, channelType)
						CacheUserConver(appkey, conver.UserId, targetId, channelType, &UserConversationItem{
							key:               key,
							UndisturbType:     conver.UndisturbType,
							UnreadIndex:       conver.LatestUnreadIndex,
							LatestReadMsgTime: conver.LatestReadMsgTime,
							ConverTags:        conver.ConverTags,
						})
					}
				}
			}
		}()
	}
	wg.Wait()
}

func ClearConversation(ctx context.Context, userId, targetId string, channelType pbobjs.ChannelType) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := fmt.Sprintf("%s_%s_%s_%d", appkey, userId, targetId, channelType)
	converCache.Remove(key)
}

func QryConversation(ctx context.Context, userId, targetId string, channelType pbobjs.ChannelType) *UserConversationItem {
	code, resp, err := bases.SyncRpcCall(ctx, "qry_conver", userId, &pbobjs.QryConverReq{
		TargetId:    targetId,
		ChannelType: channelType,
		IsInner:     true,
	}, func() proto.Message {
		return &pbobjs.Conversation{}
	})
	if err == nil && code == errs.IMErrorCode_SUCCESS && resp != nil {
		conver, ok := resp.(*pbobjs.Conversation)
		if ok {
			return &UserConversationItem{
				UndisturbType:     conver.UndisturbType,
				UnreadIndex:       conver.LatestUnreadIndex,
				ConverTags:        conver.ConverTags,
				LatestReadMsgTime: conver.LatestReadMsgTime,
			}
		}
	}
	return &UserConversationItem{
		UndisturbType: 0,
		UnreadIndex:   0,
	}
}

// handle conversation check, such as undisturb, unread index
// userId is the receiver
func HandleDownMsgByConver(ctx context.Context, userId, targetId string, channelType pbobjs.ChannelType, downMsg *pbobjs.DownMsg) {
	conver := GetConversation(ctx, userId, targetId, channelType)
	if conver.UndisturbType == UndisturbType_Normal {
		downMsg.UndisturbType = UndisturbType_Normal
		if !commonservices.IsMentionedMe(userId, downMsg) {
			downMsg.Flags = msgdefines.SetUndisturbMsg(downMsg.Flags)
		}
	} else {
		userSettings := commonservices.GetTargetUserSettings(ctx, userId)
		if userSettings != nil && userSettings.UndisturbObj != nil {
			if userSettings.UndisturbObj.CheckUndisturb(ctx, userId) {
				if !commonservices.IsMentionedMe(userId, downMsg) {
					downMsg.Flags = msgdefines.SetUndisturbMsg(downMsg.Flags)
				}
			}
		}
	}
	if msgdefines.IsCountMsg(downMsg.Flags) {
		downMsg.UnreadIndex = conver.GetUnreadIndex()
	}
	downMsg.ConverTags = append(downMsg.ConverTags, conver.ConverTags...)
}

type BatchConverItem struct {
	Appkey string
	UserId string
	Msg    *pbobjs.DownMsg
}

func batchSaveConver(tasks []interface{}) {
	grp4Appkey := map[string][]*pbobjs.Conversation{}
	for _, task := range tasks {
		item, ok := task.(*BatchConverItem)
		if ok && item != nil {
			var items []*pbobjs.Conversation
			if existItems, ok := grp4Appkey[item.Appkey]; ok {
				items = existItems
			} else {
				items = []*pbobjs.Conversation{}
			}
			items = append(items, &pbobjs.Conversation{
				UserId: item.UserId,
				Msg:    item.Msg,
			})
			grp4Appkey[item.Appkey] = items
		}
	}
	for appkey, convers := range grp4Appkey {
		//build context
		ctx := context.Background()
		ctx = context.WithValue(ctx, bases.CtxKey_AppKey, appkey)
		ctx = context.WithValue(ctx, bases.CtxKey_Session, tools.GenerateUUIDShort11())
		commonservices.BatchSaveConversation(ctx, convers)
	}
}
