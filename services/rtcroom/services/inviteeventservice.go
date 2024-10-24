package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/rtcroom/storages"
	"im-server/services/rtcroom/storages/models"
	"time"
)

var (
	RtcPingTimeOut int64 = 30 * 1000
)

func RtcInvite(ctx context.Context, req *pbobjs.RtcInviteReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := req.RoomId
	userId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := bases.GetDeviceIdFromCtx(ctx)
	currTime := time.Now().UnixMilli()

	memberKey := getRtcMemberKey(appkey, userId)
	lock := memberLocks.GetLocks(memberKey)
	lock.Lock()
	defer lock.Unlock()

	memberStorage := storages.NewRtcRoomMemberStorage()
	memberRooms, err := memberStorage.QueryRoomsByMember(appkey, userId, 100)
	if err == nil {
		hasInRoom := false
		hasExpire := false
		for _, memberRoom := range memberRooms {
			if currTime-memberRoom.LatestPingTime < RtcPingTimeOut {
				hasInRoom = true
			} else {
				hasExpire = true
			}
		}
		if hasExpire {
			//del expire item
			go delExpireRtcMember(appkey, roomId, currTime-RtcPingTimeOut)
		}
		if req.RoomType == pbobjs.RtcRoomType_OneOne {
			if hasInRoom {
				return errs.IMErrorCode_RTCINVITE_BUSY
			}
		}
	}
	roomStorage := storages.NewRtcRoomStorage()
	err = roomStorage.Create(models.RtcRoom{
		RoomId:   roomId,
		RoomType: req.RoomType,
		OwnerId:  userId,
		AppKey:   appkey,
	})
	if err != nil {
		return errs.IMErrorCode_RTCROOM_CREATEROOMFAILED
	}
	memberStorage.Insert(models.RtcRoomMember{
		RoomId:         roomId,
		MemberId:       userId,
		DeviceId:       deviceId,
		RtcState:       pbobjs.RtcState_RtcOutgoing,
		LatestPingTime: currTime,
		AppKey:         appkey,
	})
	for _, targetId := range req.TargetIds {
		memberStorage.Insert(models.RtcRoomMember{
			RoomId:         roomId,
			MemberId:       targetId,
			InviterId:      userId,
			RtcState:       pbobjs.RtcState_RtcIncoming,
			LatestPingTime: currTime,
			AppKey:         appkey,
		})
		//send event
		SendInviteEvent(ctx, targetId, &pbobjs.RtcInviteEvent{
			InviteType: pbobjs.InviteType_RtcInvite,
			TargetUser: &pbobjs.UserInfo{
				UserId: userId,
			},
			Room: &pbobjs.RtcRoom{
				RoomId:   roomId,
				RoomType: req.RoomType,
			},
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func RtcDecline(ctx context.Context, req *pbobjs.RtcAnswerReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := req.RoomId
	userId := bases.GetRequesterIdFromCtx(ctx)
	memberKey := getRtcMemberKey(appkey, userId)
	lock := memberLocks.GetLocks(memberKey)
	lock.Lock()
	defer lock.Unlock()

	roomStorage := storages.NewRtcRoomStorage()
	room, err := roomStorage.FindById(appkey, roomId)
	if err != nil || room == nil {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST
	}
	memberStorage := storages.NewRtcRoomMemberStorage()
	if room.RoomType == pbobjs.RtcRoomType_OneOne {
		roomStorage.Delete(appkey, roomId)
		memberStorage.DeleteByRoomId(appkey, roomId)
	} else {
		memberStorage.Delete(appkey, roomId, userId)
	}
	SendInviteEvent(ctx, req.TargetId, &pbobjs.RtcInviteEvent{
		InviteType: pbobjs.InviteType_RtcDecline,
		TargetUser: &pbobjs.UserInfo{
			UserId: userId,
		},
		Room: &pbobjs.RtcRoom{
			RoomId:   roomId,
			RoomType: room.RoomType,
		},
	})
	return errs.IMErrorCode_SUCCESS
}

func RtcAccept(ctx context.Context, req *pbobjs.RtcAnswerReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := req.RoomId
	userId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := bases.GetDeviceIdFromCtx(ctx)
	memberKey := getRtcMemberKey(appkey, userId)
	lock := memberLocks.GetLocks(memberKey)
	lock.Lock()
	defer lock.Unlock()
	storage := storages.NewRtcRoomMemberStorage()
	memberRooms, err := storage.QueryRoomsByMember(appkey, userId, 10)
	if err == nil {
		var currRoom *models.RtcRoomMember
		for _, memberRoom := range memberRooms {
			if memberRoom.RoomId == roomId {
				currRoom = memberRoom
			} else {
				if memberRoom.RtcState == pbobjs.RtcState_RtcConnecting || memberRoom.RtcState == pbobjs.RtcState_RtcConnected {
					return errs.IMErrorCode_RTCINVITE_BUSY
				}
			}
		}
		if currRoom == nil {
			return errs.IMErrorCode_RTCINVITE_CANCEL
		} else {
			if currRoom.RtcState != pbobjs.RtcState_RtcIncoming {
				return errs.IMErrorCode_RTCINVITE_HASACCEPT
			}
		}
	}
	storage.UpdateState(appkey, roomId, userId, pbobjs.RtcState_RtcConnecting, deviceId)
	SendInviteEvent(ctx, req.TargetId, &pbobjs.RtcInviteEvent{
		InviteType: pbobjs.InviteType_RtcAccept,
		TargetUser: &pbobjs.UserInfo{
			UserId: userId,
		},
		Room: &pbobjs.RtcRoom{
			RoomId: roomId,
		},
	})
	return errs.IMErrorCode_SUCCESS
}

func RtcHangup(ctx context.Context, req *pbobjs.RtcAnswerReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := req.RoomId
	userId := bases.GetRequesterIdFromCtx(ctx)
	memberKey := getRtcMemberKey(appkey, userId)
	lock := memberLocks.GetLocks(memberKey)
	lock.Lock()
	defer lock.Unlock()

	roomStorage := storages.NewRtcRoomStorage()
	room, err := roomStorage.FindById(appkey, roomId)
	if err != nil || room == nil {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST
	}
	memberStorage := storages.NewRtcRoomMemberStorage()
	if room.RoomType == pbobjs.RtcRoomType_OneOne {
		memberStorage.DeleteByRoomId(appkey, roomId)
		roomStorage.Delete(appkey, roomId)
		SendInviteEvent(ctx, req.TargetId, &pbobjs.RtcInviteEvent{
			InviteType: pbobjs.InviteType_RtcHangup,
			TargetUser: &pbobjs.UserInfo{
				UserId: userId,
			},
			Room: &pbobjs.RtcRoom{
				RoomType: pbobjs.RtcRoomType_OneOne,
				RoomId:   roomId,
			},
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func SendInviteEvent(ctx context.Context, targetId string, event *pbobjs.RtcInviteEvent) {
	msg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "rtc_invite_event", event)
	msg.Qos = 0
	bases.UnicastRouteWithNoSender(msg)
}

func BatchSendInviteEvent(ctx context.Context, targetIds []string, event *pbobjs.RtcInviteEvent) {
	for _, targetId := range targetIds {
		msg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "rtc_invite_event", event)
		msg.Qos = 0
		bases.UnicastRouteWithNoSender(msg)
	}
}

func delExpireRtcMember(appkey, roomId string, baseTime int64) {
	storage := storages.NewRtcRoomMemberStorage()
	storage.DelteByRoomIdBaseTime(appkey, roomId, baseTime)
}
