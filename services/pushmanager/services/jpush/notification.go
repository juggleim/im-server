package jpush

type Notification struct {
	Alert    string                `json:"alert,omitempty"`
	Android  *AndroidNotification  `json:"android,omitempty"`
	Ios      *IosNotification      `json:"ios,omitempty"`
	WinPhone *WinPhoneNotification `json:"winphone,omitempty"`
}

type AndroidNotification struct {
	Alert      string                 `json:"alert"`
	Title      string                 `json:"title,omitempty"`
	BuilderId  int                    `json:"builder_id,omitempty"`
	Priority   int                    `json:"priority,omitempty"`
	Category   string                 `json:"category,omitempty"`
	Style      int                    `json:"style,omitempty"`
	AlertType  int                    `json:"alert_type,omitempty"`
	BigText    string                 `json:"big_text,omitempty"`
	Inbox      map[string]interface{} `json:"inbox,omitempty"`
	BigPicPath string                 `json:"big_pic_path,omitempty"`
	Extras     map[string]interface{} `json:"extras,omitempty"`
}

type IosNotification struct {
	Alert            interface{}            `json:"alert"`
	Sound            string                 `json:"sound,omitempty"`
	Badge            interface{}            `json:"badge,omitempty"`
	ContentAvailable bool                   `json:"content-available,omitempty"`
	MutableContent   bool                   `json:"mutable-content,omitempty"`
	Category         string                 `json:"category,omitempty"`
	Extras           map[string]interface{} `json:"extras,omitempty"`
}

type WinPhoneNotification struct {
	Alert    string                 `json:"alert"`
	Title    string                 `json:"title,omitempty"`
	OpenPage string                 `json:"_open_page,omitempty"`
	Extras   map[string]interface{} `json:"extras,omitempty"`
}
