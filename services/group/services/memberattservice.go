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
	memberAttCache = caches.NewLruCacheWithAddReadTimeout(100000, nil, 8*time.Minute, 10*time.Minute)
	memberLocks = tools.NewSegmentatedLocks(512)
}

var notExistGrpMemberAtts MemberAtts = MemberAtts{
	ExtFields: make(map[string]string),
	Settings:  &commonservices.GrpMemberSettings{},
}

type MemberAtts struct {
	MemberId    string
	ExtFields   map[string]string
	Settings    *commonservices.GrpMemberSettings
	UpdatedTime int64
}

func AddGrpMemberAtts2Cache(ctx context.Context, appkey, groupId, memberId string, atts MemberAtts) {
	key := getGrpMemberKey(appkey, groupId, memberId)
	memberAttCache.Add(key, &atts)
}

func GetGrpMemberAttsFromCache(ctx context.Context, appkey, groupId, memberId string) (*MemberAtts, bool) {
	key := getGrpMemberKey(appkey, groupId, memberId)
	if val, exist := memberAttCache.Get(key); exist {
		memberAtts := val.(*MemberAtts)
		if memberAtts == &notExistGrpMemberAtts {
			return nil, false
		} else {
			return memberAtts, true
		}
	} else {
		l := memberLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()

		if val, exist := memberAttCache.Get(key); exist {
			memberAtts := val.(*MemberAtts)
			if memberAtts == &notExistGrpMemberAtts {
				return nil, false
			} else {
				return memberAtts, true
			}
		} else {
			memberAtts := getGrpMemberAttsFromDb(ctx, appkey, groupId, memberId)
			memberAttCache.Add(key, memberAtts)
			return memberAtts, memberAtts != &notExistGrpMemberAtts
		}
	}
}

func getGrpMemberAttsFromDb(ctx context.Context, appkey, groupId, memberId string) *MemberAtts {
	dao := dbs.GroupMemberExtDao{}
	exts, err := dao.QryExtFields(appkey, groupId, memberId)
	if err != nil || len(exts) <= 0 {
		return &notExistGrpMemberAtts
	} else {
		ret := &MemberAtts{
			MemberId:  memberId,
			ExtFields: make(map[string]string),
			Settings:  &commonservices.GrpMemberSettings{},
		}
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
			} else if ext.ItemType == int(commonservices.AttItemType_Att) {
				ret.ExtFields[ext.ItemKey] = ext.ItemValue
			}
		}
		commonservices.FillObjField(ret.Settings, valMap)
		return ret
	}
}

func getGrpMemberKey(appkey, groupId, memberId string) string {
	return strings.Join([]string{appkey, groupId, memberId}, "_")
}
