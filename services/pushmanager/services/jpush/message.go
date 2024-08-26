package jpush

type Message struct {
	Content     string                 `json:"msg_content"`
	Title       string                 `json:"title,omitempty"`
	ContentType string                 `json:"content_type,omitempty"`
	Extras      map[string]interface{} `json:"extras,omitempty"`
}

func NewMessage() *Message {
	return &Message{
		Extras: make(map[string]interface{}),
	}
}

type SmsMessage struct {
	TempPara  interface{} `json:"temp_para,omitempty"`
	TempId    int64       `json:"temp_id"`
	DelayTime int         `json:"delay_time"`
}
