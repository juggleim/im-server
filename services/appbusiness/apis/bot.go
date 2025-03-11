package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/apimodels"
	"im-server/services/appbusiness/httputils"
	"im-server/services/appbusiness/services"
	"net/http"
	"strconv"
)

func QryBots(ctx *httputils.HttpContext) {
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
	code, bots := services.QryAiBots(ctx.ToRpcCtx(), &pbobjs.QryAiBotsReq{
		Limit:  int64(count),
		Offset: offset,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ctx.ResponseErr(code)
		return
	}
	ctx.ResponseSucc(bots)
}

func BotMsgListener(ctx *httputils.HttpContext) {
	req := apimodels.BotMsg{}
	if err := ctx.BindJson(&req); err != nil {
		ctx.ResponseErr(errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	prompt := `你是我的个人助理，帮我解答一切问题。`
	if req.Stream {
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		assistantInfo := services.GetAssistantInfo(ctx.ToRpcCtx())
		if assistantInfo != nil && assistantInfo.AiEngine != nil {
			idIndex := 1
			assistantInfo.AiEngine.StreamChat(ctx.ToRpcCtx(), req.SenderId, "assistant", prompt, req.Messages[0].Content, func(answerPart string, isEnd bool) {
				if !isEnd {
					item := &apimodels.BotResponsePartData{
						Id:      tools.Int642String(int64(idIndex)),
						Type:    "message",
						Content: answerPart,
					}
					idIndex++
					ctx.Writer.Write([]byte("data: " + tools.ToJson(item) + "\n"))
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
