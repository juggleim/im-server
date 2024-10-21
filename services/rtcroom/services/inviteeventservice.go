package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/rtcroom/storages/models"
	"time"
)

func HandleInviteEvent(ctx context.Context, req *pbobjs.RtcInviteReq) errs.IMErrorCode {
	switch req.InviteType {
	case pbobjs.InviteType_RtcInvite:
		return doRtcInvite(ctx, req)
	case pbobjs.InviteType_RtcAccept:
		return doRtcAccept(ctx, req)
	case pbobjs.InviteType_RtcReject:
		return doRtcReject(ctx, req)
	case pbobjs.InviteType_RtcCancel:
		return doRtcCancel(ctx, req)
	}
	return errs.IMErrorCode_SUCCESS
}

func doRtcInvite(ctx context.Context, req *pbobjs.RtcInviteReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_NOTEXIST
	}
	isMember := container.MemberExist(userId)
	if !isMember {
		return errs.IMErrorCode_RTCROOM_NOTMEMBER
	}
	needInviteIds := []string{}
	for _, targetId := range req.TargetIds {
		_, exist := container.GetMember(targetId)
		var member *models.RtcRoomMember
		if !exist {
			member = &models.RtcRoomMember{
				RoomId:    roomId,
				MemberId:  targetId,
				RtcState:  pbobjs.RtcState_RtcIncoming,
				InviterId: userId,
				CallTime:  time.Now().UnixMilli(),
			}
			container.ForceJoinRoom(member)
			needInviteIds = append(needInviteIds, targetId)
		}
	}
	//send event to targets
	SendInviteEvent(ctx, needInviteIds, &pbobjs.RtcInviteEvent{
		InviteType: pbobjs.InviteType_RtcInvite,
		TargetUser: &pbobjs.UserInfo{
			UserId: userId,
		},
		Room: &pbobjs.RtcRoom{
			RoomId: container.RoomId,
			Owner:  container.Owner,
		},
	})
	return errs.IMErrorCode_SUCCESS
}

func doRtcAccept(ctx context.Context, req *pbobjs.RtcInviteReq) errs.IMErrorCode {
	return errs.IMErrorCode_SUCCESS
}

func doRtcReject(ctx context.Context, req *pbobjs.RtcInviteReq) errs.IMErrorCode {
	return errs.IMErrorCode_SUCCESS
}

func doRtcCancel(ctx context.Context, req *pbobjs.RtcInviteReq) errs.IMErrorCode {
	return errs.IMErrorCode_SUCCESS
}

func SendInviteEvent(ctx context.Context, targetIds []string, event *pbobjs.RtcInviteEvent) {
	for _, targetId := range targetIds {
		msg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "rtc_invite_event", event)
		msg.Qos = 0
		bases.UnicastRouteWithNoSender(msg)
	}
}
