package botengines

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/botmsg/storages"
	"im-server/services/botmsg/storages/models"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"net/http"
	"strings"
	"time"
)

type CozeBotEngine struct {
	Token string `json:"token"`
	Url   string `json:"url"`
	BotId string `json:"bot_id"`
}

func (engine *CozeBotEngine) Chat(ctx context.Context, senderId, converKey string, channelType pbobjs.ChannelType, question string) string {
	return ""
}

func (engine *CozeBotEngine) StreamChat(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, question string, f func(part string, sectionStart, sectionEnd, isEnd bool)) {
	converKey := ""
	if channelType == pbobjs.ChannelType_Private {
		converKey = commonservices.GetConversationId(senderId, targetId, pbobjs.ChannelType_Private)
		converKey = fmt.Sprintf("%s_%d", converKey, pbobjs.ChannelType_Private)
	} else if channelType == pbobjs.ChannelType_Group {
		converKey = commonservices.GetConversationId(senderId, targetId, pbobjs.ChannelType_Group)
		converKey = fmt.Sprintf("%s_%d", converKey, pbobjs.ChannelType_Group)
	}
	url := engine.Url
	cozeConverItem := GetCozeConverId(ctx, converKey, engine.Token)
	if cozeConverItem != nil && cozeConverItem.ConverId != "" {
		url = fmt.Sprintf("%s?conversation_id=%s", url, cozeConverItem.ConverId)
	}
	req := &CozeChatMsgReq{
		BotId:  engine.BotId,
		UserId: senderId,
		Stream: true,
		AdditionalMessages: []*CozeChatMsgItem{
			{
				Content:     question,
				ContentType: "text",
				Role:        "user",
				Name:        senderId,
			},
		},
	}
	body := tools.ToJson(req)
	headers := map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Bearer %s", engine.Token)
	headers["Content-Type"] = "application/json"
	stream, code, err := tools.CreateStream(http.MethodPost, url, headers, body)
	if err != nil || code != http.StatusOK {
		logs.WithContext(ctx).Errorf("call coze api failed. http_code:%d,err:%v", code, err)
		return
	}
	event := ""
	sectionStart := true
	for {
		line, err := stream.Receive()
		if err != nil {
			f("", false, false, true)
			return
		}
		if strings.TrimSpace(string(line)) == "\"[DONE]\"" {
			f("", false, false, true)
			return
		}
		if strings.HasPrefix(line, "event:") {
			event = strings.TrimPrefix(line, "event:")
			continue
		}
		line = strings.TrimPrefix(line, "data:")
		item := CozeChatMsgRespItem{}
		err = tools.JsonUnMarshal([]byte(line), &item)
		if err != nil {
			fmt.Println("unmarshal_err:", err, string(line))
			continue
		}
		if item.Type == "answer" && item.CreatedAt == 0 {
			if event == "conversation.message.delta" {
				f(item.Content, sectionStart, false, false)
				sectionStart = false
			} else if event == "conversation.message.completed" {
				f(item.Content, false, true, false)
				sectionStart = true
			}
		}
	}
}

type CozeChatMsgReq struct {
	BotId              string             `json:"bot_id"`
	UserId             string             `json:"user_id"`
	Stream             bool               `json:"stream"`
	AdditionalMessages []*CozeChatMsgItem `json:"additional_messages"`
}

type CozeChatMsgItem struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Role        string `json:"role"`
	Name        string `json:"name"`
}

type CozeChatMsgRespItem struct {
	Id             string `json:"id"`
	ConversationId string `json:"conversation_id"`
	BotId          string `json:"bot_id"`
	Role           string `json:"role"`
	Type           string `json:"type"`
	Content        string `json:"content"`
	ContentType    string `json:"content_type"`
	ChatId         string `json:"chat_id"`
	SectionId      string `json:"section_id"`
	CreatedAt      int64  `json:"created_at"`
}

var cozeConverCache *caches.LruCache
var cozeConverLock *tools.SegmentatedLocks

type CozeConverItem struct {
	AppKey    string
	ConverKey string
	ConverId  string
}

func init() {
	cozeConverCache = caches.NewLruCacheWithAddReadTimeout("coze_conver_cache", 10000, nil, 10*time.Minute, 10*time.Minute)
	cozeConverLock = tools.NewSegmentatedLocks(128)
}

func GetCozeConverId(ctx context.Context, converKey, token string) *CozeConverItem {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := strings.Join([]string{appkey, converKey}, "_")
	if val, exist := cozeConverCache.Get(key); exist {
		return val.(*CozeConverItem)
	} else {
		l := cozeConverLock.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if val, exist := cozeConverCache.Get(key); exist {
			return val.(*CozeConverItem)
		} else {
			item := &CozeConverItem{
				AppKey:    appkey,
				ConverKey: converKey,
			}
			storage := storages.NewBotConverStorage()
			botConver, err := storage.Find(appkey, models.BotConverType_Coze, converKey)
			if err == nil && botConver != nil && botConver.ConverId != "" {
				item.ConverId = botConver.ConverId
			} else {
				converId := createCozeConver(token)
				if converId != "" {
					item.ConverId = converId
					storage.Upsert(models.BotConver{
						AppKey:     appkey,
						ConverType: models.BotConverType_Coze,
						ConverKey:  converKey,
						ConverId:   converId,
					})
				}
			}
			cozeConverCache.Add(key, item)
			return item
		}
	}
}

func createCozeConver(token string) string {
	url := "https://api.coze.cn/v1/conversation/create"
	headers := map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	headers["Content-Type"] = "application/json"
	respBs, _, err := tools.HttpDoBytes(http.MethodPost, url, headers, "")
	if err == nil && len(respBs) > 0 {
		var resp CozeResp
		err = tools.JsonUnMarshal(respBs, &resp)
		if err == nil && resp.Data != nil && resp.Data.Id != "" {
			return resp.Data.Id
		}
	}
	return ""
}

type CozeResp struct {
	Code int             `json:"coze"`
	Data *CozeConverResp `json:"data"`
}
type CozeConverResp struct {
	Id            string `json:"id"`
	LastSectionId string `json:"last_section_id"`
	CreatedAt     int64  `json:"created_at"`
}
