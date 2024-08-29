package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/group/dbs"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var groupInfoCache *caches.LruCache
var groupLocks *tools.SegmentatedLocks

func init() {
	groupInfoCache = caches.NewLruCacheWithAddReadTimeout(10000, nil, 10*time.Minute, 10*time.Minute)
	groupLocks = tools.NewSegmentatedLocks(128)
}

type GroupInfo struct {
	GroupId       string
	GroupName     string
	GroupPortrait string
	IsMute        int32
	ExtFields     map[string]string
	UpdatedTime   int64
	Settings      *commonservices.GroupSettings
}

var notExistGroup GroupInfo = GroupInfo{
	ExtFields: make(map[string]string),
	Settings:  &commonservices.GroupSettings{},
}

func DissolveGroup(ctx context.Context, groupId string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//remove from db
	groupDao := dbs.GroupDao{}
	groupDao.Delete(appkey, groupId)
	memberDao := dbs.GroupMemberDao{}
	memberDao.DeleteByGroupId(appkey, groupId)

	//remove from cache
	DelGroupInfo(ctx, appkey, groupId)
	//dissolve group members
	clearMembersFromCache(ctx, appkey, groupId)
}

func DelGroupInfo(ctx context.Context, appkey, groupId string) {
	key := getGroupKey(appkey, groupId)
	groupInfoCache.Add(key, &notExistGroup)
}

func getGroupKey(appkey, groupId string) string {
	return strings.Join([]string{appkey, groupId}, "_")
}
func AddGroupInfo2Cache(ctx context.Context, appkey, groupId string, grpInfo *GroupInfo) {
	key := getGroupKey(appkey, groupId)
	groupInfoCache.Add(key, grpInfo)
}
func GetGroupInfoFromCache(ctx context.Context, appkey, groupId string) (*GroupInfo, bool) {
	key := getGroupKey(appkey, groupId)
	if val, exist := groupInfoCache.Get(key); exist {
		groupInfo := val.(*GroupInfo)
		if groupInfo == &notExistGroup {
			return nil, false
		} else {
			return groupInfo, true
		}
	} else {
		l := groupLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()

		if val, exist := groupInfoCache.Get(key); exist {
			groupInfo := val.(*GroupInfo)
			if groupInfo == &notExistGroup {
				return nil, false
			} else {
				return groupInfo, true
			}
		} else {
			groupInfo := getGroupInfoFromDb(ctx, appkey, groupId)
			if groupInfo != nil {
				groupInfoCache.Add(key, groupInfo)
				if groupInfo == &notExistGroup {
					return nil, false
				} else {
					return groupInfo, true
				}
			}
			return nil, false
		}
	}
}

func getGroupInfoFromDb(ctx context.Context, appkey, groupId string) *GroupInfo {
	dao := dbs.GroupDao{}
	dbGroup, err := dao.FindById(appkey, groupId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &notExistGroup
		}
		return nil
	} else {
		groupInfo := &GroupInfo{
			GroupId:       dbGroup.GroupId,
			GroupPortrait: dbGroup.GroupPortrait,
			GroupName:     dbGroup.GroupName,
			IsMute:        int32(dbGroup.IsMute),
			UpdatedTime:   dbGroup.UpdatedTime.UnixMilli(),
			ExtFields:     map[string]string{},
			Settings:      &commonservices.GroupSettings{},
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
	groupInfo, exist := GetGroupInfoFromCache(ctx, appkey, req.GroupId)
	if exist {
		groupInfo.IsMute = req.IsMute
	}
	return errs.IMErrorCode_SUCCESS
}

func QryGroupInfo(ctx context.Context, req *pbobjs.GroupInfoReq) (errs.IMErrorCode, *pbobjs.GroupInfo) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	if len(req.CareFields) > 0 {
		grpInfo := &pbobjs.GroupInfo{
			GroupId:   req.GroupId,
			ExtFields: []*pbobjs.KvItem{},
		}
		memberContainer, exist := GetGroupMembersFromCache(ctx, appkey, commonservices.GroupField_MemberCount)
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
			return errs.IMErrorCode_SUCCESS, &pbobjs.GroupInfo{
				GroupId:       groupInfo.GroupId,
				GroupName:     groupInfo.GroupName,
				GroupPortrait: groupInfo.GroupPortrait,
				IsMute:        groupInfo.IsMute,
				UpdatedTime:   groupInfo.UpdatedTime,
				ExtFields:     commonservices.Map2KvItems(fields),
			}
		}
	}
	return errs.IMErrorCode_GROUP_GROUPNOTEXIST, nil
}

func UpdGroupInfo(ctx context.Context, groupInfo *pbobjs.GroupInfo) errs.IMErrorCode {
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
		})
		rvCache = rvCache || true
	}
	if rvCache { //remove cache
		grpInfo := getGroupInfoFromDb(ctx, appkey, groupId)
		if grpInfo != nil {
			key := getGroupKey(appkey, groupId)
			l := groupLocks.GetLocks(key)
			l.Lock()
			defer l.Unlock()
			groupInfoCache.Remove(key)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func SetGroupSettings(ctx context.Context, groupId string, settings []*pbobjs.KvItem) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.GroupExtDao{}

	groupInfo, exist := GetGroupInfoFromCache(ctx, appkey, groupId)
	valMap := make(map[string]string)
	for _, setting := range settings {
		dao.Upsert(dbs.GroupExtDao{
			AppKey:    appkey,
			GroupId:   groupId,
			ItemKey:   setting.Key,
			ItemValue: setting.Value,
			ItemType:  int(commonservices.AttItemType_Setting),
		})
		valMap[setting.Key] = setting.Value
	}
	if exist && len(valMap) > 0 {
		commonservices.FillObjFieldWithIgnore(groupInfo.Settings, valMap, true)
	}
	return errs.IMErrorCode_SUCCESS
}
