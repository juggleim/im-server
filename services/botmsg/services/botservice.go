package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
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
	BotId     string
	Nickname  string
	Portrait  string
	ExtFields []*pbobjs.KvItem
	Webhook   string
	BotType   string
	APIKey    string
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
			bInfo := commonservices.GetUserInfoFromRpcWithAttTypes(ctx, botId, []int32{int32(commonservices.AttItemType_Setting)})
			botInfo := &BotInfo{
				BotId:     bInfo.UserId,
				Nickname:  bInfo.Nickname,
				Portrait:  bInfo.UserPortrait,
				ExtFields: bInfo.ExtFields,
			}
			if len(bInfo.Settings) > 0 {
				settingMap := commonservices.Kvitems2Map(bInfo.Settings)
				if webhook, exist := settingMap[string(commonservices.AttItemKey_Bot_WebHook)]; exist && webhook != "" {
					botInfo.Webhook = webhook
				}
				if apiKey, exist := settingMap[string(commonservices.AttItemKey_Bot_ApiKey)]; exist && apiKey != "" {
					botInfo.APIKey = apiKey
				}
				if botType, exist := settingMap[string(commonservices.AttItemKey_Bot_Type)]; exist && botType != "" {
					botInfo.BotType = botType
				}
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
	botId := bases.GetTargetIdFromCtx(ctx)
	botInfo := GetBotInfo(ctx, botId)
	if botInfo.BotType == "dify" {
		if botInfo.Webhook == "" || botInfo.APIKey == "" {
			logs.WithContext(ctx).Infof("no webhook/apikey")
			return
		}
		//https://api.dify.ai/v1/chat-messages
		//app-UD0yqEwQykpxA8hMbtzP0ktz
		Chat2Dify(ctx, botId, msg, botInfo.Webhook, botInfo.APIKey)
	} else {
		SyncMsg2Bot(ctx, botId, msg)
	}
}
