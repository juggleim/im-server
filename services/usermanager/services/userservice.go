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
	userStorage "im-server/services/usermanager/storages"
	"im-server/services/usermanager/storages/dbs"
	"im-server/services/usermanager/storages/models"
	"strings"
	"time"
)

var (
	userCache *caches.LruCache
	userLocks *tools.SegmentatedLocks
)

type UserInfo struct {
	AppKey       string
	UserType     int
	UserId       string
	Nickname     string
	UserPortrait string
	ExtFields    map[string]string
	UpdatedTime  int64

	SettingFields map[string]string

	Statuses map[string]*StatusItem
}

type StatusItem struct {
	ItemKey     string
	ItemValue   string
	UpdatedTime int64
}

func (u *UserInfo) AddStatus(key, val string, updTime int64) {
	lock := userLocks.GetLocks(u.AppKey, u.UserId)
	lock.Lock()
	defer lock.Unlock()
	u.Statuses[key] = &StatusItem{
		ItemKey:     key,
		ItemValue:   val,
		UpdatedTime: updTime,
	}
}

func (u *UserInfo) GetStatus() map[string]*StatusItem {
	lock := userLocks.GetLocks(u.AppKey, u.UserId)
	lock.RLock()
	defer lock.RUnlock()
	return u.Statuses
}

func (u *UserInfo) CheckPrivateGlobalMute() bool {
	if val, exist := u.SettingFields[string(commonservices.AttItemKey_PriGlobalMute)]; exist {
		intVal, err := tools.String2Int64(val)
		if err == nil {
			if intVal == 0 || intVal > time.Now().UnixMilli() {
				return true
			} else {
				return false
			}
		}
		return true
	}
	return false
}

func (u *UserInfo) SetPriGlobalMute(isDelete bool, endTime int64) {
	lock := userLocks.GetLocks(u.AppKey, u.UserId)
	lock.Lock()
	defer lock.Unlock()
	if isDelete {
		delete(u.SettingFields, string(commonservices.AttItemKey_PriGlobalMute))
	} else {
		u.SettingFields[string(commonservices.AttItemKey_PriGlobalMute)] = tools.Int642String(endTime)
	}
}

var notExistUser *UserInfo

func init() {
	notExistUser = &UserInfo{}
	userCache = caches.NewLruCacheWithAddReadTimeout("user_cache", 100000, nil, 10*time.Minute, 10*time.Minute)
	userLocks = tools.NewSegmentatedLocks(512)
}

func AddUser(ctx context.Context, userId, nickname, userPortrait string, extFields []*pbobjs.KvItem, settings []*pbobjs.KvItem, userType pbobjs.UserType) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := strings.Join([]string{appkey, userId}, "_")
	userInfo, exist := GetUserInfo(appkey, userId)
	storage := userStorage.NewUserStorage()
	if exist && userInfo != nil {
		if nickname != userInfo.Nickname || userPortrait != userInfo.UserPortrait {
			err := storage.Upsert(models.User{
				UserId:       userId,
				UserType:     userType,
				Nickname:     nickname,
				UserPortrait: userPortrait,
				AppKey:       appkey,
			})
			if err == nil {
				userExtStorage := userStorage.NewUserExtStorage()
				for _, ext := range extFields {
					itemKey := ext.Key
					itemValue := ext.Value
					err = userExtStorage.Upsert(models.UserExt{
						AppKey:    appkey,
						UserId:    userId,
						ItemKey:   itemKey,
						ItemValue: itemValue,
						ItemType:  int(commonservices.AttItemType_Att),
					})
					if err != nil {
						logs.NewLogEntity().Error(err.Error())
					}
				}
				for _, set := range settings {
					err = userExtStorage.Upsert(models.UserExt{
						AppKey:    appkey,
						UserId:    userId,
						ItemKey:   set.Key,
						ItemValue: set.Value,
						ItemType:  int(commonservices.AttItemType_Setting),
					})
					if err != nil {
						logs.NewLogEntity().Error(err.Error())
					}
				}
			} else {
				logs.NewLogEntity().Error(err.Error())
			}
			userCache.Remove(key)
		}
	} else {
		err := storage.Upsert(models.User{
			UserId:       userId,
			UserType:     userType,
			Nickname:     nickname,
			UserPortrait: userPortrait,
			AppKey:       appkey,
		})
		if err == nil {
			userExtStorage := userStorage.NewUserExtStorage()
			for _, ext := range extFields {
				itemKey := ext.Key
				itemValue := ext.Value
				err = userExtStorage.Upsert(models.UserExt{
					AppKey:    appkey,
					UserId:    userId,
					ItemKey:   itemKey,
					ItemValue: itemValue,
					ItemType:  int(commonservices.AttItemType_Att),
				})
				if err != nil {
					logs.NewLogEntity().Error(err.Error())
				}
			}
			for _, set := range settings {
				err = userExtStorage.Upsert(models.UserExt{
					AppKey:    appkey,
					UserId:    userId,
					ItemKey:   set.Key,
					ItemValue: set.Value,
					ItemType:  int(commonservices.AttItemType_Setting),
				})
				if err != nil {
					logs.NewLogEntity().Error(err.Error())
				}
			}
		} else {
			logs.NewLogEntity().Error(err.Error())
		}
		userCache.Remove(key)
	}
	return errs.IMErrorCode_SUCCESS
}

