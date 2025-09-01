package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/rtcroom/storages"
	"im-server/services/rtcroom/storages/models"
	"time"

	"google.golang.org/protobuf/proto"
)

func RtcInvite(ctx context.Context, req *pbobjs.RtcInviteReq) (errs.IMErrorCode, *pbobjs.RtcAuth) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := req.RoomId
	userId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := bases.GetDeviceIdFromCtx(ctx)
	if roomId == "" || len(req.TargetIds) <= 0 {
		return errs.IMErrorCode_RTCROOM_PARAMILLIGAL, nil
	}
	//check user block
	if req.RoomType == pbobjs.RtcRoomType_OneOne {
		targetId := req.TargetIds[0]
		code, resp, err := bases.SyncRpcCall(ctx, "check_block_user", targetId, &pbobjs.Nil{}, func() proto.Message {
			return &pbobjs.CheckBlockUserResp{}
		})
		if err == nil && code == errs.IMErrorCode_SUCCESS && resp != nil {
			checkResp, ok := resp.(*pbobjs.CheckBlockUserResp)
			if ok && checkResp.IsBlock {
				return errs.IMErrorCode_MSG_BLOCK, nil
			}
		}
	}
	//auth
	code, auth := GenerateAuth(appkey, userId, roomId, req.RtcChannel)
	if code != errs.IMErrorCode_SUCCESS {
		return code, auth
	}
	rtcRoom := &models.RtcRoom{
		RoomId:       req.RoomId,
		RoomType:     req.RoomType,
		RtcChannel:   req.RtcChannel,
		RtcMediaType: req.RtcMediaType,
		OwnerId:      userId,
		Ext:          req.Ext,
		AppKey:       appkey,
	}
	if req.AttachedConver != nil && req.AttachedConver.TargetId != "" {
		converId := commonservices.GetConversationId(userId, req.AttachedConver.TargetId, req.AttachedConver.ChannelType)
		rtcRoom.ConverId = &converId
		rtcRoom.ChannelType = req.AttachedConver.ChannelType
		rtcRoom.SubChannel = req.AttachedConver.SubChannel
	}
	container, succ := createRtcRoomContainer2Cache(ctx, appkey, rtcRoom)
	if succ {
		if rtcRoom.ConverId != nil && *rtcRoom.ConverId != "" {
			code, _, err := bases.SyncRpcCall(ctx, "rtc_bind_conver", *rtcRoom.ConverId, &pbobjs.RtcConverBindReq{
				Action:      pbobjs.RtcConverBindAction_RtcBind,
				ConverId:    *rtcRoom.ConverId,
				ChannelType: rtcRoom.ChannelType,
				SubChannel:  rtcRoom.SubChannel,
				RtcRoomId:   roomId,
			}, nil)
			if err == nil && code != errs.IMErrorCode_SUCCESS {
				rtcroomCache.Remove(getRoomKey(appkey, roomId))
				return errs.IMErrorCode_RTCROOM_BINDCONVERFAIL, nil
			}
		}
		// add to db
		storage := storages.NewRtcRoomStorage()
		err := storage.Create(*rtcRoom)
		if err != nil {
			logs.WithContext(ctx).Errorf("create rtc room failed:%v", err)
			rtcroomCache.Remove(getRoomKey(appkey, roomId))
			return errs.IMErrorCode_RTCROOM_ROOMHASEXIST, nil
		}
	}
	eventTime := time.Now().UnixMilli()
	if req.RoomType == pbobjs.RtcRoomType_OneOne {
		if !succ {
			return errs.IMErrorCode_RTCROOM_ROOMHASEXIST, nil
		}
		grabCode := MemberGrabState(ctx, userId, &pbobjs.MemberState{
			RoomId:   roomId,
			RoomType: container.RoomType,
			MemberId: userId,
			DeviceId: deviceId,
			RtcState: pbobjs.RtcState_RtcOutgoing,
		})
		if grabCode != errs.IMErrorCode_SUCCESS {
			return grabCode, nil
		}
		fromMember := &models.RtcRoomMember{
			RoomId:         roomId,
			RoomType:       container.RoomType,
			OwnerId:        userId,
			MemberId:       userId,
			DeviceId:       deviceId,
			RtcState:       pbobjs.RtcState_RtcOutgoing,
			LatestPingTime: eventTime,
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
			CallTime:       eventTime,
			LatestPingTime: eventTime,
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
		members := []*pbobjs.RtcMember{}
		container.ForeachMembers(func(member *models.RtcRoomMember) {
			members = append(members, &pbobjs.RtcMember{
				Member:   commonservices.GetTargetDisplayUserInfo(ctx, member.MemberId),
				RtcState: member.RtcState,
			})
		})
		//send invite event
		SendInviteEvent(ctx, targetId, &pbobjs.RtcInviteEvent{
			InviteType: pbobjs.InviteType_RtcInvite,
			User:       commonservices.GetTargetDisplayUserInfo(ctx, userId),
			Room: &pbobjs.RtcRoom{
				RoomType:     container.RoomType,
				RoomId:       container.RoomId,
				Owner:        container.Owner,
				RtcChannel:   container.RtcChannel,
				RtcMediaType: container.RtcMediaType,
				Ext:          container.Ext,
				Members:      members,
			},
			TargetUsers: []*pbobjs.UserInfo{
				commonservices.GetTargetDisplayUserInfo(ctx, targetId),
			},
			EventTime: eventTime,
		})
		//trigger push
		msg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "user_push", &pbobjs.DownMsg{
			TargetId:       userId,
			ChannelType:    pbobjs.ChannelType_Private,
			SenderId:       userId,
			MsgType:        msgdefines.InnerMsgType_VoiceCall,
			TargetUserInfo: commonservices.GetTargetDisplayUserInfo(ctx, userId),
			PushData: &pbobjs.PushData{
				IsVoip:       true,
				RtcRoomId:    roomId,
				RtcInviterId: userId,
				RtcRoomType:  int32(req.RoomType),
				RtcMediaType: int32(req.RtcMediaType),
			},
		})
		bases.UnicastRouteWithNoSender(msg)
	} else if req.RoomType == pbobjs.RtcRoomType_OneMore {
		memberStorage := storages.NewRtcRoomMemberStorage()
		if succ {
			grabCode := MemberGrabState(ctx, userId, &pbobjs.MemberState{
				RoomId:   roomId,
				RoomType: container.RoomType,
				MemberId: userId,
				DeviceId: deviceId,
				RtcState: pbobjs.RtcState_RtcConnecting,
			})
			if grabCode != errs.IMErrorCode_SUCCESS {
				return grabCode, nil
			}
			ownerMember := &models.RtcRoomMember{
				RoomId:         roomId,
				RoomType:       container.RoomType,
				OwnerId:        userId,
				MemberId:       userId,
				DeviceId:       deviceId,
				RtcState:       pbobjs.RtcState_RtcConnecting,
				LatestPingTime: eventTime,
				AppKey:         appkey,
			}
			container.JoinRoom(ownerMember)
			memberStorage.Upsert(*ownerMember)
		} else {
			if exist := container.MemberExist(userId); !exist {
				return errs.IMErrorCode_RTCROOM_NOTMEMBER, nil
			}
		}
		rtcRoom := &pbobjs.RtcRoom{
			RoomType:     container.RoomType,
			RoomId:       container.RoomId,
			Owner:        container.Owner,
			RtcChannel:   container.RtcChannel,
			RtcMediaType: container.RtcMediaType,
			Ext:          container.Ext,
			Members:      []*pbobjs.RtcMember{},
		}
		container.ForeachMembers(func(member *models.RtcRoomMember) {
			rtcRoom.Members = append(rtcRoom.Members, &pbobjs.RtcMember{
				Member:   commonservices.GetTargetDisplayUserInfo(ctx, member.MemberId),
				RtcState: member.RtcState,
			})
		})
		// add target members
		for _, targetId := range req.TargetIds {
			targetMember := &models.RtcRoomMember{
				RoomId:         roomId,
				RoomType:       container.RoomType,
				OwnerId:        container.Owner.UserId,
				MemberId:       targetId,
				RtcState:       pbobjs.RtcState_RtcIncoming,
				InviterId:      userId,
				CallTime:       eventTime,
				LatestPingTime: eventTime,
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
		}
		targetUsers := commonservices.GetTargetDisplayUserInfos(ctx, req.TargetIds)
		container.ForeachMembers(func(member *models.RtcRoomMember) {
			SendInviteEvent(ctx, member.MemberId, &pbobjs.RtcInviteEvent{
				InviteType:  pbobjs.InviteType_RtcInvite,
				User:        commonservices.GetTargetDisplayUserInfo(ctx, userId),
				Room:        rtcRoom,
				TargetUsers: targetUsers,
			})
		})
		//notify target conver
		if req.AttachedConver != nil && req.AttachedConver.TargetId != "" && req.AttachedConver.ChannelType == pbobjs.ChannelType_Group {
			syncMsg2Conver(ctx, container)
		}
	}
	return errs.IMErrorCode_SUCCESS, auth
}

