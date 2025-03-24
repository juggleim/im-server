package services

import (
	"context"
	"encoding/json"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/conversation/convercallers"
	converStorages "im-server/services/conversation/storages"
	"im-server/services/historymsg/storages"
	"im-server/services/historymsg/storages/models"
	"time"

	"google.golang.org/protobuf/proto"
)

func SavePrivateHisMsg(ctx context.Context, converId, senderId, receiverId string, downMsg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	privateHisMsgStorage := storages.NewPrivateHisMsgStorage()
	var index int64 = downMsg.MsgSeqNo
	msgBs, _ := tools.PbMarshal(downMsg)

	privateHisMsgStorage.SavePrivateHisMsg(models.PrivateHisMsg{
		HisMsg: models.HisMsg{
			ConverId:    converId,
			SenderId:    senderId,
			ReceiverId:  receiverId,
			ChannelType: pbobjs.ChannelType_Private,
			MsgType:     downMsg.MsgType,
			MsgId:       downMsg.MsgId,
			SendTime:    downMsg.MsgTime,
			MsgSeqNo:    index,
			MsgBody:     msgBs,
			AppKey:      appkey,
		},
	})
}
func SaveGroupHisMsg(ctx context.Context, converId string, downMsg *pbobjs.DownMsg, groupMemberCount int) {
	grpHisMsgStorage := storages.NewGroupHisMsgStorage()
	msgBs, _ := tools.PbMarshal(downMsg)

	err := grpHisMsgStorage.SaveGroupHisMsg(models.GroupHisMsg{
		HisMsg: models.HisMsg{
			ConverId:    converId,
			SenderId:    bases.GetRequesterIdFromCtx(ctx),
			ReceiverId:  bases.GetTargetIdFromCtx(ctx),
			ChannelType: pbobjs.ChannelType_Group,
			MsgType:     downMsg.MsgType,
			MsgId:       downMsg.MsgId,
			SendTime:    downMsg.MsgTime,
			MsgSeqNo:    downMsg.MsgSeqNo,
			MsgBody:     msgBs,
			AppKey:      bases.GetAppKeyFromCtx(ctx),
		},
		MemberCount: groupMemberCount,
	})
	if err != nil {
		logs.WithContext(ctx).Errorf("msg_id:%s\terr:%v", downMsg.MsgId, err)
	}
}

func SaveSystemHisMsg(ctx context.Context, converId, senderId, receiverId string, downMsg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewSystemHisMsgStorage()
	msgBs, _ := tools.PbMarshal(downMsg)

	storage.SaveSystemHisMsg(models.SystemHisMsg{
		HisMsg: models.HisMsg{
			ConverId:    converId,
			SenderId:    senderId,
			ReceiverId:  receiverId,
			ChannelType: pbobjs.ChannelType_System,
			MsgType:     downMsg.MsgType,
			MsgId:       downMsg.MsgId,
			SendTime:    downMsg.MsgTime,
			MsgSeqNo:    downMsg.MsgSeqNo,
			MsgBody:     msgBs,
			AppKey:      appkey,
		},
	})
}

func SaveGroupCastHisMsg(ctx context.Context, converId, senderId, receiverId string, downMsg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewGrpCastHisMsgStorage()
	msgBs, _ := tools.PbMarshal(downMsg)
	storage.SaveGrpCastHisMsg(models.GrpCastHisMsg{
		ConverId:    converId,
		SenderId:    senderId,
		ReceiverId:  receiverId,
		ChannelType: pbobjs.ChannelType_GroupCast,
		MsgType:     downMsg.MsgType,
		MsgId:       downMsg.MsgId,
		SendTime:    downMsg.MsgTime,
		MsgSeqNo:    downMsg.MsgSeqNo,
		MsgBody:     msgBs,
		AppKey:      appkey,
	})
}

func SaveBroadCastHisMsg(ctx context.Context, converId, senderId string, downMsg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewBrdCastHisMsgStorage()
	msgBs, _ := tools.PbMarshal(downMsg)
	storage.SaveBrdCastHisMsg(models.BrdCastHisMsg{
		ConverId:    converId,
		SenderId:    senderId,
		ChannelType: pbobjs.ChannelType_BroadCast,
		MsgType:     downMsg.MsgType,
		MsgId:       downMsg.MsgId,
		SendTime:    downMsg.MsgTime,
		MsgSeqNo:    downMsg.MsgSeqNo,
		MsgBody:     msgBs,
		AppKey:      appkey,
	})
}

