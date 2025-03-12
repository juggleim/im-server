package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/configures"
	"im-server/commons/tools"
	"im-server/services/appbusiness/services/imsdk"
	"im-server/services/commonservices"

	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

func InitUserAssistant(ctx context.Context, userId, nickname, portrait string) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	sdk := imsdk.GetImSdk(appkey)
	if sdk != nil {
		if appinfo, exist := commonservices.GetAppInfo(appkey); exist {
			apikey, err := GenerateApiKey(appkey, appinfo.AppSecureKey)
			if err == nil {
				botId := GetAssistantId(userId)
				sdk.AddBot(juggleimsdk.BotInfo{
					BotId:    botId,
					Nickname: GetAssistantNickname(nickname),
					Portrait: portrait,
					BotType:  int(commonservices.BotType_Custom),
					BotConf:  fmt.Sprintf(`{"url":"http://127.0.0.1:%d/jim/bots/messages/listener","api_key":"%s","bot_id":"%s"}`, configures.Config.ConnectManager.WsPort, apikey, botId),
				})
				sdk.SendPrivateMsg(juggleimsdk.Message{
					SenderId:       botId,
					TargetIds:      []string{userId},
					MsgType:        "jg:text",
					MsgContent:     `{"content":"欢迎注册JuggleIM，我是您的私人助理，任何问题都可以问我！"}`,
					IsNotifySender: tools.BoolPtr(false),
				})
			}
		}
	}
}

func GetAssistantId(userId string) string {
	botId := fmt.Sprintf("ass_%s", userId)
	return botId
}

func GetAssistantNickname(nickname string) string {
	return fmt.Sprintf("%s 的助理", nickname)
}
