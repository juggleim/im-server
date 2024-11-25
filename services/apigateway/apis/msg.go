package apis

import (
	"im-server/commons/errs"
	"im-server/commons/gmicro/utils"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"time"

	"github.com/gin-gonic/gin"
)

var defaultFlag int32

func init() {
	defaultFlag = commonservices.SetCountMsg(0)
	defaultFlag = commonservices.SetStoreMsg(defaultFlag)
}

func SendPrivateMsg(ctx *gin.Context) {
	var sendMsgReq models.SendMsgReq
	if err := ctx.BindJSON(&sendMsgReq); err != nil || sendMsgReq.SenderId == "" || sendMsgReq.MsgType == "" || sendMsgReq.MsgContent == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	targetIds := []string{}
	if sendMsgReq.TargetId != "" {
		targetIds = append(targetIds, sendMsgReq.TargetId)
	}
	if sendMsgReq.ReceiverId != "" {
		targetIds = append(targetIds, sendMsgReq.ReceiverId)
	}
	targetIds = append(targetIds, sendMsgReq.TargetIds...)
	targetIds = purgeTargetIds(targetIds)
	ret := []*models.SendMsgRespItem{}
	msgIdMap := map[string]string{}
	for _, targetId := range targetIds {
		msgId := tools.GenerateMsgId(time.Now().UnixMilli(), int32(pbobjs.ChannelType_Private), targetId)
		ret = append(ret, &models.SendMsgRespItem{
			TargetId: targetId,
			MsgId:    msgId,
		})
		msgIdMap[targetId] = msgId
	}
	utils.SafeGo(func() {
		for _, targetId := range targetIds {
			msgId := msgIdMap[targetId]
			services.AsyncSendMsg(ctx, "p_msg", sendMsgReq.SenderId, targetId, &pbobjs.UpMsg{
				MsgType:     sendMsgReq.MsgType,
				MsgContent:  []byte(sendMsgReq.MsgContent),
				Flags:       handleFlag(sendMsgReq),
				MentionInfo: handleMentionInfo(sendMsgReq.MentionInfo),
				ReferMsg:    handleReferMsg(sendMsgReq.ReferMsg),
				PushData:    handlePushData(sendMsgReq.PushData),
			}, isNotifySender(sendMsgReq), msgId)
			time.Sleep(10 * time.Millisecond)
		}
	})
	tools.SuccessHttpResp(ctx, ret)
}

func SendSystemMsg(ctx *gin.Context) {
	var sendMsgReq models.SendMsgReq
	if err := ctx.BindJSON(&sendMsgReq); err != nil || sendMsgReq.SenderId == "" || sendMsgReq.MsgType == "" || sendMsgReq.MsgContent == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	targetIds := []string{}
	if sendMsgReq.TargetId != "" {
		targetIds = append(targetIds, sendMsgReq.TargetId)
	}
	if sendMsgReq.ReceiverId != "" {
		targetIds = append(targetIds, sendMsgReq.ReceiverId)
	}
	targetIds = append(targetIds, sendMsgReq.TargetIds...)
	targetIds = purgeTargetIds(targetIds)
	ret := []*models.SendMsgRespItem{}
	msgIdMap := map[string]string{}
	for _, targetId := range targetIds {
		msgId := tools.GenerateMsgId(time.Now().UnixMilli(), int32(pbobjs.ChannelType_System), targetId)
		ret = append(ret, &models.SendMsgRespItem{
			TargetId: targetId,
			MsgId:    msgId,
		})
		msgIdMap[targetId] = msgId
	}
	utils.SafeGo(func() {
		for _, targetId := range targetIds {
			msgId := msgIdMap[targetId]
			services.AsyncSendMsg(ctx, "s_msg", sendMsgReq.SenderId, targetId, &pbobjs.UpMsg{
				MsgType:    sendMsgReq.MsgType,
				MsgContent: []byte(sendMsgReq.MsgContent),
				Flags:      handleFlag(sendMsgReq),
				PushData:   handlePushData(sendMsgReq.PushData),
			}, false, msgId)
			time.Sleep(10 * time.Millisecond)
		}
	})
	tools.SuccessHttpResp(ctx, ret)
}

