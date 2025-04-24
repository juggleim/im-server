package models

// user register
type UserRegResp struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfo struct {
	UserId       string            `json:"user_id"`
	Nickname     string            `json:"nickname"`
	UserPortrait string            `json:"user_portrait"`
	ExtFields    map[string]string `json:"ext_fields"`
	UpdatedTime  int64             `json:"updated_time"`
}

type KickUserReq struct {
	UserId    string   `json:"user_id"`
	Platforms []string `json:"platforms"`
	DeviceIds []string `json:"device_ids"`
	Ext       string   `json:"ext"`
}

type UserIds struct {
	UserIds []string `json:"user_ids"`
}

// user online
type UserOnlineStatusReq struct {
	UserIds []string `json:"user_ids"`
}
type UserOnlineStatusResp struct {
	Items []*UserOnlineStatusItem `json:"items"`
}
type UserOnlineStatusItem struct {
	UserId   string `json:"user_id"`
	IsOnline bool   `json:"is_online"`
}

// user ban
type BanUsersReq struct {
	Items []*BanUser `json:"items"`
}

type QryBanUsersResp struct {
	Items  []*BanUser `json:"items"`
	Offset string     `json:"offset"`
}

type BanUser struct {
	UserId        string `json:"user_id"`
	CreatedTime   int64  `json:"created_time"`
	EndTime       int64  `json:"end_time"`
	EndTimeOffset int64  `json:"end_time_offset"`
	ScopeKey      string `json:"scope_key"`
	ScopeValue    string `json:"scope_value"`
	Ext           string `json:"ext,omitempty"`
}

func (user *BanUser) GetTargetId() string {
	return user.UserId
}

// user block
type BlockUsersReq struct {
	UserId       string   `json:"user_id"`
	BlockUserIds []string `json:"block_user_ids"`
}

type QryBlockUsersResp struct {
	UserId string       `json:"user_id"`
	Items  []*BlockUser `json:"items"`
	Offset string       `json:"offset"`
}

type BlockUser struct {
	BlockUserId string `json:"block_user_id"`
	CreatedTime int64  `json:"createed_time"`
}

// private mute
type MuteUser struct {
	UserId      string `json:"user_id"`
	CreatedTime int64  `json:"created_time"`
}
type QryMuteUsersResp struct {
	Items  []*MuteUser `json:"items"`
	Offset string      `json:"offset"`
}

// undisturb convers
type UndisturbConversReq struct {
	UserId string                 `json:"user_id"`
	Items  []*UndisturbConverItem `json:"items"`
}

type UndisturbConverItem struct {
	TargetId      string `json:"target_id"`
	ChannelType   int    `json:"channel_type"`
	UndisturbType int32  `json:"undisturb_type"`
}

type UserSettings struct {
	UserId   string                 `json:"user_id"`
	Settings map[string]interface{} `json:"settings"`
}

// top convers
type TopConversReq struct {
	UserId string              `json:"user_id"`
	Items  []*TopConverReqItem `json:"items"`
}

type TopConverReqItem struct {
	TargetId    string `json:"target_id"`
	ChannelType int    `json:"channel_type"`
	IsTop       bool   `json:"is_top"`
}
