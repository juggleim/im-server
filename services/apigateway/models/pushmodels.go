package models

type PushPayload struct {
	FromUserId string `json:"from_user_id"`
	Condition  struct {
		TagsAnd []string `json:"tags_and"`
		TagsOr  []string `json:"tags_or"`
	} `json:"condition"`
	MsgBody struct {
		MsgType    string `json:"msg_type"`
		MsgContent string `json:"msg_content"`
	} `json:"msg_body,omitempty"`
	Notification struct {
		Title    string `json:"title"`
		PushText string `json:"push_text"`
	} `json:"notification,omitempty"`
}

func (p *PushPayload) Validate() bool {
	if p.FromUserId == "" {
		return false
	}
	if p.Condition.TagsAnd == nil && p.Condition.TagsOr == nil {
		return false
	}
	if len(p.Condition.TagsAnd) > 0 && len(p.Condition.TagsOr) > 0 {
		return false
	}

	return true
}

type UserTag struct {
	UserID string   `json:"user_id"`
	Tags   []string `json:"tags"`
}

type UserTagsPayload struct {
	UserTags []UserTag `json:"user_tags"`
}
