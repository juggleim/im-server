package apimodels

type AssistantAnswerReq struct {
	ConverId    string            `json:"conver_id"`
	ChannelType int               `json:"channel_type"`
	PromptId    string            `json:"prompt_id"`
	Msgs        []*MsgContentItem `json:"msgs"`
}

type MsgContentItem struct {
	SenderId string `json:"sender_id"`
	Content  string `json:"content"`
	MsgTime  int64  `json:"msg_time"`
}

type AssistantAnswerResp struct {
	Answer      string `json:"answer"`
	StreamMsgId string `json:"stream_msg_id"`
}

type Prompt struct {
	Id          string `json:"id"`
	Prompts     string `json:"prompts"`
	CreatedTime int64  `json:"created_time"`
}

type PromptIds struct {
	Ids []string `json:"ids"`
}

type Prompts struct {
	Items  []*Prompt `json:"items"`
	Offset string    `json:"offset"`
}
