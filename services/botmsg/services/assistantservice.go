package services

import (
	"bytes"
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/tools"
	"im-server/services/botmsg/services/botengines"
	"im-server/services/botmsg/storages"
	"im-server/services/botmsg/storages/models"
	"time"
)

var assistantCache *caches.LruCache
var assistantLocks *tools.SegmentatedLocks

func init() {
	assistantCache = caches.NewLruCacheWithAddReadTimeout(10000, nil, 5*time.Minute, 5*time.Minute)
	assistantLocks = tools.NewSegmentatedLocks(128)
}

type AssistantInfo struct {
	AppKey      string
	AssistantId string
	OwnerId     string
	Nickname    string
	Portrait    string
	BotType     models.BotType
	BotEngine   botengines.IBotEngine
}

func GetAssistantInfo(ctx context.Context, ownerId string) *AssistantInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getKey(appkey, ownerId)
	if val, exist := assistantCache.Get(key); exist {
		return val.(*AssistantInfo)
	} else {
		l := assistantLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if val, exist := assistantCache.Get(key); exist {
			return val.(*AssistantInfo)
		} else {
			assInfo := &AssistantInfo{
				AppKey:    appkey,
				OwnerId:   ownerId,
				BotEngine: &botengines.NilBotEngine{},
			}
			storage := storages.NewAssistantStorage()
			ass, err := storage.FindByOwnerId(appkey, ownerId)
			if err == nil {
				assInfo.AssistantId = ass.AssistantId
				assInfo.Nickname = ass.Nickname
				assInfo.Portrait = ass.Portrait
				assInfo.BotType = ass.BotType
				switch assInfo.BotType {
				case models.BotType_Dify:
					difyBot := &botengines.DifyBotEngine{}
					err = tools.JsonUnMarshal([]byte(ass.BotConf), difyBot)
					if err == nil && difyBot.ApiKey != "" && difyBot.Url != "" {
						assInfo.BotEngine = difyBot
					}
				case models.BotType_Coze:
					cozeBot := &botengines.CozeBotEngine{}
					err = tools.JsonUnMarshal([]byte(ass.BotConf), cozeBot)
					if err == nil && cozeBot.BotId != "" && cozeBot.Url != "" && cozeBot.Token != "" {
						assInfo.BotEngine = cozeBot
					}
				}
			}
			assistantCache.Add(key, assInfo)
			return assInfo
		}
	}
}

func GenerateAnswer(ctx context.Context, content string) string {
	userId := bases.GetRequesterIdFromCtx(ctx)
	assistantInfo := GetAssistantInfo(ctx, userId)
	if assistantInfo != nil && assistantInfo.BotEngine != nil {
		buf := bytes.NewBuffer([]byte{})
		assistantInfo.BotEngine.StreamChat(ctx, userId, "", content, func(answerPart string, sectionStart, sectionEnd, isEnd bool) {
			if !isEnd {
				buf.WriteString(answerPart)
			}
		})
		return buf.String()
	}
	return "No Answer"
}
