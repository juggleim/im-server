package aiengines

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/tools"
	"im-server/services/appbusiness/storages"
	"im-server/services/appbusiness/storages/models"
	"time"
)

type AssistantEngineType int

var (
	AssistantEngineType_SiliconFlow AssistantEngineType = 1
)

type IAiEngine interface {
	StreamChat(ctx context.Context, senderId, converId string, prompt string, question string, f func(answerPart string, isEnd bool))
}

var aiEngineCache *caches.LruCache
var aiEngineLocks *tools.SegmentatedLocks

func init() {
	aiEngineCache = caches.NewLruCacheWithAddReadTimeout("assistant_cache", 1000, nil, 5*time.Minute, 5*time.Minute)
	aiEngineLocks = tools.NewSegmentatedLocks(32)
}

type AiEngineInfo struct {
	AppKey   string
	AiEngine IAiEngine
}

func GetAiEngineInfo(ctx context.Context) *AiEngineInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := appkey
	if val, exist := aiEngineCache.Get(key); exist {
		return val.(*AiEngineInfo)
	} else {
		l := aiEngineLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if val, exist := aiEngineCache.Get(key); exist {
			return val.(*AiEngineInfo)
		} else {
			aiEngineInfo := &AiEngineInfo{
				AppKey: appkey,
			}
			storage := storages.NewAiEngineStorage()
			ass, err := storage.FindEnableAiEngine(appkey)
			if err == nil {
				switch ass.EngineType {
				case models.EngineType_SiliconFlow:
					sfBot := &SiliconFlowEngine{}
					err = tools.JsonUnMarshal([]byte(ass.EngineConf), sfBot)
					if err == nil && sfBot.ApiKey != "" && sfBot.Url != "" && sfBot.Model != "" {
						aiEngineInfo.AiEngine = sfBot
					}
				}
			}
			aiEngineCache.Add(key, aiEngineInfo)
			return aiEngineInfo
		}
	}
}
