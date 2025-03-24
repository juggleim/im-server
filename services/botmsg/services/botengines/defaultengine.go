package botengines

import (
	"context"
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/msgdefines"
	"net/http"
	"time"
)

type DefaultEngine struct {
	Webhook string `json:"webhook"`
	ApiKey  string `json:"api_key"`
}

func (engine *DefaultEngine) StreamChat(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, question string, f func(part string, sectionStart, sectionEnd, isFinish bool)) {
	engine.Chat(ctx, senderId, targetId, channelType, question)
}

func (engine *DefaultEngine) Chat(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, question string) string {
	req := &Event{
		EventType: EventType_Message,
		Timestamp: time.Now().UnixMilli(),
		Payload:   []interface{}{},
	}
	msgEvent := &MsgEvent{
		Sender:     senderId,
		Receiver:   targetId,
		ConverType: int(channelType),
		MsgType:    msgdefines.InnerMsgType_Text,
		MsgContent: question,
	}
	fmt.Println("wehbook:", engine.Webhook, engine.ApiKey)
	fmt.Println("target:", targetId)
	req.Payload = append(req.Payload, msgEvent)
	body := tools.ToJson(req)
	headers := map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Bearer %s", engine.ApiKey)
	headers["Content-Type"] = "application/json"
	resp, code, err := tools.HttpDo(http.MethodPost, engine.Webhook, headers, body)
	fmt.Println("resp:", resp)
	fmt.Println("code:", code)
	fmt.Println("err:", err)
	return ""
}

type EventType string

const (
	EventType_Message EventType = "message"
)

type Event struct {
	EventType EventType     `json:"event_type"`
	Timestamp int64         `json:"timestamp"`
	Payload   []interface{} `json:"payload"`
}

type MsgEvent struct {
	Sender      string       `json:"sender"`
	Receiver    string       `json:"receiver"`
	ConverType  int          `json:"conver_type"`
	MsgType     string       `json:"msg_type"`
	MsgContent  string       `json:"msg_content"`
	MsgId       string       `json:"msg_id"`
	MsgTime     int64        `json:"msg_time"`
	MentionInfo *MentionInfo `json:"mention_info"`
}
type MentionInfo struct {
	MentionType   string   `json:"mention_type"`
	TargetUserIds []string `json:"target_user_ids"`
}
