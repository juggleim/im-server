package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/connectmanager/dbs"
	"sort"
	"strings"
	"time"
)

var banUserCache *caches.LruCache
var userLocks *tools.SegmentatedLocks
var notBanUser *BanUserItem

func init() {
	banUserCache = caches.NewLruCacheWithAddReadTimeout("banuser_cache", 100000, nil, 8*time.Minute, 10*time.Minute)
	userLocks = tools.NewSegmentatedLocks(512)

	notBanUser = &BanUserItem{}
}

func GetBanUserFromCache(appkey, userId string) (*BanUserItem, bool) {
	key := strings.Join([]string{appkey, userId}, "_")
	if val, exist := banUserCache.Get(key); exist {
		banUser := val.(*BanUserItem)
		if banUser == notBanUser {
			return nil, false
		} else {
			return banUser, true
		}
	} else {
		l := userLocks.GetLocks(appkey, userId)
		l.Lock()
		defer l.Unlock()

		if val, exist := banUserCache.Get(key); exist {
			banUser := val.(*BanUserItem)
			if banUser == notBanUser {
				return nil, false
			} else {
				return banUser, true
			}
		} else {
			var banUser *BanUserItem = notBanUser
			var exist bool
			banUserDao := dbs.BanUserDao{}
			legacyBanUserDao := dbs.LegacyBanUserDao{}
			legacyBanUsers, legacyErr := legacyBanUserDao.FindById(appkey, userId)
			dbBanUsers, err := banUserDao.FindById(appkey, userId)
			if len(legacyBanUsers) > 0 || len(dbBanUsers) > 0 {
				banUser = &BanUserItem{
					UserId: userId,
					Items:  make(map[string]map[string]*BanItem),
				}
				needCleanLegacy := false
				needCleanRules := false
				curr := time.Now().UnixMilli()
				for _, legacyBanUser := range legacyBanUsers {
					legacyScopeKey := normalizeBanScopeKey(legacyBanUser.ScopeKey)
					for _, scopeValue := range dbs.SplitLegacyBanScopeValues(legacyScopeKey, legacyBanUser.ScopeValue) {
						addBanItem(banUser, legacyScopeKey, scopeValue, legacyBanUser.EndTime, legacyBanUser.Ext)
					}
					if legacyBanUser.EndTime > 0 && legacyBanUser.EndTime <= curr {
						needCleanLegacy = true
					}
				}
				// New rules are loaded last so they take precedence over matching legacy values.
				for _, dbBanUser := range dbBanUsers {
					addBanItem(banUser, dbBanUser.ScopeKey, dbBanUser.ScopeValue, dbBanUser.EndTime, dbBanUser.Ext)
					if dbBanUser.EndTime > 0 && dbBanUser.EndTime <= curr {
						needCleanRules = true
					}
				}
				if needCleanLegacy || needCleanRules {
					go func() {
						if needCleanLegacy {
							if cleanErr := legacyBanUserDao.CleanBaseTime(appkey, userId, curr); cleanErr != nil {
								logs.NewLogEntity().Error(cleanErr.Error())
							}
						}
						if needCleanRules {
							if cleanErr := banUserDao.CleanBaseTime(appkey, userId, curr); cleanErr != nil {
								logs.NewLogEntity().Error(cleanErr.Error())
							}
						}
					}()
				}
				exist = true
			}
			if legacyErr != nil {
				logs.NewLogEntity().Error(legacyErr.Error())
			}
			if err != nil {
				logs.NewLogEntity().Error(err.Error())
			}
			banUserCache.Add(key, banUser)
			return banUser, exist
		}
	}
}

type BanUserItem struct {
	UserId string
	Items  map[string]map[string]*BanItem
}

type BanItem struct {
	EndTime    int64
	ScopeKey   string
	ScopeValue string
	Ext        string
}

func addBanItem(user *BanUserItem, scopeKey, scopeValue string, endTime int64, ext string) {
	if _, exists := user.Items[scopeKey]; !exists {
		user.Items[scopeKey] = make(map[string]*BanItem)
	}
	user.Items[scopeKey][scopeValue] = &BanItem{
		EndTime:    endTime,
		ScopeKey:   scopeKey,
		ScopeValue: scopeValue,
		Ext:        ext,
	}
}

