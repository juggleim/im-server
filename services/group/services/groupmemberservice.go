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
	"im-server/services/commonservices/msgdefines"
	"im-server/services/group/dbs"
	"sync"
	"time"
)

var groupMembersCache *caches.LruCache

const (
	groupMembersCacheIdleTimeout = 8 * time.Minute
	groupMembersCacheMaxLifetime = 60 * time.Minute
	groupMembersCacheHardTTL     = 2 * time.Hour
	groupMembersRefreshJitter    = 5 * time.Minute
)

// loadGroupMembersFromDB is a variable so cache-load concurrency can be tested
// without depending on an external database.
var loadGroupMembersFromDB = GetGroupMembersFromDb

type groupMembersCacheTrace struct {
	cacheHit   bool
	load       time.Duration
	sharedWait time.Duration
}

func init() {
	groupMembersCache = caches.NewLruCacheWithAddReadTimeout(
		"grpmembers_cache", 10000, nil,
		groupMembersCacheIdleTimeout, groupMembersCacheMaxLifetime,
	)
	err := groupMembersCache.EnableStaleWhileRevalidate(caches.StaleWhileRevalidateOptions{
		SoftTTL:            groupMembersCacheMaxLifetime,
		HardTTL:            groupMembersCacheHardTTL,
		RefreshJitter:      groupMembersRefreshJitter,
		LoadTimeout:        30 * time.Second,
		MaxConcurrentLoads: 2,
		MaxPendingLoads:    16,
		ServeStaleOnError:  true,
		ApplyGuard: func(key interface{}) func() {
			cacheKey := key.(groupCacheKey)
			unlock, wait := groupMemberStateLocks.lock(cacheKey.lockKey())
			groupMemberCacheMetrics.stateLockWait.observe(wait)
			return unlock
		},
	}, func(ctx context.Context, key interface{}) (interface{}, bool, error) {
		cacheKey := key.(groupCacheKey)
		groupMemberCacheMetrics.leaderLoads.Add(1)
		loadStarted := time.Now()
		members, err := loadGroupMembersFromDB(ctx, cacheKey.appKey, cacheKey.groupID)
		groupMemberCacheMetrics.load.observe(time.Since(loadStarted))
		if err != nil {
			return nil, false, err
		}
		return &GroupMemberContainer{
			Appkey:  cacheKey.appKey,
			GroupId: cacheKey.groupID,
			Members: members,
		}, true, nil
	})
	if err != nil {
		panic(err)
	}
}

