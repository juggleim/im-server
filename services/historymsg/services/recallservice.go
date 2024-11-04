package services

import (
	"context"
	"encoding/json"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	converStorages "im-server/services/conversation/storages"
	"im-server/services/historymsg/storages"
)

var RecallCmdType string = "jg:recall"
var RecallInfoType string = "jg:recallinfo"

type RecallMsgContent struct {
	MsgId   string            `json:"msg_id"`
	MsgTime int64             `json:"msg_time"`
	Exts    map[string]string `json:"exts"`
}

func RecallMsg(ctx context.Context, recallMsg *pbobjs.RecallMsgReq) errs.IMErrorCode {
	userId := bases.GetRequesterIdFromCtx(ctx)
	targetId := recallMsg.TargetId
	appkey := bases.GetAppKeyFromCtx(ctx)
	converId := commonservices.GetConversationId(bases.GetRequesterIdFromCtx(ctx), targetId, recallMsg.ChannelType)

	//send recall msg
	var flag int32 = 0
	flag = commonservices.SetCmdMsg(flag)

	contentBs, _ := json.Marshal(RecallMsgContent{
		MsgId:   recallMsg.MsgId,
		MsgTime: recallMsg.MsgTime,
		Exts:    commonservices.Kvitems2Map(recallMsg.Exts),
	})

	upMsg := &pbobjs.UpMsg{
		MsgType:    RecallCmdType,
		MsgContent: contentBs,
		Flags:      flag,
	}
	if recallMsg.ChannelType == pbobjs.ChannelType_Private {
		//replace history msg
		storage := storages.NewPrivateHisMsgStorage()
		flag := commonservices.SetStoreMsg(0)
		flag = commonservices.SetCountMsg(flag)
		replaceMsg := &pbobjs.DownMsg{
			SenderId:    userId,
			TargetId:    targetId,
			ChannelType: recallMsg.ChannelType,
			MsgType:     RecallInfoType,
			MsgContent:  contentBs,
			Flags:       flag,
			MsgId:       recallMsg.MsgId,
			MsgTime:     recallMsg.MsgTime,
		}

		replaceMsgBs, _ := tools.PbMarshal(replaceMsg)
		storage.UpdateMsgBody(appkey, converId, recallMsg.MsgId, RecallInfoType, replaceMsgBs)
		//send cmd msg
		commonservices.AsyncPrivateMsg(ctx, userId, targetId, upMsg)
		return errs.IMErrorCode_SUCCESS
	} else if recallMsg.ChannelType == pbobjs.ChannelType_Group {
		//replace history msg
		storage := storages.NewGroupHisMsgStorage()
		flag := commonservices.SetStoreMsg(0)
		flag = commonservices.SetCountMsg(flag)
		replaceMsg := &pbobjs.DownMsg{
			SenderId:    userId,
			TargetId:    targetId,
			ChannelType: recallMsg.ChannelType,
			MsgType:     RecallInfoType,
			MsgContent:  contentBs,
			Flags:       flag,
			MsgId:       recallMsg.MsgId,
			MsgTime:     recallMsg.MsgTime,
		}
		replaceMsgBs, _ := tools.PbMarshal(replaceMsg)
		storage.UpdateMsgBody(appkey, converId, recallMsg.MsgId, RecallInfoType, replaceMsgBs)

		//send cmd msg
		commonservices.AsyncGroupMsg(ctx, userId, targetId, upMsg)
		//delete mention msg
		mentionStorage := converStorages.NewMentionMsgStorage()
		mentionStorage.DelOnlyByMsgIds(appkey, []string{recallMsg.MsgId})
		return errs.IMErrorCode_SUCCESS
	}

	return errs.IMErrorCode_MSG_DEFAULT
}
