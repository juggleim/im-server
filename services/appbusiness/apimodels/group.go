package apimodels

import "im-server/commons/pbdefines/pbobjs"

type Group struct {
	GroupId         string           `json:"group_id"`
	GroupName       string           `json:"group_name"`
	GroupPortrait   string           `json:"group_portrait"`
	GrpMembers      []*GroupMember   `json:"members,omitempty"`
	MemberIds       []string         `json:"member_ids,omitempty"`
	MemberCount     int              `json:"member_count"`
	Owner           *pbobjs.UserObj  `json:"owner,omitempty"`
	MyRole          int              `json:"my_role"`
	GroupManagement *GroupManagement `json:"group_management"`
}

type GroupManagement struct {
	GroupMute       int `json:"group_mute"`
	MaxAdminCount   int `json:"max_admin_count"`
	AdminCount      int `json:"admin_count"`
	GroupVerifyType int `json:"group_verify_type"`
}

type Groups struct {
	Items  []*Group `json:"items"`
	Offset string   `json:"offset,omitempty"`
}

type GroupAnnouncement struct {
	GroupId string `json:"group_id"`
	Content string `json:"content"`
}

type GroupMember struct {
	pbobjs.UserObj
}

type GroupMembersResp struct {
	Items  []*GroupMember `json:"items"`
	Offset string         `json:"offset"`
}

type CheckGroupMembersReq struct {
	GroupId   string   `json:"group_id"`
	MemberIds []string `json:"member_ids"`
}

type CheckGroupMembersResp struct {
	GroupId        string          `json:"group_id"`
	MemberExistMap map[string]bool `json:"member_exist_map"`
}
