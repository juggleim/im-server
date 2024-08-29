package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"time"

	"google.golang.org/protobuf/proto"
)

var (
	converCache *caches.LruCache

	UndisturbType_None   int32 = 0
	UndisturbType_Normal int32 = 1
)

func init() {
	converCache = caches.NewLruCacheWithAddReadTimeout(100000, nil, 10*time.Minute, 10*time.Minute)
}

type UserConversationItem struct {
	key           string
	UndisturbType int32
	UnreadIndex   int64
}

func (conver *UserConversationItem) GetUnreadIndex() int64 {
	lock := userLocks.GetLocks(conver.key)
	lock.Lock()
	defer lock.Unlock()
	conver.UnreadIndex++
	return conver.UnreadIndex
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
				UndisturbType: conver.UndisturbType,
				UnreadIndex:   conver.LatestUnreadIndex,
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
		downMsg.Flags = commonservices.SetUndisturbMsg(downMsg.Flags)
	} else {
		userSettings := commonservices.GetTargetUserSettings(ctx, userId)
		if userSettings != nil && userSettings.UndisturbObj != nil {
			if userSettings.UndisturbObj.CheckUndisturb(ctx, userId) {
				downMsg.Flags = commonservices.SetUndisturbMsg(downMsg.Flags)
			}
		}
	}
	if commonservices.IsCountMsg(downMsg.Flags) {
		downMsg.UnreadIndex = conver.GetUnreadIndex()
	}
}
