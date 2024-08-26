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

var defaultModifyMsgType string = "jg:modify"

type ModifyMsgContent struct {
	MsgType    string `json:"msg_type"`
	MsgId      string `json:"msg_id"`
	MsgTime    int64  `json:"msg_time"`
	MsgSeq     int64  `json:"msg_seq"`
	MsgContent []byte `json:"msg_content"`
}

func ModifyMsg(ctx context.Context, modifyReq *pbobjs.ModifyMsgReq) errs.IMErrorCode {
	fromUserId := bases.GetRequesterIdFromCtx(ctx)
	targetId := modifyReq.TargetId
	appkey := bases.GetAppKeyFromCtx(ctx)
	converId := commonservices.GetConversationId(fromUserId, targetId, modifyReq.ChannelType)

	//send modify msg
	flag := commonservices.SetCmdMsg(0)

	contentBs, _ := json.Marshal(ModifyMsgContent{
		MsgType:    modifyReq.MsgType,
		MsgId:      modifyReq.MsgId,
		MsgTime:    modifyReq.MsgTime,
		MsgSeq:     modifyReq.MsgSeqNo,
		MsgContent: modifyReq.MsgContent,
	})

	upMsg := &pbobjs.UpMsg{
		MsgType:    defaultModifyMsgType,
		MsgContent: contentBs,
		Flags:      flag,
	}
	if modifyReq.ChannelType == pbobjs.ChannelType_Private {
		//update msg in history
		storage := storages.NewPrivateHisMsgStorage()
		dbMsg, err := storage.FindById(appkey, converId, modifyReq.MsgId)
		if err == nil {
			newDownMsg := &pbobjs.DownMsg{}
			err = tools.PbUnMarshal(dbMsg.MsgBody, newDownMsg)
			if err == nil {
				newDownMsg.MsgContent = modifyReq.MsgContent
				newDownMsg.Flags = commonservices.SetModifiedMsg(newDownMsg.Flags)
				newDownMsgBs, _ := tools.PbMarshal(newDownMsg)
				storage.UpdateMsgBody(appkey, converId, modifyReq.MsgId, newDownMsg.MsgType, newDownMsgBs)
				//upd latest msg for conversation
				if IsLatestMsg(ctx, converId, modifyReq.ChannelType, modifyReq.MsgId, modifyReq.MsgTime, modifyReq.MsgSeqNo) {
					//get latest msg
					convercallers.UpdLatestMsgBody(ctx, fromUserId, targetId, modifyReq.ChannelType, modifyReq.MsgId, newDownMsg)
					convercallers.UpdLatestMsgBody(ctx, targetId, fromUserId, modifyReq.ChannelType, modifyReq.MsgId, newDownMsg)
				}
			}
		}
		//send cmd msg
		commonservices.AsyncPrivateMsg(ctx, fromUserId, targetId, upMsg)
		return errs.IMErrorCode_SUCCESS
	} else if modifyReq.ChannelType == pbobjs.ChannelType_Group {
		//update history msg
		storage := storages.NewGroupHisMsgStorage()
		dbMsg, err := storage.FindById(appkey, converId, modifyReq.MsgId)
		if err == nil {
			newDownMsg := &pbobjs.DownMsg{}
			err = tools.PbUnMarshal(dbMsg.MsgBody, newDownMsg)
			if err == nil {
				newDownMsg.MsgContent = modifyReq.MsgContent
				newDownMsg.Flags = commonservices.SetModifiedMsg(newDownMsg.Flags)
				newDownMsgBs, _ := tools.PbMarshal(newDownMsg)
				storage.UpdateMsgBody(appkey, converId, modifyReq.MsgId, newDownMsg.MsgType, newDownMsgBs)
				//upd latest msg for conversation
				if IsLatestMsg(ctx, converId, modifyReq.ChannelType, modifyReq.MsgId, modifyReq.MsgTime, modifyReq.MsgSeqNo) {
					bases.AsyncRpcCall(ctx, "upd_grp_conver", targetId, &pbobjs.UpdLatestMsgReq{
						TargetId:    targetId,
						ChannelType: modifyReq.ChannelType,
						LatestMsgId: modifyReq.MsgId,
						Action:      pbobjs.UpdLatestMsgAction_UpdMsg,
						Msg:         newDownMsg,
					})
				}
			}
		}
		//send cmd msg
		commonservices.AsyncGroupMsg(ctx, fromUserId, targetId, upMsg)
		return errs.IMErrorCode_SUCCESS
	}
	return errs.IMErrorCode_MSG_DEFAULT
}
