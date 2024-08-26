package services

import (
	"context"
	"encoding/json"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"net/http"
)

func Chat2Dify(ctx context.Context, botId string, msg *pbobjs.DownMsg, webhook, apiKey string) {
	if msg.MsgType == "jg:text" {
		var txtMsg commonservices.TextMsg
		err := tools.JsonUnMarshal(msg.MsgContent, txtMsg)
		if err != nil {
			logs.WithContext(ctx).Errorf("text msg illigal. content:%s", string(msg.MsgContent))
			return
		}
		req := &ChatMsgReq{
			Inputs:         map[string]string{},
			Query:          txtMsg.Content,
			ResponseMode:   "streaming",
			ConversationId: "",
			User:           "",
		}
		bs, _ := json.Marshal(req)
		body := string(bs)
		headers := map[string]string{}
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
		headers["Content-Type"] = "application/json"
		stream, code, err := tools.CreateStream(http.MethodPost, webhook, headers, body)
		if err != nil || code != http.StatusOK {
			logs.WithContext(ctx).Errorf("call dify api failed. http_code:%d,err:%v", code, err)
			return
		}
		for {
			line, err := stream.Receive()
			if err != nil {
				fmt.Println("xxxx:", err)
				return
			}
			fmt.Println(line)
		}
	}
}

func TestDify() {
	url := "https://api.dify.ai/v1/chat-messages"
	req := &ChatMsgReq{
		Inputs:         map[string]string{},
		Query:          "What are the specs of the iPhone 13 Pro Max?",
		ResponseMode:   "streaming",
		ConversationId: "",
		User:           "userid1",
	}
	bs, _ := json.Marshal(req)
	body := string(bs)
	headers := map[string]string{}
	headers["Authorization"] = "Bearer app-UD0yqEwQykpxA8hMbtzP0ktz"
	headers["Content-Type"] = "application/json"

	stream, code, err := tools.CreateStream("POST", url, headers, body)
	fmt.Println("code:", code)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		line, err := stream.Receive()
		if err != nil {
			fmt.Println("xxxx:", err)
			return
		}
		fmt.Println(line)
	}
}

type ChatMsgReq struct {
	Inputs         map[string]string `json:"inputs"`
	Query          string            `json:"query"`
	ResponseMode   string            `json:"response_mode"`
	ConversationId string            `json:"conversation_id"`
	User           string            `json:"user"`
}

type StreamRespItem struct {
	Event          string `json:"event"`
	ConversationId string `json:"conversation_id"`
	MessageId      string `json:"message_id"`
	CreatedAt      int64  `json:"created_at"`
	TaskId         string `json:"task_id"`
	Id             string `json:"id"`
	Answer         string `json:"answer"`

	Audio string `json:"audio"`
}
type MetaData struct {
	Usage *Usage `json:"usage"`
}
type Usage struct {
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