func syncMsg2Conver(ctx context.Context, container *RtcRoomContainer) {
	if container.ConverId != nil && *container.ConverId != "" && container.ChannelType == pbobjs.ChannelType_Group {
		userId := bases.GetRequesterIdFromCtx(ctx)
		activedCallMsg := &msgdefines.ActivedCallMsg{
			RoomType:     int32(container.RoomType),
			RoomId:       container.RoomId,
			Owner:        &msgdefines.UserInfo{},
			RtcChannel:   int32(container.RtcChannel),
			RtcMediaType: int32(container.RtcMediaType),
			Members:      []*msgdefines.UserInfo{},
		}
		if container.Owner != nil {
			activedCallMsg.Owner.UserId = container.Owner.UserId
			activedCallMsg.Owner.Nickname = container.Owner.Nickname
			activedCallMsg.Owner.UserPortrait = container.Owner.UserPortrait
		}
		if container.MemberCount() > 0 {
			container.ForeachMembers(func(member *models.RtcRoomMember) {
				uInfo := commonservices.GetTargetDisplayUserInfo(ctx, member.MemberId)
				activedCallMsg.Members = append(activedCallMsg.Members, &msgdefines.UserInfo{
					UserId:       uInfo.UserId,
					Nickname:     uInfo.Nickname,
					UserPortrait: uInfo.UserPortrait,
				})
			})
		} else {
			activedCallMsg.Finished = true
			//un bind conver
			bases.AsyncRpcCall(ctx, "rtc_bind_conver", *container.ConverId, &pbobjs.RtcConverBindReq{
				Action:      pbobjs.RtcConverBindAction_RtcUnBind,
				ConverId:    *container.ConverId,
				ChannelType: container.ChannelType,
				SubChannel:  container.SubChannel,
				RtcRoomId:   container.RoomId,
			})
		}
		flag := msgdefines.SetCmdMsg(0)
		upMsg := &pbobjs.UpMsg{
			MsgType:    msgdefines.CmdMsgType_ActivedCall,
			MsgContent: tools.ToJsonBs(activedCallMsg),
			Flags:      flag,
			LifeTime:   10 * 60 * 1000,
		}
		commonservices.AsyncGroupMsg(ctx, userId, *container.ConverId, upMsg, &bases.ReGenerateSessionOption{})
	}
}

