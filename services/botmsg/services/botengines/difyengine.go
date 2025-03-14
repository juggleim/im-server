package botengines

import (
	"context"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"net/http"
	"strings"
)

type DifyBotEngine struct {
	ApiKey string `json:"api_key"`
	Url    string `json:"url"`
}

func (engine *DifyBotEngine) Chat(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, question string) string {
	return ""
}

func (engine *DifyBotEngine) StreamChat(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, question string, f func(answerPart string, sectionStart, sectionEnd, isEnd bool)) {
	req := &DifyChatMsgReq{
		Inputs:         map[string]string{},
		Query:          question,
		ResponseMode:   "streaming",
		ConversationId: "",
		User:           senderId,
	}
	body := tools.ToJson(req)
	headers := map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Bearer %s", engine.ApiKey)
	headers["Content-Type"] = "application/json"
	stream, code, err := tools.CreateStream(http.MethodPost, engine.Url, headers, body)
	if err != nil || code != http.StatusOK {
		logs.WithContext(ctx).Errorf("call dify api failed. http_code:%d,err:%v", code, err)
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
		item := DifyStreamRespItem{}
		err = tools.JsonUnMarshal([]byte(line), &item)
		if err != nil {
			f("", false, false, true)
			return
		}
		if item.Event == "message" {
			f(item.Answer, sectionStart, false, false)
			sectionStart = false
		} else if item.Event == "message_end" {
			f(item.Answer, false, false, true)
			return
		}
	}
}

type DifyChatMsgReq struct {
	Inputs         map[string]string `json:"inputs"`
	Query          string            `json:"query"`
	ResponseMode   string            `json:"response_mode"`
	ConversationId string            `json:"conversation_id"`
	User           string            `json:"user"`
}

type DifyStreamRespItem struct {
	Event          string `json:"event"`
	ConversationId string `json:"conversation_id"`
	MessageId      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
	TaskId         string `json:"task_id"`
	Id             string `json:"id"`
	Answer         string `json:"answer"`

	Audio string `json:"audio"`
}

type DifyMetaData struct {
	Usage *DifyUsage `json:"usage"`
}

type DifyUsage struct {
	PromptTokens        int32   `json:"prompt_tokens"`
	PromptUnitPrice     string  `json:"prompt_price_unit"`
	PromptPrice         string  `json:"prompt_price"`
	CompletionTokens    int32   `json:"completion_tokens"`
	CompletionUnitPrice string  `json:"completion_unit_price"`
	CompletionPriceUnit string  `json:"completion_price_unit"`
	CompletionPrice     string  `json:"completion_price"`
	TotalTokens         int32   `json:"total_tokens"`
	TotalPrice          string  `json:"total_price"`
	Currency            string  `json:"currency"`
	Latency             float64 `json:"latency"`
}
