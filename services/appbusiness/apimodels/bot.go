package apimodels

type BotMsg struct {
	SenderId    string        `json:"sender_id"`
	BotId       string        `json:"bot_id"`
	ChannelType int           `json:"channel_type"`
	Stream      bool          `json:"stream"`
	Messages    []*BotMsgItem `json:"messages"`
}

type BotMsgItem struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type BotResponsePartData struct {
	Id             string `json:"id"`
	ConversationId string `json:"conversation_id"`
	Type           string `json:"type"`
	BotId          string `json:"bot_id"`
	Content        string `json:"content"`
	ContentType    string `json:"content_type"`
	SectionId      string `json:"section_id"`
	CreatedTime    int64  `json:"created_time"`
}
