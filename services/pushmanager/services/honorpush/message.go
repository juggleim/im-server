package honorpush

// 以下结构与荣耀 Push Kit「下行消息」REST 请求体字段对齐，详见：
// https://developer.honor.com/cn/docs/11002/reference/downlink-message

type StyleType = byte

const (
	StyleTypeNormal  StyleType = iota // 标准样式
	StyleTypeBigText                  // 大文本样式
)

type ButtonActionType = byte

const (
	ButtonActionTypeLaunchApp  ButtonActionType = iota // 启动应用
	ButtonActionTypeJumpCustom                         // 打开应用自定义页面
	ButtonActionTypeJumpSite                           // 跳转网页
)

type ClickActionType = byte

const (
	ClickActionTypeOpenCustom ClickActionType = iota + 1 // 打开应用自定义页面
	ClickActionTypeOpenUrl                               // 打开特定 URL
	ClickActionTypeLaunchApp                             // 启动应用
)

type AndroidIntentType = byte

const (
	AndroidIntentTypeIntent AndroidIntentType = iota // 通过 intent URI
	AndroidIntentTypeAction                          // 通过 action
)

type PushCategory = string

const (
	PushCategoryLow    PushCategory = "LOW"    // 资讯营销类
	PushCategoryNormal PushCategory = "NORMAL" // 服务与通讯类
)

// PushAndroidButton 通知栏按钮（最多 3 个）
type PushAndroidButton struct {
	Name       string            `json:"name"`
	ActionType ButtonActionType  `json:"actionType,omitempty"`
	IntentType AndroidIntentType `json:"intentType,omitempty"`
	IntentUri  string            `json:"intent,omitempty"`
	IntentData string            `json:"data,omitempty"`
}

// PushAndroidClickAction 点击通知行为
type PushAndroidClickAction struct {
	ActionType   ClickActionType `json:"type"`
	URL          string          `json:"url,omitempty"`
	IntentUri    string          `json:"intent,omitempty"`
	IntentAction string          `json:"action,omitempty"`
}

// PushBadge 角标
type PushBadge struct {
	AddNum     uint16 `json:"addNum,omitempty"`
	BadgeClass string `json:"badgeClass,omitempty"`
	SetNum     uint16 `json:"setNum,omitempty"`
}

// PushAndroidNotification Android 通知栏内容
type PushAndroidNotification struct {
	NotifyID    int64                   `json:"notifyId,omitempty"`
	Style       StyleType               `json:"style,omitempty"`
	Title       string                  `json:"title,omitempty"`
	Body        string                  `json:"body,omitempty"`
	ImageURL    string                  `json:"image,omitempty"`
	ClickAction *PushAndroidClickAction `json:"clickAction,omitempty"`
	BigTitle    string                  `json:"bigTitle,omitempty"`
	BigBody     string                  `json:"bigBody,omitempty"`
	// RFC3339Nano 时间串，例如：2014-10-02T15:01:23.045123456Z
	When     string              `json:"when,omitempty"`
	Buttons  []PushAndroidButton `json:"buttons,omitempty"`
	Badge    *PushBadge          `json:"badge,omitempty"`
	Tag      string              `json:"tag,omitempty"`
	Group    string              `json:"group,omitempty"`
	Category PushCategory        `json:"importance,omitempty"`
}

// PushAndroidConfig Android 侧控制参数
type PushAndroidConfig struct {
	TTL   string `json:"ttl,omitempty"`
	BiTag string `json:"biTag,omitempty"`
	// Data 为 android 级别的透传负载，和顶层 data 并存时以业务需要为准。
	Data           string                   `json:"data,omitempty"`
	Notification   *PushAndroidNotification `json:"notification,omitempty"`
	TargetUserType int                      `json:"targetUserType,omitempty"`
}

// PushNotification 顶层通知体
type PushNotification struct {
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
	ImageURL string `json:"image,omitempty"`
}

// SendMessageReq 下行消息请求体（单播/按 token 列表）
type SendMessageReq struct {
	// Payload 自定义负载：通知栏消息可为 JSON 字符串；透传可为普通字符串或 JSON
	Payload      string             `json:"data,omitempty"`
	Token        []string           `json:"token"`
	Notification *PushNotification  `json:"notification,omitempty"`
	Android      *PushAndroidConfig `json:"android,omitempty"`
}

// SendMessageResp 下行消息响应
type SendMessageResp struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    struct {
		SendResult   bool     `json:"sendResult"`
		RequestID    string   `json:"requestId,omitempty"`
		FailTokens   []string `json:"failTokens,omitempty"`
		ExpireTokens []string `json:"expireTokens,omitempty"`
	} `json:"data,omitempty"`
}

type authTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type authErrorResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