func QryLatestHisMsg(ctx context.Context, appkey, converId string, channelType pbobjs.ChannelType) *LatestMsgItem {
	latestMsgItem := GetLatestMsg(ctx, converId, channelType)
	return latestMsgItem
}

func qryLatestReadMsgTime(ctx context.Context, userId, targetId string, channelType pbobjs.ChannelType) (bool, int64) {
	code, resp, err := bases.SyncRpcCall(ctx, "qry_conver", userId, &pbobjs.QryConverReq{
		TargetId:    targetId,
		ChannelType: channelType,
		IsInner:     true,
	}, func() proto.Message {
		return &pbobjs.Conversation{}
	})
	if err == nil && code == errs.IMErrorCode_SUCCESS && resp != nil {
		conver, ok := resp.(*pbobjs.Conversation)
		if ok && conver != nil {
			if conver.LatestReadIndex >= conver.LatestUnreadIndex {
				return false, 0
			} else {
				return true, conver.LatestReadMsgTime
			}
		}
	}
	return false, 0
}

func QryFirstUnreadMsg(ctx context.Context, req *pbobjs.QryFirstUnreadMsgReq) (errs.IMErrorCode, *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	channelType := req.ChannelType
	targetId := req.TargetId
	converId := commonservices.GetConversationId(userId, targetId, channelType)
	var count int32 = 1
	var startTime int64 = 0
	hasUnread, sTime := qryLatestReadMsgTime(ctx, userId, targetId, channelType)
	if hasUnread {
		startTime = sTime
	} else {
		return errs.IMErrorCode_SUCCESS, nil
	}

	if channelType == pbobjs.ChannelType_Private {
		storage := storages.NewPrivateHisMsgStorage()
		msgs, err := storage.QryHisMsgs(appkey, converId, startTime, count, true, 0, []string{}, []string{})
		if err == nil {
			for _, msg := range msgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(msg.MsgBody, downMsg)
				if err == nil {
					if downMsg.MsgSeqNo <= 0 {
						downMsg.MsgSeqNo = msg.MsgSeqNo
					}
					if downMsg.MsgId == "" {
						downMsg.MsgId = msg.MsgId
					}
					if downMsg.MsgTime <= 0 {
						downMsg.MsgTime = msg.SendTime
					}
					if userId == msg.SenderId {
						downMsg.IsSend = true
						downMsg.TargetId = targetId
					}
					downMsg.IsRead = msg.IsRead > 0
					if msg.IsExt > 0 {
						downMsg.Flags = msgdefines.SetExtMsg(downMsg.Flags)
					}
					return errs.IMErrorCode_SUCCESS, downMsg
				}
			}
		}
	} else if channelType == pbobjs.ChannelType_System {
		storage := storages.NewSystemHisMsgStorage()
		dbMsgs, err := storage.QryHisMsgs(appkey, converId, startTime, count, true, 0, []string{})
		if err == nil {
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					if downMsg.MsgSeqNo <= 0 {
						downMsg.MsgSeqNo = dbMsg.MsgSeqNo
					}
					if downMsg.MsgId == "" {
						downMsg.MsgId = dbMsg.MsgId
					}
					if downMsg.MsgTime <= 0 {
						downMsg.MsgTime = dbMsg.SendTime
					}
					return errs.IMErrorCode_SUCCESS, downMsg
				}
			}
		}
	} else if channelType == pbobjs.ChannelType_Group {
		storage := storages.NewGroupHisMsgStorage()
		dbMsgs, err := storage.QryHisMsgs(appkey, converId, startTime, count, true, 0, []string{}, []string{})
		if err == nil {
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					if downMsg.MsgSeqNo <= 0 {
						downMsg.MsgSeqNo = dbMsg.MsgSeqNo
					}
					if downMsg.MsgTime <= 0 {
						downMsg.MsgTime = dbMsg.SendTime
					}
					if downMsg.MsgId == "" {
						downMsg.MsgId = dbMsg.MsgId
					}
					if userId == dbMsg.SenderId {
						downMsg.IsSend = true
					}
					downMsg.MemberCount = int32(dbMsg.MemberCount)
					downMsg.ReadCount = int32(dbMsg.ReadCount)
					if dbMsg.IsExt > 0 {
						downMsg.Flags = msgdefines.SetExtMsg(downMsg.Flags)
					}
					return errs.IMErrorCode_SUCCESS, downMsg
				}
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, nil
}

