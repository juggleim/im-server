package oppopush

const (
	Host = "https://api.push.oppomobile.com"

	AuthURL = "/server/v1/auth"                         // 鉴权
	SendURL = "/server/v1/message/notification/unicast" // 单推-通知栏消息推送
)

// SendReq 单推-通知栏消息推送
// https://open.oppomobile.com/wiki/doc#id=10692
type SendReq struct {
	TargetType           int           `json:"target_type,omitempty"`            // 目标类型 2: registration_id  5:别名
	TargetValue          string        `json:"target_value,omitempty"`           // 推送目标用户: registration_id或alias
	Notification         *Notification `json:"notification,omitempty"`           // 请参见通知栏消息
	VerifyRegistrationId bool          `json:"verify_registration_id,omitempty"` // 消息到达客户端后是否校验registration_id。 true表示推送目标与客户端registration_id进行比较，如果一致则继续展示，不一致则就丢弃；false表示不校验
}

// 问：OPPO Push 推送消息是否可以提供声音、震动等提醒选项设置？
// 答：消息推送时目前没有提供提醒方式的选择，所以通知都是使用系统默认的提醒方式。有
// 用户有需要，可在“设置”-“通知与状态栏”-“通知管理”中，找到应用并开启对应的消息提醒方式。
// https://pfs.oppomobile.com/static/document/oppo_push_faq.pdf

