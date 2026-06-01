package apis

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type BotInfo struct {
	BotId       string            `json:"bot_id"`
	Nickname    string            `json:"nickname"`
	Portrait    string            `json:"portrait"`
	BotType     int               `json:"bot_type"`
	BotConf     string            `json:"bot_conf"`
	Webhook     string            `json:"webhook"`
	ExtFields   map[string]string `json:"ext_fields"`
	UpdatedTime int64             `json:"updated_time"`
}

func RegisterBot(ctx *gin.Context) {
	var botInfo models.BotInfo
	if err := ctx.BindJSON(&botInfo); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	appkey := services.GetAppkeyFromCtx(ctx)
	appinfo, exist := commonservices.GetAppInfo(appkey)
	if exist && appinfo != nil {
		settings := []*pbobjs.KvItem{}
		settings = append(settings, &pbobjs.KvItem{
			Key:   string(commonservices.AttItemKey_Bot_Type),
			Value: tools.Int642String(int64(commonservices.BotType_Default)),
		})
		if botInfo.BotConf != nil {
			settings = append(settings, &pbobjs.KvItem{
				Key:   string(commonservices.AttItemKey_Bot_BotConf),
				Value: tools.ToJson(botInfo.BotConf),
			})
		}
		if botInfo.BotSettings != nil {
			settings = append(settings, &pbobjs.KvItem{
				Key:   string(commonservices.AttItemKey_Bot_Settings),
				Value: tools.ToJson(botInfo.BotSettings),
			})
		}
		code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "add_bot", botInfo.BotId, &pbobjs.UserInfo{
			UserId:       botInfo.BotId,
			Nickname:     botInfo.Nickname,
			UserPortrait: botInfo.Portrait,
			ExtFields:    commonservices.Map2KvItems(botInfo.ExtFields),
			Settings:     settings,
		}, func() proto.Message {
			return &pbobjs.UserRegResp{}
		})
		if err != nil {
			tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
			return
		}
		if code > 0 {
			tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
			return
		}
		token := tokens.ImToken{
			AppKey:    appkey,
			UserId:    botInfo.BotId,
			DeviceId:  "",
			TokenTime: time.Now().UnixMilli(),
			UserType:  pbobjs.UserType_Bot,
		}
		tokenStr, _ := token.ToTokenString([]byte(appinfo.AppSecureKey))
		tools.SuccessHttpResp(ctx, models.UserRegResp{
			UserId: botInfo.BotId,
			Token:  tokenStr,
		})
	} else {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_APP_NOT_EXISTED)
	}
}

func UpdateBot(ctx *gin.Context) {
	var botInfo models.BotInfo
	if err := ctx.BindJSON(&botInfo); err != nil || botInfo.BotId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	settings := []*pbobjs.KvItem{}
	if botInfo.BotConf != nil {
		settings = append(settings, &pbobjs.KvItem{
			Key:   string(commonservices.AttItemKey_Bot_BotConf),
			Value: tools.ToJson(botInfo.BotConf),
		})
	}
	if botInfo.BotSettings != nil {
		settings = append(settings, &pbobjs.KvItem{
			Key:   string(commonservices.AttItemKey_Bot_Settings),
			Value: tools.ToJson(botInfo.BotSettings),
		})
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, ""), "upd_user_info", botInfo.BotId, &pbobjs.UserInfo{
		UserId:       botInfo.BotId,
		Nickname:     botInfo.Nickname,
		UserPortrait: botInfo.Portrait,
		ExtFields:    commonservices.Map2KvItems(botInfo.ExtFields),
		Settings:     settings,
	})
	tools.SuccessHttpResp(ctx, nil)
}