func GetUserInfo(appkey, userId string) (*UserInfo, bool) {
	key := strings.Join([]string{appkey, userId}, "_")
	if userObj, exist := userCache.Get(key); exist {
		userInfo := userObj.(*UserInfo)
		if userInfo == notExistUser {
			return nil, false
		} else {
			return userInfo, true
		}
	} else {
		lock := userLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if userObj, exist := userCache.Get(key); exist {
			userInfo := userObj.(*UserInfo)
			if userInfo == notExistUser {
				return nil, false
			} else {
				return userInfo, true
			}
		} else {
			storage := userStorage.NewUserStorage()
			user, err := storage.FindByUserId(appkey, userId)
			if err == nil && user != nil {
				userInfo := &UserInfo{
					UserId:        userId,
					UserType:      int(user.UserType),
					Nickname:      user.Nickname,
					UserPortrait:  user.UserPortrait,
					ExtFields:     make(map[string]string),
					SettingFields: make(map[string]string),
					UpdatedTime:   user.UpdatedTime.UnixMilli(),
					Statuses:      make(map[string]*StatusItem),
				}
				//load extfields
				userExts, err := dbs.UserExtDao{}.QryExtFields(appkey, userId)
				if err == nil {
					for _, ext := range userExts {
						if ext.ItemType == 0 {
							userInfo.ExtFields[ext.ItemKey] = ext.ItemValue
							updTime := ext.UpdatedTime.UnixMilli()
							if updTime > userInfo.UpdatedTime {
								userInfo.UpdatedTime = updTime
							}
						} else if ext.ItemType == 1 { //setting
							userInfo.SettingFields[ext.ItemKey] = ext.ItemValue
						} else if ext.ItemType == 2 {
							userInfo.Statuses[ext.ItemKey] = &StatusItem{
								ItemKey:     ext.ItemKey,
								ItemValue:   ext.ItemValue,
								UpdatedTime: ext.UpdatedTime.UnixMilli(),
							}
						}
					}
				} else {
					logs.NewLogEntity().Error(err.Error())
				}
				userCache.Add(key, userInfo)
				return userInfo, true
			} else {
				userCache.Add(key, notExistUser)
				return nil, false
			}
		}
	}
}

func UpdUserInfo(ctx context.Context, userinfo *pbobjs.UserInfo) errs.IMErrorCode {
	userinfo.Nickname = tools.TruncateText(userinfo.Nickname, 32)
	appkey := bases.GetAppKeyFromCtx(ctx)
	_, exist := GetUserInfo(appkey, userinfo.UserId)
	if !exist {
		return errs.IMErrorCode_USER_NOT_EXIST
	}
	rvCache := false
	//upd db
	dao := dbs.UserDao{}
	err := dao.Update(appkey, userinfo.UserId, userinfo.Nickname, userinfo.UserPortrait)
	if err == nil {
		rvCache = rvCache || true
	}
	extStorage := userStorage.NewUserExtStorage()
	for _, ext := range userinfo.ExtFields {
		extStorage.Upsert(models.UserExt{
			AppKey:    appkey,
			UserId:    userinfo.UserId,
			ItemKey:   ext.Key,
			ItemValue: ext.Value,
			ItemType:  int(commonservices.AttItemType_Att),
		})
		rvCache = rvCache || true
	}
	for _, setting := range userinfo.Settings {
		extStorage.Upsert(models.UserExt{
			AppKey:    appkey,
			UserId:    userinfo.UserId,
			ItemKey:   setting.Key,
			ItemValue: setting.Value,
			ItemType:  int(commonservices.AttItemType_Setting),
		})
		rvCache = rvCache || true
	}
	//upd cache
	if rvCache {
		key := strings.Join([]string{appkey, userinfo.UserId}, "_")
		userCache.Remove(key)
	}
	return errs.IMErrorCode_SUCCESS
}

