package services

import (
	"context"
	"errors"
	"expvar"
	"strconv"
	"time"

	"gorm.io/gorm"

	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/group/dbs"
)

var groupInfoCache *caches.LruCache
var groupLocks *tools.SegmentatedLocks

const (
	groupInfoCacheSoftTTL       = 10 * time.Minute
	groupInfoCacheHardTTL       = time.Hour
	groupInfoCacheRefreshJitter = time.Minute
)

// loadGroupInfoFromDB is a variable so cache-load concurrency can be tested
// without depending on an external database.
var loadGroupInfoFromDB = getGroupInfoFromDb

func init() {
	groupInfoCache = caches.NewLruCacheWithAddReadTimeout("groupinfo_cache", 10000, nil, groupInfoCacheSoftTTL, groupInfoCacheSoftTTL)
	groupLocks = tools.NewSegmentatedLocks(128)
	err := groupInfoCache.EnableStaleWhileRevalidate(caches.StaleWhileRevalidateOptions{
		SoftTTL:            groupInfoCacheSoftTTL,
		HardTTL:            groupInfoCacheHardTTL,
		RefreshJitter:      groupInfoCacheRefreshJitter,
		LoadTimeout:        10 * time.Second,
		MaxConcurrentLoads: 8,
		MaxPendingLoads:    64,
		ServeStaleOnError:  true,
		ApplyGuard: func(key interface{}) func() {
			cacheKey := key.(groupCacheKey)
			lock := groupLocks.GetLocks(cacheKey.appKey, cacheKey.groupID)
			lock.Lock()
			return lock.Unlock
		},
	}, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		cacheKey := key.(groupCacheKey)
		groupInfo := loadGroupInfoFromDB(cacheKey.appKey, cacheKey.groupID)
		if groupInfo == nil {
			return nil, false, errors.New("load group info from database failed")
		}
		return groupInfo, true, nil
	})
	if err != nil {
		panic(err)
	}
	expvar.Publish("group_info_cache", expvar.Func(func() interface{} {
		stats, _ := groupInfoCache.StaleWhileRevalidateStats()
		return map[string]interface{}{
			"cache_size":             groupInfoCache.Len(),
			"stale_while_revalidate": stats,
		}
	}))
}

type groupCacheKey struct {
	appKey  string
	groupID string
}

func (key groupCacheKey) lockKey() string {
	return strconv.Itoa(len(key.appKey)) + ":" + key.appKey + key.groupID
}

type GroupInfo struct {
	AppKey        string
	GroupId       string
	GroupName     string
	GroupPortrait string
	IsMute        int32
	ExtFields     map[string]string
	UpdatedTime   int64
	Settings      *commonservices.GroupSettings
	SettingFields map[string]string
}

func cloneGroupInfo(groupInfo *GroupInfo) *GroupInfo {
	if groupInfo == nil {
		return nil
	}
	cloned := *groupInfo
	if groupInfo.ExtFields != nil {
		cloned.ExtFields = make(map[string]string, len(groupInfo.ExtFields))
		for key, value := range groupInfo.ExtFields {
			cloned.ExtFields[key] = value
		}
	}
	if groupInfo.SettingFields != nil {
		cloned.SettingFields = make(map[string]string, len(groupInfo.SettingFields))
		for key, value := range groupInfo.SettingFields {
			cloned.SettingFields[key] = value
		}
	}
	if groupInfo.Settings != nil {
		settings := *groupInfo.Settings
		cloned.Settings = &settings
	}
	return &cloned
}

func (grp *GroupInfo) SetSetting(itemKey, itemValue string) {
	key := getGroupKey(grp.AppKey, grp.GroupId)
	lock := groupLocks.GetLocks(key.appKey, key.groupID)
	lock.Lock()
	defer lock.Unlock()
	if current, cached := currentGroupInfoFromCache(key); cached {
		grp = current
	}
	grp = cloneGroupInfo(grp)
	if grp.SettingFields == nil {
		grp.SettingFields = make(map[string]string)
	}
	grp.SettingFields[itemKey] = itemValue
	groupInfoCache.Add(key, grp)
}

