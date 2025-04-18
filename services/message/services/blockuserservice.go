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
	blockUserCache = caches.NewLruCacheWithAddReadTimeout("blockuser_cache", 100000, nil, 10*time.Minute, 10*time.Minute)
	blockUserLock = tools.NewSegmentatedLocks(256)
}

type BlockUserItem struct {
	Appkey      string
	UserId      string
	BlockUserId string
	IsBlock     bool
	BlockTime   int64
}

func (block *BlockUserItem) SetBlock() (bool, int64) {
	if block.IsBlock {
		return false, 0
	}
	lock := blockUserLock.GetLocks(block.Appkey, block.UserId, block.BlockUserId)
	lock.Lock()
	defer lock.Unlock()
	if block.IsBlock {
		return false, 0
	} else {
		t := time.Now().UnixMilli()
		block.IsBlock = true
		block.BlockTime = t
		return true, t
	}
}

func (block *BlockUserItem) DelBlock() bool {
	if !block.IsBlock {
		return false
	}
	lock := blockUserLock.GetLocks(block.Appkey, block.UserId, block.BlockUserId)
	lock.Lock()
	defer lock.Unlock()
	if !block.IsBlock {
		return false
	} else {
		block.IsBlock = false
		block.BlockTime = 0
		return true
	}
}

func getBlockUserCacheKey(appkey, userId, blockUserId string) string {
	return strings.Join([]string{appkey, userId, blockUserId}, "_")
}

func GetBlockUserItem(appkey, userId, blockUserId string) *BlockUserItem {
	key := getBlockUserCacheKey(appkey, userId, blockUserId)
	if obj, exist := blockUserCache.Get(key); exist {
		return obj.(*BlockUserItem)
	} else {
		lock := blockUserLock.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()

		if obj, exist := blockUserCache.Get(key); exist {
			return obj.(*BlockUserItem)
		} else {
			blockUserItem := &BlockUserItem{
				Appkey:      appkey,
				UserId:      userId,
				BlockUserId: blockUserId,
				IsBlock:     false,
				BlockTime:   0,
			}
			dao := dbs.BlockDao{}
			item, err := dao.Find(appkey, userId, blockUserId)
			if err == nil && item != nil {
				blockUserItem.IsBlock = true
				blockUserItem.BlockTime = item.CreatedTime.UnixMilli()
			}
			blockUserCache.Add(key, blockUserItem)
			return blockUserItem
		}
	}
}

func AddBlockUser(ctx context.Context, userId, blockUserId string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//check from cache
	item := GetBlockUserItem(appkey, userId, blockUserId)
	succ, t := item.SetBlock()
	if succ {
		//add to db
		dao := dbs.BlockDao{}
		dao.Create(dbs.BlockDao{
			AppKey:      appkey,
			UserId:      userId,
			BlockUserId: blockUserId,
			CreatedTime: time.UnixMilli(t),
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func DelBlockUser(ctx context.Context, userId, blockUserId string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	//check from cache
	item := GetBlockUserItem(appkey, userId, blockUserId)
	succ := item.DelBlock()
	if succ {
		//del from db
		dao := dbs.BlockDao{}
		dao.DelBlockUser(appkey, userId, blockUserId)
	}
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