func SetUserUndisturb(ctx context.Context, userId string, req *pbobjs.UserUndisturb) errs.IMErrorCode {
	if req.Switch && req.Timezone != "" {
		_, err := time.LoadLocation(req.Timezone)
		if err != nil {
			logs.WithContext(ctx).Errorf("err:%v\treq:%s", err, tools.ToJson(req))
			return errs.IMErrorCode_USER_TIMEZONE_ILLGAL
		}
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	userInfo, exist := GetUserInfo(appkey, userId)
	if !exist || userInfo == nil {
		return errs.IMErrorCode_USER_NOT_EXIST
	}
	lock := userLocks.GetLocks(appkey, userId)
	lock.Lock()
	defer lock.Unlock()
	userUndisturb := &commonservices.UserUndisturb{
		Switch:   req.Switch,
		Timezone: req.Timezone,
		Rules:    []*commonservices.UserUndisturbItem{},
	}
	for _, rule := range req.Rules {
		userUndisturb.Rules = append(userUndisturb.Rules, &commonservices.UserUndisturbItem{
			Start: rule.Start,
			End:   rule.End,
		})
	}
	jsonStr := tools.ToJson(userUndisturb)
	oldUndisturb, exist := userInfo.SettingFields[string(commonservices.AttItemKey_Undisturb)]
	if exist && oldUndisturb == jsonStr {
		return errs.IMErrorCode_SUCCESS
	}
	userInfo.SettingFields[string(commonservices.AttItemKey_Undisturb)] = jsonStr
	//upd db
	extStorage := userStorage.NewUserExtStorage()
	extStorage.Upsert(models.UserExt{
		AppKey:    appkey,
		UserId:    userId,
		ItemKey:   string(commonservices.AttItemKey_Undisturb),
		ItemValue: jsonStr,
		ItemType:  int(commonservices.AttItemType_Setting),
	})
	return errs.IMErrorCode_SUCCESS
}

func GetUserUndisturb(ctx context.Context, userId string) (*pbobjs.UserUndisturb, errs.IMErrorCode) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userInfo, exist := GetUserInfo(appkey, userId)
	if !exist || userInfo == nil {
		return nil, errs.IMErrorCode_USER_NOT_EXIST
	}
	lock := userLocks.GetLocks(appkey, userId)
	lock.Lock()
	defer lock.Unlock()

	resp := &pbobjs.UserUndisturb{
		Rules: []*pbobjs.UserUndisturbItem{},
	}
	oldUndisturb, exist := userInfo.SettingFields[string(commonservices.AttItemKey_Undisturb)]
	if exist {
		var userUndisturb commonservices.UserUndisturb
		err := tools.JsonUnMarshal([]byte(oldUndisturb), &userUndisturb)
		if err != nil {
			logs.WithContext(ctx).Errorf("data format error:%s", err)
		}
		resp.Switch = userUndisturb.Switch
		resp.Timezone = userUndisturb.Timezone
		for _, rule := range userUndisturb.Rules {
			resp.Rules = append(resp.Rules, &pbobjs.UserUndisturbItem{
				Start: rule.Start,
				End:   rule.End,
			})
		}
	}
	return resp, errs.IMErrorCode_SUCCESS
}

func SetUserSettings(ctx context.Context, userId string, userinfo *pbobjs.UserInfo) errs.IMErrorCode {
	//check setting keys
	for _, item := range userinfo.Settings {
		if !commonservices.CheckUserSettingKey(item.Key) {
			return errs.IMErrorCode_USER_NOT_SUPPROT_SETTING
		}
	}

	appkey := bases.GetAppKeyFromCtx(ctx)
	_, exist := GetUserInfo(appkey, userId)
	if !exist {
		return errs.IMErrorCode_USER_NOT_EXIST
	}
	rvCache := false
	//upd db
	extStorage := userStorage.NewUserExtStorage()
	for _, item := range userinfo.Settings {
		extStorage.Upsert(models.UserExt{
			AppKey:    appkey,
			UserId:    userId,
			ItemKey:   item.Key,
			ItemValue: item.Value,
			ItemType:  int(commonservices.AttItemType_Setting),
		})
		rvCache = rvCache || true
	}
	//upd cache
	if rvCache {
		key := strings.Join([]string{appkey, userId}, "_")
		userCache.Remove(key)
	}
	return errs.IMErrorCode_SUCCESS
}

func SetUserStatus(ctx context.Context, userId string, userinfo *pbobjs.UserInfo) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userInfo, exist := GetUserInfo(appkey, userId)
	if !exist {
		return errs.IMErrorCode_USER_NOT_EXIST
	}
	uInfo := &pbobjs.UserInfo{
		UserId:   userId,
		Statuses: []*pbobjs.KvItem{},
	}
	extStorage := userStorage.NewUserExtStorage()
	for _, item := range userinfo.Statuses {
		//upd db
		extStorage.Upsert(models.UserExt{
			AppKey:    appkey,
			UserId:    userId,
			ItemKey:   item.Key,
			ItemValue: item.Value,
			ItemType:  int(commonservices.AttItemType_Status),
		})
		//upd cache
		updTime := time.Now().UnixMilli()
		userInfo.AddStatus(item.Key, item.Value, updTime)
		uInfo.Statuses = append(uInfo.Statuses, &pbobjs.KvItem{
			Key:     item.Key,
			Value:   item.Value,
			UpdTime: updTime,
		})
	}
	//publish to subscribers
	bases.AsyncRpcCall(ctx, "pub_status", userId, uInfo)
	return errs.IMErrorCode_SUCCESS
}

func CheckUserExist(ctx context.Context, userId string) bool {
	appkey := bases.GetAppKeyFromCtx(ctx)
	_, exist := GetUserInfo(appkey, userId)
	return exist
}
