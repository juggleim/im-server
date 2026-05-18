package botengines

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CustomBotEngine struct {
	Url      string `json:"url"`
	ApiKey   string `json:"api_key"`
	BotId    string `json:"bot_id"`
	IsStream bool   `json:"is_stream"`

	lock                sync.Mutex
	consecutiveFailures int
	pauseUntil          time.Time
}

func (engine *CustomBotEngine) IsStreamChat() bool {
	return engine.IsStream
}

func (engine *CustomBotEngine) StreamChat(ctx context.Context, senderId, targetId string, msg *pbobjs.DownMsg, f func(part string, sectionStart, sectionEnd, isFinish bool)) {
	if msg.MsgType != msgdefines.InnerMsgType_Text {
		return
	}
	if !engine.allowRequest(ctx) {
		return
	}
	txtMsg := &msgdefines.TextMsg{}
	err := tools.JsonUnMarshal(msg.MsgContent, txtMsg)
	if err != nil {
		logs.WithContext(ctx).Errorf("text msg illigal. content:%s", string(msg.MsgContent))
		return
	}
	url := engine.Url
	req := &CustomChatStreamMsgReq{
		Sender:   senderId,
		Receiver: engine.BotId,
		Stream:   true,
		Messages: []*CustomChatMsgItem{
			{
				Content: txtMsg.Content,
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
		engine.recordFailure(ctx)
		logs.WithContext(ctx).Errorf("call custom api failed. http_code%d,err:%v", code, err)
		return
	}
	engine.recordSuccess()
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

func (engine *CustomBotEngine) Chat(ctx context.Context, senderId, targetId string, msg *pbobjs.DownMsg) string {
	if !engine.allowRequest(ctx) {
		return ""
	}
	receiver := engine.BotId
	if receiver == "" {
		receiver = targetId
	}
	req := &CustomChatMsgReq{
		Sender:      senderId,
		Receiver:    receiver,
		ConverType:  int(msg.ChannelType),
		MsgType:     msg.MsgType,
		MsgContent:  string(msg.MsgContent),
		MsgId:       msg.MsgId,
		MsgTime:     msg.MsgTime,
		MentionInfo: convertMentionInfo(msg.MentionInfo),
	}
	body := tools.ToJson(req)
	headers := map[string]string{}
	headers["appkey"] = bases.GetAppKeyFromCtx(ctx)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", engine.ApiKey)
	headers["Content-Type"] = "application/json"
	resp, code, err := tools.HttpDo(http.MethodPost, engine.Url, headers, body)
	if err != nil || code != http.StatusOK {
		engine.recordFailure(ctx)
		logs.WithContext(ctx).Errorf("call custom api failed. http_code:%d,err:%v", code, err)
		return ""
	}
	engine.recordSuccess()
	return resp
}

func (engine *CustomBotEngine) allowRequest(ctx context.Context) bool {
	engine.lock.Lock()
	defer engine.lock.Unlock()

	if !engine.pauseUntil.IsZero() {
		if time.Now().Before(engine.pauseUntil) {
			logs.WithContext(ctx).Infof("ignore custom api request during cooldown. bot_id:%s\tresume_time:%d", engine.BotId, engine.pauseUntil.UnixMilli())
			return false
		}
		engine.pauseUntil = time.Time{}
		engine.consecutiveFailures = 0
	}
	return true
}

func (engine *CustomBotEngine) recordFailure(ctx context.Context) {
	engine.lock.Lock()
	defer engine.lock.Unlock()

	if !engine.pauseUntil.IsZero() && time.Now().Before(engine.pauseUntil) {
		return
	}
	engine.consecutiveFailures++
	if engine.consecutiveFailures > 10 {
		engine.pauseUntil = time.Now().Add(time.Minute)
		logs.WithContext(ctx).Infof("pause custom api requests for 1 minute after continuous failures. bot_id:%s", engine.BotId)
	}
}

func (engine *CustomBotEngine) recordSuccess() {
	engine.lock.Lock()
	defer engine.lock.Unlock()

	engine.consecutiveFailures = 0
	engine.pauseUntil = time.Time{}
}

func convertMentionInfo(mentionInfo *pbobjs.MentionInfo) *MentionInfo {
	if mentionInfo == nil {
		return nil
	}
	mentionType := msgdefines.ToMentionTypeStr(mentionInfo.MentionType)
	targetUserIds := make([]string, 0, len(mentionInfo.TargetUsers))
	for _, user := range mentionInfo.TargetUsers {
		if user != nil && user.UserId != "" {
			targetUserIds = append(targetUserIds, user.UserId)
		}
	}
	if mentionType == "" && len(targetUserIds) == 0 {
		return nil
	}
	return &MentionInfo{
		MentionType:   mentionType,
		TargetUserIds: targetUserIds,
	}
}

type CustomChatStreamMsgReq struct {
	Sender     string               `json:"sender"`
	Receiver   string               `json:"receiver"`
	ConverType int                  `json:"conver_type"`
	Stream     bool                 `json:"stream"`
	Messages   []*CustomChatMsgItem `json:"messages"`
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

type CustomChatMsgReq struct {
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
