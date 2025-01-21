package services

import (
	"context"
	"encoding/json"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	mentionStorages "im-server/services/conversation/storages"
	"im-server/services/historymsg/storages"
)

func MarkRead(ctx context.Context, userId string, req *pbobjs.MarkReadReq) errs.IMErrorCode {
	if len(req.Msgs) <= 0 {
		return errs.IMErrorCode_SUCCESS
	}
	if userId == req.TargetId {
		return errs.IMErrorCode_MSG_PARAM_ILLEGAL
	}
	if !checkMarkReadReq(req) {
		return errs.IMErrorCode_MSG_PARAM_ILLEGAL
	}
	if req.ChannelType == pbobjs.ChannelType_Private {
		return markReadPrivateMsgs(ctx, userId, req)
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		return markReadGroupMsgs(ctx, req)
	}
	return errs.IMErrorCode_SUCCESS
}

func checkMarkReadReq(req *pbobjs.MarkReadReq) bool {
	if len(req.Msgs) > 100 {
		return false
	}
	if len(req.IndexScopes) > 10 {
		return false
	}
	for _, scope := range req.IndexScopes {
		start := scope.StartIndex
		end := scope.EndIndex
		if start > end {
			start, end = end, start
		}
		if end-start > 100 {
			return false
		}
	}
	return true
}
func markReadGroupMsgs(ctx context.Context, req *pbobjs.MarkReadReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	groupId := req.TargetId
	markReadMsg := &ReadNtf{
		Msgs:        []*ReadMsg{},
		IndexScopes: []*IndexScope{},
	}
	//TODO check group msg
	if len(req.Msgs) > 0 {
		var latestReadMsgIndex int64 = 0
		msgIds := []string{}
		for _, msg := range req.Msgs {
			msgIds = append(msgIds, msg.MsgId)
			if msg.MsgReadIndex > latestReadMsgIndex {
				latestReadMsgIndex = msg.MsgReadIndex
			}
			markReadMsg.Msgs = append(markReadMsg.Msgs, &ReadMsg{
				MsgId:    msg.MsgId,
				MsgIndex: msg.MsgReadIndex,
			})
		}

		DispatchGroupMsgMarkRead(ctx, groupId, userId, req.ChannelType, msgIds)

		flag := msgdefines.SetCmdMsg(0)
		//notify for other device
		bs, _ := json.Marshal(markReadMsg)
		commonservices.AsyncGroupMsg(ctx, userId, groupId, &pbobjs.UpMsg{
			MsgType:    ReadNtfType,
			MsgContent: bs,
			Flags:      flag,
		}, &bases.OnlySendboxOption{})

		//update mention msg's read state
		mentionStorage := mentionStorages.NewMentionMsgStorage()
		mentionStorage.MarkRead(appkey, userId, req.TargetId, req.ChannelType, msgIds)
	}
	return errs.IMErrorCode_SUCCESS
}

func markReadPrivateMsgs(ctx context.Context, userId string, req *pbobjs.MarkReadReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	targetId := req.TargetId
	var latestReadMsgIndex int64 = 0
	markReadMsg := &ReadNtf{
		Msgs:        []*ReadMsg{},
		IndexScopes: []*IndexScope{},
	}
	converId := commonservices.GetConversationId(userId, targetId, pbobjs.ChannelType_Private)
	storage := storages.NewPrivateHisMsgStorage()
	if len(req.Msgs) > 0 {
		msgIds := []string{}
		for _, msg := range req.Msgs {
			msgIds = append(msgIds, msg.MsgId)
			if msg.MsgReadIndex > latestReadMsgIndex {
				latestReadMsgIndex = msg.MsgReadIndex
			}
			markReadMsg.Msgs = append(markReadMsg.Msgs, &ReadMsg{
				MsgId:    msg.MsgId,
				MsgIndex: msg.MsgReadIndex,
			})
		}
		if len(msgIds) > 0 {
			storage.MarkReadByMsgIds(appkey, converId, msgIds)
		}
	}
	if len(req.IndexScopes) > 0 {
		//scope
		for _, scope := range req.IndexScopes {
			start := scope.StartIndex
			end := scope.EndIndex
			if start > end {
				start, end = end, start
			}
			if end > latestReadMsgIndex {
				latestReadMsgIndex = end
			}
			storage.MarkReadByScope(appkey, converId, start, end)
			markReadMsg.IndexScopes = append(markReadMsg.IndexScopes, &IndexScope{
				StartIndex: start,
				EndIndex:   end,
			})
		}
	}
	//Notify msg's sender
	bs, _ := json.Marshal(markReadMsg)
	flag := msgdefines.SetCmdMsg(0)
	commonservices.AsyncPrivateMsg(ctx, userId, targetId, &pbobjs.UpMsg{
		MsgType:    ReadNtfType,
		Flags:      flag,
		MsgContent: bs,
	})

	return errs.IMErrorCode_SUCCESS
}

var ReadNtfType string = msgdefines.CmdMsgType_ReadNtf

type ReadNtf struct {
	Msgs        []*ReadMsg    `json:"msgs"`
	IndexScopes []*IndexScope `json:"index_scopes"`
}
type ReadMsg struct {
	MsgId    string `json:"msg_id"`
	MsgIndex int64  `json:"msg_index"`
}

type IndexScope struct {
	StartIndex int64 `json:"start_index"`
	EndIndex   int64 `json:"end_index"`
}
