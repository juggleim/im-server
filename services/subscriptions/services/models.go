package services

type EventType string

const (
	EventType_Message EventType = "message"
	EventType_Online  EventType = "online"
	EventType_Offline EventType = "offline"
)

type SubEvent struct {
	EventType EventType     `json:"event_type"`
	Timestamp int64         `json:"timestamp"`
	Payload   []interface{} `json:"payload"`
}

type SubEventResp struct {
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
