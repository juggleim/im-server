package apis

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"

	"github.com/gin-gonic/gin"
)

type BotInfo struct {
	BotId       string            `json:"bot_id"`
	Nickname    string            `json:"nickname"`
	Portrait    string            `json:"portrait"`
	Webhook     string            `json:"webhook"`
	ExtFields   map[string]string `json:"ext_fields"`
	UpdatedTime int64             `json:"updated_time"`
}

func AddBot(ctx *gin.Context) {
	var botInfo BotInfo
	if err := ctx.BindJSON(&botInfo); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	settings := []*pbobjs.KvItem{}
	if botInfo.Webhook != "" {
		settings = append(settings, &pbobjs.KvItem{
			Key:   string(commonservices.AttItemKey_Bot_WebHook),
			Value: botInfo.Webhook,
		})
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, ""), "add_bot", botInfo.BotId, &pbobjs.UserInfo{
		UserId:       botInfo.BotId,
		Nickname:     botInfo.Nickname,
		UserPortrait: botInfo.Portrait,
		ExtFields:    commonservices.Map2KvItems(botInfo.ExtFields),
		Settings:     settings,
	})
	tools.SuccessHttpResp(ctx, nil)
}
