package services

import (
	"bytes"
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/apimodels"
	"im-server/services/appbusiness/storages"
	"im-server/services/botmsg/services/botengines"
	botModels "im-server/services/botmsg/storages/models"
	"time"
)

func AssistantAnswer(ctx context.Context, req *apimodels.AssistantAnswerReq) (errs.IMErrorCode, *apimodels.AssistantAnswerResp) {
	if req == nil || len(req.Msgs) <= 0 {
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	buf := bytes.NewBuffer([]byte{})
	userId := bases.GetRequesterIdFromCtx(ctx)
	for _, msg := range req.Msgs {
		if msg.SenderId != userId {
			buf.WriteString(fmt.Sprintf("对方:%s\n", msg.Content))
		} else {
			buf.WriteString(fmt.Sprintf("我:%s\n", msg.Content))
		}
	}
	buf.WriteString("帮我生成回复")
	answer := GenerateAnswer(ctx, buf.String())
	return errs.IMErrorCode_SUCCESS, &apimodels.AssistantAnswerResp{
		Answer: answer,
	}
}

var assistantCache *caches.LruCache
var assistantLocks *tools.SegmentatedLocks

func init() {
	assistantCache = caches.NewLruCacheWithAddReadTimeout(1000, nil, 5*time.Minute, 5*time.Minute)
	assistantLocks = tools.NewSegmentatedLocks(32)
}

type AssistantInfo struct {
	AppKey      string
	AssistantId string
	OwnerId     string
	Nickname    string
	Portrait    string
	BotType     botModels.BotType
	BotEngine   botengines.IBotEngine
}

func GetAssistantInfo(ctx context.Context) *AssistantInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := appkey
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
				BotEngine: &botengines.NilBotEngine{},
			}
			storage := storages.NewAssistantStorage()
			ass, err := storage.FindEnableAssistant(appkey)
			if err == nil {
				assInfo.AssistantId = ass.AssistantId
				assInfo.Nickname = ass.Nickname
				assInfo.Portrait = ass.Portrait
				assInfo.BotType = ass.BotType
				switch assInfo.BotType {
				case botModels.BotType_Dify:
					difyBot := &botengines.DifyBotEngine{}
					err = tools.JsonUnMarshal([]byte(ass.BotConf), difyBot)
					if err == nil && difyBot.ApiKey != "" && difyBot.Url != "" {
						assInfo.BotEngine = difyBot
					}
				case botModels.BotType_Coze:
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
	assistantInfo := GetAssistantInfo(ctx)
	if assistantInfo != nil && assistantInfo.BotEngine != nil {
		buf := bytes.NewBuffer([]byte{})
		assistantInfo.BotEngine.StreamChat(ctx, bases.GetRequesterIdFromCtx(ctx), "", content, func(answerPart string, sectionStart, sectionEnd, isEnd bool) {
			if !isEnd {
				buf.WriteString(answerPart)
			}
		})
		return buf.String()
	}
	return "No Answer"
}
