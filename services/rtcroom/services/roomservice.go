package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/rtcroom/storages"
	"im-server/services/rtcroom/storages/models"
	"strings"
	"time"

	"github.com/rfyiamcool/go-timewheel"
)

var (
	rtcroomCache *caches.LruCache
	rtcroomLocks *tools.SegmentatedLocks
	checkTimer   *timewheel.TimeWheel
)

const (
	PingTimeout       int64 = 30
	PingCheckInterval int64 = 10
)

func init() {
	rtcroomCache = caches.NewLruCacheWithReadTimeout(10000, nil, time.Hour)
	rtcroomLocks = tools.NewSegmentatedLocks(128)
	checkTimer, _ = timewheel.NewTimeWheel(1*time.Second, 360)
	checkTimer.Start()
}

type RtcRoomStatus int

const (
	RtcRoomStatus_Normal   = 0
	RtcRoomStatus_Destroy  = 1
	RtcRoomStatus_NotExist = 2
)

type RtcRoomContainer struct {
	Appkey   string
	RoomId   string
	RoomType pbobjs.RtcRoomType
	Owner    *pbobjs.UserInfo
	Status   RtcRoomStatus // 0:normal; 1: destroy; 2: not exist;

	Members map[string]*models.RtcRoomMember
}

func (container *RtcRoomContainer) UpdPingTime(memberId string) errs.IMErrorCode {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	if member, exist := container.Members[memberId]; exist {
		member.LatestPingTime = time.Now().UnixMilli()
		return errs.IMErrorCode_SUCCESS
	} else {
		return errs.IMErrorCode_RTCROOM_NOTMEMBER
	}
}

func (container *RtcRoomContainer) ForceJoinRoom(member *models.RtcRoomMember) {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	container.innerJoinRoom(member, true)
}

func (container *RtcRoomContainer) JoinRoom(member *models.RtcRoomMember) errs.IMErrorCode {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	return container.innerJoinRoom(member, false)
}

func (container *RtcRoomContainer) innerJoinRoom(member *models.RtcRoomMember, isForce bool) errs.IMErrorCode {
	if !isForce {
		if oldMember, exist := container.Members[member.MemberId]; exist {
			oldMember.LatestPingTime = time.Now().UnixMilli()
			return errs.IMErrorCode_RTCROOM_HASMEMBER
		}
	}
	member.LatestPingTime = time.Now().UnixMilli()
	container.Members[member.MemberId] = member
	// appkey := container.Appkey
	// roomId := container.RoomId
	// memberId := member.MemberId
	// checkTimer.Add(time.Duration(PingCheckInterval)*time.Second, func() {
	// 	checkRtcMemberTimeout(appkey, roomId, memberId)
	// })
	return errs.IMErrorCode_SUCCESS
}

func checkRtcMemberTimeout(appkey, roomId, memberId string) {
	container, exist := getRtcRoomContainer(appkey, roomId)
	if exist && container != nil {
		member, exist := container.GetMember(memberId)
		if exist {
			curr := time.Now().UnixMilli()
			if curr-member.LatestPingTime > (PingTimeout * 1000) {
				innerQuitRtcRoom(appkey, roomId, memberId)
			} else {
				// checkTimer.Add(time.Duration(PingCheckInterval)*time.Second, func() {
				// 	checkRtcMemberTimeout(appkey, roomId, memberId)
				// })
			}
		}
	}
}

func (container *RtcRoomContainer) QuitRoom(memberId string) (errs.IMErrorCode, *models.RtcRoomMember) {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if member, exist := container.Members[memberId]; exist {
		delete(container.Members, memberId)
		return errs.IMErrorCode_SUCCESS, member
	} else {
		return errs.IMErrorCode_RTCROOM_NOTMEMBER, nil
	}
}

func (container *RtcRoomContainer) ForeachMembers(f func(member *models.RtcRoomMember)) {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	for _, member := range container.Members {
		f(member)
	}
}

