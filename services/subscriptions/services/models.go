package services

type EventType string

const (
	EventType_Message EventType = "message"
	EventType_Online  EventType = "online"
	EventType_Offline EventType = "offline" //discarded
)

type SubEvent struct {
	EventType EventType     `json:"event_type"`
	Timestamp int64         `json:"timestamp"`
	Payload   []interface{} `json:"payload"`
}

type SubEventResp struct {
}

type MsgEvent struct {
	Platform    string       `json:"platform"`
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

type OnlineEvent struct {
	Type int32 `json:"type"`

	DepUserId string `json:"userId"`
	UserId    string `json:"user_id"`

	DepDeviceId string `json:"deviceId"`
	DeviceId    string `json:"device_id"`

	Platform string `json:"platform"`

	DepClientIp string `json:"clientIp"`
	ClientIp    string `json:"client_ip"`

	DepSessionId string `json:"sessionId"`
	SessionId    string `json:"session_id"`

	Timestamp int64 `json:"timestamp"`

	DepConnectionExt string `json:"connectionExt"`
	ConnectionExt    string `json:"connection_ext"`

	InstanceId string `json:"instance_id"`
}