func (grp *GroupInfo) SetExt(itemKey, itemValue string) {
	key := getGroupKey(grp.AppKey, grp.GroupId)
	lock := groupLocks.GetLocks(key.appKey, key.groupID)
	lock.Lock()
	defer lock.Unlock()
	if current, cached := currentGroupInfoFromCache(key); cached {
		grp = current
	}
	grp = cloneGroupInfo(grp)
	if grp.ExtFields == nil {
		grp.ExtFields = make(map[string]string)
	}
	grp.ExtFields[itemKey] = itemValue
	groupInfoCache.Add(key, grp)
}

func (grp *GroupInfo) GetMemberLimiter(memberId string) *commonservices.GrpMemberLimiter {
	key := getGroupKey(grp.AppKey, grp.GroupId)
	lock := groupLocks.GetLocks(key.appKey, key.groupID)
	lock.Lock()
	defer lock.Unlock()
	if grp.Settings != nil && (grp.Settings.GrpMsgSecondLimiter > 0 || grp.Settings.GrpMsgMinuteLimiter > 0 || grp.Settings.GrpMsgHourLimiter > 0) {
		if grp.Settings.MemberLimiterMap == nil {
			grp.Settings.MemberLimiterMap = make(map[string]*commonservices.GrpMemberLimiter)
		}
		if limiter, exist := grp.Settings.MemberLimiterMap[memberId]; exist {
			return limiter
		} else {
			limiter := &commonservices.GrpMemberLimiter{}
			if grp.Settings.GrpMsgSecondLimiter > 0 {
				limiter.GrpMsgSecondLimiter = tools.NewPerSecondLimiter(grp.Settings.GrpMsgSecondLimiter)
			}
			if grp.Settings.GrpMsgMinuteLimiter > 0 {
				limiter.GrpMsgMinuteLimiter = tools.NewPerMinuteLimiter(grp.Settings.GrpMsgMinuteLimiter)
			}
			if grp.Settings.GrpMsgHourLimiter > 0 {
				limiter.GrpMsgHourLimiter = tools.NewPerHourLimiter(grp.Settings.GrpMsgHourLimiter)
			}
			grp.Settings.MemberLimiterMap[memberId] = limiter
			return limiter
		}
	}
	return nil
}

var notExistGroup GroupInfo = GroupInfo{
	ExtFields:     make(map[string]string),
	Settings:      &commonservices.GroupSettings{},
	SettingFields: make(map[string]string),
}

func DissolveGroup(ctx context.Context, groupId string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//remove from db
	groupDao := dbs.GroupDao{}
	groupDao.Delete(appkey, groupId)
	memberDao := dbs.GroupMemberDao{}
	memberDao.DeleteByGroupId(appkey, groupId)
	memberExtDao := dbs.GroupMemberExtDao{}
	memberExtDao.DeleteByGroupId(appkey, groupId)

	//remove from cache
	DelGroupInfo(ctx, appkey, groupId)
	//dissolve group members
	clearMembersFromCache(ctx, appkey, groupId)
}

func DelGroupInfo(ctx context.Context, appkey, groupId string) {
	key := getGroupKey(appkey, groupId)
	groupInfoCache.Add(key, &notExistGroup)
}

func getGroupKey(appkey, groupId string) groupCacheKey {
	return groupCacheKey{appKey: appkey, groupID: groupId}
}
func AddGroupInfo2Cache(ctx context.Context, appkey, groupId string, grpInfo *GroupInfo) {
	key := getGroupKey(appkey, groupId)
	groupInfoCache.Add(key, grpInfo)
}
func GetGroupInfoFromCache(ctx context.Context, appkey, groupId string) (*GroupInfo, bool) {
	key := getGroupKey(appkey, groupId)
	value, found, err := groupInfoCache.GetOrLoad(ctx, key)
	if err != nil {
		return nil, false
	}
	if !found || value == nil {
		return nil, false
	}
	return resolveCachedGroupInfo(value)
}

func currentGroupInfoFromCache(key groupCacheKey) (*GroupInfo, bool) {
	entries := groupInfoCache.SnapshotForKeys([]interface{}{key})
	if len(entries) == 0 || !entries[0].Found || entries[0].Value == nil {
		return nil, false
	}
	groupInfo := entries[0].Value.(*GroupInfo)
	if groupInfo == &notExistGroup {
		return nil, false
	}
	return groupInfo, true
}

func resolveCachedGroupInfo(value interface{}) (*GroupInfo, bool) {
	groupInfo := value.(*GroupInfo)
	if groupInfo == &notExistGroup {
		return nil, false
	}
	return groupInfo, true
}

