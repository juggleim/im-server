package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
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
				mergedContent := ""
				for _, item := range req.MsgItems {
					newDownMsg.MsgItems = append(newDownMsg.MsgItems, item)
					//merged content
					if req.MsgType == msgdefines.InnerMsgType_StreamText {
						streamMsgBody := &msgdefines.StreamMsg{}
						if len(item.PartialContent) > 0 {
							err = tools.JsonUnMarshal(item.PartialContent, streamMsgBody)
							if err == nil {
								mergedContent = mergedContent + streamMsgBody.Content
							}
						}
					}
				}
				if mergedContent != "" {
					newDownMsg.MsgContent, _ = tools.JsonMarshal(&msgdefines.StreamMsg{Content: mergedContent})
				}
				newDownMsgBs, _ := tools.PbMarshal(newDownMsg)
				storage.UpdateMsgBody(appkey, converId, req.MsgId, newDownMsg.MsgType, newDownMsgBs)
			}
		}
	}
	return errs.IMErrorCode_SUCCESS
}
