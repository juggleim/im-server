package models

type FriendIds struct {
	UserId    string        `json:"user_id"`
	FriendIds []string      `json:"friend_ids"`
	Friends   []*FriendItem `json:"friends"`
}

type FriendItem struct {
	UserId      string `json:"user_id,omitempty"`
	FriendId    string `json:"friend_id"`
	DisplayName string `json:"display_name"`
}

type FriendsResp struct {
	Items  []*FriendItem `json:"items"`
	Offset string        `json:"offset"`
}