func QryHisMsgs(ctx context.Context, appkey, targetId string, channelType pbobjs.ChannelType, startTime int64, count int32, isPositive bool, msgTypes []string) (errs.IMErrorCode, *pbobjs.DownMsgSet) {
	userId := bases.GetRequesterIdFromCtx(ctx)
	resp := &pbobjs.DownMsgSet{
		Msgs: []*pbobjs.DownMsg{},
	}
	if channelType == pbobjs.ChannelType_Private {
		cleanTime := GetCleanTime(appkey, userId, targetId, channelType)
		storage := storages.NewPrivateHisMsgStorage()
		converId := commonservices.GetConversationId(userId, targetId, channelType)
		dbMsgs, err := storage.QryHisMsgsExcludeDel(appkey, converId, userId, targetId, startTime, count, isPositive, cleanTime, msgTypes)
		if err == nil {
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					if downMsg.ClientUid == "" {
						downMsg.ClientUid = tools.GenerateUUIDShort22()
					}
					if downMsg.MsgSeqNo <= 0 {
						downMsg.MsgSeqNo = dbMsg.MsgSeqNo
					}
					if downMsg.MsgId == "" {
						downMsg.MsgId = dbMsg.MsgId
					}
					downMsg.TargetUserInfo = nil
					if downMsg.MsgTime <= 0 {
						downMsg.MsgTime = dbMsg.SendTime
					}
					if userId == dbMsg.SenderId {
						downMsg.IsSend = true
						downMsg.TargetId = targetId
					}
					downMsg.IsRead = dbMsg.IsRead > 0
					//msg ext
					if len(dbMsg.MsgExt) > 0 {
						extItems := &pbobjs.MsgExtItems{
							MsgId: dbMsg.MsgId,
						}
						err = tools.PbUnMarshal(dbMsg.MsgExt, extItems)
						if err == nil {
							downMsg.MsgExts = extItems.Exts
						}
					}
					//msg exset
					if len(dbMsg.MsgExset) > 0 {
						extItems := &pbobjs.MsgExtItems{
							MsgId: dbMsg.MsgId,
						}
						err = tools.PbUnMarshal(dbMsg.MsgExset, extItems)
						if err == nil {
							downMsg.MsgExSet = extItems.Exts
							fillUserInfos(ctx, []*pbobjs.MsgExtItems{
								{
									MsgId: downMsg.MsgId,
									Exts:  downMsg.MsgExSet,
								},
							})
						}
					}
					resp.Msgs = append(resp.Msgs, downMsg)
				}
			}
			//add userinfo
			targetUserInfo := commonservices.GetTargetDisplayUserInfo(ctx, targetId)
			resp.TargetUserInfo = targetUserInfo
		}
	} else if channelType == pbobjs.ChannelType_System {
		storage := storages.NewSystemHisMsgStorage()
		converId := commonservices.GetConversationId(userId, targetId, channelType)
		cleanTime := GetCleanTime(appkey, userId, targetId, channelType)
		dbMsgs, err := storage.QryHisMsgs(appkey, converId, startTime, count, isPositive, cleanTime, msgTypes)
		if err == nil {
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					if downMsg.ClientUid == "" {
						downMsg.ClientUid = tools.GenerateUUIDShort22()
					}
					if downMsg.MsgSeqNo <= 0 {
						downMsg.MsgSeqNo = dbMsg.MsgSeqNo
					}
					if downMsg.MsgId == "" {
						downMsg.MsgId = dbMsg.MsgId
					}
					if downMsg.MsgTime <= 0 {
						downMsg.MsgTime = dbMsg.SendTime
					}
					downMsg.TargetUserInfo = nil
					resp.Msgs = append(resp.Msgs, downMsg)
				}
			}
		}
		//add userinfo
		targetUserInfo := commonservices.GetTargetDisplayUserInfo(ctx, targetId)
		resp.TargetUserInfo = targetUserInfo
	} else if channelType == pbobjs.ChannelType_Group {
		var cleanTime int64 = 0
		if !bases.GetIsFromApiFromCtx(ctx) {
			appInfo, exist := commonservices.GetAppInfo(appkey)
			if exist && appInfo != nil {
				if !appInfo.NotCheckGrpMember {
					memberSettings := qryGrpMemberSettings(ctx, targetId, userId)
					if memberSettings == nil || !memberSettings.IsMember { // not group member
						return errs.IMErrorCode_GROUP_NOTGROUPMEMBER, nil
					}
					if memberSettings.GrpMemberSetting.HideGrpMsg {
						cleanTime = memberSettings.JoinTime
					}
				}
			}
		}
		dbCleanTime := GetCleanTime(appkey, userId, targetId, channelType)
		if dbCleanTime > cleanTime {
			cleanTime = dbCleanTime
		}
		storage := storages.NewGroupHisMsgStorage()
		converId := commonservices.GetConversationId(userId, targetId, channelType)
		dbMsgs, err := storage.QryHisMsgsExcludeDel(appkey, converId, userId, targetId, startTime, count, isPositive, cleanTime, msgTypes)
		if err == nil {
			msgMap := map[string]*pbobjs.DownMsg{}
			msgIds := []string{}
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					msgMap[dbMsg.MsgId] = downMsg
					msgIds = append(msgIds, dbMsg.MsgId)

					if downMsg.ClientUid == "" {
						downMsg.ClientUid = tools.GenerateUUIDShort22()
					}
					if downMsg.MsgSeqNo <= 0 {
						downMsg.MsgSeqNo = dbMsg.MsgSeqNo
					}
					if downMsg.MsgTime <= 0 {
						downMsg.MsgTime = dbMsg.SendTime
					}
					if downMsg.MsgId == "" {
						downMsg.MsgId = dbMsg.MsgId
					}
					downMsg.GroupInfo = nil
					if userId == dbMsg.SenderId {
						downMsg.IsSend = true
					}
					downMsg.MemberCount = int32(dbMsg.MemberCount)
					downMsg.ReadCount = int32(dbMsg.ReadCount)
					//msg ext
					if len(dbMsg.MsgExt) > 0 {
						extItems := &pbobjs.MsgExtItems{
							MsgId: dbMsg.MsgId,
						}
						err = tools.PbUnMarshal(dbMsg.MsgExt, extItems)
						if err == nil {
							downMsg.MsgExts = extItems.Exts
						}
					}
					//msg exset
					if len(dbMsg.MsgExset) > 0 {
						extItems := &pbobjs.MsgExtItems{
							MsgId: dbMsg.MsgId,
						}
						err = tools.PbUnMarshal(dbMsg.MsgExset, extItems)
						if err == nil {
							downMsg.MsgExSet = extItems.Exts
							fillUserInfos(ctx, []*pbobjs.MsgExtItems{
								{
									MsgId: downMsg.MsgId,
									Exts:  downMsg.MsgExSet,
								},
							})
						}
					}
					resp.Msgs = append(resp.Msgs, downMsg)
				}
			}
			//readed status of group msg
			readStorage := storages.NewReadInfoStorage()
			readMap, err := readStorage.CheckMsgsRead(appkey, targetId, userId, pbobjs.ChannelType_Group, msgIds)
			if err == nil {
				for msgId, readStatus := range readMap {
					if msg, exist := msgMap[msgId]; exist {
						msg.IsRead = readStatus
					}
				}
			}
			//add groupinfo
			groupInfo := commonservices.GetGroupInfoFromCache(ctx, targetId)
			resp.GroupInfo = groupInfo
		}
	} else if channelType == pbobjs.ChannelType_GroupCast {
		storage := storages.NewGrpCastHisMsgStorage()
		converId := commonservices.GetConversationId(userId, targetId, channelType)
		dbMsgs, err := storage.QryHisMsgs(appkey, converId, startTime, count, isPositive, 0, msgTypes)
		if err == nil {
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					downMsg.GroupInfo = nil
					if userId == dbMsg.SenderId {
						downMsg.IsSend = true
					}
					if downMsg.ClientUid == "" {
						downMsg.ClientUid = tools.GenerateUUIDShort22()
					}
					resp.Msgs = append(resp.Msgs, downMsg)
				}
			}
		}
	} else if channelType == pbobjs.ChannelType_BroadCast {
		storage := storages.NewBrdCastHisMsgStorage()
		converId := commonservices.GetConversationId(userId, targetId, channelType)
		dbMsgs, err := storage.QryHisMsgs(appkey, converId, startTime, count, isPositive, 0, msgTypes)
		if err == nil {
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					if downMsg.ClientUid == "" {
						downMsg.ClientUid = tools.GenerateUUIDShort22()
					}
					resp.Msgs = append(resp.Msgs, downMsg)
				}
			}
		}
	}
	if len(resp.Msgs) < int(count) {
		resp.IsFinished = true
	}
	//statistic
	if len(resp.Msgs) > 0 {
		for _, msg := range resp.Msgs {
			commonservices.ReportDownMsg(appkey, msg.ChannelType, 1)
		}
	}
	return errs.IMErrorCode_SUCCESS, resp
}

