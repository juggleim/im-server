package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/historymsg/storages"
	"im-server/services/historymsg/storages/models"
	"time"
)

func MergeMsg(ctx context.Context, mergeMsgReq *pbobjs.MergeMsgReq) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	if len(mergeMsgReq.MergedMsgs.Msgs) > 0 {
		parentMsgId := mergeMsgReq.ParentMsgId
		channelType := mergeMsgReq.MergedMsgs.ChannelType
		userId := mergeMsgReq.MergedMsgs.UserId
		targetId := mergeMsgReq.MergedMsgs.TargetId

		msgIds := []string{}
		for _, msg := range mergeMsgReq.MergedMsgs.Msgs {
			msgIds = append(msgIds, msg.MsgId)
		}
		// find msg body from db and save to merged msgs
		mergedMsgStorage := storages.NewMergedMsgStorage()
		mergedMsgs := []models.MergedMsg{}
		converId := commonservices.GetConversationId(userId, targetId, channelType)
		if channelType == pbobjs.ChannelType_Private {
			storage := storages.NewPrivateHisMsgStorage()
			msgs, err := storage.FindByIds(appkey, converId, mergeMsgReq.SubChannel, msgIds, 0)
			if err == nil {
				for _, msg := range msgs {
					mergedMsgs = append(mergedMsgs, models.MergedMsg{
						ParentMsgId: parentMsgId,
						FromId:      userId,
						TargetId:    targetId,
						ChannelType: channelType,
						SubChannel:  mergeMsgReq.SubChannel,
						MsgId:       msg.MsgId,
						MsgTime:     msg.SendTime,
						MsgBody:     msg.MsgBody,
						AppKey:      appkey,
					})
				}
			}
		} else if channelType == pbobjs.ChannelType_Group {
			storage := storages.NewGroupHisMsgStorage()
			msgs, err := storage.FindByIds(appkey, converId, mergeMsgReq.SubChannel, msgIds, 0)
			if err == nil {
				for _, msg := range msgs {
					mergedMsgs = append(mergedMsgs, models.MergedMsg{
						ParentMsgId: parentMsgId,
						FromId:      userId,
						TargetId:    targetId,
						ChannelType: channelType,
						SubChannel:  mergeMsgReq.SubChannel,
						MsgId:       msg.MsgId,
						MsgTime:     msg.SendTime,
						MsgBody:     msg.MsgBody,
						AppKey:      appkey,
					})
				}
			}
		}
		if len(mergedMsgs) > 0 {
			mergedMsgStorage.BatchSaveMergedMsgs(mergedMsgs)
		}
	}
}

func QryMergedMsgs(ctx context.Context, req *pbobjs.QryMergedMsgsReq) (errs.IMErrorCode, *pbobjs.DownMsgSet) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	parentMsgId := bases.GetTargetIdFromCtx(ctx)
	resp := &pbobjs.DownMsgSet{
		Msgs: []*pbobjs.DownMsg{},
	}
	startTime := req.StartTime
	isPositiveOrder := false
	if req.Order == 0 { //0:倒序;1:正序;
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
	} else {
		isPositiveOrder = true
	}
	dbMsgs, err := storages.NewMergedMsgStorage().QryMergedMsgs(appkey, parentMsgId, startTime, req.Count+1, isPositiveOrder)
	if err == nil {
		for _, dbMsg := range dbMsgs {
			downMsg := &pbobjs.DownMsg{}
			err = tools.PbUnMarshal(dbMsg.MsgBody, downMsg)
			if err == nil {
				resp.Msgs = append(resp.Msgs, downMsg)
			}
		}
	}
	if len(resp.Msgs) < int(req.Count)+1 {
		resp.IsFinished = true
	} else {
		resp.Msgs = resp.Msgs[0:req.Count]
	}
	return errs.IMErrorCode_SUCCESS, resp
}