func getGroupInfoFromDb(appkey, groupId string) *GroupInfo {
	dao := dbs.GroupDao{}
	dbGroup, err := dao.FindById(appkey, groupId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &notExistGroup
		}
		return nil
	} else {
		groupInfo := &GroupInfo{
			AppKey:        appkey,
			GroupId:       dbGroup.GroupId,
			GroupPortrait: dbGroup.GroupPortrait,
			GroupName:     dbGroup.GroupName,
			IsMute:        int32(dbGroup.IsMute),
			UpdatedTime:   dbGroup.UpdatedTime.UnixMilli(),
			ExtFields:     map[string]string{},
			Settings:      &commonservices.GroupSettings{},
			SettingFields: map[string]string{},
		}
		groupExts, err := dbs.GroupExtDao{}.QryExtFields(appkey, groupId)
		settingValMap := make(map[string]string)
		//get default from appinfo
		appinfo, exist := commonservices.GetAppInfo(appkey)
		if exist && appinfo != nil {
			if appinfo.HideGrpMsg {
				settingValMap[string(commonservices.AttItemKey_HideGrpMsg)] = "1"
			} else {
				settingValMap[string(commonservices.AttItemKey_HideGrpMsg)] = "0"
			}
		}
		if err == nil {
			for _, ext := range groupExts {
				if ext.ItemType == int(commonservices.AttItemType_Att) {
					groupInfo.ExtFields[ext.ItemKey] = ext.ItemValue
					extUpdTime := ext.UpdatedTime.UnixMilli()
					if extUpdTime > groupInfo.UpdatedTime {
						groupInfo.UpdatedTime = extUpdTime
					}
				} else if ext.ItemType == int(commonservices.AttItemType_Setting) {
					settingValMap[ext.ItemKey] = ext.ItemValue
				}
			}
			groupInfo.SettingFields = settingValMap
			commonservices.FillObjField(groupInfo.Settings, settingValMap)
		}
		return groupInfo
	}
}

func SetGroupMute(ctx context.Context, req *pbobjs.GroupMuteReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//update db
	dao := dbs.GroupDao{}
	dao.UpdateGroupMuteStatus(appkey, req.GroupId, req.IsMute)

	//update cache
	_, _ = GetGroupInfoFromCache(ctx, appkey, req.GroupId)
	key := getGroupKey(appkey, req.GroupId)
	lock := groupLocks.GetLocks(key.appKey, key.groupID)
	lock.Lock()
	if groupInfo, exist := currentGroupInfoFromCache(key); exist {
		groupInfo = cloneGroupInfo(groupInfo)
		groupInfo.IsMute = req.IsMute
		groupInfoCache.Add(key, groupInfo)
	}
	lock.Unlock()
	return errs.IMErrorCode_SUCCESS
}

func QryGroupInfo(ctx context.Context, req *pbobjs.GroupInfoReq) (errs.IMErrorCode, *pbobjs.GroupInfo) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	if len(req.CareFields) > 0 {
		grpInfo := &pbobjs.GroupInfo{
			GroupId:   req.GroupId,
			ExtFields: []*pbobjs.KvItem{},
		}
		memberContainer, exist := GetGroupMembersFromCache(ctx, appkey, req.GroupId)
		for _, field := range req.CareFields {
			if field == commonservices.GroupField_MemberCount {
				if exist {
					grpInfo.ExtFields = append(grpInfo.ExtFields, &pbobjs.KvItem{
						Key:   commonservices.GroupField_MemberCount,
						Value: strconv.Itoa(memberContainer.GroupMemberCount()),
					})
				} else {
					grpInfo.ExtFields = append(grpInfo.ExtFields, &pbobjs.KvItem{
						Key:   commonservices.GroupField_MemberCount,
						Value: "0",
					})
				}
			}
		}
		return errs.IMErrorCode_SUCCESS, grpInfo
	} else {
		groupInfo, exist := GetGroupInfoFromCache(ctx, appkey, req.GroupId)
		if exist && groupInfo != nil {
			fields := make(map[string]string)
			for k, v := range groupInfo.ExtFields {
				fields[k] = v
			}
			var memberCount int32 = 0
			memberContainer, exist := GetGroupMembersFromCache(ctx, appkey, req.GroupId)
			if exist && memberContainer != nil {
				memberCount = int32(memberContainer.GroupMemberCount())
			}
			return errs.IMErrorCode_SUCCESS, &pbobjs.GroupInfo{
				GroupId:       groupInfo.GroupId,
				GroupName:     groupInfo.GroupName,
				GroupPortrait: groupInfo.GroupPortrait,
				IsMute:        groupInfo.IsMute,
				UpdatedTime:   groupInfo.UpdatedTime,
				ExtFields:     commonservices.Map2KvItems(fields),
				Settings:      commonservices.Map2KvItems(groupInfo.SettingFields),
				MemberCount:   memberCount,
			}
		}
	}
	return errs.IMErrorCode_GROUP_GROUPNOTEXIST, nil
}

