package apis

import (
	"im-server/commons/bases"
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
	code, msgId, msgTime, msgSeq := commonservices.SyncPrivateMsgOverUpstream(services.ToRpcCtx(ctx, ""), sendMsgReq.SenderId, sendMsgReq.TargetId, &pbobjs.UpMsg{
		MsgType:     sendMsgReq.MsgType,
		MsgContent:  []byte(sendMsgReq.MsgContent),
		Flags:       commonservices.SetStreamMsg(msgFlag),
		MentionInfo: handleMentionInfo(sendMsgReq.MentionInfo),
		ReferMsg:    handleReferMsg(sendMsgReq.ReferMsg),
	}, &bases.NoNotifySenderOption{})
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, code)
		return
	}
	tools.SuccessHttpResp(ctx, &models.SendMsgResp{
		MsgId:   msgId,
		MsgTime: msgTime,
		MsgSeq:  msgSeq,
	})
}

func AppendPrivateStreamMsg(ctx *gin.Context) {
	var req models.AppendStreamMsgReq
	if err := ctx.BindJSON(&req); err != nil || req.SenderId == "" || req.TargetId == "" || req.MsgId == "" || len(req.Items) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	streamDown := &pbobjs.StreamDownMsg{
		MsgId:    req.MsgId,
		MsgItems: []*pbobjs.StreamMsgItem{},
	}
	for _, item := range req.Items {
		pbEvent := pbobjs.StreamEvent_StreamComplete
		event := item.Event
		if event == "msg" {
			pbEvent = pbobjs.StreamEvent_StreamMessage
		} else {
			pbEvent = pbobjs.StreamEvent_StreamComplete
		}
		streamDown.MsgItems = append(streamDown.MsgItems, &pbobjs.StreamMsgItem{
			Event:          pbEvent,
			SubSeq:         item.SubSeq,
			PartialContent: []byte(item.PartialContent),
		})
		if pbEvent == pbobjs.StreamEvent_StreamComplete {
			break
		}
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, req.SenderId), "pri_stream", req.TargetId, streamDown)
	tools.SuccessHttpResp(ctx, nil)
}

func CompletePrivateStreamMsg(ctx *gin.Context) {
	var req models.CompleteStreamMsgReq
	if err := ctx.BindJSON(&req); err != nil || req.SenderId == "" || req.TargetId == "" || req.MsgId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	streamDown := &pbobjs.StreamDownMsg{
		MsgId: req.MsgId,
		MsgItems: []*pbobjs.StreamMsgItem{
			{
				Event:          pbobjs.StreamEvent_StreamComplete,
				SubSeq:         math.MaxInt16,
				PartialContent: []byte{},
			},
		},
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, req.SenderId), "pri_stream", req.TargetId, streamDown)
	tools.SuccessHttpResp(ctx, nil)
}
