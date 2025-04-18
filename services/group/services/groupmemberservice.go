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
	"time"
)

var groupMembersCache *caches.LruCache

func init() {
	groupMembersCache = caches.NewLruCacheWithAddReadTimeout("grpmembers_cache", 10000, nil, 8*time.Minute, 10*time.Minute)
}

func AddGroupMembers(ctx context.Context, groupId, groupName, groupPortrait string, memberIds []string, extFields []*pbobjs.KvItem, settings []*pbobjs.KvItem) errs.IMErrorCode {
	if groupId == "" {
		return errs.IMErrorCode_API_PARAM_REQUIRED
	}
	groupName = tools.TruncateText(groupName, 64)
	appKey := bases.GetAppKeyFromCtx(ctx)
	//check group exist
	grpInfo, exist := GetGroupInfoFromCache(ctx, appKey, groupId)
	if !exist {
		//create group
		if groupName == "" {
			return errs.IMErrorCode_API_PARAM_REQUIRED
		}
		grpDao := dbs.GroupDao{}
		err := grpDao.Create(dbs.GroupDao{
			GroupId:       groupId,
			GroupName:     groupName,
			GroupPortrait: groupPortrait,
			UpdatedTime:   time.Now(),
			CreatedTime:   time.Now(),
			AppKey:        appKey,
		})
		//group ext
		extDao := dbs.GroupExtDao{}
		for _, ext := range extFields {
			itemKey, itemValue := ext.Key, ext.Value
			extDao.Upsert(dbs.GroupExtDao{
				AppKey:    appKey,
				GroupId:   groupId,
				ItemKey:   itemKey,
				ItemValue: itemValue,
			})
		}
		for _, setting := range settings {
			itemKey, itemValue := setting.Key, setting.Value
			extDao.Upsert(dbs.GroupExtDao{
				AppKey:    appKey,
				GroupId:   groupId,
				ItemKey:   itemKey,
				ItemValue: itemValue,
				ItemType:  int(commonservices.AttItemType_Setting),
			})
		}
		if err == nil {
			AddGroupInfo2Cache(ctx, appKey, groupId, &GroupInfo{
				GroupId:       groupId,
				GroupPortrait: groupPortrait,
				GroupName:     groupName,
				IsMute:        int32(0),
				UpdatedTime:   time.Now().UnixMilli(),
				ExtFields:     commonservices.Kvitems2Map(extFields),
				SettingFields: commonservices.Kvitems2Map(settings),
			})
		}
	} else {
		grpName := ""
		grpPortrait := ""
		if grpInfo.GroupName == "" && groupName != "" {
			grpName = groupName
			grpInfo.GroupName = groupName
		}
		if grpInfo.GroupPortrait == "" && groupPortrait != "" {
			grpPortrait = groupPortrait
			grpInfo.GroupPortrait = groupPortrait
		}
		if grpName != "" || grpPortrait != "" {
			grpInfo.GroupName = groupName
			grpDao := dbs.GroupDao{}
			grpDao.UpdateGrpName(appKey, groupId, grpName, grpPortrait)
		}
	}

	needAddMemberIds := []string{}
	memberContainer, exist := GetGroupMembersFromCache(ctx, appKey, groupId)
	currentMemberCount := memberContainer.GroupMemberCount()
	if exist {
		memberMap := memberContainer.GetMemberMap()
		for _, memberId := range memberIds {
			if _, exist := memberMap[memberId]; !exist {
				needAddMemberIds = append(needAddMemberIds, memberId)
			}
		}
	} else {
		needAddMemberIds = append(needAddMemberIds, memberIds...)
	}
	//check group member count limit
	appInfo, exist := commonservices.GetAppInfo(appKey)
	if exist && appInfo != nil {
		needAddCount := len(needAddMemberIds)
		if (needAddCount + currentMemberCount) > appInfo.MaxGrpMemberCount {
			return errs.IMErrorCode_GROUP_GROUPMEMBERCOUNTEXCEED
		}
	} else {
		return errs.IMErrorCode_API_APP_NOT_EXISTED
	}
	groupInfo, exist := GetGroupInfoFromCache(ctx, appKey, groupId)
	grpHideGrpMsg := false
	if exist && groupInfo != nil && groupInfo.Settings != nil {
		grpHideGrpMsg = groupInfo.Settings.HideGrpMsg
	}
	for _, memberId := range needAddMemberIds {
		memberContainer.AddMember(GroupMember{
			MemberId:    memberId,
			CreatedTime: time.Now().UnixMilli(),
		})
		AddGrpMemberAtts2Cache(ctx, appKey, groupId, memberId, MemberAtts{
			MemberId: memberId,
			Settings: &commonservices.GrpMemberSettings{
				HideGrpMsg: grpHideGrpMsg,
			},
			ExtFields:     make(map[string]string),
			SettingFields: make(map[string]string),
		})
	}
	//添加群成员
	memberDao := dbs.GroupMemberDao{}
	memberExtDao := dbs.GroupMemberExtDao{}

	members := []dbs.GroupMemberDao{}
	memberExts := []dbs.GroupMemberExtDao{}
	for _, memberId := range needAddMemberIds {
		members = append(members, dbs.GroupMemberDao{
			AppKey:   appKey,
			GroupId:  groupId,
			MemberId: memberId,
		})
		memberExts = append(memberExts, dbs.GroupMemberExtDao{
			AppKey:    appKey,
			GroupId:   groupId,
			MemberId:  memberId,
			ItemKey:   string(commonservices.AttItemKey_HideGrpMsg),
			ItemValue: tools.Bool2String(grpHideGrpMsg),
			ItemType:  int(commonservices.AttItemType_Setting),
		})
	}
	err := memberDao.BatchCreate(members)
	if err == nil {
		existMemberIds := []string{}
		memberMap := memberContainer.GetMemberMap()
		for memberId := range memberMap {
			existMemberIds = append(existMemberIds, memberId)
		}
		GenerateGroupSnapshot(appKey, groupId, existMemberIds)
	}
	memberExtDao.BatchCreate(memberExts)
	return errs.IMErrorCode_SUCCESS
}