func UpdGroupInfo(ctx context.Context, groupInfo *pbobjs.GroupInfo) errs.IMErrorCode {
	groupInfo.GroupName = tools.TruncateText(groupInfo.GroupName, 64)
	appkey := bases.GetAppKeyFromCtx(ctx)
	rvCache := false
	groupId := groupInfo.GroupId
	dao := dbs.GroupDao{}
	err := dao.UpdateGrpName(appkey, groupId, groupInfo.GroupName, groupInfo.GroupPortrait)
	if err == nil {
		rvCache = rvCache || true
	}
	extDao := dbs.GroupExtDao{}
	for _, ext := range groupInfo.ExtFields {
		itemKey, itemValue := ext.Key, ext.Value
		extDao.Upsert(dbs.GroupExtDao{
			AppKey:    appkey,
			GroupId:   groupId,
			ItemKey:   itemKey,
			ItemValue: itemValue,
			ItemType:  int(commonservices.AttItemType_Att),
		})
		rvCache = rvCache || true
	}
	for _, setting := range groupInfo.Settings {
		extDao.Upsert(dbs.GroupExtDao{
			AppKey:    appkey,
			GroupId:   groupId,
			ItemKey:   setting.Key,
			ItemValue: setting.Value,
			ItemType:  int(commonservices.AttItemType_Setting),
		})
		rvCache = rvCache || true
	}
	if rvCache { //remove cache
		key := getGroupKey(appkey, groupId)
		groupInfoCache.Remove(key)
	}
	return errs.IMErrorCode_SUCCESS
}

func SetGroupSettings(ctx context.Context, groupId string, settings []*pbobjs.KvItem) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.GroupExtDao{}

	_, exist := GetGroupInfoFromCache(ctx, appkey, groupId)
	for _, setting := range settings {
		err := dao.Upsert(dbs.GroupExtDao{
			AppKey:    appkey,
			GroupId:   groupId,
			ItemKey:   setting.Key,
			ItemValue: setting.Value,
			ItemType:  int(commonservices.AttItemType_Setting),
		})
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
	}
	if exist && len(settings) > 0 {
		key := getGroupKey(appkey, groupId)
		lock := groupLocks.GetLocks(key.appKey, key.groupID)
		lock.Lock()
		if groupInfo, cached := currentGroupInfoFromCache(key); cached {
			groupInfo = cloneGroupInfo(groupInfo)
			if groupInfo.SettingFields == nil {
				groupInfo.SettingFields = make(map[string]string)
			}
			if groupInfo.Settings == nil {
				groupInfo.Settings = &commonservices.GroupSettings{}
			}
			for _, setting := range settings {
				groupInfo.SettingFields[setting.Key] = setting.Value
			}
			commonservices.FillObjFieldWithIgnore(groupInfo.Settings, groupInfo.SettingFields, true)
			groupInfoCache.Add(key, groupInfo)
		}
		lock.Unlock()
	}
	return errs.IMErrorCode_SUCCESS
}

func GetGroupSettings(ctx context.Context, groupId string) (errs.IMErrorCode, *pbobjs.GroupInfo) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.GroupExtDao{}
	groupExts, err := dao.QryExtFields(appkey, groupId)
	if err != nil {
		return errs.IMErrorCode_API_INTERNAL_RESP_FAIL, nil
	}

	groupInfo := &pbobjs.GroupInfo{
		GroupId:  groupId,
		Settings: []*pbobjs.KvItem{},
	}

	for _, item := range groupExts {
		groupInfo.ExtFields = append(groupInfo.ExtFields, &pbobjs.KvItem{
			Key:     item.ItemKey,
			Value:   item.ItemValue,
			UpdTime: item.UpdatedTime.UnixMilli(),
		})
	}

	return errs.IMErrorCode_SUCCESS, groupInfo

}
