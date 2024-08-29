package services

// import (
// 	"context"
// 	"im-server/commons/bases"
// 	"im-server/commons/caches"
// 	"im-server/commons/pbdefines/pbobjs"
// 	"im-server/commons/tools"
// 	"im-server/services/commonservices"
// 	"time"
// )

// var targetUserCache *caches.LruCache

// type TargetUserInfo struct {
// 	UserId       string
// 	Nickname     string
// 	UserPortrait string
// 	ExtFields    []*pbobjs.KvItem
// 	UpdatedTime  int64
// 	Settings     *commonservices.UserSettings
// 	UserType     pbobjs.UserType
// }

// func init() {
// 	targetUserCache = caches.NewLruCacheWithAddReadTimeout(100000, nil, 5*time.Second, 5*time.Second)
// }

// func GetTargetUserInfo(ctx context.Context, userId string) *TargetUserInfo {
// 	appkey := bases.GetAppKeyFromCtx(ctx)
// 	key := getKey(appkey, userId)
// 	if val, exist := targetUserCache.Get(key); exist {
// 		return val.(*TargetUserInfo)
// 	} else {
// 		l := userLocks.GetLocks(appkey, userId)
// 		l.Lock()
// 		defer l.Unlock()
// 		if val, exist := targetUserCache.Get(key); exist {
// 			return val.(*TargetUserInfo)
// 		} else {
// 			uInfo := commonservices.GetUserInfoFromRpcWithAttTypes(ctx, userId, []int32{int32(commonservices.AttItemType_Att), int32(commonservices.AttItemType_Setting)})
// 			targetUserInfo := &TargetUserInfo{
// 				UserId:       userId,
// 				Nickname:     uInfo.Nickname,
// 				UserPortrait: uInfo.UserPortrait,
// 				ExtFields:    uInfo.ExtFields,
// 				UpdatedTime:  uInfo.UpdatedTime,
// 				Settings:     &commonservices.UserSettings{},
// 				UserType:     uInfo.UserType,
// 			}
// 			commonservices.FillObjField(targetUserInfo.Settings, commonservices.Kvitems2Map(uInfo.Settings))
// 			if targetUserInfo.Settings.Undisturb != "" {
// 				var userUndisturb commonservices.UserUndisturb
// 				err := tools.JsonUnMarshal([]byte(targetUserInfo.Settings.Undisturb), &userUndisturb)
// 				if err == nil {
// 					targetUserInfo.Settings.UndisturbObj = &userUndisturb
// 				}
// 			}
// 			targetUserCache.Add(key, targetUserInfo)
// 			return targetUserInfo
// 		}
// 	}
// }

// func GetTargetDisplayUserInfo(ctx context.Context, userId string) *pbobjs.UserInfo {
// 	tUserInfo := GetTargetUserInfo(ctx, userId)
// 	return &pbobjs.UserInfo{
// 		UserId:       tUserInfo.UserId,
// 		Nickname:     tUserInfo.Nickname,
// 		UserPortrait: tUserInfo.UserPortrait,
// 		ExtFields:    tUserInfo.ExtFields,
// 		UpdatedTime:  tUserInfo.UpdatedTime,
// 	}
// }

// func GetTargetUserSettings(ctx context.Context, userId string) *commonservices.UserSettings {
// 	tUserInfo := GetTargetUserInfo(ctx, userId)
// 	return tUserInfo.Settings
// }