func AddGroupMembers(ctx context.Context, req *pbobjs.GroupMembersReq) (result errs.IMErrorCode) {
	started := time.Now()
	phases := groupAddPhaseDurations{}
	cacheChecked := false
	memberCacheHit := false
	groupMemberCount := 0
	addedMemberCount := 0
	defer func() {
		observeGroupAdd(ctx, req.GroupId, time.Since(started), phases, cacheChecked, memberCacheHit,
			groupMemberCount, len(req.MemberIds), addedMemberCount, result)
	}()

	groupId := req.GroupId
	groupName := req.GroupName
	groupPortrait := req.GroupPortrait
	memberIds := req.MemberIds
	extFields := req.ExtFields
	settings := req.Settings

	groupName = tools.TruncateText(groupName, 64)
	appKey := bases.GetAppKeyFromCtx(ctx)
	//check group exist
	phaseStarted := time.Now()
	grpInfo, exist := GetGroupInfoFromCache(ctx, appKey, groupId)
	if !exist {
		//create group
		grpDao := dbs.GroupDao{}
		err := grpDao.Create(dbs.GroupDao{
			GroupId:       groupId,
			GroupName:     groupName,
			GroupPortrait: groupPortrait,
			UpdatedTime:   time.Now(),
			CreatedTime:   time.Now(),
			AppKey:        appKey,
		})
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
		//group ext
		extDao := dbs.GroupExtDao{}
		for _, ext := range extFields {
			itemKey, itemValue := ext.Key, ext.Value
			err := extDao.Upsert(dbs.GroupExtDao{
				AppKey:    appKey,
				GroupId:   groupId,
				ItemKey:   itemKey,
				ItemValue: itemValue,
			})
			if err != nil {
				logs.WithContext(ctx).Error(err.Error())
			}
		}
		for _, setting := range settings {
			itemKey, itemValue := setting.Key, setting.Value
			err := extDao.Upsert(dbs.GroupExtDao{
				AppKey:    appKey,
				GroupId:   groupId,
				ItemKey:   itemKey,
				ItemValue: itemValue,
				ItemType:  int(commonservices.AttItemType_Setting),
			})
			if err != nil {
				logs.WithContext(ctx).Error(err.Error())
			}
		}
		if err == nil {
			AddGroupInfo2Cache(ctx, appKey, groupId, &GroupInfo{
				AppKey:        appKey,
				GroupId:       groupId,
				GroupPortrait: groupPortrait,
				GroupName:     groupName,
				IsMute:        int32(0),
				UpdatedTime:   time.Now().UnixMilli(),
				ExtFields:     commonservices.Kvitems2Map(extFields),
				Settings:      &commonservices.GroupSettings{},
				SettingFields: commonservices.Kvitems2Map(settings),
			})
		}
	} else {
		cacheKey := getGroupKey(appKey, groupId)
		lock := groupLocks.GetLocks(cacheKey.appKey, cacheKey.groupID)
		lock.Lock()
		if current, cached := currentGroupInfoFromCache(cacheKey); cached {
			grpInfo = current
		}
		grpInfo = cloneGroupInfo(grpInfo)
		grpName := ""
		grpPortrait := ""
		if grpInfo.GroupName == "" && groupName != "" {
			grpName = groupName
		}
		if grpInfo.GroupPortrait == "" && groupPortrait != "" {
			grpPortrait = groupPortrait
		}
		if grpName != "" || grpPortrait != "" {
			grpDao := dbs.GroupDao{}
			err := grpDao.UpdateGrpName(appKey, groupId, grpName, grpPortrait)
			if err != nil {
				logs.WithContext(ctx).Error(err.Error())
			}
			if grpName != "" {
				grpInfo.GroupName = grpName
			}
			if grpPortrait != "" {
				grpInfo.GroupPortrait = grpPortrait
			}
			groupInfoCache.Add(cacheKey, grpInfo)
		}
		lock.Unlock()
	}
	phases.groupInfo = time.Since(phaseStarted)

	appInfo, appInfoExists := commonservices.GetAppInfo(appKey)
	if !appInfoExists || appInfo == nil {
		return errs.IMErrorCode_API_APP_NOT_EXISTED
	}

	// Snapshot-enabled applications still need the complete member set. Other
	// applications use the targeted cold path below and avoid loading the whole
	// group merely to add a small number of members.
	var memberContainer *GroupMemberContainer
	if appInfo.OpenGrpSnapshot {
		phaseStarted = time.Now()
		loadedContainer, loaded, cacheTrace := getGroupMembersFromCache(ctx, appKey, groupId)
		memberCacheHit = cacheTrace.cacheHit
		cacheChecked = true
		phases.memberCache = time.Since(phaseStarted)
		phases.memberCacheLoad = cacheTrace.load
		phases.memberCacheSharedWait = cacheTrace.sharedWait
		if !loaded || loadedContainer == nil {
			return errs.IMErrorCode_API_INTERNAL_RESP_FAIL
		}
		memberContainer = loadedContainer
	}

	stateKey := getGroupKey(appKey, groupId)
	stateUnlock, stateLockWait := groupMemberStateLocks.lock(stateKey.lockKey())
	phases.memberStateLockWait = stateLockWait
	defer func() {
		if stateUnlock != nil {
			stateUnlock()
		}
	}()

	phaseStarted = time.Now()
	if cached, ok := currentGroupMembersFromCache(stateKey); ok {
		memberContainer = cached
		memberCacheHit = true
		groupMemberCacheMetrics.cacheHits.Add(1)
	} else if memberContainer == nil {
		groupMemberCacheMetrics.cacheMisses.Add(1)
	}
	if !cacheChecked || memberContainer == nil {
		cacheChecked = true
		phases.memberCache = time.Since(phaseStarted)
	}

	needAddMemberIds := []string{}
	currentMemberCount := 0
	if memberContainer != nil {
		var lockWait time.Duration
		needAddMemberIds, currentMemberCount, lockWait, phases.memberFilter = memberContainer.FilterMissingMembers(memberIds)
		phases.memberLockWait += lockWait
	} else {
		phases.coldPath = true
		var coldTrace groupMemberColdPathTrace
		var err error
		needAddMemberIds, currentMemberCount, coldTrace, err = prepareColdGroupMemberAdd(appKey, groupId, memberIds)
		phases.coldCountDB = coldTrace.countDB
		phases.coldExistingDB = coldTrace.existingDB
		if err != nil {
			phases.coldPathFailed = true
			logs.WithContext(ctx).Errorf("[G_ADD_MEMBERS_COLD_PATH_FAILED] group_id:%s\terr:%v", groupId, err)
			return errs.IMErrorCode_API_INTERNAL_RESP_FAIL
		}
	}
	addedMemberCount = len(needAddMemberIds)
	groupMemberCount = currentMemberCount + addedMemberCount
	//check group member count limit
	if groupMemberCount > appInfo.MaxGrpMemberCount {
		return errs.IMErrorCode_GROUP_GROUPMEMBERCOUNTEXCEED
	}

	groupInfo, exist := GetGroupInfoFromCache(ctx, appKey, groupId)
	grpHideGrpMsg := false
	if exist && groupInfo != nil && groupInfo.Settings != nil {
		grpHideGrpMsg = groupInfo.Settings.HideGrpMsg
	}
	phaseStarted = time.Now()
	if memberContainer != nil {
		phases.memberLockWait += memberContainer.AddMembers(needAddMemberIds, time.Now().UnixMilli())
	}
	for _, memberId := range needAddMemberIds {
		AddGrpMemberAtts2Cache(ctx, appKey, groupId, memberId, MemberAtts{
			MemberId: memberId,
			Settings: &commonservices.GrpMemberSettings{
				HideGrpMsg: grpHideGrpMsg,
			},
			ExtFields: make(map[string]string),
			SettingFields: map[string]string{
				string(commonservices.AttItemKey_HideGrpMsg): tools.Bool2String(grpHideGrpMsg),
			},
		})
	}
	phases.cacheUpdate = time.Since(phaseStarted)
	//添加群成员
	memberDao := dbs.GroupMemberDao{}
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
	if len(members) <= 0 {
		return errs.IMErrorCode_SUCCESS
	}
	phaseStarted = time.Now()
	err := memberDao.BatchCreate(members)
	phases.memberDB = time.Since(phaseStarted)
	if err != nil {
		logs.WithContext(ctx).Error(err.Error())
		if phases.coldPath {
			phases.coldPathFailed = true
			return errs.IMErrorCode_API_INTERNAL_RESP_FAIL
		}
	}
	if memberContainer != nil {
		// Replace the LRU wrapper after the DB write so a refresh that started
		// before or during this mutation cannot apply an older DB snapshot.
		groupMembersCache.Add(stateKey, memberContainer)
	}

	generateSnapshot := err == nil && appInfo.OpenGrpSnapshot && memberContainer != nil
	var snapshotMemberIDs []string
	if generateSnapshot {
		snapshotMemberIDs = memberContainer.GetMemberIds()
	}
	stateUnlock()
	stateUnlock = nil

	phaseStarted = time.Now()
	submitGroupMemberPostPersist(ctx, appKey, groupId, memberExts, generateSnapshot, snapshotMemberIDs)
	phases.postPersistSchedule = time.Since(phaseStarted)
	//add conver for new grp members
	phaseStarted = time.Now()
	if req.UndisturbType > 0 && len(needAddMemberIds) > 0 {
		convers := make([]*pbobjs.Conversation, 0, len(needAddMemberIds))
		for _, memberId := range needAddMemberIds {
			convers = append(convers, &pbobjs.Conversation{
				UserId:      memberId,
				TargetId:    req.GroupId,
				ChannelType: pbobjs.ChannelType_Group,
				Msg: &pbobjs.DownMsg{
					TargetId:    req.GroupId,
					ChannelType: pbobjs.ChannelType_Group,
					Flags:       msgdefines.SetStoreMsg(0),
				},
				SortTime:      0,
				SyncTime:      0,
				UndisturbType: req.UndisturbType,
				IsDelete:      1,
			})
		}
		bases.AsyncRpcCall(ctx, "batch_add_conver", req.GroupId, &pbobjs.BatchAddConvers{
			Convers:    convers,
			UserOnline: true,
		})
	}
	phases.conversationSchedule = time.Since(phaseStarted)
	return errs.IMErrorCode_SUCCESS
}

