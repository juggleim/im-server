package apis

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"time"

	"github.com/gin-gonic/gin"
)

func SendPrivateStreamMsg(ctx *gin.Context) {
	var msgReq models.StreamMsg
	if err := ctx.BindJSON(&msgReq); err != nil || msgReq.FromId == "" || msgReq.TargetId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	msgId := msgReq.MsgId
	if msgId == "" {
		msgId = tools.GenerateMsgId(time.Now().UnixMilli(), int32(pbobjs.ChannelType_Private), msgReq.TargetId)
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, msgReq.FromId), "send_stream_msg", msgId, &pbobjs.StreamMsg{
		StreamMsgId:    msgId,
		PartialContent: []byte(msgReq.PartialContent),
		Seq:            int64(msgReq.Seq),
		IsFinished:     msgReq.IsFinished,

		TargetId: msgReq.TargetId,
	})
	tools.SuccessHttpResp(ctx, &models.SendMsgRespItem{
		MsgId: msgId,
	})
}
