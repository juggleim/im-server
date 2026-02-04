package models

type Conversations struct {
	UserId     string          `json:"user_id,omitempty"`
	Items      []*Conversation `json:"items"`
	IsFinished bool            `json:"is_finished"`
}

type Conversation struct {
	Id          string `json:"id"`
	UserId      string `json:"user_id"`
	TargetId    string `json:"target_id"`
	ChannelType int    `json:"channel_type"`
	SubChannel  string `json:"sub_channel"`
	Time        int64  `json:"time"`
}

// top convers
type TopConversReq struct {
	UserId string              `json:"user_id"`
	Items  []*TopConverReqItem `json:"items"`
}

type TopConverReqItem struct {
	TargetId    string `json:"target_id"`
	ChannelType int    `json:"channel_type"`
	IsTop       bool   `json:"is_top"`
}

// tag convers
type TagConversReq struct {
	UserId  string `json:"user_id"`
	Tag     string `json:"tag"`
	TagName string `json:"tag_name"`

	Convers []*Conversation `json:"convers"`
}
