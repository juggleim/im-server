package models

type GroupMuteReq struct {
	GrouopId string `json:"group_id"`
	IsMute   int    `json:"is_mute"`
}

type GroupInfo struct {
	GroupId       string            `json:"group_id"`
	GroupName     string            `json:"group_name"`
	GroupPortrait string            `json:"group_portrait"`
	IsMute        int               `json:"is_mute"`
	UpdatedTime   int64             `json:"updated_time"`
	ExtFields     map[string]string `json:"ext_fields"`
	Settings      map[string]string `json:"settings"`
}

type GroupMembersReq struct {
	GroupId       string            `json:"group_id"`
	GroupName     string            `json:"group_name"`
	GroupPortrait string            `json:"group_portrait"`
	ExtFields     map[string]string `json:"ext_fields"`
	MemberIds     []string          `json:"member_ids"`
}

type GroupMemberUpdateReq struct {
	GroupId        string            `json:"group_id"`
	MemberId       string            `json:"member_id"`
	GrpDisplayName string            `json:"grp_display_name"`
	ExtFields      map[string]string `json:"ext_fields"`
}

type GroupMember struct {
	MemberId       string            `json:"member_id"`
	IsMute         int               `json:"is_mute"`
	IsAllow        int               `json:"is_allow"`
	GrpDisplayName string            `json:"grp_display_name"`
	ExtFields      map[string]string `json:"ext_fields"`
}

type GroupMembersResp struct {
	Items  []*GroupMember `json:"items"`
	Offset string         `json:"offset"`
}

type GroupMemberMuteReq struct {
	GroupId    string   `json:"group_id"`
	MemberIds  []string `json:"member_ids"`
	IsMute     int      `json:"is_mute"`
	MuteMinute int      `json:"mute_minute"`
}

type GroupMemberAllowReq struct {
	GroupId   string   `json:"group_id"`
	MemberIds []string `json:"member_ids"`
	IsAllow   int      `json:"is_allow"`
}

type SetGroupSettingReq struct {
	GroupId  string                 `json:"group_id"`
	Settings map[string]interface{} `json:"settings"`
}
