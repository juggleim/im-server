package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	chrmService "im-server/services/chatroom/services"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bytedance/gopkg/collection/skipmap"
	"github.com/bytedance/gopkg/collection/zset"
	"google.golang.org/protobuf/proto"
)

var (
	chrmCache *caches.LruCache
	chrmLocks *tools.SegmentatedLocks
)

func init() {
	chrmCache = caches.NewLruCacheWithReadTimeout(10000, nil, time.Hour)
	chrmLocks = tools.NewSegmentatedLocks(256)
}

type ChatroomContainer struct {
	Appkey string
	ChatId string
	Status chrmService.ChatroomStatus

	PartialMembers map[string]*ChatroomMember

	MsgSet *skipmap.Int64Map

	Atts      map[string]*pbobjs.ChatAttItem
	AttsIndex *zset.Float64Set
}

type ChatroomMember struct {
	MemberId       string
	UnReadCount    int32
	LatestSyncTime int64
}

func (container *ChatroomContainer) Destroy() {
	container.Status = chrmService.ChatroomStatus_Destroy
	container.PartialMembers = make(map[string]*ChatroomMember)
	container.MsgSet = skipmap.NewInt64()
	container.Atts = make(map[string]*pbobjs.ChatAttItem)
	container.AttsIndex = zset.NewFloat64()
}

func (container *ChatroomContainer) AppendMsg(ctx context.Context, msg *pbobjs.DownMsg) {
	key := getChrmKey(container.Appkey, container.ChatId)
	lock := chrmLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	container.MsgSet.Store(msg.MsgTime, msg)
	maxMsgCount := getMaxMsgCount(container.Appkey)
	if container.MsgSet.Len() > maxMsgCount {
		evictCount := container.MsgSet.Len() - maxMsgCount
		deletedArr := []int64{}
		container.MsgSet.Range(func(key int64, value interface{}) bool {
			evictCount--
			deletedArr = append(deletedArr, key)
			if evictCount <= 0 {
				return false
			} else {
				return true
			}
		})
		for _, k := range deletedArr {
			container.MsgSet.Delete(k)
		}
	}
}

func (container *ChatroomContainer) GetMsgsBaseTime(ctx context.Context, userId string, start int64) ([]*pbobjs.DownMsg, errs.IMErrorCode) {
	key := getChrmKey(container.Appkey, container.ChatId)
	lock := chrmLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	retMsgs := []*pbobjs.DownMsg{}
	container.MsgSet.Range(func(key int64, value interface{}) bool {
		if key > start {
			msg := value.(*pbobjs.DownMsg)
			if msg.SenderId == userId {
				msg.IsSend = true
			}
			retMsgs = append(retMsgs, msg)
		}
		return true
	})
	return retMsgs, errs.IMErrorCode_SUCCESS
}

func (container *ChatroomContainer) AppendAtt(ctx context.Context, att *pbobjs.ChatAttItem) {
	key := getChrmKey(container.Appkey, container.ChatId)
	lock := chrmLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	container.Atts[att.Key] = att
	container.AttsIndex.Add(float64(att.AttTime), att.Key)
}

func (container *ChatroomContainer) GetAttsBaseTime(ctx context.Context, start int64) ([]*pbobjs.ChatAttItem, errs.IMErrorCode) {
	key := getChrmKey(container.Appkey, container.ChatId)
	lock := chrmLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	retAtts := []*pbobjs.ChatAttItem{}
	nodes := container.AttsIndex.RangeByScoreWithOpt(float64(start), math.MaxFloat64, zset.RangeOpt{
		ExcludeMin: true,
	})
	for _, node := range nodes {
		if att, exist := container.Atts[node.Value]; exist {
			retAtts = append(retAtts, &pbobjs.ChatAttItem{
				Key:     att.Key,
				Value:   att.Value,
				UserId:  att.UserId,
				AttTime: att.AttTime,
				OptType: att.OptType,
			})
		}
	}
	return retAtts, errs.IMErrorCode_SUCCESS
}

func (container *ChatroomContainer) ForeachMembers(f func(member *ChatroomMember) bool) {
	key := getChrmKey(container.Appkey, container.ChatId)
	lock := chrmLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	for _, member := range container.PartialMembers {
		isContinue := f(member)
		if !isContinue {
			break
		}
	}
}