func (container *RtcRoomContainer) MemberExist(memberId string) bool {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	_, exist := container.Members[memberId]
	return exist
}

func (container *RtcRoomContainer) RemoveMember(memberId string) (*models.RtcRoomMember, bool) {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	member, exist := container.Members[memberId]
	if exist {
		delete(container.Members, memberId)
		return member, true
	} else {
		return nil, false
	}
}

func (container *RtcRoomContainer) GetMember(memberId string) (*models.RtcRoomMember, bool) {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	member, exist := container.Members[memberId]
	return member, exist
}

func (container *RtcRoomContainer) ExistAndChgState(memberId string, state pbobjs.RtcState) bool {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	member, exist := container.Members[memberId]
	if exist {
		member.RtcState = state
	}
	return exist
}

func (container *RtcRoomContainer) CompareAndSetState(memberId string, except pbobjs.RtcState, state pbobjs.RtcState) bool {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	member, exist := container.Members[memberId]
	if exist && member.RtcState == except {
		member.RtcState = state
		return true
	}
	return false
}

func (container *RtcRoomContainer) UpdMemberState(memberId string, state pbobjs.RtcState, deviceId string) errs.IMErrorCode {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if member, exist := container.Members[memberId]; exist {
		member.RtcState = state
		member.DeviceId = deviceId
		member.LatestPingTime = time.Now().UnixMilli()
		return errs.IMErrorCode_SUCCESS
	} else {
		return errs.IMErrorCode_RTCROOM_NOTMEMBER
	}
}

func getRtcRoomContainer(appkey, roomId string) (*RtcRoomContainer, bool) {
	key := getRoomKey(appkey, roomId)
	if cacheContainer, exist := rtcroomCache.Get(key); exist {
		container := cacheContainer.(*RtcRoomContainer)
		if container.Status == RtcRoomStatus_Normal {
			return container, true
		}
		return container, false
	} else {
		lock := rtcroomLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if cacheContainer, exist := rtcroomCache.Get(key); exist {
			container := cacheContainer.(*RtcRoomContainer)
			if container.Status == RtcRoomStatus_Normal {
				return container, true
			}
			return container, false
		} else {
			container := getRtcRoomContainerFromDb(appkey, roomId)
			rtcroomCache.Add(key, container)
			return container, container.Status == RtcRoomStatus_Normal
		}
	}
}

func getRtcRoomContainerWithInit(appkey, roomId, ownerId string, roomType pbobjs.RtcRoomType) (*RtcRoomContainer, bool) {
	key := getRoomKey(appkey, roomId)
	if cacheContainer, exist := rtcroomCache.Get(key); exist {
		container := cacheContainer.(*RtcRoomContainer)
		if container.Status == RtcRoomStatus_Normal {
			return container, false
		}
	}
	lock := rtcroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if cacheContainer, exist := rtcroomCache.Get(key); exist {
		container := cacheContainer.(*RtcRoomContainer)
		if container.Status == RtcRoomStatus_Normal {
			return container, false
		}
	}
	container := getRtcRoomContainerFromDb(appkey, roomId)
	if container.Status == RtcRoomStatus_Normal {
		rtcroomCache.Add(key, container)
		return container, false
	}
	container = &RtcRoomContainer{
		Appkey:   appkey,
		RoomId:   roomId,
		RoomType: roomType,
		Owner: &pbobjs.UserInfo{
			UserId: ownerId,
		},
		Status:  RtcRoomStatus_Normal,
		Members: make(map[string]*models.RtcRoomMember),
	}
	rtcroomCache.Add(key, container)
	return container, true
}

