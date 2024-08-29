package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

func GetSenderUserInfo(ctx context.Context) *pbobjs.UserInfo {
	userInfo := bases.GetSenderInfoFromCtx(ctx)
	if userInfo == nil {
		userId := bases.GetRequesterIdFromCtx(ctx)
		userInfo = GetTargetDisplayUserInfo(ctx, userId)
	}
	return userInfo
}

func GetUserInfoFromRpcWithAttTypes(ctx context.Context, userId string, attTypes []int32) *pbobjs.UserInfo {
	_, respObj, err := bases.SyncRpcCall(ctx, "qry_user_info", userId, &pbobjs.UserIdReq{
		UserId:   userId,
		AttTypes: attTypes,
	}, func() proto.Message {
		return &pbobjs.UserInfo{}
	})
	if err == nil && respObj != nil {
		return respObj.(*pbobjs.UserInfo)
	}
	return &pbobjs.UserInfo{
		UserId: userId,
	}
}

func Map2KvItems(m map[string]string) []*pbobjs.KvItem {
	items := []*pbobjs.KvItem{}
	for k, v := range m {
		items = append(items, &pbobjs.KvItem{
			Key:   k,
			Value: v,
		})
	}
	return items
}

func Kvitems2Map(items []*pbobjs.KvItem) map[string]string {
	m := make(map[string]string)
	for _, item := range items {
		m[item.Key] = item.Value
	}
	return m
}

var targetUserCache *caches.LruCache
var targetUserLocks *tools.SegmentatedLocks

type TargetUserInfo struct {
	UserId       string
	Nickname     string
	UserPortrait string
	ExtFields    []*pbobjs.KvItem
	UpdatedTime  int64
	Settings     *UserSettings
	UserType     pbobjs.UserType
}

func init() {
	targetUserCache = caches.NewLruCacheWithAddReadTimeout(100000, nil, 5*time.Second, 5*time.Second)
	targetUserLocks = tools.NewSegmentatedLocks(256)
}

func GetTargetUserInfo(ctx context.Context, userId string) *TargetUserInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getKey(appkey, userId)
	if val, exist := targetUserCache.Get(key); exist {
		return val.(*TargetUserInfo)
	} else {
		l := targetUserLocks.GetLocks(appkey, userId)
		l.Lock()
		defer l.Unlock()
		if val, exist := targetUserCache.Get(key); exist {
			return val.(*TargetUserInfo)
		} else {
			uInfo := GetUserInfoFromRpcWithAttTypes(ctx, userId, []int32{int32(AttItemType_Att), int32(AttItemType_Setting)})
			targetUserInfo := &TargetUserInfo{
				UserId:       userId,
				Nickname:     uInfo.Nickname,
				UserPortrait: uInfo.UserPortrait,
				ExtFields:    uInfo.ExtFields,
				UpdatedTime:  uInfo.UpdatedTime,
				Settings:     &UserSettings{},
				UserType:     uInfo.UserType,
			}
			FillObjField(targetUserInfo.Settings, Kvitems2Map(uInfo.Settings))
			if targetUserInfo.Settings.Undisturb != "" {
				var userUndisturb UserUndisturb
				err := tools.JsonUnMarshal([]byte(targetUserInfo.Settings.Undisturb), &userUndisturb)
				if err == nil {
					targetUserInfo.Settings.UndisturbObj = &userUndisturb
				}
			}
			targetUserCache.Add(key, targetUserInfo)
			return targetUserInfo
		}
	}
}

func GetTargetDisplayUserInfo(ctx context.Context, userId string) *pbobjs.UserInfo {
	tUserInfo := GetTargetUserInfo(ctx, userId)
	return &pbobjs.UserInfo{
		UserId:       tUserInfo.UserId,
		Nickname:     tUserInfo.Nickname,
		UserPortrait: tUserInfo.UserPortrait,
		ExtFields:    tUserInfo.ExtFields,
		UpdatedTime:  tUserInfo.UpdatedTime,
	}
}

func GetTargetUserSettings(ctx context.Context, userId string) *UserSettings {
	tUserInfo := GetTargetUserInfo(ctx, userId)
	return tUserInfo.Settings
}

func getKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}
