package services

import (
	"context"
	"encoding/json"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/historymsg/storages"
	"im-server/services/historymsg/storages/models"

	"time"

	"google.golang.org/protobuf/proto"
)

var (
	GrpMsgReadInfoLimit int64 = 1000
)

func MarkGrpMsgRead(ctx context.Context, req *pbobjs.MarkGrpMsgReadReq) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	addedTime := time.Now()

	grpHisStorage := storages.NewGroupHisMsgStorage()
	grpDelMsgs := []models.GroupDelHisMsg{}
	allReadedMsgIds := []string{}
	for _, msgId := range req.MsgIds {
		readInfo := GetMsgInfo(appkey, req.GroupId, req.SubChannel, msgId, req.ChannelType)
		if readInfo.SenderId == userId {
			continue
		}
		isSucc, readCount := readInfo.AddReadMember(userId, addedTime.UnixMilli())
		if isSucc {
			//add read info
			readInfoStorage := storages.NewReadInfoStorage()
			readInfoStorage.Create(models.ReadInfo{
				AppKey:      appkey,
				MsgId:       msgId,
				ChannelType: req.ChannelType,
				GroupId:     req.GroupId,
				MemberId:    userId,
				CreatedTime: addedTime.UnixMilli(),
			})
			//update hismsg
			grpHisStorage.UpdateReadCount(appkey, req.GroupId, req.SubChannel, msgId, readCount)
			//update msg destroy time after read
			if readInfo.LifeTimeAfterRead > 0 {
				grpDelMsgs = append(grpDelMsgs, models.GroupDelHisMsg{
					UserId:        userId,
					TargetId:      req.GroupId,
					MsgId:         msgId,
					MsgTime:       readInfo.MsgTime,
					MsgSeq:        readInfo.MsgSeq,
					EffectiveTime: addedTime.UnixMilli() + readInfo.LifeTimeAfterRead,
					AppKey:        appkey,
				})
				if readCount >= readInfo.MemberCount {
					allReadedMsgIds = append(allReadedMsgIds, msgId)
				}
			}
			ntf := &GrpReadNtf{
				Msgs: []*GrpReadMsg{},
			}
			ntf.Msgs = append(ntf.Msgs, &GrpReadMsg{
				MsgId:       msgId,
				ReadCount:   readCount,
				MemberCount: readInfo.MemberCount,
			})
			bs, _ := json.Marshal(ntf)
			commonservices.AsyncGroupMsg(ctx, req.GroupId, req.GroupId, &pbobjs.UpMsg{
				MsgType:    GrpReadNtfType,
				MsgContent: bs,
				Flags:      msgdefines.SetCmdMsg(0),
				ToUserIds:  []string{readInfo.SenderId},
			}, &bases.NoNotifySenderOption{}, &bases.MarkFromApiOption{})
		}
	}
	if len(grpDelMsgs) > 0 {
		grpDelMsgStorage := storages.NewGroupDelHisMsgStorage()
		grpDelMsgStorage.BatchCreate(grpDelMsgs)
	}
	if len(allReadedMsgIds) > 0 {
		grpHisStorage.UpdateDestroyTimeAfterReadByMsgIds(appkey, req.GroupId, req.SubChannel, allReadedMsgIds)
	}
}

func DispatchGroupMsgMarkRead(ctx context.Context, groupId, memberId string, channelType pbobjs.ChannelType, msgIds []string) {
	if len(msgIds) > 0 {
		method := "mark_grp_msg_read"
		groups := bases.GroupTargets(method, msgIds)
		for _, ids := range groups {
			bases.AsyncRpcCall(ctx, method, ids[0], &pbobjs.MarkGrpMsgReadReq{
				GroupId:     groupId,
				ChannelType: channelType,
				MsgIds:      ids,
			})
		}
	}
}

