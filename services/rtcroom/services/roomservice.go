package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
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
	CallTimeout       int64 = 30
	PingTimeout       int64 = 10
	PingCheckInterval int64 = 5
)

func init() {
	rtcroomCache = caches.NewLruCacheWithReadTimeout("rtcroom_cache", 10000, nil, time.Hour)
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
	Appkey       string
	RoomId       string
	RoomType     pbobjs.RtcRoomType
	RtcChannel   pbobjs.RtcChannel
	RtcMediaType pbobjs.RtcMediaType
	Owner        *pbobjs.UserInfo
	CreatedTime  int64
	AcceptedTime int64
	Status       RtcRoomStatus // 0:normal; 1: destroy; 2: not exist;

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
	return errs.IMErrorCode_SUCCESS
}

func checkRtcMemberTimeout(appkey, roomId string) {
	container, exist := getRtcRoomContainer(appkey, roomId)
	if exist && container != nil {
		if container.MemberCount() <= 0 {
			rtcroomCache.Remove(getRoomKey(appkey, roomId))
			return
		}
		curr := time.Now().UnixMilli()
		pingTimeoutMemberIds := []string{}
		callTimeOutMemberIds := []string{}
		container.ForeachMembers(func(member *models.RtcRoomMember) {
			if member.RtcState == pbobjs.RtcState_RtcIncoming || member.RtcState == pbobjs.RtcState_RtcOutgoing {
				if curr-member.LatestPingTime > (CallTimeout * 1000) {
					callTimeOutMemberIds = append(callTimeOutMemberIds, member.MemberId)
				}
			} else {
				if curr-member.LatestPingTime > (PingTimeout * 1000) {
					pingTimeoutMemberIds = append(pingTimeoutMemberIds, member.MemberId)
				}
			}
		})
		for _, memberId := range callTimeOutMemberIds {
			innerQuitRtcRoom(bases.CreateRpcCtx(appkey, memberId), appkey, roomId, memberId, pbobjs.RtcRoomQuitReason_CallTimeout, true)
			if container.RoomType == pbobjs.RtcRoomType_OneOne {
				break
			}
		}
		for _, memberId := range pingTimeoutMemberIds {
			innerQuitRtcRoom(bases.CreateRpcCtx(appkey, memberId), appkey, roomId, memberId, pbobjs.RtcRoomQuitReason_PingTimeout, true)
			if container.RoomType == pbobjs.RtcRoomType_OneOne {
				break
			}
		}
		checkTimer.Add(time.Duration(PingCheckInterval)*time.Second, func() {
			checkRtcMemberTimeout(appkey, roomId)
		})
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

func (container *RtcRoomContainer) MemberCount() int {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	return len(container.Members)
}

func (container *RtcRoomContainer) MemberExist(memberId string) bool {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	_, exist := container.Members[memberId]
	return exist
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

func (container *RtcRoomContainer) UpdAcceptedTime(acceptedTime int64) {
	key := getRoomKey(container.Appkey, container.RoomId)
	lock := rtcroomLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	container.AcceptedTime = acceptedTime
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
			if container.Status == RtcRoomStatus_Normal {
				checkTimer.Add(time.Duration(PingCheckInterval)*time.Second, func() {
					checkRtcMemberTimeout(appkey, roomId)
				})
			}
			return container, container.Status == RtcRoomStatus_Normal
		}
	}
}

func createRtcRoomContainer2Cache(ctx context.Context, appkey string, room *models.RtcRoom) (*RtcRoomContainer, bool) {
	key := getRoomKey(appkey, room.RoomId)
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
	container := getRtcRoomContainerFromDb(appkey, room.RoomId)
	if container.Status == RtcRoomStatus_Normal {
		rtcroomCache.Add(key, container)
		checkTimer.Add(time.Duration(PingCheckInterval)*time.Second, func() {
			checkRtcMemberTimeout(appkey, room.RoomId)
		})
		return container, false
	}
	container = &RtcRoomContainer{
		Appkey:       appkey,
		RoomId:       room.RoomId,
		RoomType:     room.RoomType,
		RtcChannel:   room.RtcChannel,
		RtcMediaType: room.RtcMediaType,
		CreatedTime:  time.Now().UnixMilli(),
		Owner:        commonservices.GetTargetDisplayUserInfo(ctx, room.OwnerId),
		Status:       RtcRoomStatus_Normal,
		Members:      make(map[string]*models.RtcRoomMember),
	}
	rtcroomCache.Add(key, container)
	checkTimer.Add(time.Duration(PingCheckInterval)*time.Second, func() {
		checkRtcMemberTimeout(appkey, room.RoomId)
	})
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
		container.RoomType = room.RoomType
		container.RtcChannel = room.RtcChannel
		container.RtcMediaType = room.RtcMediaType
		container.Owner = &pbobjs.UserInfo{
			UserId: room.OwnerId,
		}
		container.CreatedTime = room.CreatedTime
		container.AcceptedTime = room.AcceptedTime
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
				startId = member.ID
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

	rtcRoom := &models.RtcRoom{
		RoomId:       req.RoomId,
		RoomType:     req.RoomType,
		RtcChannel:   req.RtcChannel,
		RtcMediaType: req.RtcMediaType,
		OwnerId:      userId,
		AppKey:       appkey,
	}
	container, succ := createRtcRoomContainer2Cache(ctx, appkey, rtcRoom)
	if succ {
		// add to db
		storage := storages.NewRtcRoomStorage()
		err := storage.Create(*rtcRoom)
		if err != nil {
			logs.WithContext(ctx).Errorf("create rtc room failed:%v", err)
		}
	} else {
		return errs.IMErrorCode_RTCROOM_ROOMHASEXIST, generatePbRtcRoom(ctx, container)
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
	return errs.IMErrorCode_SUCCESS, generatePbRtcRoom(ctx, container)
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

func generatePbRtcRoom(ctx context.Context, container *RtcRoomContainer) *pbobjs.RtcRoom {
	members := []*pbobjs.RtcMember{}
	container.ForeachMembers(func(member *models.RtcRoomMember) {
		members = append(members, &pbobjs.RtcMember{
			Member:      commonservices.GetTargetDisplayUserInfo(ctx, member.MemberId),
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
		RoomId:       container.RoomId,
		Owner:        container.Owner,
		RtcChannel:   container.RtcChannel,
		RtcMediaType: container.RtcMediaType,
		RoomType:     container.RoomType,
		Members:      members,
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
		return errs.IMErrorCode_RTCROOM_HASMEMBER, generatePbRtcRoom(ctx, container)
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
		return code, generatePbRtcRoom(ctx, container)
	}
	//add to db
	storage := storages.NewRtcRoomMemberStorage()
	err := storage.Upsert(*member)
	if err != nil {
		logs.WithContext(ctx).Errorf("failed to join rtc room. err:%v", err)
	}
	return errs.IMErrorCode_SUCCESS, generatePbRtcRoom(ctx, container)
}

func QuitRtcRoom(ctx context.Context) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)

	return innerQuitRtcRoom(ctx, appkey, roomId, userId, pbobjs.RtcRoomQuitReason_Active, true)
}

func innerQuitRtcRoom(ctx context.Context, appkey, roomId, userId string, quitReason pbobjs.RtcRoomQuitReason, isSendEvent bool) errs.IMErrorCode {
	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST
	}
	storage := storages.NewRtcRoomMemberStorage()
	code, delMember := container.QuitRoom(userId)
	if code == errs.IMErrorCode_SUCCESS {
		//delete from db
		storage.Delete(appkey, roomId, userId)
		//sync other node
		MemberSyncState(ctx, userId, &pbobjs.SyncMemberStateReq{
			IsDelete: true,
			Member: &pbobjs.MemberState{
				RoomId:   roomId,
				RoomType: container.RoomType,
				MemberId: userId,
				DeviceId: delMember.DeviceId,
			},
		})
		eventTime := time.Now().UnixMilli()
		//send room event
		if isSendEvent {
			SendRoomEvent(ctx, userId, &pbobjs.RtcRoomEvent{
				RoomEventType: pbobjs.RtcRoomEventType_RtcQuit,
				Members: []*pbobjs.RtcMember{
					{
						Member: &pbobjs.UserInfo{
							UserId: userId,
						},
					},
				},
				Room: &pbobjs.RtcRoom{
					RoomType: container.RoomType,
					RoomId:   container.RoomId,
				},
				Reason:    quitReason,
				EventTime: eventTime,
			})
		}

		if container.RoomType == pbobjs.RtcRoomType_OneOne {
			var calleeId string = ""
			container.ForeachMembers(func(member *models.RtcRoomMember) {
				MemberSyncState(ctx, member.MemberId, &pbobjs.SyncMemberStateReq{
					IsDelete: true,
					Member: &pbobjs.MemberState{
						RoomId:   roomId,
						RoomType: container.RoomType,
						MemberId: userId,
						DeviceId: delMember.DeviceId,
					},
				})
				if isSendEvent {
					SendRoomEvent(ctx, member.MemberId, &pbobjs.RtcRoomEvent{
						RoomEventType: pbobjs.RtcRoomEventType_RtcQuit,
						Members: []*pbobjs.RtcMember{
							{
								Member: &pbobjs.UserInfo{
									UserId: userId,
								},
							},
						},
						Room: &pbobjs.RtcRoom{
							RoomType: container.RoomType,
							RoomId:   container.RoomId,
						},
						Reason:    quitReason,
						EventTime: eventTime,
					})
				}
				if member.MemberId != container.Owner.UserId {
					calleeId = member.MemberId
				}
			})
			//send notify msg
			var reason CallFinishReasonType
			var duration int64 = 0
			if container.AcceptedTime > 0 {
				reason = CallFinishReasonType_Complete
				duration = eventTime - container.AcceptedTime
			} else {
				reason = CallFinishReasonType_NoAnswer
			}
			SendFinishNtf(ctx, container.Owner.UserId, calleeId, reason, duration, container.RtcMediaType)
			//destroy room
			storage.DeleteByRoomId(appkey, roomId)
			roomStorage := storages.NewRtcRoomStorage()
			roomStorage.Delete(appkey, roomId)
			rtcroomCache.Remove(getRoomKey(appkey, roomId))
		} else if container.RoomType == pbobjs.RtcRoomType_OneMore {
			if container.MemberCount() > 0 {
				if isSendEvent {
					container.ForeachMembers(func(member *models.RtcRoomMember) {
						SendRoomEvent(ctx, member.MemberId, &pbobjs.RtcRoomEvent{
							RoomEventType: pbobjs.RtcRoomEventType_RtcQuit,
							Members: []*pbobjs.RtcMember{
								{
									Member: &pbobjs.UserInfo{
										UserId: userId,
									},
								},
							},
							Room: &pbobjs.RtcRoom{
								RoomType: container.RoomType,
								RoomId:   container.RoomId,
							},
							Reason:    quitReason,
							EventTime: eventTime,
						})
					})
				}
			} else {
				//desctroy room
				storage.DeleteByRoomId(appkey, roomId)
				roomStorage := storages.NewRtcRoomStorage()
				roomStorage.Delete(appkey, roomId)
				rtcroomCache.Remove(getRoomKey(appkey, roomId))
			}
		}
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
	return errs.IMErrorCode_SUCCESS, generatePbRtcRoom(ctx, container)
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
	eventTime := time.Now().UnixMilli()
	container.ForeachMembers(func(member *models.RtcRoomMember) {
		if member.MemberId != userId && member.RtcState != pbobjs.RtcState_RtcIncoming {
			SendRoomEvent(ctx, member.MemberId, &pbobjs.RtcRoomEvent{
				RoomEventType: pbobjs.RtcRoomEventType_RtcStateChg,
				Room: &pbobjs.RtcRoom{
					RoomId:   container.RoomId,
					RoomType: container.RoomType,
					Owner:    container.Owner,
				},
				Members: []*pbobjs.RtcMember{
					{
						Member: &pbobjs.UserInfo{
							UserId: userId,
						},
						RtcState: req.RtcState,
					},
				},
				EventTime: eventTime,
			})
		}
	})

	return errs.IMErrorCode_SUCCESS
}

func SendRoomEvent(ctx context.Context, targetId string, event *pbobjs.RtcRoomEvent) {
	msg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "rtc_room_event", event)
	msg.Qos = 0
	msg.PublishType = int32(commonservices.PublishType_AllSessionExceptSelf)
	bases.UnicastRouteWithNoSender(msg)
}