func (user *BanUserItem) ActiveBan(scopeKey, scopeValue string, currentTime int64) (*BanItem, bool) {
	if user == nil {
		return nil, false
	}
	rules, exists := user.Items[scopeKey]
	if !exists {
		return nil, false
	}
	rule, exists := rules[scopeValue]
	if !exists || (rule.EndTime > 0 && rule.EndTime <= currentTime) {
		return nil, false
	}
	return rule, true
}

func BanUsers(ctx context.Context, banUsers []*pbobjs.BanUser) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.BanUserDao{}
	for _, user := range banUsers {
		scopeKey := normalizeBanScopeKey(user.ScopeKey)
		scopeValues := normalizeBanScopeValues(scopeKey, user.ScopeValue, user.ScopeValues)
		if len(scopeValues) == 0 {
			continue
		}
		createdTime := time.Now()
		rules := make([]dbs.BanUserDao, 0, len(scopeValues))
		for _, scopeValue := range scopeValues {
			rules = append(rules, dbs.BanUserDao{
				UserId:      user.UserId,
				ScopeKey:    scopeKey,
				ScopeValue:  scopeValue,
				Ext:         user.Ext,
				CreatedTime: createdTime,
				EndTime:     user.EndTime,
				AppKey:      appkey,
			})
		}
		if err := dao.UpsertMany(rules); err != nil {
			logs.NewLogEntity().Error(err.Error())
			continue
		}
		key := strings.Join([]string{appkey, user.UserId}, "_")
		banUserCache.Remove(key)

		//kick user
		kicReq := &pbobjs.KickUserReq{
			Ext: user.Ext,
		}
		if user.UserId != "" {
			switch scopeKey {
			case string(dbs.UserBanScopeDefault):
				kicReq.UserId = user.UserId
			case string(dbs.UserBanScopeDevice):
				kicReq.UserId = user.UserId
				kicReq.DeviceIds = scopeValues
			case string(dbs.UserBanScopePlatform):
				kicReq.UserId = user.UserId
				kicReq.Platforms = scopeValues
			case string(dbs.UserBanScopeIp):
				kicReq.UserId = user.UserId
				kicReq.Ips = scopeValues
			}
			if kicReq.UserId != "" {
				KickUser(ctx, kicReq, errs.IMErrorCode_CONNECT_USER_BLOCK)
			}
		}
	}
}

func UnBanUsers(ctx context.Context, banUsers []*pbobjs.BanUser) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.BanUserDao{}
	legacyDao := dbs.LegacyBanUserDao{}
	for _, user := range banUsers {
		scopeKey := strings.TrimSpace(user.ScopeKey)
		var scopeValues []string
		if scopeKey != "" {
			scopeValues = normalizeBanScopeValues(scopeKey, user.ScopeValue, user.ScopeValues)
		}
		if err := dao.DelBanUser(appkey, user.UserId, scopeKey, scopeValues); err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
		if err := legacyDao.DelBanUser(appkey, user.UserId, scopeKey, scopeValues); err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
		key := strings.Join([]string{appkey, user.UserId}, "_")
		banUserCache.Remove(key)
	}
}

func QryBanUsers(ctx context.Context, limit int64, startIdStr string, userIds []string) (errs.IMErrorCode, []*pbobjs.BanUser, string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	banUsers := []*pbobjs.BanUser{}
	dao := dbs.BanUserDao{}
	legacyDao := dbs.LegacyBanUserDao{}
	if len(userIds) > 0 {
		legacyBanUsers, legacyErr := legacyDao.FindByUserIds(appkey, userIds)
		if legacyErr != nil {
			logs.NewLogEntity().Error(legacyErr.Error())
		} else {
			for _, legacyBanUser := range legacyBanUsers {
				banUsers = append(banUsers, legacyBanUserToProto(legacyBanUser))
			}
		}
		dbBanUsers, dbErr := dao.FindByUserIds(appkey, userIds)
		if dbErr != nil {
			logs.NewLogEntity().Error(dbErr.Error())
		} else {
			for _, dbBanUser := range dbBanUsers {
				banUsers = append(banUsers, banRuleToProto(dbBanUser))
			}
		}
		return errs.IMErrorCode_SUCCESS, banUsers, ""
	}

	phase, startId := decodeBanQueryOffset(startIdStr)
	if phase == banQueryPhaseLegacy {
		legacyBanUsers, legacyErr := legacyDao.QryBanUsers(appkey, limit, startId)
		if legacyErr != nil {
			logs.NewLogEntity().Error(legacyErr.Error())
			return errs.IMErrorCode_SUCCESS, banUsers, ""
		}
		if len(legacyBanUsers) > 0 {
			for _, legacyBanUser := range legacyBanUsers {
				banUsers = append(banUsers, legacyBanUserToProto(legacyBanUser))
			}
			return errs.IMErrorCode_SUCCESS, banUsers, encodeBanQueryOffset(banQueryPhaseLegacy, legacyBanUsers[len(legacyBanUsers)-1].ID)
		}
		startId = 0
	}

	dbBanUsers, dbErr := dao.QryBanUsers(appkey, limit, startId)
	if dbErr != nil {
		logs.NewLogEntity().Error(dbErr.Error())
		return errs.IMErrorCode_SUCCESS, banUsers, ""
	}
	for _, dbBanUser := range dbBanUsers {
		banUsers = append(banUsers, banRuleToProto(dbBanUser))
	}
	offset := ""
	if len(dbBanUsers) > 0 {
		offset = encodeBanQueryOffset(banQueryPhaseRules, dbBanUsers[len(dbBanUsers)-1].ID)
	}
	return errs.IMErrorCode_SUCCESS, banUsers, offset
}

