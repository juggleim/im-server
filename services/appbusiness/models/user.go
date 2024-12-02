package models

type User struct {
	UserId   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Status   int    `json:"status"`
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty"`
	Language string `json:"language,omitempty"`
	Province string `json:"province,omitempty"`
	IsFriend bool   `json:"is_friend"`
}

type Users struct {
	Items []*User `json:"items"`
}

type Friends struct {
	Items  []*User `json:"items"`
	Offset string  `json:"offset,omitempty"`
}

type Friend struct {
	UserId   string `json:"user_id"`
	FriendId string `json:"friend_id"`
}