func SendGroupMsg(ctx *gin.Context) {
	var sendMsgReq models.SendMsgReq
	if err := ctx.BindJSON(&sendMsgReq); err != nil || sendMsgReq.SenderId == "" || sendMsgReq.MsgType == "" || sendMsgReq.MsgContent == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	targetIds := []string{}
	if sendMsgReq.TargetId != "" {
		targetIds = append(targetIds, sendMsgReq.TargetId)
	}
	if sendMsgReq.ReceiverId != "" {
		targetIds = append(targetIds, sendMsgReq.ReceiverId)
	}
	targetIds = append(targetIds, sendMsgReq.TargetIds...)
	targetIds = purgeTargetIds(targetIds)
	ret := []*models.SendMsgRespItem{}
	msgIdMap := map[string]string{}
	for _, targetId := range targetIds {
		msgId := tools.GenerateMsgId(time.Now().UnixMilli(), int32(pbobjs.ChannelType_Group), targetId)
		ret = append(ret, &models.SendMsgRespItem{
			TargetId: targetId,
			MsgId:    msgId,
		})
		msgIdMap[targetId] = msgId
	}
	utils.SafeGo(func() {
		for _, targetId := range targetIds {
			msgId := msgIdMap[targetId]
			services.AsyncSendMsg(ctx, "g_msg", sendMsgReq.SenderId, targetId, &pbobjs.UpMsg{
				MsgType:     sendMsgReq.MsgType,
				MsgContent:  []byte(sendMsgReq.MsgContent),
				Flags:       handleFlag(sendMsgReq),
				MentionInfo: handleMentionInfo(sendMsgReq.MentionInfo),
				ReferMsg:    handleReferMsg(sendMsgReq.ReferMsg),
				PushData:    handlePushData(sendMsgReq.PushData),
			}, isNotifySender(sendMsgReq), msgId)
			time.Sleep(10 * time.Millisecond)
		}
	})
	tools.SuccessHttpResp(ctx, ret)
}

func SendGroupCastMsg(ctx *gin.Context) {
	var req models.SendGrpCastMsgReq
	if err := ctx.BindJSON(&req); err != nil || req.SenderId == "" || req.MsgType == "" || req.MsgContent == "" || len(req.TargetConvers) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if req.TargetId != "" {
		//add group cast msg
		services.AsyncApiCall(ctx, "gc_msg", req.SenderId, req.TargetId, &pbobjs.UpMsg{
			MsgType:    req.MsgType,
			MsgContent: []byte(req.MsgContent),
			Flags:      commonservices.SetStoreMsg(0),
		})
	}
	//dispatch for target conversations
	utils.SafeGo(func() {
		flag := commonservices.SetStoreMsg(0)
		flag = commonservices.SetCountMsg(flag)
		flag = commonservices.SetNoAffectSenderConver(flag)
		upMsg := &pbobjs.UpMsg{
			MsgType:    req.MsgType,
			MsgContent: []byte(req.MsgContent),
			Flags:      flag,
		}
		for _, conver := range req.TargetConvers {
			if conver.ChannelType == int(pbobjs.ChannelType_Private) {
				services.AsyncSendMsg(ctx, "p_msg", req.SenderId, conver.TargetId, upMsg, false, "")
			} else if conver.ChannelType == int(pbobjs.ChannelType_Group) {
				services.AsyncSendMsg(ctx, "g_msg", req.SenderId, conver.TargetId, upMsg, false, "")
			}
			time.Sleep(50 * time.Millisecond)
		}
	})
	tools.SuccessHttpResp(ctx, nil)
}

