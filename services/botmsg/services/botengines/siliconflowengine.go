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

type SiliconFlowEngine struct {
	ApiKey string `json:"api_key"`
	Url    string `json:"url"`
	Model  string `json:"model"`
}

func (engine *SiliconFlowEngine) StreamChat(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, question string, f func(part string, sectionStart, sectionend, isFinish bool)) {
	headers := map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Bearer %s", engine.ApiKey)
	headers["Content-Type"] = "application/json"
	req := &SiliconFlowChatReq{
		Model: engine.Model,
		Messages: []*SiliconFlowChatMsg{
			{
				Role:    "user",
				Content: question,
			},
		},
		Stream:    true,
		MaxTokens: 512,
		// Tools:     getTools(),
	}
	body := tools.ToJson(req)
	stream, code, err := tools.CreateStream(http.MethodPost, engine.Url, headers, body)
	if err != nil || code != http.StatusOK {
		logs.WithContext(ctx).Errorf("call siliconflow api failed. http_code:%d,err:%v", code, err)
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
		item := SiliconFlowChatResp{}
		err = tools.JsonUnMarshal([]byte(line), &item)
		if err != nil {
			f("", false, false, true)
			return
		}
		if len(item.Choices) > 0 {
			for _, choice := range item.Choices {
				if choice.Delta != nil {
					if choice.FinishReason != "stop" {
						f(choice.Delta.Content, sectionStart, false, false)
						sectionStart = false
					} else {
						f(choice.Delta.Content, false, false, true)
					}
				}
			}
			fmt.Println("x:", line)
		}
	}
}

func (engine *SiliconFlowEngine) Chat(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, question string) string {
	headers := map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Bearer %s", engine.ApiKey)
	headers["Content-Type"] = "application/json"
	req := &SiliconFlowChatReq{
		Model: engine.Model,
		Messages: []*SiliconFlowChatMsg{
			{
				Role:    "user",
				Content: question,
			},
		},
		Stream:    false,
		MaxTokens: 512,
		Tools:     getTools(),
	}
	body := tools.ToJson(req)
	fmt.Println(body)
	resp, code, err := tools.HttpDoBytesWithTimeout(http.MethodPost, engine.Url, headers, body, 0)
	if err != nil || code != http.StatusOK {
		logs.WithContext(ctx).Errorf("call siliconflow api failed. http_code:%d,err:%v", code, err)
		return ""
	}
	return string(resp)
}

func getTools() []*SiliconFlowTool {
	ret := []*SiliconFlowTool{}
	ret = append(ret, &SiliconFlowTool{
		Type: "function",
		Function: &SiliconFlowFunction{
			Name:        "compare",
			Description: "Compare two number, which one is bigger",
			Parameters: &SfParameters{
				Type: "object",
				Properties: map[string]*SfParamProperty{
					"a": {
						Type:        "int",
						Description: "A number",
					},
					"b": {
						Type:        "int",
						Description: "A number",
					},
				},
				Required: []string{"a", "b"},
			},
		},
	})
	return ret
}

type SiliconFlowChatReq struct {
	Model     string                `json:"model"`
	Messages  []*SiliconFlowChatMsg `json:"messages"`
	Stream    bool                  `json:"stream"`
	MaxTokens int                   `json:"max_tokens"`
	Tools     []*SiliconFlowTool    `json:"tools,omitempty"`
}

type SiliconFlowChatMsg struct {
	Role      string        `json:"role"`
	Content   string        `json:"content"`
	ToolCalls []*SfToolCall `json:"tool_calls,omitempty"`
}

type SfToolCall struct {
	Id       string               `json:"id"`
	Type     string               `json:"type"`
	Function *SiliconFlowFunction `json:"function,omitempty"`
}

type SiliconFlowChatResp struct {
	Id                string               `json:"id"`
	Object            string               `json:"object"`
	Created           int64                `json:"created"`
	Model             string               `json:"model"`
	Choices           []*SiliconFlowChoice `json:"choices"`
	SystemFingerprint string               `json:"system_fingerprint"`
	Usage             *SiliconFlowUsage    `json:"usage"`
}

type SiliconFlowChoice struct {
	Index        int                     `json:"index"`
	Delta        *SiliconFlowChoiceDelta `json:"delta"`
	Message      *SiliconFlowChatMsg     `json:"message"`
	FinishReason string                  `json:"finish_reason"`
}

type SiliconFlowChoiceDelta struct {
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
}

type SiliconFlowUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type SiliconFlowTool struct {
	Type     string               `json:"type"`
	Function *SiliconFlowFunction `json:"function"`
}

type SiliconFlowFunction struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Parameters  *SfParameters `json:"parameters"`
	Arguments   string        `json:"arguments,omitempty"`
}

type SfParameters struct {
	Type       string                      `json:"type"`
	Properties map[string]*SfParamProperty `json:"properties"`
	Required   []string                    `json:"required"`
}

type SfParamProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}