// Notification 通知栏消息
// https://open.oppomobile.com/wiki/doc#id=10688
type Notification struct {
	AppMessageID        string `json:"app_message_id,omitempty"`        // App开发者自定义消息Id，OPPO推送平台根据此ID做去重处理，对于广播推送相同app_message_id只会保存一次，对于单推相同app_message_id只会推送一次。
	Style               int    `json:"style,omitempty"`                 // 通知栏样式 1. 标准样式 2. 长文本样式（ColorOS版本>5.0可用，通知栏第一条消息可展示全部内容，非第一条消息只展示一行内容） 3. 大图样式（ColorOS版本>5.0可用，通知栏第一条消息展示大图，非第一条消息不显示大图，推送方式仅支持广播，且不支持定速功能）
	BigPictureId        string `json:"big_picture_id,omitempty"`        // 大图id【style为3时，必填】,可在大图上传接口获取
	SmallPictureId      string `json:"small_picture_id,omitempty"`      // 通知图标id,可在图标上传接口获取
	Title               string `json:"title,omitempty"`                 // 设置在通知栏展示的通知栏标题, 【字数限制1~50，中英文均以一个计算】
	SubTitle            string `json:"sub_title,omitempty"`             // 子标题，设置在通知栏展示的通知栏标题, 【字数限制1~10，中英文均以一个计算】
	Content             string `json:"content,omitempty"`               // 设置在通知栏展示的通知的内容,必填 1）标准样式（style 为 1）：字数限制200以内（兼容API文档以前定义，实际手机端通知栏消息只能展示50字数） 2）长文本样式（style 为 2）限制128个以内 3）大图样式（style 为 3）字数限制50以内，中英文均以一个计算】
	ClickActionType     int    `json:"click_action_type,omitempty"`     // 点击动作类型： 0，启动应用； 1，打开应用内页（activity的intent action）； 2，打开网页； 4，打开应用内页（activity）； 【非必填，默认值为0】; 5,Intent scheme URL
	ClickActionActivity string `json:"click_action_activity,omitempty"` // 应用内页地址【click_action_type为1/4/时必填，长度500】
	ClickActionURL      string `json:"click_action_url,omitempty"`      // 网页地址或【click_action_type为2与5时必填，长度500】示例： click_action_type为2时http://oppo.com?key1=val1&key2=val2
	ActionParameters    string `json:"action_parameters,omitempty"`     // 动作参数，打开应用内页或网页时传递给应用或网页【JSON格式，非必填】，字符数不能超过4K 示例： {"key1":"value1","key2":"value2"}
	ShowTimeType        int    `json:"show_time_type,omitempty"`        // 展示类型 (0, “即时”),(1, “定时”)
	ShowStartTime       int64  `json:"show_start_time,omitempty"`       // 定时展示开始时间（根据time_zone转换成当地时间），时间的毫秒数
	ShowEndTime         int64  `json:"show_end_time,omitempty"`         // 定时展示结束时间（根据time_zone转换成当地时间），时间的毫秒数
	OffLine             bool   `json:"off_line,omitempty"`              // 是否进离线消息,【非必填，默认为True】
	OffLineTTL          int    `json:"off_line_ttl,omitempty"`          // 离线消息的存活时间(time_to_live) (单位：秒), 【最长10天】
	PushTimeType        int    `json:"push_time_type,omitempty"`        // 定时推送 (0, “即时”),(1, “定时”), 【只对全部用户推送生效】
	PushStartTime       int64  `json:"push_start_time,omitempty"`       // 定时推送开始时间（根据time_zone转换成当地时间）, 【push_time_type 为1必填】，时间的毫秒数
	TimeZone            string `json:"time_zone,omitempty"`             // 时区，默认值：（GMT+08:00）北京，香港，新加坡
	FixSpeed            bool   `json:"fix_speed,omitempty"`             // 是否定速推送,【非必填，默认值为false】
	FixSpeedRate        int64  `json:"fix_speed_rate,omitempty"`        // 定速速率 【fixSpeed为true时，必填】
	NetworkType         int    `json:"network_type,omitempty"`          // 0：不限联网方式, 1：仅wifi推送；
	CallBackURL         string `json:"call_back_url,omitempty"`         // 仅支持registrationId推送方式 应用接收消息到达回执的回调URL，字数限制200以内，中英文均以一个计算。 OPPO Push服务器POST一个JSON数据到call_back_url； Content-Type为application/json的方式提交数据。
	CallBackParameter   string `json:"call_back_parameter,omitempty"`   // App开发者自定义回执参数，字数限制100以内，中英文均以一个计算。
	ChannelID           string `json:"channel_id,omitempty"`            // 通知栏通道（NotificationChannel），从Android9开始发送通知消息必须要指定通道Id（如果是快应用，必须带置顶的通道Id:OPPO PUSH推送）
	ShowTtl             int    `json:"show_ttl,omitempty"`              // 限时展示(单位：秒)，消息在通知栏展示后开始计时，到达填写的相对应时间后自动从通知栏消失，默认是1天。时间范围6 * 60 * 60 s -- 48 * 60 * 60 s
	NotifyId            int    `json:"notify_id,omitempty"`             // 每条消息在通知显示时的唯一标识。不携带时，PUSH自动为给每条消息生成一个唯一标识；不同的通知栏消息可以拥有相同的notifyId，实现新的消息覆盖上一条消息功能。
}

type SendRes struct {
	Code    int    `json:"code"`    // 返回码,请参考公共返回码与接口返回码
	Message string `json:"message"` // 错误详细信息，不存在则不填
	Data    struct {
		MessageID string `json:"messageId"` // 消息 ID
	} `json:"data"` // 返回值，JSON类型
}

type AuthReq struct {
	AppKey    string `json:"app_key,omitempty"`   // OPPO-OPEN 分配给应用的AppKey，内部开放API是PUSH分配给应用的AppKey
	Sign      string `json:"sign,omitempty"`      // sha256(appkey+timestamp+mastersecret) mastersecret为注册应用时生成
	Timestamp string `json:"timestamp,omitempty"` // 时间戳，时间毫秒数，时区为GMT+8。PUSH API服务端允许客户端请求最大时间误差为10分钟。
}

type AuthRes struct {
	Code    int    `json:"code"`    // 返回码,请参考公共返回码与接口返回码
	Message string `json:"message"` // 错误详细信息，不存在则不填
	Data    struct {
		AuthToken  string `json:"auth_token"`  // 权限令牌，推送消息时，需要提供auth_token，有效期默认为24小时，过期后无法使用
		CreateTime int64  `json:"create_time"` // 时间毫秒数
	} `json:"data"` // 返回值，JSON类型，包含响应结构体
}
