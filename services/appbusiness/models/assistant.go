package models

type AssistantAnswerReq struct {
	Msgs []*MsgContentItem `json:"msgs"`
}

type MsgContentItem struct {
	SenderId string `json:"sender_id"`
	Content  string `json:"content"`
	MsgTime  int64  `json:"msg_time"`
}

type AssistantAnswerResp struct {
	Answer string `json:"answer"`
}
