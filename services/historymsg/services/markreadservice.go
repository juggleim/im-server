package services

import (
	"context"
	"encoding/json"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/conversation/convercallers"
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
	converId := commonservices.GetConversationId(userId, groupId, req.ChannelType)
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
			//upd latest msg for conversation
			if IsLatestMsg(ctx, converId, req.ChannelType, msg.MsgId, msg.MsgTime, 0) {
				//find msg
				storage := storages.NewGroupHisMsgStorage()
				latestMsg, err := storage.FindById(appkey, converId, msg.MsgId)
				if err == nil && latestMsg != nil {
					newDownMsg := &pbobjs.DownMsg{}
					err = tools.PbUnMarshal(latestMsg.MsgBody, newDownMsg)
					if err == nil {
						newDownMsg.IsRead = true
						convercallers.UpdLatestMsgBody(ctx, userId, groupId, req.ChannelType, msg.MsgId, newDownMsg)
					}
				}
			}
		}

		DispatchGroupMsgMarkRead(ctx, groupId, userId, req.ChannelType, msgIds)

		flag := commonservices.SetCmdMsg(0)
		//notify for other device
		bs, _ := json.Marshal(markReadMsg)
		commonservices.AsyncGroupMsg(ctx, userId, groupId, &pbobjs.UpMsg{
			MsgType:    ReadNtfType,
			MsgContent: bs,
			Flags:      flag,
		}, &bases.OnlySendboxOption{})

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
			for _, msg := range msgIds {
				if IsLatestMsg(ctx, converId, req.ChannelType, msg, 0, 0) {
					//find msg
					latestMsg, err := storage.FindById(appkey, converId, msg)
					if err == nil && latestMsg != nil {
						newDownMsg := &pbobjs.DownMsg{}
						err = tools.PbUnMarshal(latestMsg.MsgBody, newDownMsg)
						if err == nil {
							newDownMsg.IsRead = true
							convercallers.UpdLatestMsgBody(ctx, userId, targetId, req.ChannelType, msg, newDownMsg)
							convercallers.UpdLatestMsgBody(ctx, targetId, userId, req.ChannelType, msg, newDownMsg)
						}
					}
					break
				}
			}
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
	flag := commonservices.SetCmdMsg(0)
	commonservices.AsyncPrivateMsg(ctx, userId, targetId, &pbobjs.UpMsg{
		MsgType:    ReadNtfType,
		Flags:      flag,
		MsgContent: bs,
	})

	return errs.IMErrorCode_SUCCESS
}

var ReadNtfType string = "jg:readntf"

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
