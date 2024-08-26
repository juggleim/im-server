package xiaomipush

const (
	Host       = "https://api.xmpush.xiaomi.com"        // 国内
	HostGlobal = "https://api.xmpush.global.xiaomi.com" // 海外

	HostRu   = "https://ru-api.xmpush.global.xiaomi.com"   // 机房：莫斯科，服务范围：俄罗斯
	HostIdmb = "https://idmb-api.xmpush.global.xiaomi.com" // 机房：孟买，服务范围：印度
	HostFr   = "https://fr-api.xmpush.global.xiaomi.com"   // 机房：法兰克福，服务范围：欧洲经济区以及英国
	HostSgp  = "https://sgp-api.xmpush.global.xiaomi.com"  // 机房：新加坡，服务范围：以上各地区之外的其他地区

	SendURL = "/v3/message/regid"
)

type SendReq struct {
	Payload               string      `json:"payload,omitempty"`                 // 消息的内容。（注意：需要对payload字符串做urlencode处理）
	RestrictedPackageName string      `json:"restricted_package_name,omitempty"` // App的包名。备注：V2版本支持一个包名，V3版本支持多包名（中间用逗号分割）。
	PassThrough           int32       `json:"pass_through,omitempty"`            // pass_through的值可以为： 0 表示通知栏消息 1 表示透传消息
	Title                 string      `json:"title,omitempty"`                   // 通知栏展示的通知的标题，不允许全是空白字符，长度小于50， 一个中英文字符均计算为1（通知栏消息必填）。
	Description           string      `json:"description,omitempty"`             // 通知栏展示的通知的描述，不允许全是空白字符，长度小于128，一个中英文字符均计算为1（通知栏消息必填）。
	NotifyType            int32       `json:"notify_type,omitempty"`             // notify_type的值可以是DEFAULT_ALL或者以下其他几种的OR组合： DEFAULT_ALL = -1; DEFAULT_SOUND  = 1;  // 使用默认提示音提示； DEFAULT_VIBRATE = 2;  // 使用默认震动提示； DEFAULT_LIGHTS = 4;   // 使用默认led灯光提示；
	TimeToLive            int64       `json:"time_to_live,omitempty"`            // 可选项。如果用户离线，设置消息在服务器保存的时间，单位：ms。服务器默认最长保留两周。
	TimeToSend            int64       `json:"time_to_send,omitempty"`            // 可选项。定时发送消息。用自1970年1月1日以来00:00:00.0 UTC时间表示（以毫秒为单位的时间）。注：仅支持七天内的定时消息。
	NotifyID              int64       `json:"notify_id,omitempty"`               // 可选项。默认情况下，通知栏只显示一条推送消息。如果通知栏要显示多条推送消息，需要针对不同的消息设置不同的notify_id（相同notify_id的通知栏消息会覆盖之前的）。
	Extra                 *Extra      `json:"extra,omitempty"`                   // 可选项，对app提供一些扩展的功能，请参考2.2。除了这些扩展功能，开发者还可以通过 ExtraCustom 定义一些key和value来控制客户端的行为。注：key和value的字符数不能超过1024，至多可以设置10个key-value键值对。
	ExtraCustom           interface{} `json:"extra,omitempty"`                   // 可选项，对app提供一些扩展的功能，请参考2.2。开发者自己定义的一些key和value来控制客户端的行为。传入类型应该为结构体指针，参考 *Extra。

	RegistrationId string `json:"registration_id,omitempty"` // 根据registration_id，发送消息到指定设备上。可以提供多个registration_id，发送给一组设备，不同的registration_id之间用“,”分割。参数仅适用于“/message/regid”HTTP API。（注意：需要对payload字符串做urlencode处理）
	Alias          string `json:"alias,omitempty"`           // 根据alias，发送消息到指定的设备。可以提供多个alias，发送给一组设备，不同的alias之间用“,”分割。参数仅适用于“/message/alias”HTTP API。
	UserAccount    string `json:"user_account,omitempty"`    // 根据user_account，发送消息给设置了该user_account的所有设备。可以提供多个user_account，user_account之间用“,”分割。参数仅适用于“/message/user_account”HTTP API。
	Topic          string `json:"topic,omitempty"`           // 根据topic，发送消息给订阅了该topic的所有设备。参数仅适用于“/message/topic”HTTP API。
	Topics         string `json:"topics,omitempty"`          // topic列表，使用;$;分割。注: topics参数需要和topic_op参数配合使用，另外topic的数量不能超过5。参数仅适用于“/message/multi_topic”HTTP API。
	TopicOp        string `json:"topic_op,omitempty"`        // topic之间的操作关系。支持以下三种： UNION并集 INTERSECTION交集 EXCEPT差集 例如： topics的列表元素是[A, B, C, D]，则并集结果是A∪B∪C∪D，交集的结果是A ∩B ∩C ∩D，差集的结果是A-B-C-D。 参数仅适用于“/message/multi_topic”HTTP API。
}

