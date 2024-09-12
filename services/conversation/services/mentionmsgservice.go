package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages"
	hisService "im-server/services/historymsg/services"
	"time"
)

func QryMentionedMsgs(ctx context.Context, userId string, req *pbobjs.QryMentionMsgsReq) *pbobjs.QryMentionMsgsResp {
	ret := &pbobjs.QryMentionMsgsResp{
		MentionMsgs: []*pbobjs.DownMsg{},
		IsFinished:  true,
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	startTime := req.StartTime
	isPositiveOrder := false
	if req.Order == 0 { //0:倒序;1:正序
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
	} else {
		isPositiveOrder = true
	}

	mentionMsgStorage := storages.NewMentionMsgStorage()
	// dbMentionMsgs, err := mentionMsgStorage.QryMentionMsgs(appkey, userId, req.TargetId, req.ChannelType, startTime, int(req.Count), isPositiveOrder, req.LatestReadIndex)
	dbMentionMsgs, err := mentionMsgStorage.QryUnreadMentionMsgs(appkey, userId, req.TargetId, req.ChannelType, startTime, int(req.Count), isPositiveOrder)
	if err == nil {
		msgIds := []string{}
		for _, dbMentionMsg := range dbMentionMsgs {
			msgIds = append(msgIds, dbMentionMsg.MsgId)
		}
		downMsgs := QryHisMsgByIds(ctx, userId, req.TargetId, req.ChannelType, msgIds)
		for _, downMsg := range downMsgs {
			if downMsg.MsgType != hisService.RecallInfoType {
				ret.MentionMsgs = append(ret.MentionMsgs, downMsg)
			}
		}
		if len(ret.MentionMsgs) >= int(req.Count) {
			ret.IsFinished = false
		}
	}

	return ret
}