func (container *ChatroomContainer) CleanUnread(memberId string) {
	key := getChrmKey(container.Appkey, container.ChatId)
	lock := chrmLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	if member, exist := container.PartialMembers[memberId]; exist {
		atomic.StoreInt32(&member.UnReadCount, 0)
		atomic.StoreInt64(&member.LatestSyncTime, time.Now().UnixMilli())
	}
}

func (container *ChatroomContainer) AddMember(memberId string) bool {
	key := getChrmKey(container.Appkey, container.ChatId)
	lock := chrmLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if _, exist := container.PartialMembers[memberId]; !exist {
		container.PartialMembers[memberId] = &ChatroomMember{
			MemberId: memberId,
		}
		return true
	}
	return false
}

func (container *ChatroomContainer) DelMember(memberId string) bool {
	key := getChrmKey(container.Appkey, container.ChatId)
	lock := chrmLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if _, exist := container.PartialMembers[memberId]; exist {
		delete(container.PartialMembers, memberId)
		return true
	}
	return false
}

func getChrmKey(appkey, chatId string) string {
	return strings.Join([]string{appkey, chatId}, "_")
}

func GetChrmContainer(ctx context.Context, appkey, chatId string) (*ChatroomContainer, bool) {
	key := getChrmKey(appkey, chatId)
	if obj, exist := chrmCache.Get(key); exist {
		container := obj.(*ChatroomContainer)
		if container.Status == chrmService.ChatroomStatus_Normal {
			return container, true
		}
		return container, false
	} else {
		lock := chrmLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if obj, exist := chrmCache.Get(key); exist {
			container := obj.(*ChatroomContainer)
			if container.Status == chrmService.ChatroomStatus_Normal {
				return container, true
			}
			return container, false
		} else {
			chrmContainer := qryChatroomPartialInfo(ctx, chatId)
			chrmCache.Add(key, chrmContainer)
			return chrmContainer, chrmContainer.Status == chrmService.ChatroomStatus_Normal
		}
	}
}

func initChatroomContainer(appkey, chatId string) {
	key := getChrmKey(appkey, chatId)
	lock := chrmLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if !chrmCache.Contains(key) {
		chrmCache.Add(key, &ChatroomContainer{
			Appkey:         appkey,
			ChatId:         chatId,
			Status:         chrmService.ChatroomStatus_Normal,
			PartialMembers: make(map[string]*ChatroomMember),
			MsgSet:         skipmap.NewInt64(),
			Atts:           make(map[string]*pbobjs.ChatAttItem),
			AttsIndex:      zset.NewFloat64(),
		})
	}
}

func qryChatroomPartialInfo(ctx context.Context, chatId string) *ChatroomContainer {
	appkey := bases.GetAppKeyFromCtx(ctx)
	code, resp, err := bases.SyncRpcCall(ctx, "c_sync_partial", chatId, &pbobjs.ChatMsgNode{
		NodeName: bases.GetCluster().GetCurrentNode().Name,
		Method:   "c_members_dispatch",
	}, func() proto.Message {
		return &pbobjs.ChatroomInfo{}
	})
	chrmContainer := &ChatroomContainer{
		Appkey:         appkey,
		ChatId:         chatId,
		Status:         chrmService.ChatroomStatus_Normal,
		PartialMembers: make(map[string]*ChatroomMember),
		MsgSet:         skipmap.NewInt64(),
		Atts:           make(map[string]*pbobjs.ChatAttItem),
		AttsIndex:      zset.NewFloat64(),
	}
	if err == nil {
		if code == errs.IMErrorCode_CHATROOM_NOTEXIST {
			chrmContainer.Status = chrmService.ChatroomStatus_NotExist
		} else {
			chrmInfo := resp.(*pbobjs.ChatroomInfo)
			// partial members
			for _, member := range chrmInfo.Members {
				chrmContainer.PartialMembers[member.MemberId] = &ChatroomMember{
					MemberId: member.MemberId,
				}
			}
			// atts
			for _, att := range chrmInfo.Atts {
				chrmContainer.Atts[att.Key] = att
				chrmContainer.AttsIndex.Add(float64(att.AttTime), att.Key)
			}
		}
		return chrmContainer
	}
	return chrmContainer
}