func getRtcRoomContainerFromDb(appkey, roomId string) *RtcRoomContainer {
	container := &RtcRoomContainer{
		Appkey:  appkey,
		RoomId:  roomId,
		Members: make(map[string]*models.RtcRoomMember),
	}
	storage := storages.NewRtcRoomStorage()
	room, err := storage.FindById(appkey, roomId)
	if err == nil && room != nil {
		container.Status = RtcRoomStatus_Normal
		//init rtc member relations
		memberStorage := storages.NewRtcRoomMemberStorage()
		var startId int64 = 0
		var limit int64 = 1000
		curr := time.Now().UnixMilli()
		for {
			members, err := memberStorage.QueryMembers(appkey, roomId, startId, limit)
			if err != nil {
				break
			}
			for _, member := range members {
				member.LatestPingTime = curr
				container.Members[member.MemberId] = member
				// startId = member.ID
				// memberId := member.MemberId
				// checkTimer.Add(time.Duration(PingCheckInterval)*time.Second, func() {
				// 	checkRtcMemberTimeout(appkey, roomId, memberId)
				// })
			}
			if len(members) < int(limit) {
				break
			}
		}
	} else {
		container.Status = RtcRoomStatus_NotExist
	}
	return container
}

func getRoomKey(appkey, roomId string) string {
	return strings.Join([]string{appkey, roomId}, "_")
}

func CreateRtcRoom(ctx context.Context, req *pbobjs.RtcRoomReq) (errs.IMErrorCode, *pbobjs.RtcRoom) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := bases.GetDeviceIdFromCtx(ctx)
	//add to cache
	container, exist := getRtcRoomContainer(appkey, req.RoomId)
	if exist {
		return errs.IMErrorCode_RTCROOM_ROOMHASEXIST, generatePbRtcRoom(container)
	}
	container, succ := initRtcRoomContainer(appkey, req.RoomId, userId, req.RoomType)
	if succ {
		// add to db
		storage := storages.NewRtcRoomStorage()
		err := storage.Create(models.RtcRoom{
			RoomId:   req.RoomId,
			RoomType: req.RoomType,
			OwnerId:  userId,
			AppKey:   appkey,
		})
		if err != nil {
			logs.WithContext(ctx).Errorf("create rtc room failed:%v", err)
		}
	} else {
		return errs.IMErrorCode_RTCROOM_ROOMHASEXIST, generatePbRtcRoom(container)
	}
	container.JoinRoom(&models.RtcRoomMember{
		RoomId:   req.RoomId,
		MemberId: userId,
		DeviceId: deviceId,
		RtcState: req.JoinMember.RtcState,
		AppKey:   appkey,
	})
	memberStorage := storages.NewRtcRoomMemberStorage()
	err := memberStorage.Upsert(models.RtcRoomMember{
		RoomId:   req.RoomId,
		MemberId: userId,
		DeviceId: deviceId,
		RtcState: req.JoinMember.RtcState,
		AppKey:   appkey,
	})
	if err != nil {
		logs.WithContext(ctx).Errorf("join rtc room failed:%v", err)
	}
	return errs.IMErrorCode_SUCCESS, generatePbRtcRoom(container)
}

func DestroyRtcRoom(ctx context.Context) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := bases.GetTargetIdFromCtx(ctx)
	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST
	}
	storage := storages.NewRtcRoomStorage()
	err := storage.Delete(appkey, roomId)
	if err == nil {
		memberStorage := storages.NewRtcRoomMemberStorage()
		err = memberStorage.DeleteByRoomId(appkey, roomId)
		if err == nil {
			container.ForeachMembers(func(member *models.RtcRoomMember) {
				fmt.Println("member:", member.MemberId)
			})
			rtcroomCache.Remove(getRoomKey(appkey, roomId))
		} else {
			logs.WithContext(ctx).Errorf("failed to clean rtc room members:%v", err)
		}
	} else {
		logs.WithContext(ctx).Errorf("failed to del rtc room:%v", err)
	}

	return errs.IMErrorCode_SUCCESS
}

