package services

import (
	"context"
	"encoding/json"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/historymsg/storages"
	"im-server/services/historymsg/storages/models"

	"time"

	"google.golang.org/protobuf/proto"
)

var (
	msgCache   *caches.LruCache
	msgLocks   *tools.SegmentatedLocks
	noExistMsg = &GrpMsgReadInfo{
		ReadMembers: make(map[string]int64),
	}
	GrpMsgReadInfoLimit int64 = 1000
)

func init() {
	msgCache = caches.NewLruCacheWithReadTimeout(100000, nil, 10*time.Minute)
	msgLocks = tools.NewSegmentatedLocks(512)
}

type GrpMsgReadInfo struct {
	Appkey      string
	ChannelType pbobjs.ChannelType
	GroupId     string
	MsgId       string
	SenderId    string
	ReadMembers map[string]int64
	MemberCount int
}

func (info *GrpMsgReadInfo) AddReadMembers(members map[string]int64) (bool, int) {
	isAdded := false
	for memberId, addedTime := range members {
		if _, exist := info.ReadMembers[memberId]; !exist {
			info.ReadMembers[memberId] = addedTime
			isAdded = true
		}
	}
	return isAdded, len(info.ReadMembers)
}

func (info *GrpMsgReadInfo) AddReadMember(memberId string, addedTime int64) (bool, int) {
	key := getCacheKey(info.Appkey, info.GroupId, info.MsgId, info.ChannelType)
	lock := msgLocks.GetLocks(key)
	lock.Lock()
	defer lock.Unlock()
	if _, exist := info.ReadMembers[memberId]; exist {
		return false, len(info.ReadMembers)
	} else {
		info.ReadMembers[memberId] = addedTime
		return true, len(info.ReadMembers)
	}
}

func (info *GrpMsgReadInfo) GetReadMemberCount() int {
	key := getCacheKey(info.Appkey, info.GroupId, info.MsgId, info.ChannelType)
	lock := msgLocks.GetLocks(key)
	lock.RLock()
	defer lock.RUnlock()
	return len(info.ReadMembers)
}

func getCacheKey(appkey, groupId, msgId string, channelType pbobjs.ChannelType) string {
	return fmt.Sprintf("%s_%d_%s_%s", appkey, channelType, groupId, msgId)
}
func GetGroupMsgReadInfo(appkey, groupId, msgId string, channelType pbobjs.ChannelType) *GrpMsgReadInfo {
	key := getCacheKey(appkey, groupId, msgId, channelType)
	if info, exist := msgCache.Get(key); exist {
		return info.(*GrpMsgReadInfo)
	} else {
		lock := msgLocks.GetLocks(key)
		lock.Lock()
		defer lock.Unlock()
		if info, exist := msgCache.Get(key); exist {
			return info.(*GrpMsgReadInfo)
		} else {
			hisStorage := storages.NewGroupHisMsgStorage()
			grpMsg, err := hisStorage.FindById(appkey, groupId, msgId)
			if err != nil || grpMsg == nil {
				msgCache.Add(key, noExistMsg)
				return noExistMsg
			}
			readInfo := &GrpMsgReadInfo{
				Appkey:      appkey,
				ChannelType: channelType,
				GroupId:     groupId,
				MsgId:       msgId,
				ReadMembers: make(map[string]int64),
			}
			readInfoStorage := storages.NewReadInfoStorage()
			infos, err := readInfoStorage.QryReadInfosByMsgId(appkey, groupId, pbobjs.ChannelType_Group, msgId, 0, GrpMsgReadInfoLimit)
			if err == nil {
				members := make(map[string]int64)
				for _, info := range infos {
					members[info.MemberId] = info.CreatedTime
				}
				readInfo.AddReadMembers(members)

				readInfo.MemberCount = grpMsg.MemberCount
				readInfo.SenderId = grpMsg.SenderId
			}
			msgCache.Add(key, readInfo)
			return readInfo
		}
	}
}

func MarkGrpMsgRead(ctx context.Context, req *pbobjs.MarkGrpMsgReadReq) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	addedTime := time.Now()
	for _, msgId := range req.MsgIds {
		readInfo := GetGroupMsgReadInfo(appkey, req.GroupId, msgId, req.ChannelType)
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
			grpHisStorage := storages.NewGroupHisMsgStorage()
			grpHisStorage.UpdateReadCount(appkey, req.GroupId, msgId, readCount)
			ntf := &GrpReadNtf{
				Msgs: []*GrpReadMsg{},
			}
			ntf.Msgs = append(ntf.Msgs, &GrpReadMsg{
				MsgId:       msgId,
				ReadCount:   readCount,
				MemberCount: readInfo.MemberCount,
			})
			bs, _ := json.Marshal(ntf)
			commonservices.GroupMsgFromApi(ctx, req.GroupId, req.GroupId, &pbobjs.UpMsg{
				MsgType:    GrpReadNtfType,
				MsgContent: bs,
				Flags:      commonservices.SetCmdMsg(0),
				ToUserIds:  []string{readInfo.SenderId},
			}, true)

		}
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
	resp := &pbobjs.QryReadInfosResp{
		Items: []*pbobjs.ReadInfoItem{},
	}
	if req.ChannelType == pbobjs.ChannelType_Private {

	} else if req.ChannelType == pbobjs.ChannelType_Group {
		storage := storages.NewReadInfoStorage()
		for _, msgId := range req.MsgIds {
			readCount := storage.CountReadInfosByMsgId(appkey, req.TargetId, req.ChannelType, msgId)
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
		infos, err := readInfoStorage.QryReadInfosByMsgId(appkey, req.TargetId, req.ChannelType, req.MsgId, 0, GrpMsgReadInfoLimit)
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
		msg, err := storage.FindById(appkey, req.TargetId, req.MsgId)
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

var GrpReadNtfType string = "jg:grpreadntf"

type GrpReadNtf struct {
	Msgs []*GrpReadMsg `json:"msgs"`
}
type GrpReadMsg struct {
	MsgId       string `json:"msg_id"`
	ReadCount   int    `json:"read_count"`
	MemberCount int    `json:"member_count"`
}