func SendBroadCastMsg(ctx *gin.Context) {
	var req models.SendBrdCastMsgReq
	if err := ctx.BindJSON(&req); err != nil || req.SenderId == "" || req.MsgType == "" || req.MsgContent == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	flag := commonservices.SetStoreMsg(0)
	flag = commonservices.SetCountMsg(flag)
	services.AsyncApiCall(ctx, "bc_msg", req.SenderId, req.SenderId, &pbobjs.UpMsg{
		MsgType:    req.MsgType,
		MsgContent: []byte(req.MsgContent),
		Flags:      flag,
	})
	tools.SuccessHttpResp(ctx, nil)
}

func handleFlag(sendMsgReq models.SendMsgReq) int32 {
	var flag int32 = 0
	if sendMsgReq.IsStorage == nil || *sendMsgReq.IsStorage {
		flag = commonservices.SetStoreMsg(flag)
	}
	if sendMsgReq.IsCount == nil || *sendMsgReq.IsCount {
		flag = commonservices.SetCountMsg(flag)
	}
	if sendMsgReq.IsCmd != nil && *sendMsgReq.IsCmd {
		flag = commonservices.SetCmdMsg(flag)
	}
	if sendMsgReq.IsState != nil && *sendMsgReq.IsState {
		flag = commonservices.SetStateMsg(flag)
	}
	return flag
}

func handleMentionInfo(mentionInfo *models.MentionInfo) *pbobjs.MentionInfo {
	retMention := &pbobjs.MentionInfo{}
	if mentionInfo != nil {
		if mentionInfo.MentionType == commonservices.MentionType_All {
			retMention.MentionType = pbobjs.MentionType_All
		} else if mentionInfo.MentionType == commonservices.MentionType_Someone {
			retMention.MentionType = pbobjs.MentionType_Someone
		} else if mentionInfo.MentionType == commonservices.MentionType_AllSomeone {
			retMention.MentionType = pbobjs.MentionType_AllAndSomeone
		}
		if mentionInfo.TargetUserIds != nil && len(mentionInfo.TargetUserIds) > 0 {
			retMention.TargetUsers = []*pbobjs.UserInfo{}
			for _, userId := range mentionInfo.TargetUserIds {
				retMention.TargetUsers = append(retMention.TargetUsers, &pbobjs.UserInfo{
					UserId: userId,
				})
			}
		} else if mentionInfo.TargetUsers != nil && len(mentionInfo.TargetUsers) > 0 {
			retMention.TargetUsers = []*pbobjs.UserInfo{}
			for _, user := range mentionInfo.TargetUsers {
				retMention.TargetUsers = append(retMention.TargetUsers, &pbobjs.UserInfo{
					UserId:       user.UserId,
					Nickname:     user.Nickname,
					UserPortrait: user.UserPortrait,
					ExtFields:    commonservices.Map2KvItems(user.ExtFields),
				})
			}
		}
	}
	return retMention
}

func handleReferMsg(referMsg *models.ReferMsg) *pbobjs.DownMsg {
	if referMsg != nil {
		downMsg := &pbobjs.DownMsg{
			TargetId:    referMsg.TargetId,
			SenderId:    referMsg.SenderId,
			MsgId:       referMsg.MsgId,
			MsgTime:     referMsg.MsgTime,
			ChannelType: pbobjs.ChannelType(referMsg.ChannelType),
			MsgType:     referMsg.MsgType,
			MsgContent:  []byte(referMsg.MsgContent),
		}
		return downMsg
	}
	return nil
}

func handlePushData(data *models.PushData) *pbobjs.PushData {
	if data != nil {
		return &pbobjs.PushData{
			Title:         data.PushTitle,
			PushText:      data.PushText,
			PushExtraData: data.PushExtra,
			PushLevel:     pbobjs.PushLevel(data.PushLevel),
		}
	}
	return nil
}

func isNotifySender(sendMsgReq models.SendMsgReq) bool {
	if sendMsgReq.IsNotifySender != nil {
		return *sendMsgReq.IsNotifySender
	}
	return true
}

func purgeTargetIds(targetIds []string) []string {
	tmp := map[string]bool{}
	ret := []string{}
	for _, id := range targetIds {
		if _, exist := tmp[id]; !exist {
			tmp[id] = true
			ret = append(ret, id)
		}
	}
	return ret
}