type Extra struct {
	SoundUri           string    `json:"sound_uri,omitempty"`           // 可选项，自定义通知栏消息铃声。extra.sound_uri的值设置为铃声的URI。（请参见《服务端Java SDK文档》中的“自定义铃声”一节。） 注：铃声文件放在Android app的raw目录下。
	Ticker             string    `json:"ticker,omitempty"`              // 可选项，开启通知消息在状态栏滚动显示。
	NotifyForeground   string    `json:"notify_foreground,omitempty"`   // 可选项，开启/关闭app在前台时的通知弹出。当extra.notify_foreground值为”1″时，app会弹出通知栏消息；当extra.notify_foreground值为”0″时，app不会弹出通知栏消息。注：默认情况下会弹出通知栏消息。 （请参见《服务端Java SDK文档》中的“控制App前台通知弹出”一节。）
	NotifyEffect       string    `json:"notify_effect,omitempty"`       // 可选项，预定义通知栏消息的点击行为。通过设置extra.notify_effect的值以得到不同的预定义点击行为。 “1″：通知栏点击后打开app的Launcher Activity。 “2″：通知栏点击后打开app的任一Activity（开发者还需要传入extra.intent_uri）。 “3″：通知栏点击后打开网页（开发者还需要传入extra.web_uri）。 （请参见《服务端Java SDK文档》中的“预定义通知栏通知的点击行为”一节。）
	IntentUri          string    `json:"intent_uri,omitempty"`          // 可选项，打开当前app的任一组件。
	WebUri             string    `json:"web_uri,omitempty"`             // 可选项，打开某一个网页。
	FlowControl        string    `json:"flow_control,omitempty"`        // 可选项，控制是否需要进行平缓发送。默认不支持。value代表平滑推送的速度。注：服务端支持最低3000/s的qps。也就是说，如果将平滑推送设置为1000，那么真实的推送速度是3000/s。
	Jobkey             string    `json:"jobkey,omitempty"`              // 可选项，使用推送批次（JobKey）功能聚合消息。客户端会按照jobkey去重，即相同jobkey的消息只展示第一条，其他的消息会被忽略。合法的jobkey由数字（[0-9]），大小写字母（[a-zA-Z]），下划线（_）和中划线（-）组成，长度不大于12个字符，且不能以下划线(_)开头。
	Callback           *Callback `json:"callback,omitempty"`            // 可选项，开启消息回执。消息发送后，推送系统能发送回执给开发者，告知开发者这些消息的送达、点击或发送失败状态。 将extra.callback的值设置为第三方接收回执的http接口。接口协议请参见《服务端Java SDK文档》中的“消息回执”一节。
	Locale             string    `json:"locale,omitempty"`              // 可选项，可以接收消息的设备的语言范围，用逗号分隔。当前设备的local的获取方法： Locale.getDefault().toString() 比如，中国大陆用zh_CN表示。
	LocaleNotIn        string    `json:"locale_not_in,omitempty"`       // 可选项，无法收到消息的设备的语言范围，逗号分隔。
	Model              string    `json:"model,omitempty"`               // 可选项，model支持三种用法。 可以收到消息的设备的机型范围，逗号分隔。当前设备的model的获取方法： Build.MODEL 比如，小米手机4移动版用”MI 4LTE”表示。 可以收到消息的设备的品牌范围，逗号分割。比如，三星手机用”samsung”表示。 （目前覆盖42个主流品牌，对应关系见附录”品牌表”） 可以收到消息的设备的价格范围，逗号分隔。比如，0-999价格的设备用”0-999″表示。 （目前覆盖4个价格区间，对应关系见附录”价格表”）
	ModelNotIn         string    `json:"model_not_in,omitempty"`        // 可选项，无法收到消息的设备的机型范围，逗号分隔。
	AppVersion         string    `json:"app_version,omitempty"`         // 可以接收消息的app版本号，用逗号分割。安卓app版本号来源于manifest文件中的”android:versionName”的值。注：目前支持MiPush_SDK_Client_2_2_12_sdk.jar（及以后）的版本。
	AppVersionNotIn    string    `json:"app_version_not_in,omitempty"`  // 无法接收消息的app版本号，用逗号分割。
	Connpt             string    `json:"connpt,omitempty"`              // 可选项，指定在特定的网络环境下才能接收到消息。目前仅支持指定”wifi”。
	OnlySendOnce       string    `json:"only_send_once,omitempty"`      // 可选项，extra.only_send_once的值设置为1，表示该消息仅在设备在线时发送一次，不缓存离线消息进行多次下发。
	ChannelId          string    `json:"channel_id,omitempty"`          // 必填，通知类别的ID。
	ChannelName        string    `json:"channel_name,omitempty"`        // 可选，通知类别的名称。
	ChannelDescription string    `json:"channel_description,omitempty"` // 可选，通知类别的描述。
}

type Callback struct {
	Param string `json:"param,omitempty"` // 可选项，第三方自定义回执参数，最大长度64字节。
	Type  string `json:"type,omitempty"`  // 可选项，表示回执类型。详细用法请参见《服务端Java SDK文档》中“消息回执”一节中的callback.type字段。
}

type Result struct {
	Result      string `json:"result,omitempty"`
	TraceId     string `json:"trace_id,omitempty"`
	Code        int64  `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Info        string `json:"info,omitempty"`
	Reason      string `json:"reason,omitempty"`
}
type SendRes struct {
	Result
	Data struct {
		ID                 string `json:"id,omitempty"`
		DayQuota           string `json:"day_quota,omitempty"`
		DayAcked           string `json:"day_acked,omitempty"`
		ChannelExceedQuota string `json:"channel_exceed_quota,omitempty"`
	} `json:"data,omitempty"`
}