func DelGroupMembers(ctx context.Context, groupId string, memberIds []string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	isAffected := false
	memberContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	if exist {
		succ := memberContainer.DelMembers(memberIds)
		if succ {
			isAffected = true
		}
	}
	//update for db
	memberDao := dbs.GroupMemberDao{}
	err := memberDao.BatchDelete(appkey, groupId, memberIds)
	if err == nil && isAffected {
		existMemberIds := []string{}
		memberMap := memberContainer.GetMemberMap()
		for memberId := range memberMap {
			existMemberIds = append(existMemberIds, memberId)
		}
		GenerateGroupSnapshot(appkey, groupId, existMemberIds)
		//delete ext from db
		memberExtDao := dbs.GroupMemberExtDao{}
		memberExtDao.BatchDelete(appkey, groupId, memberIds)
	}
}

func clearMembersFromCache(ctx context.Context, appkey, groupId string) {
	memberContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	if exist {
		memberContainer.ClearMembers()
	}
}

func GetGroupMembersFromCache(ctx context.Context, appkey, groupId string) (*GroupMemberContainer, bool) {
	key := getGroupKey(appkey, groupId)
	if groupContainer, exist := groupMembersCache.Get(key); exist {
		return groupContainer.(*GroupMemberContainer), true
	} else {
		lock := groupLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if groupContainer, exist := groupMembersCache.Get(key); exist {
			return groupContainer.(*GroupMemberContainer), true
		} else {
			members := GetGroupMembersFromDb(ctx, appkey, groupId)
			groupContainer := &GroupMemberContainer{
				Appkey:  appkey,
				GroupId: groupId,
				Members: members,
			}
			groupMembersCache.Add(key, groupContainer)
			return groupContainer, true
		}
	}
}

func GetGroupMembersFromDb(ctx context.Context, appkey, groupId string) map[string]*GroupMember {
	memberDao := dbs.GroupMemberDao{}
	var startId int64 = 0
	var limit int64 = 10000
	members := map[string]*GroupMember{}
	// grpInfo := getGroupInfoFromDb(ctx, appkey, groupId)
	// hideGrpMsg := "0"
	// if grpInfo != nil && grpInfo.Settings != nil {
	// 	if grpInfo.Settings.HideGrpMsg {
	// 		hideGrpMsg = "1"
	// 	}
	// }
	for {
		dbMembers, err := memberDao.QueryMembers(appkey, groupId, startId, limit)
		if err == nil {
			for _, dbMember := range dbMembers {
				grpMember := &GroupMember{
					MemberId:    dbMember.MemberId,
					IsMute:      dbMember.IsMute,
					MuteEndAt:   dbMember.MuteEndAt,
					IsAllow:     dbMember.IsAllow,
					CreatedTime: dbMember.CreatedTime.UnixMilli(),
					// Settings:    &commonservices.GrpMemberSettings{},
				}
				// valMap := make(map[string]string)
				// valMap[string(commonservices.AttItemKey_HideGrpMsg)] = hideGrpMsg
				// memberExtDao := dbs.GroupMemberExtDao{}
				// exts, err := memberExtDao.QryExtFields(appkey, groupId, grpMember.MemberId)
				// if err == nil {
				// 	for _, ext := range exts {
				// 		if ext.ItemType == int(commonservices.AttItemType_Setting) {
				// 			valMap[ext.ItemKey] = ext.ItemValue
				// 		}
				// 	}
				// }
				// commonservices.FillObjField(grpMember.Settings, valMap)
				members[dbMember.MemberId] = grpMember
				if startId < dbMember.ID {
					startId = dbMember.ID
				}
			}
			if len(dbMembers) < int(limit) {
				break
			}
		} else {
			break
		}
	}
	return members
}