type GrpMemberSettings struct {
	GroupId          string
	MemberId         string
	IsMember         bool
	JoinTime         int64
	GrpMemberSetting *commonservices.GrpMemberSettings
}

func qryGrpMemberSettings(ctx context.Context, groupId, memberId string) *GrpMemberSettings {
	code, resp, err := bases.SyncRpcCall(ctx, "qry_grp_member_settings", groupId, &pbobjs.QryGrpMemberSettingsReq{
		MemberId: memberId,
	}, func() proto.Message {
		return &pbobjs.QryGrpMemberSettingsResp{}
	})
	if err != nil || code != 0 {
		return nil
	}
	memberSettings, ok := resp.(*pbobjs.QryGrpMemberSettingsResp)
	if !ok || memberSettings == nil {
		return nil
	}
	retMemberSettings := &GrpMemberSettings{
		GroupId:          memberSettings.GroupId,
		MemberId:         memberSettings.MemberId,
		IsMember:         memberSettings.IsMember,
		JoinTime:         memberSettings.JoinTime,
		GrpMemberSetting: &commonservices.GrpMemberSettings{},
	}
	commonservices.FillObjField(retMemberSettings.GrpMemberSetting, memberSettings.MemberSettings)
	return retMemberSettings
}

func QryHisMsgByIds(ctx context.Context, req *pbobjs.QryHisMsgByIdsReq) *pbobjs.DownMsgSet {
	userId := bases.GetRequesterIdFromCtx(ctx)
	appkey := bases.GetAppKeyFromCtx(ctx)
	targetId := req.TargetId
	converId := commonservices.GetConversationId(userId, targetId, req.ChannelType)
	cleanTime := GetCleanTime(appkey, userId, targetId, req.ChannelType)
	resp := &pbobjs.DownMsgSet{
		Msgs:       []*pbobjs.DownMsg{},
		IsFinished: true,
	}
	if req.ChannelType == pbobjs.ChannelType_Private {
		storage := storages.NewPrivateHisMsgStorage()
		dbMsgs, err := storage.FindByIds(appkey, converId, req.MsgIds, cleanTime)
		if err == nil {
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					if userId == dbMsg.SenderId {
						downMsg.IsSend = true
						downMsg.TargetId = targetId
					}
					downMsg.IsRead = dbMsg.IsRead > 0
					resp.Msgs = append(resp.Msgs, downMsg)
				}
			}
		}
	} else if req.ChannelType == pbobjs.ChannelType_System {
		storage := storages.NewSystemHisMsgStorage()
		dbMsgs, err := storage.FindByIds(appkey, converId, req.MsgIds, cleanTime)
		if err == nil {
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					resp.Msgs = append(resp.Msgs, downMsg)
				}
			}
		}
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		storage := storages.NewGroupHisMsgStorage()
		dbMsgs, err := storage.FindByIds(appkey, converId, req.MsgIds, cleanTime)
		if err == nil {
			for _, dbMsg := range dbMsgs {
				downMsg := &pbobjs.DownMsg{}
				err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
				if err == nil {
					if userId == dbMsg.SenderId {
						downMsg.IsSend = true
					}
					resp.Msgs = append(resp.Msgs, downMsg)
				}
			}
		}
	}
	//statistic
	if len(resp.Msgs) > 0 {
		for _, msg := range resp.Msgs {
			commonservices.ReportDownMsg(appkey, msg.ChannelType, 1)
		}
	}
	return resp
}

