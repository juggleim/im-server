package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"strings"
	"sync"
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

func BatchGetUserInfoFromRpc(ctx context.Context, userIds []string) map[string]*pbobjs.UserInfo {
	targetGroups := map[string][]string{}
	method := "qry_user_info_by_ids"
	for _, userId := range userIds {
		node := bases.GetCluster().GetTargetNode(method, userId)
		if node != nil {
			var uids []string
			if existUids, ok := targetGroups[node.Name]; ok {
				uids = existUids
			} else {
				uids = []string{}
			}
			uids = append(uids, userId)
			targetGroups[node.Name] = uids
		}
	}
	wg := &sync.WaitGroup{}
	lock := &sync.RWMutex{}
	retUserMap := make(map[string]*pbobjs.UserInfo)
	for _, ids := range targetGroups {
		wg.Add(1)
		uids := ids
		if len(uids) > 0 {
			go func() {
				defer wg.Done()
				code, respObj, err := bases.SyncRpcCall(ctx, method, uids[0], &pbobjs.UserIdsReq{
					UserIds:  uids,
					AttTypes: []int32{int32(AttItemType_Att), int32(AttItemType_Setting)},
				}, func() proto.Message {
					return &pbobjs.UserInfosResp{}
				})
				if err == nil && code == errs.IMErrorCode_SUCCESS && respObj != nil {
					resp := respObj.(*pbobjs.UserInfosResp)
					if len(resp.UserInfoMap) > 0 {
						lock.Lock()
						defer lock.Unlock()
						for k, v := range resp.UserInfoMap {
							retUserMap[k] = v
						}
					}
				}
			}()
		}
	}
	wg.Wait()
	return retUserMap
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
	UserId         string
	Nickname       string
	UserPortrait   string
	ExtFields      []*pbobjs.KvItem
	UpdatedTime    int64
	Settings       *UserSettings
	SettingsFields []*pbobjs.KvItem
	UserType       pbobjs.UserType
}

func init() {
	targetUserCache = caches.NewLruCacheWithAddReadTimeout("userinfo_cache", 100000, nil, 5*time.Second, 5*time.Second)
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
				UserId:         userId,
				Nickname:       uInfo.Nickname,
				UserPortrait:   uInfo.UserPortrait,
				ExtFields:      uInfo.ExtFields,
				UpdatedTime:    uInfo.UpdatedTime,
				Settings:       &UserSettings{},
				SettingsFields: uInfo.Settings,
				UserType:       uInfo.UserType,
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

func BatchGetTargetUserInfo(ctx context.Context, userIds []string) map[string]*TargetUserInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	noCacheUserIds := []string{}
	ret := map[string]*TargetUserInfo{}
	for _, userId := range userIds {
		cacheKey := getKey(appkey, userId)
		if !targetUserCache.Contains(cacheKey) {
			noCacheUserIds = append(noCacheUserIds, userId)
		} else {
			obj, exist := targetUserCache.Get(cacheKey)
			if exist && obj != nil {
				ret[userId] = obj.(*TargetUserInfo)
			} else {
				noCacheUserIds = append(noCacheUserIds, userId)
			}
		}
	}
	if len(noCacheUserIds) > 0 {
		userMap := BatchGetUserInfoFromRpc(ctx, noCacheUserIds)
		for userId, userInfo := range userMap {
			targetUserInfo := &TargetUserInfo{
				UserId:       userId,
				Nickname:     userInfo.Nickname,
				UserPortrait: userInfo.UserPortrait,
				ExtFields:    userInfo.ExtFields,
				UpdatedTime:  userInfo.UpdatedTime,
				Settings:     &UserSettings{},
				UserType:     userInfo.UserType,
			}
			FillObjField(targetUserInfo.Settings, Kvitems2Map(userInfo.Settings))
			if targetUserInfo.Settings.Undisturb != "" {
				var userUndisturb UserUndisturb
				err := tools.JsonUnMarshal([]byte(targetUserInfo.Settings.Undisturb), &userUndisturb)
				if err == nil {
					targetUserInfo.Settings.UndisturbObj = &userUndisturb
				}
			}
			cacheKey := getKey(appkey, userId)
			targetUserCache.Add(cacheKey, targetUserInfo)
			ret[userId] = targetUserInfo
		}
	}
	return ret
}

func GetTargetDisplayUserInfo(ctx context.Context, userId string) *pbobjs.UserInfo {
	tUserInfo := GetTargetUserInfo(ctx, userId)
	return &pbobjs.UserInfo{
		UserId:       tUserInfo.UserId,
		Nickname:     tUserInfo.Nickname,
		UserPortrait: tUserInfo.UserPortrait,
		ExtFields:    tUserInfo.ExtFields,
		UpdatedTime:  tUserInfo.UpdatedTime,
		UserType:     tUserInfo.UserType,
	}
}

func GetTargetDisplayUserInfosMap(ctx context.Context, userIds []string) map[string]*pbobjs.UserInfo {
	targetUserMap := BatchGetTargetUserInfo(ctx, userIds)
	userMap := map[string]*pbobjs.UserInfo{}
	for userId, tUserInfo := range targetUserMap {
		userMap[userId] = &pbobjs.UserInfo{
			UserId:       tUserInfo.UserId,
			Nickname:     tUserInfo.Nickname,
			UserPortrait: tUserInfo.UserPortrait,
			ExtFields:    tUserInfo.ExtFields,
			UpdatedTime:  tUserInfo.UpdatedTime,
		}
	}
	return userMap
}

func GetTargetDisplayUserInfos(ctx context.Context, userIds []string) []*pbobjs.UserInfo {
	targetUserMap := BatchGetTargetUserInfo(ctx, userIds)
	users := []*pbobjs.UserInfo{}
	for _, userInfo := range targetUserMap {
		users = append(users, &pbobjs.UserInfo{
			UserId:       userInfo.UserId,
			Nickname:     userInfo.Nickname,
			UserPortrait: userInfo.UserPortrait,
			ExtFields:    userInfo.ExtFields,
			UpdatedTime:  userInfo.UpdatedTime,
		})
	}
	return users
}

func GetTargetUserSettings(ctx context.Context, userId string) *UserSettings {
	tUserInfo := GetTargetUserInfo(ctx, userId)
	return tUserInfo.Settings
}

func getKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}
