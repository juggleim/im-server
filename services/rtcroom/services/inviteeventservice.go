package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
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

	container, succ := getRtcRoomContainerWithInit(appkey, roomId, userId, req.RoomType)
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
	}
	if container.RoomType == pbobjs.RtcRoomType_OneOne {
		if !succ {
			return errs.IMErrorCode_RTCROOM_ROOMHASEXIST
		}
		grabCode := MemberGrabState(ctx, userId, &pbobjs.MemberState{
			RoomId:   roomId,
			RoomType: container.RoomType,
			MemberId: userId,
			DeviceId: deviceId,
			RtcState: pbobjs.RtcState_RtcOutgoing,
		})
		if grabCode != errs.IMErrorCode_SUCCESS {
			return grabCode
		}
		fromMember := &models.RtcRoomMember{
			RoomId:         roomId,
			RoomType:       container.RoomType,
			OwnerId:        userId,
			MemberId:       userId,
			DeviceId:       deviceId,
			RtcState:       pbobjs.RtcState_RtcOutgoing,
			LatestPingTime: time.Now().UnixMilli(),
			AppKey:         appkey,
		}
		container.JoinRoom(fromMember)
		memberStorage := storages.NewRtcRoomMemberStorage()
		memberStorage.Upsert(*fromMember)

		//target
		targetId := req.TargetIds[0]
		targetMember := &models.RtcRoomMember{
			RoomId:         roomId,
			RoomType:       container.RoomType,
			OwnerId:        container.Owner.UserId,
			MemberId:       targetId,
			RtcState:       pbobjs.RtcState_RtcIncoming,
			InviterId:      userId,
			CallTime:       time.Now().UnixMilli(),
			LatestPingTime: time.Now().UnixMilli(),
			AppKey:         appkey,
		}
		container.JoinRoom(targetMember)
		memberStorage.Upsert(*targetMember)
		//notify other node
		MemberSyncState(ctx, targetId, &pbobjs.SyncMemberStateReq{
			Member: &pbobjs.MemberState{
				RoomId:   roomId,
				RoomType: container.RoomType,
				MemberId: targetId,
				RtcState: pbobjs.RtcState_RtcIncoming,
			},
		})
		//send room event
		SendRoomEvent(ctx, targetId, &pbobjs.RtcRoomEvent{
			RoomEventType: pbobjs.RtcRoomEventType_RtcJoin,
			Member: &pbobjs.RtcMember{
				Member: &pbobjs.UserInfo{
					UserId: targetId,
				},
				RtcState: pbobjs.RtcState_RtcIncoming,
				Inviter: &pbobjs.UserInfo{
					UserId: userId,
				},
			},
			Room: &pbobjs.RtcRoom{
				RoomId:   roomId,
				RoomType: req.RoomType,
			},
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func MemberGrabState(ctx context.Context, targetId string, req *pbobjs.MemberState) errs.IMErrorCode {
	code, _, err := bases.SyncRpcCall(ctx, "rtc_grab_member", targetId, req, nil)
	if err != nil {
		return errs.IMErrorCode_RTCINVITE_BUSY
	}
	return code
}

func MemberSyncState(ctx context.Context, targetId string, req *pbobjs.SyncMemberStateReq) {
	bases.AsyncRpcCall(ctx, "rtc_sync_member", targetId, req)
}