func DelGroupMembers(ctx context.Context, groupId string, memberIds []string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getGroupKey(appkey, groupId)
	isAffected := false
	memberContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	stateUnlock, wait := groupMemberStateLocks.lock(key.lockKey())
	groupMemberCacheMetrics.stateLockWait.observe(wait)
	if current, cached := currentGroupMembersFromCache(key); cached {
		memberContainer = current
		exist = true
	}
	if exist {
		succ := memberContainer.DelMembers(memberIds)
		if succ {
			isAffected = true
		}
		groupMembersCache.Add(key, memberContainer)
	}
	go func() {
		defer func() {
			if stateUnlock != nil {
				stateUnlock()
			}
		}()
		//update for db
		memberDao := dbs.GroupMemberDao{}
		err := memberDao.BatchDelete(appkey, groupId, memberIds)
		if memberContainer != nil {
			// Invalidate refreshes that queried the DB while the asynchronous
			// delete was still pending.
			groupMembersCache.Add(key, memberContainer)
		}
		stateUnlock()
		stateUnlock = nil
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
	}()
}

func clearMembersFromCache(ctx context.Context, appkey, groupId string) {
	key := getGroupKey(appkey, groupId)
	memberContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	stateUnlock, wait := groupMemberStateLocks.lock(key.lockKey())
	groupMemberCacheMetrics.stateLockWait.observe(wait)
	defer stateUnlock()
	if current, cached := currentGroupMembersFromCache(key); cached {
		memberContainer = current
		exist = true
	}
	if exist {
		memberContainer.ClearMembers()
		groupMembersCache.Add(key, memberContainer)
	}
}