func RtcAccept(ctx context.Context) (errs.IMErrorCode, *pbobjs.RtcAuth) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	roomId := bases.GetTargetIdFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	deviceId := bases.GetDeviceIdFromCtx(ctx)

	container, exist := getRtcRoomContainer(appkey, roomId)
	if !exist {
		return errs.IMErrorCode_RTCROOM_ROOMNOTEXIST, nil
	}
	if !container.MemberExist(userId) {
		return errs.IMErrorCode_RTCROOM_NOTMEMBER, nil
	}
	code := MemberGrabState(ctx, userId, &pbobjs.MemberState{
		RoomId:   roomId,
		RoomType: container.RoomType,
		MemberId: userId,
		DeviceId: deviceId,
		RtcState: pbobjs.RtcState_RtcConnecting,
	})
	if code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	code = container.UpdMemberState(userId, pbobjs.RtcState_RtcConnecting, deviceId)
	if code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	acceptedTime := time.Now().UnixMilli()
	if container.RoomType == pbobjs.RtcRoomType_OneOne {
		//update accepted time
		container.UpdAcceptedTime(acceptedTime)
		roomStorage := storages.NewRtcRoomStorage()
		roomStorage.UpdateAcceptedTime(appkey, roomId, acceptedTime)
	}
	//update member state
	storage := storages.NewRtcRoomMemberStorage()
	err := storage.UpdateState(appkey, roomId, userId, pbobjs.RtcState_RtcConnecting, deviceId)
	if err != nil {
		return errs.IMErrorCode_RTCROOM_UPDATEFAILED, nil
	}
	container.ForeachMembers(func(member *models.RtcRoomMember) {
		SendInviteEvent(ctx, member.MemberId, &pbobjs.RtcInviteEvent{
			InviteType: pbobjs.InviteType_RtcAccept,
			User:       commonservices.GetTargetDisplayUserInfo(ctx, userId),
			Room: &pbobjs.RtcRoom{
				RoomType: container.RoomType,
				RoomId:   container.RoomId,
				Owner:    container.Owner,
			},
			EventTime: acceptedTime,
		})
	})
	//auth
	code, auth := GenerateAuth(appkey, userId, roomId, container.RtcChannel)
	return code, auth
}

