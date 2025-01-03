package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/historymsg/storages"
)

func UpdStreamMsg(ctx context.Context, req *pbobjs.StreamDownMsg) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	requestId := bases.GetRequesterIdFromCtx(ctx)
	converId := commonservices.GetConversationId(requestId, req.TargetId, req.ChannelType)
	if req.ChannelType == pbobjs.ChannelType_Private {
		storage := storages.NewPrivateHisMsgStorage()
		dbMsg, err := storage.FindById(appkey, converId, req.MsgId)
		if err == nil {
			newDownMsg := &pbobjs.DownMsg{}
			err = tools.PbUnMarshal(dbMsg.MsgBody, newDownMsg)
			if err == nil {
				if len(newDownMsg.MsgItems) <= 0 {
					newDownMsg.MsgItems = []*pbobjs.StreamMsgItem{}
				}
				newDownMsg.MsgItems = append(newDownMsg.MsgItems, req.MsgItems...)
				newDownMsgBs, _ := tools.PbMarshal(newDownMsg)
				storage.UpdateMsgBody(appkey, converId, req.MsgId, newDownMsg.MsgType, newDownMsgBs)
			}
		}
	}
	return errs.IMErrorCode_SUCCESS
}