const (
	banQueryPhaseLegacy = "legacy"
	banQueryPhaseRules  = "rules"
)

func decodeBanQueryOffset(offset string) (string, int64) {
	phase := banQueryPhaseLegacy
	encodedId := offset
	if prefix, value, found := strings.Cut(offset, ":"); found {
		if prefix == banQueryPhaseLegacy || prefix == banQueryPhaseRules {
			phase = prefix
			encodedId = value
		}
	}
	id, err := tools.DecodeInt(encodedId)
	if err != nil {
		id = 0
	}
	return phase, id
}

func encodeBanQueryOffset(phase string, id int64) string {
	encodedId, err := tools.EncodeInt(id)
	if err != nil {
		return ""
	}
	return phase + ":" + encodedId
}

func legacyBanUserToProto(user *dbs.LegacyBanUserDao) *pbobjs.BanUser {
	scopeKey := normalizeBanScopeKey(user.ScopeKey)
	scopeValues := dbs.SplitLegacyBanScopeValues(scopeKey, user.ScopeValue)
	if len(scopeValues) == 1 && scopeValues[0] == "" {
		scopeValues = nil
	}
	return &pbobjs.BanUser{
		UserId:      user.UserId,
		EndTime:     user.EndTime,
		CreatedTime: user.CreatedTime.UnixMilli(),
		ScopeKey:    scopeKey,
		ScopeValue:  user.ScopeValue,
		ScopeValues: scopeValues,
		Ext:         user.Ext,
	}
}

func banRuleToProto(user *dbs.BanUserDao) *pbobjs.BanUser {
	var scopeValues []string
	if user.ScopeValue != "" {
		scopeValues = []string{user.ScopeValue}
	}
	return &pbobjs.BanUser{
		UserId:      user.UserId,
		EndTime:     user.EndTime,
		CreatedTime: user.CreatedTime.UnixMilli(),
		ScopeKey:    user.ScopeKey,
		ScopeValue:  user.ScopeValue,
		ScopeValues: scopeValues,
		Ext:         user.Ext,
	}
}

func normalizeBanScopeValues(scopeKey, legacyScopeValue string, scopeValues []string) []string {
	if scopeKey == string(dbs.UserBanScopeDefault) {
		return []string{""}
	}
	values := make(map[string]struct{}, len(scopeValues)+1)
	for _, scopeValue := range scopeValues {
		scopeValue = strings.TrimSpace(scopeValue)
		if scopeValue != "" {
			values[scopeValue] = struct{}{}
		}
	}
	for _, scopeValue := range strings.Split(legacyScopeValue, ",") {
		scopeValue = strings.TrimSpace(scopeValue)
		if scopeValue != "" {
			values[scopeValue] = struct{}{}
		}
	}
	result := make([]string, 0, len(values))
	for scopeValue := range values {
		result = append(result, scopeValue)
	}
	sort.Strings(result)
	return result
}

func normalizeBanScopeKey(scopeKey string) string {
	scopeKey = strings.TrimSpace(scopeKey)
	if scopeKey == "" {
		return string(dbs.UserBanScopeDefault)
	}
	return scopeKey
}