func GetGroupMembersFromCache(ctx context.Context, appkey, groupId string) (*GroupMemberContainer, bool) {
	container, exist, _ := getGroupMembersFromCache(ctx, appkey, groupId)
	return container, exist
}

func getGroupMembersFromCache(ctx context.Context, appkey, groupId string) (*GroupMemberContainer, bool, groupMembersCacheTrace) {
	key := getGroupKey(appkey, groupId)
	wasCached := false
	if entries := groupMembersCache.SnapshotForKeys([]interface{}{key}); len(entries) > 0 {
		wasCached = entries[0].Found
	}
	started := time.Now()
	value, found, err := groupMembersCache.GetOrLoad(ctx, key)
	trace := groupMembersCacheTrace{cacheHit: wasCached}
	if wasCached {
		groupMemberCacheMetrics.cacheHits.Add(1)
	} else {
		groupMemberCacheMetrics.cacheMisses.Add(1)
		trace.load = time.Since(started)
	}
	if err != nil {
		return nil, false, trace
	}
	if !found || value == nil {
		return nil, false, trace
	}
	return value.(*GroupMemberContainer), true, trace
}

func currentGroupMembersFromCache(key groupCacheKey) (*GroupMemberContainer, bool) {
	entries := groupMembersCache.SnapshotForKeys([]interface{}{key})
	if len(entries) == 0 || !entries[0].Found || entries[0].Value == nil {
		return nil, false
	}
	return entries[0].Value.(*GroupMemberContainer), true
}