func generatePbRtcRoom(container *RtcRoomContainer) *pbobjs.RtcRoom {
	members := []*pbobjs.RtcMember{}
	container.ForeachMembers(func(member *models.RtcRoomMember) {
		members = append(members, &pbobjs.RtcMember{
			Member: &pbobjs.UserInfo{
				UserId: member.MemberId,
			},
			RtcState:    member.RtcState,
			CallTime:    member.CallTime,
			ConnectTime: member.ConnectTime,
			HangupTime:  member.HangupTime,
			Inviter: &pbobjs.UserInfo{
				UserId: member.InviterId,
			},
		})
	})
	return &pbobjs.RtcRoom{
		RoomId:  container.RoomId,
		Owner:   container.Owner,
		Members: members,
	}
}

func initRtcRoomContainer(appkey, roomId, ownerId string, roomType pbobjs.RtcRoomType) (*RtcRoomContainer, bool) {
	key := getRoomKey(appkey, roomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	var container *RtcRoomContainer
	if cacheContainer, exist := rtcroomCache.Get(key); exist {
		container = cacheContainer.(*RtcRoomContainer)
	}
	if container == nil || container.Status != RtcRoomStatus_Normal {
		container = &RtcRoomContainer{
			Appkey:   appkey,
			RoomId:   roomId,
			RoomType: roomType,
			Owner: &pbobjs.UserInfo{
				UserId: ownerId,
			},
			Members: make(map[string]*models.RtcRoomMember),
		}
		rtcroomCache.Add(key, container)
		return container, true
	} else {
		return container, false
	}
}

func JoinRtcRoom(ctx context.Context, req *pbobjs.RtcRoomReq) (errs.IMErrorCode, *pbobjs.RtcRoom) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := bases.GetDeviceIdFromCtx(ctx)
	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST, nil
	}
	if container.MemberExist(userId) {
		return errs.IMErrorCode_RTCROOM_HASMEMBER, generatePbRtcRoom(container)
	}
	member := &models.RtcRoomMember{
		RoomId:         roomId,
		MemberId:       userId,
		DeviceId:       deviceId,
		RtcState:       req.JoinMember.RtcState,
		AppKey:         appkey,
		LatestPingTime: time.Now().UnixMilli(),
	}
	//add to cache
	code := container.JoinRoom(member)
	if code != errs.IMErrorCode_SUCCESS {
		return code, generatePbRtcRoom(container)
	}
	//add to db
	storage := storages.NewRtcRoomMemberStorage()
	err := storage.Upsert(*member)
	if err != nil {
		logs.WithContext(ctx).Errorf("failed to join rtc room. err:%v", err)
	}
	return errs.IMErrorCode_SUCCESS, generatePbRtcRoom(container)
}

