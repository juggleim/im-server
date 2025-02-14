package botengines

import (
	"context"
	"fmt"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"net/http"
)

type CozeBotEngine struct {
	Token string `json:"token"`
	Url   string `json:"url"`
	BotId string `json:"bot_id"`
}

func (engine *CozeBotEngine) StreamChat(ctx context.Context, senderId, converId string, question string, f func(part string, isEnd bool)) {
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
	stream, code, err := tools.CreateStream(http.MethodPost, engine.Url, headers, body)
	if err != nil || code != http.StatusOK {
		logs.WithContext(ctx).Errorf("call coze api failed. http_code:%d,err:%v", code, err)
		return
	}
	filter := map[string]bool{}
	for {
		line, err := stream.Receive()
		if err != nil {
			f("", true)
			return
		}
		item := CozeChatMsgRespItem{}
		err = tools.JsonUnMarshal([]byte(line), &item)
		if err != nil {
			f("", true)
			return
		}
		if item.Type == "answer" {
			if _, exist := filter[item.Id]; !exist {
				filter[item.Id] = true
				f(item.Content, false)
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
