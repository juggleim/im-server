package services

import (
	"context"
	"im-server/commons/caches"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/group/dbs"
	"strings"
	"time"
)

var memberAttCache *caches.LruCache
var memberLocks *tools.SegmentatedLocks

func init() {
	memberAttCache = caches.NewLruCacheWithAddReadTimeout("grpmemberatt_cache", 100000, nil, 8*time.Minute, 10*time.Minute)
	memberLocks = tools.NewSegmentatedLocks(512)
}

type MemberAtts struct {
	AppKey        string
	GroupId       string
	MemberId      string
	ExtFields     map[string]string
	Settings      *commonservices.GrpMemberSettings
	SettingFields map[string]string
	UpdatedTime   int64
}

func (member *MemberAtts) SetMemberSetting(itemKey, itemValue string) {
	key := getGrpMemberKey(member.AppKey, member.GroupId, member.MemberId)
	lock := memberLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	member.SettingFields[itemKey] = itemValue
}

func (member *MemberAtts) GetMemberSettings() map[string]string {
	key := getGrpMemberKey(member.AppKey, member.GroupId, member.MemberId)
	lock := memberLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	ret := make(map[string]string)
	for k, v := range member.SettingFields {
		ret[k] = v
	}
	return ret
}

func (member *MemberAtts) SetMemberExt(itemKey, itemValue string) {
	key := getGrpMemberKey(member.AppKey, member.GroupId, member.MemberId)
	lock := memberLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	member.ExtFields[itemKey] = itemValue
}

func (member *MemberAtts) GetMemberExts() map[string]string {
	key := getGrpMemberKey(member.AppKey, member.GroupId, member.MemberId)
	lock := memberLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	ret := make(map[string]string)
	for k, v := range member.ExtFields {
		ret[k] = v
	}
	return ret
}

func AddGrpMemberAtts2Cache(ctx context.Context, appkey, groupId, memberId string, atts MemberAtts) {
	key := getGrpMemberKey(appkey, groupId, memberId)
	memberAttCache.Add(key, &atts)
}

func GetGrpMemberAttsFromCache(ctx context.Context, appkey, groupId, memberId string) *MemberAtts {
	key := getGrpMemberKey(appkey, groupId, memberId)
	if val, exist := memberAttCache.Get(key); exist {
		memberAtts := val.(*MemberAtts)
		return memberAtts
	} else {
		l := memberLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()

		if val, exist := memberAttCache.Get(key); exist {
			memberAtts := val.(*MemberAtts)
			return memberAtts
		} else {
			memberAtts := getGrpMemberAttsFromDb(ctx, appkey, groupId, memberId)
			memberAttCache.Add(key, memberAtts)
			return memberAtts
		}
	}
}

func getGrpMemberAttsFromDb(ctx context.Context, appkey, groupId, memberId string) *MemberAtts {
	ret := &MemberAtts{
		AppKey:        appkey,
		GroupId:       groupId,
		MemberId:      memberId,
		ExtFields:     make(map[string]string),
		Settings:      &commonservices.GrpMemberSettings{},
		SettingFields: make(map[string]string),
	}
	dao := dbs.GroupMemberExtDao{}
	exts, err := dao.QryExtFields(appkey, groupId, memberId)
	if err == nil && len(exts) > 0 {
		grpInfo, exist := GetGroupInfoFromCache(ctx, appkey, groupId)
		hideGrpMsg := "0"
		if exist && grpInfo != nil && grpInfo.Settings != nil {
			if grpInfo.Settings.HideGrpMsg {
				hideGrpMsg = "1"
			}
		}
		valMap := make(map[string]string)
		valMap[string(commonservices.AttItemKey_HideGrpMsg)] = hideGrpMsg
		for _, ext := range exts {
			if ext.ItemType == int(commonservices.AttItemType_Setting) {
				valMap[ext.ItemKey] = ext.ItemValue
				ret.SettingFields[ext.ItemKey] = ext.ItemValue
			} else if ext.ItemType == int(commonservices.AttItemType_Att) {
				ret.ExtFields[ext.ItemKey] = ext.ItemValue
			}
			updTime := ext.UpdatedTime.UnixMilli()
			if updTime > ret.UpdatedTime {
				ret.UpdatedTime = updTime
			}
		}
		commonservices.FillObjField(ret.Settings, valMap)
	}
	return ret
}

func getGrpMemberKey(appkey, groupId, memberId string) string {
	return strings.Join([]string{appkey, groupId, memberId}, "_")
}
