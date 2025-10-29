package apis

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/utils"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"time"

	"github.com/gin-gonic/gin"
)

var defaultFlag int32

func init() {
	defaultFlag = msgdefines.SetCountMsg(0)
	defaultFlag = msgdefines.SetStoreMsg(defaultFlag)
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
		if sendMsgReq.MsgId != nil && *sendMsgReq.MsgId != "" && len(*sendMsgReq.MsgId) <= 20 {
			msgId = *sendMsgReq.MsgId
		}
		ret = append(ret, &models.SendMsgRespItem{
			TargetId: targetId,
			MsgId:    msgId,
		})
		msgIdMap[targetId] = msgId
	}
	utils.SafeGo(func() {
		for _, targetId := range targetIds {
			msgId := msgIdMap[targetId]
			opts := []bases.BaseActorOption{}
			if !isNotifySender(sendMsgReq) {
				opts = append(opts, &bases.NoNotifySenderOption{})
			}
			opts = append(opts, &bases.WithMsgIdOption{
				MsgId: msgId,
			})
			commonservices.AsyncPrivateMsgOverUpstream(services.ToRpcCtx(ctx, sendMsgReq.SenderId), sendMsgReq.SenderId, targetId, &pbobjs.UpMsg{
				MsgType:           sendMsgReq.MsgType,
				MsgContent:        []byte(sendMsgReq.MsgContent),
				Flags:             handleFlag(sendMsgReq),
				MentionInfo:       handleMentionInfo(sendMsgReq.MentionInfo),
				ReferMsg:          handleReferMsg(sendMsgReq.ReferMsg),
				PushData:          handlePushData(sendMsgReq.PushData),
				LifeTime:          sendMsgReq.LifeTime,
				LifeTimeAfterRead: sendMsgReq.LifeTimeAfterRead,
			}, opts...)
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
		if sendMsgReq.MsgId != nil && *sendMsgReq.MsgId != "" && len(*sendMsgReq.MsgId) <= 20 {
			msgId = *sendMsgReq.MsgId
		}
		ret = append(ret, &models.SendMsgRespItem{
			TargetId: targetId,
			MsgId:    msgId,
		})
		msgIdMap[targetId] = msgId
	}
	utils.SafeGo(func() {
		for _, targetId := range targetIds {
			msgId := msgIdMap[targetId]
			commonservices.AsyncSystemMsgOverUpstream(services.ToRpcCtx(ctx, sendMsgReq.SenderId), sendMsgReq.SenderId, targetId, &pbobjs.UpMsg{
				MsgType:    sendMsgReq.MsgType,
				MsgContent: []byte(sendMsgReq.MsgContent),
				Flags:      handleFlag(sendMsgReq),
				PushData:   handlePushData(sendMsgReq.PushData),
			}, &bases.NoNotifySenderOption{}, &bases.WithMsgIdOption{
				MsgId: msgId,
			})
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
		if sendMsgReq.MsgId != nil && *sendMsgReq.MsgId != "" && len(*sendMsgReq.MsgId) <= 20 {
			msgId = *sendMsgReq.MsgId
		}
		ret = append(ret, &models.SendMsgRespItem{
			TargetId: targetId,
			MsgId:    msgId,
		})
		msgIdMap[targetId] = msgId
	}
	utils.SafeGo(func() {
		for _, targetId := range targetIds {
			msgId := msgIdMap[targetId]
			opts := []bases.BaseActorOption{}
			if !isNotifySender(sendMsgReq) {
				opts = append(opts, &bases.NoNotifySenderOption{})
			}
			opts = append(opts, &bases.WithMsgIdOption{
				MsgId: msgId,
			})
			commonservices.AsyncGroupMsgOverUpstream(services.ToRpcCtx(ctx, sendMsgReq.SenderId), sendMsgReq.SenderId, targetId, &pbobjs.UpMsg{
				MsgType:           sendMsgReq.MsgType,
				MsgContent:        []byte(sendMsgReq.MsgContent),
				Flags:             handleFlag(sendMsgReq),
				ToUserIds:         sendMsgReq.ToUserIds,
				MentionInfo:       handleMentionInfo(sendMsgReq.MentionInfo),
				ReferMsg:          handleReferMsg(sendMsgReq.ReferMsg),
				PushData:          handlePushData(sendMsgReq.PushData),
				LifeTime:          sendMsgReq.LifeTime,
				LifeTimeAfterRead: sendMsgReq.LifeTimeAfterRead,
			}, opts...)
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
		bases.AsyncRpcCall(services.ToRpcCtx(ctx, req.SenderId), "gc_msg", req.TargetId, &pbobjs.UpMsg{
			MsgType:    req.MsgType,
			MsgContent: []byte(req.MsgContent),
			Flags:      msgdefines.SetStoreMsg(0),
		})
	}
	//dispatch for target conversations
	utils.SafeGo(func() {
		flag := msgdefines.SetStoreMsg(0)
		flag = msgdefines.SetCountMsg(flag)
		flag = msgdefines.SetNoAffectSenderConver(flag)
		upMsg := &pbobjs.UpMsg{
			MsgType:    req.MsgType,
			MsgContent: []byte(req.MsgContent),
			Flags:      flag,
		}
		for _, conver := range req.TargetConvers {
			if conver.ChannelType == int(pbobjs.ChannelType_Private) {
				commonservices.AsyncPrivateMsgOverUpstream(services.ToRpcCtx(ctx, req.SenderId), req.SenderId, conver.TargetId, upMsg, &bases.NoNotifySenderOption{})
			} else if conver.ChannelType == int(pbobjs.ChannelType_Group) {
				commonservices.AsyncGroupMsgOverUpstream(services.ToRpcCtx(ctx, req.SenderId), req.SenderId, conver.TargetId, upMsg, &bases.NoNotifySenderOption{})
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
	flag := msgdefines.SetStoreMsg(0)
	flag = msgdefines.SetCountMsg(flag)
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, req.SenderId), "bc_msg", req.SenderId, &pbobjs.UpMsg{
		MsgType:    req.MsgType,
		MsgContent: []byte(req.MsgContent),
		Flags:      flag,
	})
	tools.SuccessHttpResp(ctx, nil)
}

func handleFlag(sendMsgReq models.SendMsgReq) int32 {
	var flag int32 = 0
	if sendMsgReq.IsStorage == nil || *sendMsgReq.IsStorage {
		flag = msgdefines.SetStoreMsg(flag)
	}
	if sendMsgReq.IsCount == nil || *sendMsgReq.IsCount {
		flag = msgdefines.SetCountMsg(flag)
	}
	if sendMsgReq.IsCmd != nil && *sendMsgReq.IsCmd {
		flag = msgdefines.SetCmdMsg(flag)
	}
	if sendMsgReq.IsState != nil && *sendMsgReq.IsState {
		flag = msgdefines.SetStateMsg(flag)
	}
	return flag
}

func handleMentionInfo(mentionInfo *models.MentionInfo) *pbobjs.MentionInfo {
	retMention := &pbobjs.MentionInfo{}
	if mentionInfo != nil {
		if mentionInfo.MentionType == msgdefines.MentionType_All {
			retMention.MentionType = pbobjs.MentionType_All
		} else if mentionInfo.MentionType == msgdefines.MentionType_Someone {
			retMention.MentionType = pbobjs.MentionType_Someone
		} else if mentionInfo.MentionType == msgdefines.MentionType_AllSomeone {
			retMention.MentionType = pbobjs.MentionType_AllAndSomeone
		}
		if len(mentionInfo.TargetUserIds) > 0 {
			retMention.TargetUsers = []*pbobjs.UserInfo{}
			for _, userId := range mentionInfo.TargetUserIds {
				retMention.TargetUsers = append(retMention.TargetUsers, &pbobjs.UserInfo{
					UserId: userId,
				})
			}
		} else if len(mentionInfo.TargetUsers) > 0 {
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
