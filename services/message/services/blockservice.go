package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/message/storages/dbs"
	"strings"
	"time"
)

var (
	blockUserCache *caches.LruCache
	blockUserLock  *tools.SegmentatedLocks
)

func init() {
	blockUserCache = caches.NewLruCacheWithAddReadTimeout(10000, nil, 10*time.Minute, 10*time.Minute)
	blockUserLock = tools.NewSegmentatedLocks(128)
}

func AddBlockUsers(ctx context.Context, blockUserIds []string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	//add to cache
	blockUsers := GetBlockUsers(appkey, userId)
	blockUsers.AddBlockUsers(blockUserIds, time.Now().UnixMilli())
	//add to db
	dao := dbs.BlockDao{}
	for _, blockId := range blockUserIds {
		dao.Create(dbs.BlockDao{
			UserId:      userId,
			AppKey:      appkey,
			BlockUserId: blockId,
			CreatedTime: time.Now(),
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func RemoveBlockUsers(ctx context.Context, blockUserIds []string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	//remove from db
	dao := dbs.BlockDao{}
	dao.BatchDelBlockUsers(appkey, userId, blockUserIds)

	//remove from cache
	blockUsers := GetBlockUsers(appkey, userId)
	blockUsers.RemoveBlockUsers(blockUserIds)
	return errs.IMErrorCode_SUCCESS
}

func QryBlockUsers(ctx context.Context, userId string, limit int64, offset string) (errs.IMErrorCode, []*pbobjs.BlockUser, string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	retBlockUsers := []*pbobjs.BlockUser{}
	dao := dbs.BlockDao{}
	startId, err := tools.DecodeInt(offset)
	if err != nil {
		startId = 0
	}
	var retOffset string = ""
	items, err := dao.QryBlockUsers(appkey, userId, limit, startId)
	if err == nil {
		for _, item := range items {
			retBlockUsers = append(retBlockUsers, &pbobjs.BlockUser{
				BlockUserId: item.BlockUserId,
				CreatedTime: item.CreatedTime.UnixMilli(),
			})
			retOffset, err = tools.EncodeInt(item.ID)
			if err != nil {
				retOffset = ""
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, retBlockUsers, retOffset
}

type BlockUsers struct {
	Appkey        string
	UserId        string
	BlockUserIds  map[string]int64
	blockUserLock *tools.SegmentatedLocks
}

func (block *BlockUsers) CheckBlockUser(blockUserId string) bool {
	lock := block.blockUserLock.GetLocks(block.Appkey, block.UserId)
	lock.RLock()
	defer lock.RUnlock()
	_, exist := block.BlockUserIds[blockUserId]
	return exist
}

func (block *BlockUsers) AddBlockUsers(blockUserIds []string, addedTime int64) {
	lock := block.blockUserLock.GetLocks(block.Appkey, block.UserId)
	lock.Lock()
	defer lock.Unlock()
	for _, id := range blockUserIds {
		block.BlockUserIds[id] = addedTime
	}
}

func (block *BlockUsers) RemoveBlockUsers(blockUserIds []string) {
	lock := block.blockUserLock.GetLocks(block.Appkey, block.UserId)
	lock.Lock()
	defer lock.Unlock()
	for _, id := range blockUserIds {
		delete(block.BlockUserIds, id)
	}
}

func GetBlockUsers(appkey, userId string) *BlockUsers {
	key := strings.Join([]string{appkey, userId}, "_")
	if obj, exist := blockUserCache.Get(key); exist {
		return obj.(*BlockUsers)
	} else {
		lock := blockUserLock.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if obj, exist := blockUserCache.Get(key); exist {
			return obj.(*BlockUsers)
		} else {
			blockUsers := &BlockUsers{
				Appkey:        appkey,
				UserId:        userId,
				BlockUserIds:  make(map[string]int64),
				blockUserLock: blockUserLock,
			}
			dao := dbs.BlockDao{}
			var start int64 = 0
			var limit int64 = 200
			for {
				items, err := dao.QryBlockUsers(appkey, userId, limit, start)
				if err == nil && len(items) > 0 {
					for _, item := range items {
						blockUsers.BlockUserIds[item.BlockUserId] = item.CreatedTime.UnixMilli()
						if item.ID > start {
							start = item.ID
						}
					}
					if len(items) < int(limit) {
						break
					}
				} else {
					break
				}
			}
			blockUserCache.Add(key, blockUsers)
			return blockUsers
		}
	}
}