func QuitRtcRoom(ctx context.Context) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)

	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST
	}
	if !container.MemberExist(userId) {
		return errs.IMErrorCode_RTCROOM_NOTMEMBER
	}
	memberStorage := storages.NewRtcRoomMemberStorage()

	delMember, succ := container.RemoveMember(userId)
	if succ {
		//delete from db
		memberStorage.Delete(appkey, roomId, userId)
		//sync other node
		MemberSyncState(ctx, delMember.MemberId, &pbobjs.SyncMemberStateReq{
			IsDelete: true,
			Member: &pbobjs.MemberState{
				RoomId:   roomId,
				RoomType: container.RoomType,
				MemberId: delMember.MemberId,
				DeviceId: delMember.DeviceId,
			},
		})
	}
	if container.RoomType == pbobjs.RtcRoomType_OneOne {
		container.ForeachMembers(func(member *models.RtcRoomMember) {
			MemberSyncState(ctx, member.MemberId, &pbobjs.SyncMemberStateReq{
				IsDelete: true,
				Member: &pbobjs.MemberState{
					RoomId:   roomId,
					RoomType: container.RoomType,
					MemberId: member.MemberId,
					DeviceId: member.DeviceId,
				},
			})
			SendRoomEvent(ctx, member.MemberId, &pbobjs.RtcRoomEvent{
				RoomEventType: pbobjs.RtcRoomEventType_RtcDestroy,
				Room: &pbobjs.RtcRoom{
					RoomType: container.RoomType,
					RoomId:   container.RoomId,
				},
			})
		})
		//destroy room
		memberStorage.DeleteByRoomId(appkey, roomId)
		roomStorage := storages.NewRtcRoomStorage()
		roomStorage.Delete(appkey, roomId)
	} else {
		container.ForeachMembers(func(member *models.RtcRoomMember) {
			SendRoomEvent(ctx, member.MemberId, &pbobjs.RtcRoomEvent{
				RoomEventType: pbobjs.RtcRoomEventType_RtcQuit,
				Member: &pbobjs.RtcMember{
					Member: &pbobjs.UserInfo{
						UserId: userId,
					},
				},
				Room: &pbobjs.RtcRoom{
					RoomType: container.RoomType,
					RoomId:   container.RoomId,
				},
			})
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func innerQuitRtcRoom(appkey, roomId, userId string) errs.IMErrorCode {
	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST
	}
	code, _ := container.QuitRoom(userId)
	storage := storages.NewRtcRoomMemberStorage()
	err := storage.Delete(appkey, roomId, userId)
	if err != nil {
		logs.WithContext(context.Background()).Errorf("failed to delete rtc member:%v", err)
	}
	return code
}

func RtcPing(ctx context.Context) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST
	}
	code := container.UpdPingTime(userId)
	if code != errs.IMErrorCode_SUCCESS {
		return code
	}
	storage := storages.NewRtcRoomMemberStorage()
	storage.RefreshPingTime(appkey, roomId, userId)
	return errs.IMErrorCode_SUCCESS
}

func QryRtcRoom(ctx context.Context, roomId string) (errs.IMErrorCode, *pbobjs.RtcRoom) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST, nil
	}
	return errs.IMErrorCode_SUCCESS, generatePbRtcRoom(container)
}

func UpdRtcMemberState(ctx context.Context, req *pbobjs.RtcMember) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := bases.GetDeviceIdFromCtx(ctx)
	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST
	}
	if !container.MemberExist(userId) {
		return errs.IMErrorCode_RTCROOM_NOTMEMBER
	}

	code := MemberGrabState(ctx, userId, &pbobjs.MemberState{
		RoomId:   roomId,
		RoomType: container.RoomType,
		MemberId: userId,
		DeviceId: deviceId,
		RtcState: req.RtcState,
	})
	if code != errs.IMErrorCode_SUCCESS {
		return code
	}
	code = container.UpdMemberState(userId, req.RtcState, deviceId)
	if code != errs.IMErrorCode_SUCCESS {
		return code
	}
	storage := storages.NewRtcRoomMemberStorage()
	err := storage.UpdateState(appkey, roomId, userId, req.RtcState, deviceId)
	if err != nil {
		return errs.IMErrorCode_RTCROOM_UPDATEFAILED
	}
	container.ForeachMembers(func(member *models.RtcRoomMember) {
		if member.MemberId != userId && member.RtcState != pbobjs.RtcState_RtcIncoming {
			SendRoomEvent(ctx, member.MemberId, &pbobjs.RtcRoomEvent{
				RoomEventType: pbobjs.RtcRoomEventType_RtcStateChg,
				Room: &pbobjs.RtcRoom{
					RoomId:   container.RoomId,
					RoomType: container.RoomType,
					Owner:    container.Owner,
				},
				Member: &pbobjs.RtcMember{
					Member: &pbobjs.UserInfo{
						UserId: userId,
					},
					RtcState: req.RtcState,
				},
			})
		}
	})

	return errs.IMErrorCode_SUCCESS
}

func SendRoomEvent(ctx context.Context, targetId string, event *pbobjs.RtcRoomEvent) {
	msg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "rtc_room_event", event)
	msg.Qos = 0
	bases.UnicastRouteWithNoSender(msg)
}