func RtcHangup(ctx context.Context) errs.IMErrorCode {
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
	eventTime := time.Now().UnixMilli()
	if container.RoomType == pbobjs.RtcRoomType_OneOne {
		var calleeId string = ""
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
			SendInviteEvent(ctx, member.MemberId, &pbobjs.RtcInviteEvent{
				InviteType: pbobjs.InviteType_RtcHangup,
				User:       commonservices.GetTargetDisplayUserInfo(ctx, userId),
				Room: &pbobjs.RtcRoom{
					RoomType: container.RoomType,
					RoomId:   container.RoomId,
					Owner:    container.Owner,
				},
				EventTime: eventTime,
			})
			if member.MemberId != container.Owner.UserId {
				calleeId = member.MemberId
			}
		})
		//send notify msg
		var reason CallFinishReasonType
		var duration int64 = 0
		if container.AcceptedTime > 0 {
			reason = CallFinishReasonType_Complete
			duration = time.Now().UnixMilli() - container.AcceptedTime
		} else {
			if userId == container.Owner.UserId {
				reason = CallFinishReasonType_Cancel
			} else {
				reason = CallFinishReasonType_Decline
			}
		}
		SendFinishNtf(ctx, container.Owner.UserId, calleeId, reason, duration, container.RtcMediaType)

		//destroy room
		storage := storages.NewRtcRoomMemberStorage()
		roomStorage := storages.NewRtcRoomStorage()
		roomStorage.Delete(appkey, roomId)
		storage.DeleteByRoomId(appkey, roomId)
		rtcroomCache.Remove(getRoomKey(appkey, roomId))
	} else if container.RoomType == pbobjs.RtcRoomType_OneMore {
		//quit
		code, _ := container.QuitRoom(userId)
		if code == errs.IMErrorCode_SUCCESS {
			//delete from db
			storage := storages.NewRtcRoomMemberStorage()
			storage.Delete(appkey, roomId, userId)
		}
		MemberSyncState(ctx, userId, &pbobjs.SyncMemberStateReq{
			IsDelete: true,
			Member: &pbobjs.MemberState{
				RoomId:   roomId,
				RoomType: container.RoomType,
				MemberId: userId,
				DeviceId: deviceId,
			},
		})
		//notify other members
		container.ForeachMembers(func(member *models.RtcRoomMember) {
			SendInviteEvent(ctx, member.MemberId, &pbobjs.RtcInviteEvent{
				InviteType: pbobjs.InviteType_RtcHangup,
				User:       commonservices.GetTargetDisplayUserInfo(ctx, userId),
				Room: &pbobjs.RtcRoom{
					RoomType: container.RoomType,
					RoomId:   container.RoomId,
					Owner:    container.Owner,
				},
				EventTime: eventTime,
			})
		})
		//send msg
		if container.ConverId != nil && *container.ConverId != "" && container.ChannelType == pbobjs.ChannelType_Group {
			syncMsg2Conver(ctx, container)
		}
		if container.MemberCount() <= 0 {
			//destroy room
			roomStorage := storages.NewRtcRoomStorage()
			roomStorage.Delete(appkey, roomId)
			storage := storages.NewRtcRoomMemberStorage()
			storage.DeleteByRoomId(appkey, roomId)
			rtcroomCache.Remove(getRoomKey(appkey, roomId))
		}
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

func SendInviteEvent(ctx context.Context, targetId string, event *pbobjs.RtcInviteEvent) {
	msg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "rtc_invite_event", event)
	msg.Qos = 0
	msg.PublishType = int32(commonservices.PublishType_AllSessionExceptSelf)
	bases.UnicastRouteWithNoSender(msg)
}

var RtcMsgType_OneOneFinishNtf string = msgdefines.InnerMsgType_CallFinishNtf

type CallFinishReasonType int

var (
	CallFinishReasonType_Cancel   CallFinishReasonType = 0
	CallFinishReasonType_Decline  CallFinishReasonType = 1
	CallFinishReasonType_NoAnswer CallFinishReasonType = 2
	CallFinishReasonType_Complete CallFinishReasonType = 3
)

type CallFinishNtf struct {
	Reason    int   `json:"reason"`
	Duration  int64 `json:"duration"`
	MediaType int   `json:"media_type"`
}

func SendFinishNtf(ctx context.Context, senderId, targetId string, reason CallFinishReasonType, duration int64, mediaType pbobjs.RtcMediaType) {
	ntf := &CallFinishNtf{
		Reason:    int(reason),
		Duration:  duration,
		MediaType: int(mediaType),
	}
	contentBs, _ := tools.JsonMarshal(ntf)
	flag := msgdefines.SetStoreMsg(0)
	msg := &pbobjs.UpMsg{
		MsgType:    RtcMsgType_OneOneFinishNtf,
		MsgContent: contentBs,
		Flags:      flag,
	}
	ctx = context.WithValue(ctx, bases.CtxKey_Session, tools.GenerateUUIDShort11())
	commonservices.AsyncPrivateMsg(ctx, senderId, targetId, msg)
}
