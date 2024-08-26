package services

type Apps struct {
	Items   []*SimpleApp `json:"items"`
	HasMore bool         `json:"has_more"`
	Offset  string       `json:"offset"`
}
type SimpleApp struct {
	AppKey       string `json:"app_key"`
	AppType      int    `json:"app_type"`
	AppName      string `json:"app_name"`
	MaxUserCount int    `json:"max_user_count"`
	CurUserCount int    `json:"cur_user_count"`
	CreatedTime  int64  `json:"created_time"`
	EndedTime    int64  `json:"ended_time"`
}

type Accounts struct {
	Items   []*Account `json:"items"`
	HasMore bool       `json:"has_more"`
	Offset  string     `json:"offset"`
}
type Account struct {
	Account       string `json:"account"`
	State         int    `json:"state"`
	CreatedTime   int64  `json:"created_time"`
	UpdatedTime   int64  `json:"updated_time"`
	ParentAccount string `json:"parent_account"`
}

type AppInfo struct {
	AppType     int    `json:"app_type"`
	AppName     string `json:"app_name"`
	CreatedTime int64  `json:"created_time"`
	UpdateTime  int64  `json:"updated_time"`

	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	AppStatus int    `json:"app_status"`

	MaxUserCount int `json:"max_user_count"`
	CurUserCount int `json:"cur_user_count"`

	RestrictedFields *RestrictedFields `json:"restricted_fields"`
	ConfigFields     map[string]string `json:"config_fields"`

	ExpiredTime int64 `json:"expired_time"`
}
type RestrictedFields struct {
	MaxUserCount int32 `json:"max_user_count"`
}

type ConfigItem struct {
	Key   string      `json:"key"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