func GetCleanTime(appkey, userId, targetId string, channelType pbobjs.ChannelType) int64 {
	converId := commonservices.GetConversationId(userId, targetId, channelType)
	var destroyTime, cleanTime int64 = 0, 0
	//conver clean time
	converCleanStorage := storages.NewHisMsgConverCleanTimeStorage()
	destroy, err := converCleanStorage.FindOne(appkey, converId, channelType)
	if err == nil && destroy != nil && destroy.CleanTime > 0 {
		destroyTime = destroy.CleanTime
	}
	//user clean time
	userCleanStorage := storages.NewHisMsgUserCleanTimeStorage()
	clean, err := userCleanStorage.FindOne(appkey, userId, targetId, channelType)
	if err == nil && clean != nil && clean.CleanTime > 0 {
		cleanTime = clean.CleanTime
	}
	if destroyTime > cleanTime {
		return destroyTime
	} else {
		return cleanTime
	}
}

func CleanHisMsg(ctx context.Context, req *pbobjs.CleanHisMsgReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	targetId := req.TargetId
	channelType := req.ChannelType
	cleanTime := req.CleanMsgTime

	if cleanTime < 0 {
		return errs.IMErrorCode_SUCCESS
	}
	if cleanTime == 0 || cleanTime > time.Now().UnixMilli() {
		cleanTime = time.Now().UnixMilli()
	}
	flag := msgdefines.SetCmdMsg(0)
	bs, _ := json.Marshal(CleanMsg{
		TargetId:    targetId,
		ChannelType: int32(channelType),
		CleanTime:   cleanTime,
		SenderId:    req.SenderId,
	})
	cmdMsg := &pbobjs.UpMsg{
		MsgType:    cleanMsgType,
		MsgContent: bs,
		Flags:      flag,
	}
	latestMsg := GetLatestMsg(ctx, commonservices.GetConversationId(userId, req.TargetId, req.ChannelType), channelType)
	if req.CleanScope == 0 { //user clean time
		storage := storages.NewHisMsgUserCleanTimeStorage()
		dbClean, err := storage.FindOne(appkey, userId, targetId, channelType)
		if err == nil && dbClean.CleanTime >= cleanTime {
			return errs.IMErrorCode_SUCCESS
		}
		err = storage.UpsertCleanTime(models.HisMsgUserCleanTime{
			AppKey:      appkey,
			UserId:      userId,
			TargetId:    targetId,
			ChannelType: req.ChannelType,
			CleanTime:   cleanTime,
		})
		if err != nil {
			logs.WithContext(ctx).Error(err.Error())
		}

		//notify other device to clean msgs
		if channelType == pbobjs.ChannelType_Private {
			commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, cmdMsg, &bases.OnlySendboxOption{})
		} else if channelType == pbobjs.ChannelType_Group {
			commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, cmdMsg, &bases.OnlySendboxOption{})
		}
		//del latest msg for conversation
		if cleanTime >= latestMsg.LatestMsgTime {
			convercallers.UpdLatestMsgBody(ctx, userId, req.TargetId, req.ChannelType, latestMsg.LatestMsgId, &pbobjs.DownMsg{})
		}
	} else if req.CleanScope == 1 { //conver clean time
		converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
		if req.SenderId == "" {
			storage := storages.NewHisMsgConverCleanTimeStorage()
			dbDestroy, err := storage.FindOne(appkey, converId, req.ChannelType)
			if err == nil && dbDestroy != nil && dbDestroy.CleanTime >= cleanTime {
				return errs.IMErrorCode_SUCCESS
			}
			err = storage.UpsertDestroyTime(models.HisMsgConverCleanTime{
				AppKey:      appkey,
				ConverId:    converId,
				ChannelType: req.ChannelType,
				CleanTime:   cleanTime,
			})
			if err != nil {
				logs.WithContext(ctx).Error(err.Error())
			}
		} else {
			if channelType == pbobjs.ChannelType_Private {
				storage := storages.NewPrivateHisMsgStorage()
				storage.DelSomeoneMsgsBaseTime(appkey, converId, cleanTime, req.SenderId)
			} else if channelType == pbobjs.ChannelType_Group {
				storage := storages.NewGroupHisMsgStorage()
				storage.DelSomeoneMsgsBaseTime(appkey, converId, cleanTime, req.SenderId)
			}
		}
		//notify other device to destroy msgs
		if channelType == pbobjs.ChannelType_Private {
			commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, cmdMsg)
			if cleanTime >= latestMsg.LatestMsgTime {
				convercallers.UpdLatestMsgBody(ctx, userId, req.TargetId, req.ChannelType, latestMsg.LatestMsgId, &pbobjs.DownMsg{})
				convercallers.UpdLatestMsgBody(ctx, req.TargetId, userId, req.ChannelType, latestMsg.LatestMsgId, &pbobjs.DownMsg{})
			}
		} else if channelType == pbobjs.ChannelType_Group {
			commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, cmdMsg)
			if cleanTime >= latestMsg.LatestMsgTime {
				bases.AsyncRpcCall(ctx, "upd_grp_conver", req.TargetId, &pbobjs.UpdLatestMsgReq{
					TargetId:    req.TargetId,
					ChannelType: req.ChannelType,
					LatestMsgId: latestMsg.LatestMsgId,
					Action:      pbobjs.UpdLatestMsgAction_UpdMsg,
					Msg:         &pbobjs.DownMsg{},
				})
			}
		}
	}

	return errs.IMErrorCode_SUCCESS
}