func GetGroupMembersFromDb(ctx context.Context, appkey, groupId string) (map[string]*GroupMember, error) {
	loadStarted := time.Now()
	memberDao := dbs.GroupMemberDao{}
	var startId int64 = 0
	var limit int64 = 10000
	members := map[string]*GroupMember{}
	pageCount := 0
	rowCount := 0
	var totalQueryDuration time.Duration
	var totalBuildDuration time.Duration
	var loadErr error
	// grpInfo := getGroupInfoFromDb(ctx, appkey, groupId)
	// hideGrpMsg := "0"
	// if grpInfo != nil && grpInfo.Settings != nil {
	// 	if grpInfo.Settings.HideGrpMsg {
	// 		hideGrpMsg = "1"
	// 	}
	// }
	for {
		queryStarted := time.Now()
		dbMembers, err := memberDao.QueryMembers(appkey, groupId, startId, limit)
		queryDuration := time.Since(queryStarted)
		totalQueryDuration += queryDuration
		if err == nil {
			buildStarted := time.Now()
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
			buildDuration := time.Since(buildStarted)
			totalBuildDuration += buildDuration
			pageCount++
			rowCount += len(dbMembers)
			observeGroupMemberDBPage(queryDuration, buildDuration, len(dbMembers), false)
			if len(dbMembers) < int(limit) {
				break
			}
		} else {
			observeGroupMemberDBPage(queryDuration, 0, 0, true)
			logs.WithContext(ctx).Error(err.Error())
			loadErr = err
			break
		}
	}
	observeGroupMemberDBLoad(pageCount, rowCount)
	logSlowGroupMemberDBLoad(ctx, groupId, time.Since(loadStarted), totalQueryDuration, totalBuildDuration, pageCount, rowCount)
	return members, loadErr
}

func SetGroupMemberMute(ctx context.Context, req *pbobjs.GroupMemberMuteReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//update db
	dao := dbs.GroupMemberDao{}
	dao.UpdateMute(appkey, req.GroupId, int(req.IsMute), req.MemberIds, req.MuteEndAt)
	//update cache
	key := getGroupKey(appkey, req.GroupId)
	container, exist := GetGroupMembersFromCache(ctx, appkey, req.GroupId)
	stateUnlock, wait := groupMemberStateLocks.lock(key.lockKey())
	groupMemberCacheMetrics.stateLockWait.observe(wait)
	defer stateUnlock()
	if current, cached := currentGroupMembersFromCache(key); cached {
		container = current
		exist = true
	}
	if exist {
		container.SetMemberMute(req.IsMute, req.MemberIds, req.MuteEndAt)
		groupMembersCache.Add(key, container)
	}
	return errs.IMErrorCode_SUCCESS
}

func SetGroupMemberAllow(ctx context.Context, req *pbobjs.GroupMemberAllowReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//update db
	dao := dbs.GroupMemberDao{}
	dao.UpdateAllow(appkey, req.GroupId, int(req.IsAllow), req.MemberIds)
	//update cache
	key := getGroupKey(appkey, req.GroupId)
	container, exist := GetGroupMembersFromCache(ctx, appkey, req.GroupId)
	stateUnlock, wait := groupMemberStateLocks.lock(key.lockKey())
	groupMemberCacheMetrics.stateLockWait.observe(wait)
	defer stateUnlock()
	if current, cached := currentGroupMembersFromCache(key); cached {
		container = current
		exist = true
	}
	if exist {
		container.SetMemberAllow(req.IsAllow, req.MemberIds)
		groupMembersCache.Add(key, container)
	}
	return errs.IMErrorCode_SUCCESS
}

func SetGroupMemberSettings(ctx context.Context, groupId string, req *pbobjs.GroupMember) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.GroupMemberExtDao{}
	rvCache := false
	for _, setting := range req.Settings {
		err := dao.Upsert(appkey, groupId, req.MemberId, setting.Key, setting.Value, int(commonservices.AttItemType_Setting))
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
		rvCache = rvCache || true
	}
	for _, extField := range req.ExtFields {
		err := dao.Upsert(appkey, groupId, req.MemberId, extField.Key, extField.Value, int(commonservices.AttItemType_Att))
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}
		rvCache = rvCache || true
	}
	if rvCache {
		RemoveGrpMemberAttsFromCache(ctx, appkey, groupId, req.MemberId)
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
	mu      sync.RWMutex
}

