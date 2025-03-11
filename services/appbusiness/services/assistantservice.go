package services

import (
	"bytes"
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/apimodels"
	"im-server/services/appbusiness/services/aiengines"
	"im-server/services/appbusiness/storages"
	"im-server/services/appbusiness/storages/models"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"time"

	"google.golang.org/protobuf/proto"
)

func AssistantAnswer(ctx context.Context, req *apimodels.AssistantAnswerReq) (errs.IMErrorCode, *apimodels.AssistantAnswerResp) {
	if req == nil {
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	userId := bases.GetRequesterIdFromCtx(ctx)
	promptStr := "你是一个智能回复生成器，能够根据用户提供的聊天记录，生成精彩回复。\n生成回复的一些限制条件：\n1. 只根据提供的聊天记录和上下文，生成回复，不进行无关的话题拓展；\n2. 确保回复的语音恰当、得体，不要产生冒犯性的表达；\n3. 回答简洁，不做过多延伸；\n4. 不要给我建议，直接以我的身份生成我该回复的内容；\n"
	appkey := bases.GetAppKeyFromCtx(ctx)
	if req.PromptId != "" {
		pId, err := tools.DecodeInt(req.PromptId)
		if err == nil && pId > 0 {
			storage := storages.NewPromptStorage()
			prompt, err := storage.FindPrompt(appkey, userId, pId)
			if err == nil && prompt != nil && prompt.Prompts != "" {
				promptStr = promptStr + "回复的其他要求：" + prompt.Prompts
			}
		}
	}
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("下面是对话内容，格式是 {userid}:{聊天内容}， 请根据这个对话内容，以我的身份生成一条回复。\n")
	if len(req.Msgs) > 0 {
		for _, msg := range req.Msgs {
			if msg.SenderId != userId {
				buf.WriteString(fmt.Sprintf("%s:%s\n", msg.SenderId, msg.Content))
			} else {
				buf.WriteString(fmt.Sprintf("我:%s\n", msg.Content))
			}
		}
	} else {
		if req.ConverId == "" || req.ChannelType == int(pbobjs.ChannelType_Unknown) {
			return errs.IMErrorCode_APP_DEFAULT, nil
		}
		//qry history msg
		converId := commonservices.GetConversationId(userId, req.ConverId, pbobjs.ChannelType(req.ChannelType))
		code, resp, err := bases.SyncRpcCall(ctx, "qry_hismsgs", converId, &pbobjs.QryHisMsgsReq{
			TargetId:    req.ConverId,
			ChannelType: pbobjs.ChannelType(req.ChannelType),
			Count:       5,
			MsgTypes:    []string{"jg:text"},
		}, func() proto.Message {
			return &pbobjs.DownMsgSet{}
		})
		if err == nil && code == errs.IMErrorCode_SUCCESS && resp != nil {
			downMsgs := resp.(*pbobjs.DownMsgSet)
			for _, msg := range downMsgs.Msgs {
				txtContent := &msgdefines.TextMsg{}
				err = tools.JsonUnMarshal(msg.MsgContent, txtContent)
				if err == nil {
					if msg.SenderId != userId {
						buf.WriteString(fmt.Sprintf("对方:%s\n", txtContent.Content))
					} else {
						buf.WriteString(fmt.Sprintf("我:%s\n", txtContent.Content))
					}
				}
			}
		}
	}
	if buf.Len() <= 0 {
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	content := buf.String()
	fmt.Println("----------------------------------------------------")
	fmt.Println(promptStr)
	fmt.Println("=")
	fmt.Println(content)
	fmt.Println("----------------------------------------------------")

	answer, streamMsgId := GenerateAnswer(ctx, promptStr, content, true)
	return errs.IMErrorCode_SUCCESS, &apimodels.AssistantAnswerResp{
		Answer:      answer,
		StreamMsgId: streamMsgId,
	}
}

var assistantCache *caches.LruCache
var assistantLocks *tools.SegmentatedLocks

func init() {
	assistantCache = caches.NewLruCacheWithAddReadTimeout(1000, nil, 5*time.Minute, 5*time.Minute)
	assistantLocks = tools.NewSegmentatedLocks(32)
}

type AssistantInfo struct {
	AppKey   string
	AiEngine aiengines.IAiEngine
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
				AppKey: appkey,
			}
			storage := storages.NewAiEngineStorage()
			ass, err := storage.FindEnableAiEngine(appkey)
			if err == nil {
				switch ass.EngineType {
				case models.EngineType_SiliconFlow:
					sfBot := &aiengines.SiliconFlowEngine{}
					err = tools.JsonUnMarshal([]byte(ass.EngineConf), sfBot)
					if err == nil && sfBot.ApiKey != "" && sfBot.Url != "" && sfBot.Model != "" {
						assInfo.AiEngine = sfBot
					}
				}
			}
			assistantCache.Add(key, assInfo)
			return assInfo
		}
	}
}

func GenerateAnswer(ctx context.Context, prompt, question string, isSync bool) (string, string) {
	userId := bases.GetRequesterIdFromCtx(ctx)
	streamMsgId := tools.GenerateMsgId(time.Now().UnixMilli(), int32(pbobjs.ChannelType_System), "assistant")
	assistantInfo := GetAssistantInfo(ctx)
	if assistantInfo != nil && assistantInfo.AiEngine != nil {
		if isSync {
			buf := bytes.NewBuffer([]byte{})
			assistantInfo.AiEngine.StreamChat(ctx, bases.GetRequesterIdFromCtx(ctx), "assistant", prompt, question, func(answerPart string, isEnd bool) {
				if !isEnd {
					buf.WriteString(answerPart)
				}
			})
			return buf.String(), streamMsgId
		} else {
			go func() {
				assistantInfo.AiEngine.StreamChat(ctx, bases.GetRequesterIdFromCtx(ctx), "assistant", prompt, question, func(answerPart string, isEnd bool) {
					if !isEnd {
						commonservices.AsyncSystemMsg(ctx, "assistant", userId, &pbobjs.UpMsg{
							MsgType: "jgs:aianswer",
							MsgContent: []byte(tools.ToJson(&StreamMsg{
								Content:     answerPart,
								StreamMsgId: streamMsgId,
							})),
						})
					}
				})
			}()

			return "", streamMsgId
		}
	}
	return "No Answer", streamMsgId
}

type StreamMsg struct {
	Content     string `json:"content,omitempty"`
	StreamMsgId string `json:"stream_msg_id"`
}
