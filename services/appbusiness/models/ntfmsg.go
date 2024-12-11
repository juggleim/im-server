package models

import "im-server/commons/pbdefines/pbobjs"

var GroupNotifyMsgType string = "jgd:grpntf"

type GroupNotify struct {
	Operator *pbobjs.UserObj   `json:"operator"`
	Name     string            `json:"name"`
	Members  []*pbobjs.UserObj `json:"members"`
	Type     GroupNotifyType   `json:"type"`
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