func QryReadInfos(ctx context.Context, req *pbobjs.QryReadInfosReq) (errs.IMErrorCode, *pbobjs.QryReadInfosResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	resp := &pbobjs.QryReadInfosResp{
		Items: []*pbobjs.ReadInfoItem{},
	}
	if req.ChannelType == pbobjs.ChannelType_Private {
		storage := storages.NewPrivateHisMsgStorage()
		msgs, err := storage.FindReadTimeByIds(appkey, converId, req.SubChannel, req.MsgIds)
		if err == nil && len(msgs) > 0 {
			for _, msg := range msgs {
				resp.Items = append(resp.Items, &pbobjs.ReadInfoItem{
					MsgId:    msg.MsgId,
					ReadTime: msg.ReadTime,
				})
			}
		}
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		storage := storages.NewReadInfoStorage()
		for _, msgId := range req.MsgIds {
			readCount := storage.CountReadInfosByMsgId(appkey, req.TargetId, req.SubChannel, req.ChannelType, msgId)
			resp.Items = append(resp.Items, &pbobjs.ReadInfoItem{
				MsgId:      msgId,
				ReadCount:  readCount,
				TotalCount: 0, //TODO get total count
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, resp
}
func QryReadDetail(ctx context.Context, req *pbobjs.QryReadDetailReq) (errs.IMErrorCode, *pbobjs.QryReadDetailResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	resp := &pbobjs.QryReadDetailResp{
		ReadMembers:   []*pbobjs.MemberReadDetailItem{},
		UnreadMembers: []*pbobjs.MemberReadDetailItem{},
	}
	if req.ChannelType == pbobjs.ChannelType_Private {

	} else if req.ChannelType == pbobjs.ChannelType_Group {
		readInfoStorage := storages.NewReadInfoStorage()
		infos, err := readInfoStorage.QryReadInfosByMsgId(appkey, req.TargetId, req.SubChannel, req.ChannelType, req.MsgId, 0, GrpMsgReadInfoLimit)
		readMemberMap := map[string]bool{}
		if err == nil {
			for _, info := range infos {
				resp.ReadCount = resp.ReadCount + 1
				resp.ReadMembers = append(resp.ReadMembers, &pbobjs.MemberReadDetailItem{
					Member: commonservices.GetTargetDisplayUserInfo(ctx, info.MemberId),
					Time:   info.CreatedTime,
				})
				readMemberMap[info.MemberId] = true
			}
		}
		storage := storages.NewGroupHisMsgStorage()
		msg, err := storage.FindById(appkey, req.TargetId, req.SubChannel, req.MsgId)
		if err == nil && msg != nil {
			resp.MemberCount = int32(msg.MemberCount)
			//get all members from grp snapshot
			code, respObj, err := bases.SyncRpcCall(ctx, "qry_group_snapshot", req.TargetId, &pbobjs.QryGrpSnapshotReq{
				GroupId:    req.TargetId,
				NearlyTime: msg.SendTime,
			}, func() proto.Message {
				return &pbobjs.GroupSnapshot{}
			})
			if err == nil && code == errs.IMErrorCode_SUCCESS && respObj != nil {
				snapshot := respObj.(*pbobjs.GroupSnapshot)
				memberIds := snapshot.MemberIds
				memberCount := 0
				for _, memberId := range memberIds {
					//exclude sender
					if memberId == msg.SenderId {
						continue
					}
					memberCount++
					if _, exist := readMemberMap[memberId]; !exist {
						resp.UnreadMembers = append(resp.UnreadMembers, &pbobjs.MemberReadDetailItem{
							Member: commonservices.GetTargetDisplayUserInfo(ctx, memberId),
						})
					}
				}
				if memberCount > 0 {
					resp.MemberCount = int32(memberCount)
				}
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, resp
}

var GrpReadNtfType string = msgdefines.CmdMsgType_GrpReadNtf

type GrpReadNtf struct {
	Msgs []*GrpReadMsg `json:"msgs"`
}
type GrpReadMsg struct {
	MsgId       string `json:"msg_id"`
	ReadCount   int    `json:"read_count"`
	MemberCount int    `json:"member_count"`
}
