package apis

import (
	"net/http"
	"strconv"

	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/services/aiengines"
	"github.com/juggleim/jugglechat-server/utils"
)

func QryBots(ctx *HttpContext) {
	offset := ctx.Query("offset")
	count := 20
	var err error
	countStr := ctx.Query("count")
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			count = 20
		}
	}
	code, bots := services.QryAiBots(ctx.ToRpcCtx(), int64(count), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	ctx.ResponseSucc(bots)
}

func BotMsgListener(ctx *HttpContext) {
	req := apimodels.BotMsg{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	prompt := `你是我的个人助理，帮我解答一切问题。`
	if req.Stream {
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		assistantInfo := aiengines.GetAiEngineInfo(ctx.ToRpcCtx(), ctx.AppKey)
		if assistantInfo != nil && assistantInfo.AiEngine != nil {
			idIndex := 1
			assistantInfo.AiEngine.StreamChat(ctx.ToRpcCtx(), req.SenderId, req.BotId, prompt, req.Messages[0].Content, func(answerPart string, isEnd bool) {
				if !isEnd {
					item := &apimodels.BotResponsePartData{
						Id:      utils.Int2String(int64(idIndex)),
						Type:    "message",
						Content: answerPart,
					}
					idIndex++
					ctx.Writer.Write([]byte("data: " + utils.ToJson(item) + "\n"))
					ctx.Writer.(http.Flusher).Flush()
				} else {
					ctx.Writer.Write([]byte("[DONE]"))
					ctx.Writer.(http.Flusher).Flush()
				}
			})
		}
	} else {

	}
}
