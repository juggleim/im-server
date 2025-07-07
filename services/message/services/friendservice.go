package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/message/storages"
	"im-server/services/message/storages/models"
	"strings"
	"time"
)

type FriendStatus struct {
	IsFriend          bool
	FriendDisplayName string
}

var friendStatusCache *caches.LruCache

func init() {
	friendStatusCache = caches.NewLruCacheWithAddReadTimeout("friendstatus_cache", 100000, func(key, value interface{}) {}, 10*time.Minute, 10*time.Minute)
}

func GetFriendStatus(appkey, userId, friendId string) *FriendStatus {
	key := getFriendStatusCacheKey(appkey, userId, friendId)
	if val, exist := friendStatusCache.Get(key); exist {
		return val.(*FriendStatus)
	} else {
		l := userLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if val, exist := friendStatusCache.Get(key); exist {
			return val.(*FriendStatus)
		} else {
			status := &FriendStatus{}
			storage := storages.NewFriendRelStorage()
			rel, err := storage.GetFriendRel(appkey, userId, friendId)
			if err == nil && rel != nil {
				status.IsFriend = true
				status.FriendDisplayName = rel.DisplayName
			}
			friendStatusCache.Add(key, status)
			return status
		}
	}
}

func getFriendStatusCacheKey(appkey, userId, friendId string) string {
	return strings.Join([]string{appkey, userId, friendId}, "_")
}

func AddFriends(ctx context.Context, req *pbobjs.FriendMembersReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	friendRels := []models.FriendRel{}
	for _, friendMember := range req.FriendMembers {
		//add to cache
		key := getFriendStatusCacheKey(appkey, userId, friendMember.FriendId)
		friendStatusCache.Add(key, &FriendStatus{
			IsFriend:          true,
			FriendDisplayName: friendMember.DisplayName,
		})
		friendRels = append(friendRels, models.FriendRel{
			AppKey:      appkey,
			UserId:      userId,
			FriendId:    friendMember.FriendId,
			OrderTag:    friendMember.OrderTag,
			DisplayName: friendMember.DisplayName,
		})
	}
	storage.BatchUpsert(friendRels)
	return errs.IMErrorCode_SUCCESS
}

func DelFriends(ctx context.Context, req *pbobjs.FriendIdsReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	storage.BatchDelete(appkey, userId, req.FriendIds)
	for _, friendId := range req.FriendIds {
		key := getFriendStatusCacheKey(appkey, userId, friendId)
		friendStatusCache.Remove(key)
	}
	return errs.IMErrorCode_SUCCESS
}

func QryFriends(ctx context.Context, req *pbobjs.QryFriendsReq) (errs.IMErrorCode, *pbobjs.QryFriendsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	var startId int64 = 0
	if req.Offset != "" {
		startId, _ = tools.DecodeInt(req.Offset)
	}
	ret := &pbobjs.QryFriendsResp{
		Items: []*pbobjs.FriendMember{},
	}
	rels, err := storage.QueryFriendRels(appkey, userId, startId, req.Limit, req.Order > 0)
	if err == nil {
		for _, rel := range rels {
			ret.Offset, _ = tools.EncodeInt(rel.ID)
			ret.Items = append(ret.Items, &pbobjs.FriendMember{
				FriendId:    rel.FriendId,
				OrderTag:    rel.OrderTag,
				DisplayName: rel.DisplayName,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryFriendsWithPage(ctx context.Context, req *pbobjs.QryFriendsWithPageReq) (errs.IMErrorCode, *pbobjs.QryFriendsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	storage := storages.NewFriendRelStorage()
	ret := &pbobjs.QryFriendsResp{
		Items: []*pbobjs.FriendMember{},
	}
	rels, err := storage.QueryFriendRelsWithPage(appkey, userId, req.OrderTag, req.Page, req.Size)
	if err == nil {
		for _, rel := range rels {
			ret.Items = append(ret.Items, &pbobjs.FriendMember{
				FriendId:    rel.FriendId,
				OrderTag:    rel.OrderTag,
				DisplayName: rel.DisplayName,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func CheckFriends(ctx context.Context, req *pbobjs.CheckFriendsReq) (errs.IMErrorCode, *pbobjs.CheckFriendsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	ret := &pbobjs.CheckFriendsResp{
		CheckResults: make(map[string]bool),
	}
	if len(req.FriendIds) <= 0 {
		return errs.IMErrorCode_SUCCESS, ret
	}
	for _, friendId := range req.FriendIds {
		status := GetFriendStatus(appkey, userId, friendId)
		ret.CheckResults[friendId] = status.IsFriend
	}
	return errs.IMErrorCode_SUCCESS, ret
}
