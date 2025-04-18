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
	memberCache = caches.NewLruCacheWithReadTimeout("rtcmember_cache", 10000, nil, time.Hour)
	memberLocks = tools.NewSegmentatedLocks(128)
}

type RtcMemberContainer struct {
	Appkey   string
	MemberId string
	Rooms    map[string]*RtcMemberRoom
}

type RtcMemberRoom struct {
	RoomId       string
	RoomType     pbobjs.RtcRoomType
	RtcChannel   pbobjs.RtcChannel
	RtcMediaType pbobjs.RtcMediaType
	RtcState     pbobjs.RtcState
	DeviceId     string
}

func (container *RtcMemberContainer) Add(memberRoom RtcMemberRoom) {
	key := getRtcMemberKey(container.Appkey, container.MemberId)
	lock := memberLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	container.InnerAdd(memberRoom)
}

func (container *RtcMemberContainer) InnerAdd(memberRoom RtcMemberRoom) {
	container.Rooms[memberRoom.RoomId] = &RtcMemberRoom{
		RoomId:   memberRoom.RoomId,
		RoomType: memberRoom.RoomType,
		RtcState: memberRoom.RtcState,
		DeviceId: memberRoom.DeviceId,
	}
}

func (container *RtcMemberContainer) Del(memberRoom RtcMemberRoom) {
	key := getRtcMemberKey(container.Appkey, container.MemberId)
	lock := memberLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	delete(container.Rooms, memberRoom.RoomId)
}

func (container *RtcMemberContainer) ForeachRooms(f func(room *RtcMemberRoom)) {
	key := getRtcMemberKey(container.Appkey, container.MemberId)
	lock := memberLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	for _, room := range container.Rooms {
		f(room)
	}
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
				Rooms:    make(map[string]*RtcMemberRoom),
			}
			storage := storages.NewRtcRoomMemberStorage()
			roomMembers, err := storage.QueryRoomsByMember(appkey, memberId, 100)
			if err == nil && len(roomMembers) > 0 {
				curr := time.Now().UnixMilli()
				for _, roomMember := range roomMembers {
					if roomMember.RtcState == pbobjs.RtcState_RtcIncoming || roomMember.RtcState == pbobjs.RtcState_RtcOutgoing {
						if curr-roomMember.LatestPingTime > (CallTimeout * 1000) {
							continue
						}
					} else {
						if curr-roomMember.LatestPingTime > (PingTimeout * 1000) {
							continue
						}
					}
					container.Rooms[roomMember.RoomId] = &RtcMemberRoom{
						RoomId:   roomMember.RoomId,
						RoomType: roomMember.RoomType,
						RtcState: roomMember.RtcState,
						DeviceId: roomMember.DeviceId,
					}
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

	container := GetRtcMemberContainer(appkey, userId)
	container.ForeachRooms(func(room *RtcMemberRoom) {
		ret.Rooms = append(ret.Rooms, &pbobjs.RtcMemberRoom{
			RoomType: room.RoomType,
			RoomId:   room.RoomId,
			RtcState: room.RtcState,
			DeviceId: room.DeviceId,
		})
	})
	return errs.IMErrorCode_SUCCESS, ret
}

func GrabMemberState(ctx context.Context, req *pbobjs.MemberState) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)

	container := GetRtcMemberContainer(appkey, userId)
	key := getRtcMemberKey(appkey, userId)
	lock := memberLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()

	if req.RoomType == pbobjs.RtcRoomType_OneOne {
		for _, room := range container.Rooms {
			if room.RoomId == req.RoomId {
				continue
			}
			if room.RtcState == pbobjs.RtcState_RtcOutgoing || room.RtcState == pbobjs.RtcState_RtcConnecting || room.RtcState == pbobjs.RtcState_RtcConnected {
				return errs.IMErrorCode_RTCINVITE_BUSY
			}
		}
	}
	container.InnerAdd(RtcMemberRoom{
		RoomId:   req.RoomId,
		RoomType: req.RoomType,
		RtcState: req.RtcState,
		DeviceId: req.DeviceId,
	})
	return errs.IMErrorCode_SUCCESS
}

func SyncMemberState(ctx context.Context, req *pbobjs.SyncMemberStateReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetTargetIdFromCtx(ctx)
	container := GetRtcMemberContainer(appkey, userId)
	if req.IsDelete {
		container.Del(RtcMemberRoom{
			RoomId: req.Member.RoomId,
		})
	} else {
		container.Add(RtcMemberRoom{
			RoomId:   req.Member.RoomId,
			RoomType: req.Member.RoomType,
			RtcState: req.Member.RtcState,
			DeviceId: req.Member.DeviceId,
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func SyncQuitWhenConnectKicked(ctx context.Context, userId string) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	container := GetRtcMemberContainer(appkey, userId)
	if container != nil {
		roomIds := []string{}
		container.ForeachRooms(func(room *RtcMemberRoom) {
			if room.RtcState != pbobjs.RtcState_RtcIncoming {
				roomIds = append(roomIds, room.RoomId)
			}
		})
		for _, roomId := range roomIds {
			bases.AsyncRpcCall(ctx, "rtc_hangup", roomId, &pbobjs.Nil{})
		}
	}
	return errs.IMErrorCode_SUCCESS
}
