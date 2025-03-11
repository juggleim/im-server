package botengines

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"net/http"
	"strings"
)

type CustomBotEngine struct {
	Url    string `json:"url"`
	ApiKey string `json:"api_key"`
	BotId  string `json:"bot_id"`
}

func (engine *CustomBotEngine) StreamChat(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, question string, f func(part string, sectionStart, sectionEnd, isFinish bool)) {
	url := engine.Url
	req := &CustomChatMsgReq{
		SenderId: senderId,
		BotId:    engine.BotId,
		Stream:   true,
		Messages: []*CustomChatMsgItem{
			{
				Content: question,
				Role:    "user",
			},
		},
	}
	body := tools.ToJson(req)
	headers := map[string]string{}
	headers["appkey"] = bases.GetAppKeyFromCtx(ctx)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", engine.ApiKey)
	headers["Content-Type"] = "application/json"
	stream, code, err := tools.CreateStream(http.MethodPost, url, headers, body)
	if err != nil || code != http.StatusOK {
		logs.WithContext(ctx).Errorf("call custom api failed. http_code%d,err:%v", code, err)
		return
	}
	sectionStart := true
	for {
		line, err := stream.Receive()
		if err != nil {
			f("", false, false, true)
			return
		}
		line = strings.TrimPrefix(line, "data:")
		item := CustomPartData{}
		err = tools.JsonUnMarshal([]byte(line), &item)
		if err != nil {
			f("", false, false, true)
			return
		}
		if item.Type == "message" {
			f(item.Content, sectionStart, false, false)
			sectionStart = false
		}
	}
}

func (engine *CustomBotEngine) Chat(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, question string) string {
	return ""
}

type CustomChatMsgReq struct {
	SenderId    string               `json:"sender_id"`
	BotId       string               `json:"bot_id"`
	ChannelType int                  `json:"channel_type"`
	Stream      bool                 `json:"stream"`
	Messages    []*CustomChatMsgItem `json:"messages"`
}
type CustomChatMsgItem struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type CustomPartData struct {
	Id             string `json:"id"`
	ConversationId string `json:"conversation_id"`
	Type           string `json:"type"`
	BotId          string `json:"bot_id"`
	Content        string `json:"content"`
	ContentType    string `json:"content_type"`
	SectionId      string `json:"section_id"`
	CreatedTime    int64  `json:"created_time"`
}
