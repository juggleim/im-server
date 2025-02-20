package apimodels

import "im-server/commons/pbdefines/pbobjs"

var (
	UserExtKey_Phone            string = "phone"
	UserExtKey_Language         string = "language"
	UserExtKey_Undisturb        string = "undisturb"
	UserExtKey_FriendVerifyType string = "friend_verify_type"
	UserExtKey_GrpVerifyType    string = "grp_verify_type"
)

type Users struct {
	Items []*pbobjs.UserObj `json:"items"`
}

type Friends struct {
	Items  []*pbobjs.UserObj `json:"items"`
	Offset string            `json:"offset,omitempty"`
}

type Friend struct {
	UserId   string `json:"user_id"`
	FriendId string `json:"friend_id"`
}

type FriendIds struct {
	FriendIds []string `json:"friend_ids"`
}
