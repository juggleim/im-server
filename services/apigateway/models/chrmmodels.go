package models

type ChatroomInfo struct {
	ChatId      string            `json:"chat_id"`
	ChatName    string            `json:"chat_name"`
	Members     []*ChatroomMember `json:"members"`
	Atts        []*ChatroomAtt    `json:"atts"`
	MemberCount int32             `json:"member_count"`
	IsMute      int               `json:"is_mute"`
}
type ChatroomMember struct {
	MemberId   string `json:"member_id"`
	MemberName string `json:"member_name"`
	AddedTime  int64  `json:"added_time"`
}

type ChatroomAtt struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	AttTime int64  `json:"att_time"`
	UserId  string `json:"user_id"`
}

type ChrmBanUserReq struct {
	ChatId    string   `json:"chat_id"`
	MemberIds []string `json:"member_ids"`
}

type ChrmBanUsers struct {
	ChatId  string            `json:"chat_id"`
	Members []*ChatroomMember `json:"members"`
	Offset  string            `json:"offset"`
}
