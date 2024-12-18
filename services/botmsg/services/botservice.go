package services

import (
	"bytes"
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/botmsg/services/botengines"
	"im-server/services/botmsg/storages"
	"im-server/services/botmsg/storages/models"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"strings"
	"time"
)

var botCache *caches.LruCache
var botLocks *tools.SegmentatedLocks

func init() {
	botCache = caches.NewLruCacheWithAddReadTimeout(10000, nil, 5*time.Second, 5*time.Second)
	botLocks = tools.NewSegmentatedLocks(128)
}

type BotInfo struct {
	AppKey    string
	BotId     string
	Nickname  string
	Portrait  string
	BotType   models.BotType
	BotEngine botengines.IBotEngine
}

func GetBotInfo(ctx context.Context, botId string) *BotInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getKey(appkey, botId)
	if val, exist := botCache.Get(key); exist {
		return val.(*BotInfo)
	} else {
		l := botLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if val, exist := botCache.Get(key); exist {
			return val.(*BotInfo)
		} else {
			botInfo := &BotInfo{
				AppKey:    appkey,
				BotId:     botId,
				BotEngine: &botengines.NilBotEngine{},
			}
			storage := storages.NewBotConfStorage()
			bot, err := storage.FindById(appkey, botId)
			if err == nil {
				botInfo.Nickname = bot.Nickname
				botInfo.Portrait = bot.BotPortrait
				botInfo.BotType = bot.BotType
				switch botInfo.BotType {
				case models.BotType_Dify:
					difyBot := &botengines.DifyBotEngine{}
					err = tools.JsonUnMarshal([]byte(bot.BotConf), difyBot)
					if err == nil && difyBot.ApiKey != "" && difyBot.Url != "" {
						botInfo.BotEngine = difyBot
					}
				}
			} else {
				botInfo.BotEngine = &botengines.NilBotEngine{}
			}
			botCache.Add(key, botInfo)
			return botInfo
		}
	}
}

func getKey(appkey, botId string) string {
	return strings.Join([]string{appkey, botId}, "_")
}

func HandleBotMsg(ctx context.Context, msg *pbobjs.DownMsg) {
	if msg.MsgType != "jg:text" {
		return
	}
	txtMsg := &commonservices.TextMsg{}
	err := tools.JsonUnMarshal(msg.MsgContent, txtMsg)
	if err != nil {
		logs.WithContext(ctx).Errorf("text msg illigal. content:%s", string(msg.MsgContent))
		return
	}
	botId := bases.GetTargetIdFromCtx(ctx)
	botInfo := GetBotInfo(ctx, botId)
	if botInfo.BotEngine != nil {
		converId := ""
		buf := bytes.NewBuffer([]byte{})
		botInfo.BotEngine.StreamChat(ctx, msg.SenderId, converId, txtMsg.Content, func(answerPart string, isEnd bool) {
			if !isEnd {
				buf.WriteString(answerPart)
			}
		})
		answer := buf.String()
		if answer != "" {
			answerMsg := &commonservices.TextMsg{
				Content: answer,
			}
			flag := commonservices.SetStoreMsg(0)
			flag = commonservices.SetCountMsg(flag)
			commonservices.AsyncPrivateMsgOverUpstream(ctx, botId, msg.SenderId, &pbobjs.UpMsg{
				MsgType:    "jg:text",
				MsgContent: []byte(tools.ToJson(answerMsg)),
				Flags:      flag,
			})
		}
	}
}
