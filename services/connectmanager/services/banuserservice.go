package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/connectmanager/dbs"
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
			var exist bool = false
			banUserDao := dbs.BanUserDao{}
			dbBanUsers, err := banUserDao.FindById(appkey, userId)
			if err == nil && len(dbBanUsers) > 0 {
				banUser = &BanUserItem{
					UserId: userId,
					Items:  make(map[string]*BanItem),
				}
				needClean := false
				curr := time.Now().UnixMilli()
				for _, dbBanUser := range dbBanUsers {
					banUser.Items[dbBanUser.ScopeKey] = &BanItem{
						EndTime:    dbBanUser.EndTime,
						ScopeKey:   dbBanUser.ScopeKey,
						ScopeValue: dbBanUser.ScopeValue,
						Ext:        dbBanUser.Ext,
					}
					if dbBanUser.EndTime > 0 && dbBanUser.EndTime < curr {
						needClean = true
					}
				}
				if needClean {
					go func() {
						banUserDao.CleanBaseTime(appkey, userId, curr)
					}()
				}
				exist = true
			}
			banUserCache.Add(key, banUser)
			return banUser, exist
		}
	}
}

type BanUserItem struct {
	UserId string
	Items  map[string]*BanItem
}

type BanItem struct {
	EndTime    int64
	ScopeKey   string
	ScopeValue string
	Ext        string
}

func BanUsers(ctx context.Context, banUsers []*pbobjs.BanUser) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.BanUserDao{}
	for _, user := range banUsers {
		dao.Upsert(dbs.BanUserDao{
			UserId:      user.UserId,
			ScopeKey:    user.ScopeKey,
			ScopeValue:  user.ScopeValue,
			Ext:         user.Ext,
			CreatedTime: time.Now(),
			EndTime:     user.EndTime,
			AppKey:      appkey,
		})
		key := strings.Join([]string{appkey, user.UserId}, "_")
		banUserCache.Remove(key)

		//kick user
		kicReq := &pbobjs.KickUserReq{
			Ext: user.Ext,
		}
		if user.UserId != "" && user.ScopeKey != "" {
			switch user.ScopeKey {
			case string(dbs.UserBanScopeDefault):
				kicReq.UserId = user.UserId
			case string(dbs.UserBanScopeDevice):
				deviceIds := strings.Split(user.ScopeValue, ",")
				if len(deviceIds) > 0 {
					kicReq.UserId = user.UserId
					kicReq.DeviceIds = deviceIds
				}
			case string(dbs.UserBanScopePlatform):
				platforms := strings.Split(user.ScopeValue, ",")
				if len(platforms) > 0 {
					kicReq.UserId = user.UserId
					kicReq.Platforms = platforms
				}
			case string(dbs.UserBanScopeIp):
				ips := strings.Split(user.ScopeValue, ",")
				if len(ips) > 0 {
					kicReq.UserId = user.UserId
					kicReq.Ips = ips
				}
			}
			if kicReq.UserId != "" {
				KickUser(ctx, kicReq)
			}
		}
	}
}

func UnBanUsers(ctx context.Context, banUsers []*pbobjs.BanUser) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.BanUserDao{}
	for _, user := range banUsers {
		dao.DelBanUser(appkey, user.UserId, user.ScopeKey)
		key := strings.Join([]string{appkey, user.UserId}, "_")
		banUserCache.Remove(key) //TODO
	}
}

func QryBanUsers(ctx context.Context, limit int64, startIdStr string) (errs.IMErrorCode, []*pbobjs.BanUser, string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	dao := dbs.BanUserDao{}

	startId, err := tools.DecodeInt(startIdStr)
	if err != nil {
		startId = 0
	}
	banUsers := []*pbobjs.BanUser{}
	offset := ""
	dbBanUsers, err := dao.QryBanUsers(appkey, limit, startId)
	if err == nil {
		for _, dbBanUser := range dbBanUsers {
			user := &pbobjs.BanUser{
				UserId:      dbBanUser.UserId,
				EndTime:     0,
				CreatedTime: dbBanUser.CreatedTime.UnixMilli(),
				ScopeKey:    dbBanUser.ScopeKey,
				ScopeValue:  dbBanUser.ScopeValue,
				Ext:         dbBanUser.Ext,
			}
			user.EndTime = dbBanUser.EndTime
			banUsers = append(banUsers, user)

			offset, err = tools.EncodeInt(dbBanUser.ID)
			if err != nil {
				offset = ""
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, banUsers, offset
}
