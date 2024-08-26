package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"math"

	"github.com/gin-gonic/gin"
)

func CreatePrivateStreamMsg(ctx *gin.Context) {
	var sendMsgReq models.SendMsgReq
	if err := ctx.BindJSON(&sendMsgReq); err != nil || sendMsgReq.SenderId == "" || sendMsgReq.MsgType == "" || sendMsgReq.TargetId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	msgFlag := handleFlag(sendMsgReq)
	code, sendAck, err := services.SyncSendMsg(ctx, "p_msg", sendMsgReq.SenderId, sendMsgReq.TargetId, &pbobjs.UpMsg{
		MsgType:     sendMsgReq.MsgType,
		MsgContent:  []byte(sendMsgReq.MsgContent),
		Flags:       commonservices.SetStreamMsg(msgFlag),
		MentionInfo: handleMentionInfo(sendMsgReq.MentionInfo),
		ReferMsg:    handleReferMsg(sendMsgReq.ReferMsg),
	}, false)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, code)
		return
	}
	tools.SuccessHttpResp(ctx, sendAck)
}

func AppendPrivateStreamMsg(ctx *gin.Context) {
	var req models.AppendStreamMsgReq
	if err := ctx.BindJSON(&req); err != nil || req.SenderId == "" || req.TargetId == "" || req.MsgId == "" || len(req.Items) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	streamDown := &pbobjs.StreamDownMsg{
		TargetId:    req.TargetId,
		ChannelType: pbobjs.ChannelType_Private,
		MsgId:       req.MsgId,
		MsgItems:    []*pbobjs.StreamMsgItem{},
	}
	for _, item := range req.Items {
		streamDown.MsgItems = append(streamDown.MsgItems, &pbobjs.StreamMsgItem{
			Event:          pbobjs.StreamEvent_StreamMessage,
			SubSeq:         item.SubSeq,
			PartialContent: []byte(item.PartialContent),
		})
	}
	targetId := commonservices.GetConversationId(req.SenderId, req.TargetId, pbobjs.ChannelType_Private)
	services.AsyncApiCall(ctx, "pri_stream", req.SenderId, targetId, streamDown)
	tools.SuccessHttpResp(ctx, nil)
}

func CompletePrivateStreamMsg(ctx *gin.Context) {
	var req models.CompleteStreamMsgReq
	if err := ctx.BindJSON(&req); err != nil || req.SenderId == "" || req.TargetId == "" || req.MsgId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	streamDown := &pbobjs.StreamDownMsg{
		TargetId:    req.TargetId,
		ChannelType: pbobjs.ChannelType_Private,
		MsgId:       req.MsgId,
		MsgItems: []*pbobjs.StreamMsgItem{
			{
				Event:          pbobjs.StreamEvent_StreamComplete,
				SubSeq:         math.MaxInt16,
				PartialContent: []byte{},
			},
		},
	}
	targetId := commonservices.GetConversationId(req.SenderId, req.TargetId, pbobjs.ChannelType_Private)
	services.AsyncApiCall(ctx, "pri_stream", req.SenderId, targetId, streamDown)
	tools.SuccessHttpResp(ctx, nil)
}
