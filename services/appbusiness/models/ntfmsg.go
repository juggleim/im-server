package models

var GroupNotifyMsgType string = "jgd:grpntf"

type GroupNotify struct {
	Operator *User           `json:"operator"`
	Name     string          `json:"name"`
	Members  []*User         `json:"members"`
	Type     GroupNotifyType `json:"type"`
}

type GroupNotifyType int

const (
	GroupNotifyType_AddMember    = 1
	GroupNotifyType_RemoveMember = 2
	GroupNotifyType_Rename       = 3
)

var FriendNotifyMsgType string = "jgd:friendntf"

type FriendNotify struct {
	Type int `json:"type"`
}
