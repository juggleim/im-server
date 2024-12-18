package services

import (
	"time"

	"github.com/avast/retry-go/v4"
)

var retryStrategy = []retry.Option{
	retry.Delay(100 * time.Millisecond),
	retry.Attempts(3),
	retry.LastErrorOnly(true),
	retry.DelayType(retry.BackOffDelay),
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

/*
func SyncMsg2Bot(ctx context.Context, botId string, msg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	botInfo := GetBotInfo(ctx, botId)
	if botInfo != nil && botInfo.Webhook != "" {
		headers := map[string]string{
			"Content-Type": "application/json",
			"appkey":       appkey,
		}
		msgEvent := &MsgEvent{
			Sender:     msg.SenderId,
			Receiver:   msg.TargetId,
			ConverType: int(msg.ChannelType),
			MsgType:    msg.MsgType,
			MsgContent: string(msg.MsgContent),
			MsgId:      msg.MsgId,
			MsgTime:    msg.MsgTime,
		}
		event := &Event{
			EventType: EventType_Message,
			Timestamp: time.Now().UnixMilli(),
			Payload: []interface{}{
				msgEvent,
			},
		}
		body := tools.ToJson(event)
		err := breaker.Do(botInfo.Webhook, func() error {
			return retry.Do(func() error {
				return notify(botInfo.Webhook, headers, body)
			}, retryStrategy...)
		})
		if err == nil {
			logs.WithContext(ctx).Infof("success to sync msg to bot. bot_id:%s\tmsg_id:%s", botId, msg.MsgId)
		} else {
			logs.WithContext(ctx).Errorf("failed to sync msg to bot. bot_id:%s\tmsg_id:%s", botId, msg.MsgId)
		}
	} else {
		logs.WithContext(ctx).Warnf("no webhook. bot_id:%s", botId)
	}
}

func notify(url string, headers map[string]string, body string) error {
	_, code, err := tools.HttpDoBytes(http.MethodPost, url, headers, body)
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("sync bot msg failed. webhook:%s\tcode:%d", url, code)
	}
	return nil
}
*/
