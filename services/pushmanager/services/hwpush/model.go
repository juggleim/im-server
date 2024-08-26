package hwpush

type MessageRequest struct {
	ValidateOnly bool     `json:"validate_only"`
	Message      *Message `json:"message"`
}

type MessageResponse struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"requestId"`
}

type Message struct {
	Data         string         `json:"data,omitempty"`
	Notification *Notification  `json:"notification,omitempty"`
	Android      *AndroidConfig `json:"android,omitempty"`
	Apns         *Apns          `json:"apns,omitempty"`
	WebPush      *WebPushConfig `json:"webpush,omitempty"`
	Token        []string       `json:"token,omitempty"`
	Topic        string         `json:"topic,omitempty"`
	Condition    string         `json:"condition,omitempty"`
}

type Notification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Image string `json:"image,omitempty"`
}

type Apns struct {
	Headers    *ApnsHeaders           `json:"headers,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
	HmsOptions *ApnsHmsOptions        `json:"hms_options,omitempty"`
}

type ApnsHmsOptions struct {
	TargetUserType int `json:"target_user_type,omitempty"`
}

type ApnsHeaders struct {
	Authorization  string `json:"authorization,omitempty"`
	ApnsId         string `json:"apns-id,omitempty"`
	ApnsExpiration int64  `json:"apns-expiration,omitempty"`
	ApnsPriority   string `json:"apns-priority,omitempty"`
	ApnsTopic      string `json:"apns-topic,omitempty"`
	ApnsCollapseId string `json:"apns-collapse-id,omitempty"`
}

type Aps struct {
	Alert            interface{} `json:"alert,omitempty"` // dictionary or string
	Badge            int         `json:"badge,omitempty"`
	Sound            string      `json:"sound,omitempty"`
	ContentAvailable int         `json:"content-available,omitempty"`
	Category         string      `json:"category,omitempty"`
	ThreadId         string      `json:"thread-id,omitempty"`
}

type AlertDictionary struct {
	Title        string   `json:"title,omitempty"`
	Body         string   `json:"body,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`
	ActionLocKey string   `json:"action-loc-key,omitempty"`
	LocKey       string   `json:"loc-key,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`
}

type WebPushConfig struct {
	Data         string               `json:"data,omitempty"`
	Headers      *WebPushHeaders      `json:"headers,omitempty"`
	HmsOptions   *HmsWebPushOption    `json:"hms_options,omitempty"`
	Notification *WebPushNotification `json:"notification,omitempty"`
}

type WebPushHeaders struct {
	TTL     string `json:"ttl,omitempty"`
	Topic   string `json:"topics,omitempty"`
	Urgency string `json:"urgency,omitempty"`
}

type HmsWebPushOption struct {
	Link string `json:"link,omitempty"`
}

type WebPushNotification struct {
	Title              string           `json:"title,omitempty"`
	Body               string           `json:"body,omitempty"`
	Actions            []*WebPushAction `json:"actions,omitempty"`
	Badge              string           `json:"badge,omitempty"`
	Dir                string           `json:"dir,omitempty"`
	Icon               string           `json:"icon,omitempty"`
	Image              string           `json:"image,omitempty"`
	Lang               string           `json:"lang,omitempty"`
	Renotify           bool             `json:"renotify,omitempty"`
	RequireInteraction bool             `json:"require_interaction,omitempty"`
	Silent             bool             `json:"silent,omitempty"`
	Tag                string           `json:"tag,omitempty"`
	Timestamp          int64            `json:"timestamp,omitempty"`
	Vibrate            []int            `json:"vibrate,omitempty"`
}

type WebPushAction struct {
	Action string `json:"action,omitempty"`
	Icon   string `json:"icon,omitempty"`
	Title  string `json:"title,omitempty"`
}

type AndroidConfig struct {
	CollapseKey   int                  `json:"collapse_key,omitempty"`
	Urgency       string               `json:"urgency,omitempty"`
	Category      string               `json:"category,omitempty"`
	TTL           string               `json:"ttl,omitempty"`
	BiTag         string               `json:"bi_tag,omitempty"`
	FastAppTarget int                  `json:"fast_app_target,omitempty"`
	Data          string               `json:"data,omitempty"`
	Notification  *AndroidNotification `json:"notification,omitempty"`
}

type AndroidNotification struct {
	Title         string                 `json:"title,omitempty"`
	Body          string                 `json:"body,omitempty"`
	Icon          string                 `json:"icon,omitempty"`
	Color         string                 `json:"color,omitempty"`
	Sound         string                 `json:"sound,omitempty"`
	DefaultSound  bool                   `json:"default_sound,omitempty"`
	Tag           string                 `json:"tag,omitempty"`
	ClickAction   *ClickAction           `json:"click_action,omitempty"`
	BodyLocKey    string                 `json:"body_loc_key,omitempty"`
	BodyLocArgs   []string               `json:"body_loc_args,omitempty"`
	TitleLocKey   string                 `json:"title_loc_key,omitempty"`
	TitleLocArgs  []string               `json:"title_loc_args,omitempty"`
	MultiLangKey  map[string]interface{} `json:"multi_lang_key,omitempty"`
	ChannelId     string                 `json:"channel_id,omitempty"`
	NotifySummary string                 `json:"notify_summary,omitempty"`
	Image         string                 `json:"image,omitempty"`
	Style         int                    `json:"style,omitempty"`
	BigTitle      string                 `json:"big_title,omitempty"`
	BigBody       string                 `json:"big_body,omitempty"`

	AutoClear         int                `json:"auto_clear,omitempty"`
	NotifyId          int                `json:"notify_id,omitempty"`
	Group             string             `json:"group,omitempty"`
	Badge             *BadgeNotification `json:"badge,omitempty"`
	Ticker            string             `json:"ticker,omitempty"`
	AutoCancel        bool               `json:"auto_cancel,omitempty"`
	When              string             `json:"when,omitempty"`
	Importance        string             `json:"importance,omitempty"`
	UseDefaultVibrate bool               `json:"use_default_vibrate,omitempty"`
	UseDefaultLight   bool               `json:"use_default_light,omitempty"`
	VibrateConfig     []string           `json:"vibrate_config,omitempty"`
	Visibility        string             `json:"visibility,omitempty"`
	LightSettings     *LightSettings     `json:"light_settings,omitempty"`
	ForegroundShow    bool               `json:"foreground_show,omitempty"`
}

type ClickAction struct {
	Type         int    `json:"type"` //when the type equals to 1, At least one of intent and action is not empty
	Intent       string `json:"intent,omitempty"`
	Action       string `json:"action,omitempty"`
	Url          string `json:"url,omitempty"`
	RichResource string `json:"richresource,omitempty"`
}

type BadgeNotification struct {
	AddNum int    `json:"add_num,omitempty"`
	SetNum int    `json:"set_num,omitempty"`
	Class  string `json:"class,omitempty"`
}

type LightSettings struct {
	Color            *Color `json:"color"`
	LightOnDuration  string `json:"light_on_duration,omitempty"`
	LightOffDuration string `json:"light_off_duration,omitempty"`
}

type Color struct {
	Alpha int `json:"alpha"`
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}