func SetGroupMemberMute(ctx context.Context, req *pbobjs.GroupMemberMuteReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//update db
	dao := dbs.GroupMemberDao{}
	dao.UpdateMute(appkey, req.GroupId, int(req.IsMute), req.MemberIds, req.MuteEndAt)
	//update cache
	container, exist := GetGroupMembersFromCache(ctx, appkey, req.GroupId)
	if exist {
		container.SetMemberMute(req.IsMute, req.MemberIds, req.MuteEndAt)
	}
	return errs.IMErrorCode_SUCCESS
}

func SetGroupMemberAllow(ctx context.Context, req *pbobjs.GroupMemberAllowReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//update db
	dao := dbs.GroupMemberDao{}
	dao.UpdateAllow(appkey, req.GroupId, int(req.IsAllow), req.MemberIds)
	//update cache
	container, exist := GetGroupMembersFromCache(ctx, appkey, req.GroupId)
	if exist {
		container.SetMemberAllow(req.IsAllow, req.MemberIds)
	}
	return errs.IMErrorCode_SUCCESS
}

func SetGroupMemberSettings(ctx context.Context, groupId string, req *pbobjs.GroupMember) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.GroupMemberExtDao{}
	memberAtts := GetGrpMemberAttsFromCache(ctx, appkey, groupId, req.MemberId)
	if memberAtts != nil {
		for _, setting := range req.Settings {
			memberAtts.SetMemberSetting(setting.Key, setting.Value)
			dao.Upsert(appkey, groupId, req.MemberId, setting.Key, setting.Value, int(commonservices.AttItemType_Setting))
		}
		for _, ext := range req.ExtFields {
			memberAtts.SetMemberExt(ext.Key, ext.Value)
			dao.Upsert(appkey, groupId, req.MemberId, ext.Key, ext.Value, int(commonservices.AttItemType_Att))
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func QryGroupMembersByIds(ctx context.Context, req *pbobjs.GroupMembersReq) (errs.IMErrorCode, *pbobjs.GroupMembersResp) {
	resp := &pbobjs.GroupMembersResp{
		Items: []*pbobjs.GroupMember{},
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	container, exist := GetGroupMembersFromCache(ctx, appkey, req.GroupId)
	if exist {
		memberMap := container.GetMemberMap()
		curr := time.Now().UnixMilli()
		for _, memberId := range req.MemberIds {
			if member, ok := memberMap[memberId]; ok {
				isMute := member.IsMute
				if isMute > 0 && member.MuteEndAt < curr {
					isMute = 0
				}
				member := &pbobjs.GroupMember{
					MemberId:   memberId,
					IsMute:     int32(isMute),
					IsAllow:    int32(member.IsAllow),
					MemberType: member.MemberType,
				}
				//add extfields
				memberAtts := GetGrpMemberAttsFromCache(ctx, appkey, req.GroupId, memberId)
				if memberAtts != nil {
					member.ExtFields = commonservices.Map2KvItems(memberAtts.ExtFields)
					member.Settings = commonservices.Map2KvItems(memberAtts.SettingFields)
				}
				resp.Items = append(resp.Items, member)
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, resp
}

func QryGroupMembers(ctx context.Context, req *pbobjs.QryGroupMembersReq) (errs.IMErrorCode, *pbobjs.GroupMembersResp) {
	resp := &pbobjs.GroupMembersResp{
		Items:  []*pbobjs.GroupMember{},
		Offset: "",
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	startId, err := tools.DecodeInt(req.Offset)
	if err != nil {
		startId = 0
	}
	dao := dbs.GroupMemberDao{}
	members, err := dao.QueryMembers(appkey, req.GroupId, startId, req.Limit)
	if err == nil {
		curr := time.Now().UnixMilli()
		memberIds := []string{}
		for _, member := range members {
			memberIds = append(memberIds, member.MemberId)
			isMute := member.IsMute
			if isMute > 0 && member.MuteEndAt < curr {
				isMute = 0
			}
			resp.Items = append(resp.Items, &pbobjs.GroupMember{
				MemberId:   member.MemberId,
				IsMute:     int32(isMute),
				IsAllow:    int32(member.IsAllow),
				MemberType: pbobjs.UserType(member.MemberType),
			})
			offset, err := tools.EncodeInt(member.ID)
			if err != nil {
				offset = ""
			}
			resp.Offset = offset
		}
		//add member extfields
		memberExtDao := dbs.GroupMemberExtDao{}
		extMap, err := memberExtDao.QryExtFieldsByMemberIds(appkey, req.GroupId, memberIds)
		if err == nil {
			for _, member := range resp.Items {
				if exts, exist := extMap[member.MemberId]; exist {
					if len(exts) > 0 {
						member.ExtFields = []*pbobjs.KvItem{}
						member.Settings = []*pbobjs.KvItem{}
						for _, ext := range exts {
							if ext.ItemType == int(commonservices.AttItemType_Att) {
								member.ExtFields = append(member.ExtFields, &pbobjs.KvItem{
									Key:     ext.ItemKey,
									Value:   ext.ItemValue,
									UpdTime: ext.UpdatedTime.UnixMilli(),
								})
							} else if ext.ItemType == int(commonservices.AttItemType_Setting) {
								member.Settings = append(member.Settings, &pbobjs.KvItem{
									Key:     ext.ItemKey,
									Value:   ext.ItemValue,
									UpdTime: ext.UpdatedTime.UnixMilli(),
								})
							}
						}
					}
				}
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, resp
}

func CheckGroupMembers(ctx context.Context, req *pbobjs.CheckGroupMembersReq) (errs.IMErrorCode, *pbobjs.CheckGroupMembersResp) {
	resp := &pbobjs.CheckGroupMembersResp{
		MemberIdMap: make(map[string]int64),
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	container, exist := GetGroupMembersFromCache(ctx, appkey, req.GroupId)
	if exist {
		memberMap := container.CheckGroupMembers(req.MemberIds)
		for _, memberId := range req.MemberIds {
			addedTime, exist := memberMap[memberId]
			if exist {
				resp.MemberIdMap[memberId] = addedTime
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, resp
}

func QryMemberSettings(ctx context.Context, groupId string, memberId string) (errs.IMErrorCode, *pbobjs.QryGrpMemberSettingsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	resp := &pbobjs.QryGrpMemberSettingsResp{}
	memberContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	if exist {
		member := memberContainer.GetMember(memberId)
		if member != nil {
			resp.IsMember = true
			resp.JoinTime = member.CreatedTime
			memberAtts := GetGrpMemberAttsFromCache(ctx, appkey, groupId, memberId)
			if memberAtts != nil {
				resp.MemberSettings = memberAtts.GetMemberSettings()
				resp.MemberExts = memberAtts.GetMemberExts()
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, resp
}

type GroupMember struct {
	MemberId    string
	IsMute      int
	MuteEndAt   int64
	IsAllow     int
	CreatedTime int64 //join time
	MemberType  pbobjs.UserType
}

type GroupMemberContainer struct {
	Appkey  string
	GroupId string
	Members map[string]*GroupMember
}

func (container *GroupMemberContainer) AddMember(member GroupMember) {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	container.Members[member.MemberId] = &member
}

func (container *GroupMemberContainer) DelMembers(memberIds []string) bool {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	ret := false
	for _, memberId := range memberIds {
		if _, exist := container.Members[memberId]; exist {
			delete(container.Members, memberId)
			ret = true
		}
	}
	return ret
}

func (container *GroupMemberContainer) ClearMembers() bool {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if len(container.Members) > 0 {
		container.Members = make(map[string]*GroupMember)
		return true
	} else {
		return false
	}
}

func (container *GroupMemberContainer) GetMemberMap() map[string]*GroupMember {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	memberMap := map[string]*GroupMember{}
	for k, v := range container.Members {
		memberMap[k] = v
	}
	return memberMap
}

func (container *GroupMemberContainer) GetMember(memberId string) *GroupMember {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	if member, exist := container.Members[memberId]; exist {
		return member
	} else {
		return nil
	}
}

func (container *GroupMemberContainer) SetMemberMute(isMute int32, memberIds []string, muteEndAt int64) {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	for _, memberId := range memberIds {
		if member, exist := container.Members[memberId]; exist {
			member.IsMute = int(isMute)
			if isMute == 0 {
				member.MuteEndAt = 0
			} else {
				member.MuteEndAt = muteEndAt
			}
		}
	}
}

func (container *GroupMemberContainer) SetMemberAllow(isAllow int32, memberIds []string) {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	for _, memberId := range memberIds {
		if member, exist := container.Members[memberId]; exist {
			member.IsAllow = int(isAllow)
		}
	}
}

func (container *GroupMemberContainer) GroupMemberCount() int {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	return len(container.Members)
}

func (container *GroupMemberContainer) SetGroupMembers(members map[string]*GroupMember) {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	container.Members = members
}

func (container *GroupMemberContainer) CheckGroupMembers(memberIds []string) map[string]int64 {
	key := getGroupKey(container.Appkey, container.GroupId)
	lock := groupLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	ret := map[string]int64{}
	for _, memberId := range memberIds {
		if member, exist := container.Members[memberId]; exist {
			ret[memberId] = member.CreatedTime
		}
	}
	return ret
}
