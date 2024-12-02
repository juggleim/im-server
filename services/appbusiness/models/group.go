package models

type Group struct {
	GroupId       string  `json:"group_id"`
	GroupName     string  `json:"group_name"`
	GroupPortrait string  `json:"group_portrait"`
	GrpMembers    []*User `json:"members,omitempty"`
	IsNotify      bool    `json:"is_notify"`
}

type Groups struct {
	Items  []*Group `json:"items"`
	Offset string   `json:"offset,omitempty"`
}