var cleanMsgType string = msgdefines.CmdMsgType_CleanMsg

type CleanMsg struct {
	TargetId    string `json:"target_id"`
	ChannelType int32  `json:"channel_type"`
	CleanTime   int64  `json:"clean_time"`
	SenderId    string `json:"sender_id,omitempty"`
}

var delMsgsType string = msgdefines.CmdMsgType_DelMsgs

type DelMsgs struct {
	TargetId    string    `json:"target_id"`
	ChannelType int32     `json:"channel_type"`
	Msgs        []*DelMsg `json:"msgs"`
}
type DelMsg struct {
	MsgId string `json:"msg_id"`
}

func DelHisMsg(ctx context.Context, req *pbobjs.DelHisMsgsReq) errs.IMErrorCode {
	userId := bases.GetRequesterIdFromCtx(ctx)
	appkey := bases.GetAppKeyFromCtx(ctx)
	delMsgIds := []string{}
	delMsgs := &DelMsgs{
		TargetId:    req.TargetId,
		ChannelType: int32(req.ChannelType),
		Msgs:        []*DelMsg{},
	}
	if len(req.Msgs) <= 0 {
		return errs.IMErrorCode_SUCCESS
	}
	if req.DelScope == 0 { //one-way
		if req.ChannelType == pbobjs.ChannelType_Private {
			pDelStorage := storages.NewPrivateDelHisMsgStorage()
			items := []models.PrivateDelHisMsg{}
			for _, msg := range req.Msgs {
				items = append(items, models.PrivateDelHisMsg{
					UserId:   userId,
					TargetId: req.TargetId,
					MsgId:    msg.MsgId,
					MsgTime:  msg.MsgTime,
					AppKey:   appkey,
				})
				delMsgs.Msgs = append(delMsgs.Msgs, &DelMsg{
					MsgId: msg.MsgId,
				})
				delMsgIds = append(delMsgIds, msg.MsgId)
			}
			if len(items) > 0 {
				pDelStorage.BatchCreate(items)
			}
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			gDelStorage := storages.NewGroupDelHisMsgStorage()
			items := []models.GroupDelHisMsg{}
			for _, msg := range req.Msgs {
				items = append(items, models.GroupDelHisMsg{
					UserId:   userId,
					TargetId: req.TargetId,
					MsgId:    msg.MsgId,
					MsgTime:  msg.MsgTime,
					AppKey:   appkey,
				})
				delMsgs.Msgs = append(delMsgs.Msgs, &DelMsg{
					MsgId: msg.MsgId,
				})
			}
			if len(items) > 0 {
				gDelStorage.BatchCreate(items)
			}
		}
		//notify other device to clean msgs
		if len(delMsgs.Msgs) > 0 {
			flag := msgdefines.SetCmdMsg(0)
			bs, _ := json.Marshal(delMsgs)
			if req.ChannelType == pbobjs.ChannelType_Private {
				commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, &pbobjs.UpMsg{
					MsgType:    delMsgsType,
					MsgContent: bs,
					Flags:      flag,
				}, &bases.OnlySendboxOption{})
			} else if req.ChannelType == pbobjs.ChannelType_Group {
				commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, &pbobjs.UpMsg{
					MsgType:    delMsgsType,
					MsgContent: bs,
					Flags:      flag,
				}, &bases.OnlySendboxOption{})
			}
			//if latest msg then update conversation
			// for _, msg := range delMsgs.Msgs {
			// 	if IsLatestMsg(ctx, converId, req.ChannelType, msg.MsgId, 0, 0) {
			// 		convercallers.UpdLatestMsgBody(ctx, userId, req.TargetId, req.ChannelType, msg.MsgId, &pbobjs.DownMsg{})
			// 		break
			// 	}
			// }
			//delete mention msg
			mentionStorage := converStorages.NewMentionMsgStorage()
			mentionStorage.DelMentionMsgs(appkey, userId, req.TargetId, req.ChannelType, delMsgIds)
		}
	} else if req.DelScope == 1 { //two-way
		if req.ChannelType == pbobjs.ChannelType_Private {
			pStorage := storages.NewPrivateHisMsgStorage()
			for _, msg := range req.Msgs {
				delMsgIds = append(delMsgIds, msg.MsgId)
				delMsgs.Msgs = append(delMsgs.Msgs, &DelMsg{
					MsgId: msg.MsgId,
				})
			}
			converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
			pStorage.DelMsgs(appkey, converId, delMsgIds)
		} else if req.ChannelType == pbobjs.ChannelType_Group {
			gStorage := storages.NewGroupHisMsgStorage()
			for _, msg := range req.Msgs {
				delMsgIds = append(delMsgIds, msg.MsgId)
				delMsgs.Msgs = append(delMsgs.Msgs, &DelMsg{
					MsgId: msg.MsgId,
				})
			}
			converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
			gStorage.DelMsgs(appkey, converId, delMsgIds)
		}
		//notify all people of conversation
		if len(delMsgs.Msgs) > 0 {
			flag := msgdefines.SetCmdMsg(0)
			bs, _ := json.Marshal(delMsgs)
			if req.ChannelType == pbobjs.ChannelType_Private {
				commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, &pbobjs.UpMsg{
					MsgType:    delMsgsType,
					MsgContent: bs,
					Flags:      flag,
				})
			} else if req.ChannelType == pbobjs.ChannelType_Group {
				commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, &pbobjs.UpMsg{
					MsgType:    delMsgsType,
					MsgContent: bs,
					Flags:      flag,
				})
			}
			// for _, msg := range delMsgs.Msgs {
			// 	if IsLatestMsg(ctx, converId, req.ChannelType, msg.MsgId, 0, 0) {
			// 		if req.ChannelType == pbobjs.ChannelType_Private {
			// 			convercallers.UpdLatestMsgBody(ctx, userId, req.TargetId, req.ChannelType, msg.MsgId, &pbobjs.DownMsg{})
			// 			convercallers.UpdLatestMsgBody(ctx, req.TargetId, userId, req.ChannelType, msg.MsgId, &pbobjs.DownMsg{})
			// 		} else if req.ChannelType == pbobjs.ChannelType_Group {
			// 			bases.AsyncRpcCall(ctx, "upd_grp_conver", req.TargetId, &pbobjs.UpdLatestMsgReq{
			// 				TargetId:    req.TargetId,
			// 				ChannelType: req.ChannelType,
			// 				LatestMsgId: msg.MsgId,
			// 				Action:      pbobjs.UpdLatestMsgAction_UpdMsg,
			// 				Msg:         &pbobjs.DownMsg{},
			// 			})
			// 		}
			// 		break
			// 	}
			// }
			//delete mention msg
			mentionStorage := converStorages.NewMentionMsgStorage()
			mentionStorage.DelOnlyByMsgIds(appkey, delMsgIds)
		}
	}
	return errs.IMErrorCode_SUCCESS
}
