package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/rtcroom/storages"
	"strings"
	"time"
)

var (
	memberCache *caches.LruCache
	memberLocks *tools.SegmentatedLocks
)

func init() {
	memberCache = caches.NewLruCacheWithReadTimeout(10000, nil, time.Hour)
	memberLocks = tools.NewSegmentatedLocks(128)
}

type RtcMemberContainer struct {
	Appkey   string
	MemberId string
	Rooms    []*RtcMemberRoom
}

type RtcMemberRoom struct {
	RoomId   string
	RtcState pbobjs.RtcState
}

func GetRtcMemberContainer(appkey, memberId string) *RtcMemberContainer {
	key := getRtcMemberKey(appkey, memberId)
	if cacheContainer, exist := memberCache.Get(key); exist {
		container := cacheContainer.(*RtcMemberContainer)
		return container
	} else {
		lock := memberLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if cacheContainer, exist := memberCache.Get(key); exist {
			container := cacheContainer.(*RtcMemberContainer)
			return container
		} else {
			container := &RtcMemberContainer{
				Appkey:   appkey,
				MemberId: memberId,
				Rooms:    []*RtcMemberRoom{},
			}
			storage := storages.NewRtcRoomMemberStorage()
			roomMembers, err := storage.QueryRoomsByMember(appkey, memberId, 10)
			if err == nil && len(roomMembers) > 0 {
				for _, roomMember := range roomMembers {
					container.Rooms = append(container.Rooms, &RtcMemberRoom{
						RoomId:   roomMember.RoomId,
						RtcState: roomMember.RtcState,
					})
				}
			}
			memberCache.Add(key, container)
			return container
		}
	}
}

func getRtcMemberKey(appkey, memberId string) string {
	return strings.Join([]string{appkey, memberId}, "_")
}

func QryRtcMemberRooms(ctx context.Context) (errs.IMErrorCode, *pbobjs.RtcMemberRooms) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	ret := &pbobjs.RtcMemberRooms{
		Rooms: []*pbobjs.RtcMemberRoom{},
	}
	currTime := time.Now().UnixMilli()
	storage := storages.NewRtcRoomMemberStorage()
	roomMembers, err := storage.QueryRoomsByMember(appkey, userId, 100)
	if err == nil && len(roomMembers) > 0 {
		for _, roomMember := range roomMembers {
			if currTime-roomMember.LatestPingTime > RtcPingTimeOut {
				roomType := roomMember.RoomType
				roomId := roomMember.RoomId
				memberId := roomMember.MemberId
				go func() {
					if roomType == pbobjs.RtcRoomType_OneOne {
						storage.DeleteByRoomId(appkey, roomId)
						roomStorage := storages.NewRtcRoomStorage()
						roomStorage.Delete(appkey, roomId)
					} else {
						storage.Delete(appkey, roomId, memberId)
					}
				}()
				continue
			}
			memberRoom := &pbobjs.RtcMemberRoom{
				RoomType: roomMember.RoomType,
				RoomId:   roomMember.RoomId,
				Owner: &pbobjs.UserInfo{
					UserId: roomMember.OwnerId,
				},
				RtcState: roomMember.RtcState,
			}
			ret.Rooms = append(ret.Rooms, memberRoom)
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}