func (container *GroupMemberContainer) AddMember(member GroupMember) {
	container.mu.Lock()
	defer container.mu.Unlock()
	container.Members[member.MemberId] = &member
}

// AddMembers updates a batch while holding this group's lock once and returns
// the time spent waiting for the lock.
func (container *GroupMemberContainer) AddMembers(memberIds []string, createdTime int64) time.Duration {
	if len(memberIds) == 0 {
		return 0
	}
	waitStarted := time.Now()
	container.mu.Lock()
	lockWait := time.Since(waitStarted)
	defer container.mu.Unlock()
	for _, memberID := range memberIds {
		container.Members[memberID] = &GroupMember{
			MemberId:    memberID,
			CreatedTime: createdTime,
		}
	}
	return lockWait
}

// FilterMissingMembers checks only the requested IDs instead of copying the
// complete member map. The returned IDs are de-duplicated in request order;
// lock wait and scan time are reported separately for diagnostics.
func (container *GroupMemberContainer) FilterMissingMembers(memberIds []string) ([]string, int, time.Duration, time.Duration) {
	waitStarted := time.Now()
	container.mu.RLock()
	lockWait := time.Since(waitStarted)
	defer container.mu.RUnlock()

	scanStarted := time.Now()
	missing := make([]string, 0, len(memberIds))
	seen := make(map[string]struct{}, len(memberIds))
	for _, memberID := range memberIds {
		if _, duplicate := seen[memberID]; duplicate {
			continue
		}
		seen[memberID] = struct{}{}
		if _, exists := container.Members[memberID]; !exists {
			missing = append(missing, memberID)
		}
	}
	return missing, len(container.Members), lockWait, time.Since(scanStarted)
}

func (container *GroupMemberContainer) GetMemberIds() []string {
	container.mu.RLock()
	defer container.mu.RUnlock()

	memberIDs := make([]string, 0, len(container.Members))
	for memberID := range container.Members {
		memberIDs = append(memberIDs, memberID)
	}
	return memberIDs
}

func (container *GroupMemberContainer) DelMembers(memberIds []string) bool {
	container.mu.Lock()
	defer container.mu.Unlock()
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
	container.mu.Lock()
	defer container.mu.Unlock()
	if len(container.Members) > 0 {
		container.Members = make(map[string]*GroupMember)
		return true
	} else {
		return false
	}
}

func (container *GroupMemberContainer) GetMemberMap() map[string]*GroupMember {
	container.mu.RLock()
	defer container.mu.RUnlock()
	memberMap := map[string]*GroupMember{}
	for k, v := range container.Members {
		memberMap[k] = v
	}
	return memberMap
}

func (container *GroupMemberContainer) GetMember(memberId string) *GroupMember {
	container.mu.RLock()
	defer container.mu.RUnlock()
	if member, exist := container.Members[memberId]; exist {
		return member
	} else {
		return nil
	}
}

func (container *GroupMemberContainer) SetMemberMute(isMute int32, memberIds []string, muteEndAt int64) {
	container.mu.Lock()
	defer container.mu.Unlock()
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
	container.mu.Lock()
	defer container.mu.Unlock()
	for _, memberId := range memberIds {
		if member, exist := container.Members[memberId]; exist {
			member.IsAllow = int(isAllow)
		}
	}
}

func (container *GroupMemberContainer) GroupMemberCount() int {
	container.mu.RLock()
	defer container.mu.RUnlock()
	return len(container.Members)
}

func (container *GroupMemberContainer) SetGroupMembers(members map[string]*GroupMember) {
	container.mu.Lock()
	defer container.mu.Unlock()
	container.Members = members
}

func (container *GroupMemberContainer) CheckGroupMembers(memberIds []string) map[string]int64 {
	container.mu.RLock()
	defer container.mu.RUnlock()
	ret := map[string]int64{}
	for _, memberId := range memberIds {
		if member, exist := container.Members[memberId]; exist {
			ret[memberId] = member.CreatedTime
		}
	}
	return ret
}
