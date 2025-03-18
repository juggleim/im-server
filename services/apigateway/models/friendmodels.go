package models

type FriendIds struct {
	UserId    string   `json:"user_id"`
	FriendIds []string `json:"friend_ids"`
}
