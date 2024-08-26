package vivopush

// PUSH-UPS-API接口文档
// https://dev.vivo.com.cn/documentCenter/doc/362
const (
	Host = "https://api-push.vivo.com.cn"

	AuthURL = "/message/auth" // 推送鉴权接口
	SendURL = "/message/send" // 单推接口
)

type AuthReq struct {
	AppId     string `json:"appId,omitempty"`     // 用户申请推送业务时生成的appId
	AppKey    string `json:"appKey,omitempty"`    // 用户申请推送业务时获得的appKey
	Timestamp int64  `json:"timestamp,omitempty"` // Unix时间戳 做签名用，单位：毫秒，且在vivo服务器当前utc时间戳前后十分钟区间内。
	Sign      string `json:"sign,omitempty"`      // 签名 使用MD5算法，字符串trim后拼接（appId+appKey+timestamp+appSecret），然后通过MD5加密得到的值（字母小写）
}

type AuthRes struct {
	Result    int    `json:"result"`    // 接口调用是否成功的状态码 0成功，非0失败
	Desc      string `json:"desc"`      // 文字描述接口调用情况
	AuthToken string `json:"authToken"` // 当鉴权成功时才会有该字段，推送消息时，需要提供authToken，有效期默认为1天，过期后无法使用。一个appId可对应多个token，24小时过期，业务方做中心缓存，1-2小时更新一次。
}

type SendReq struct {
	AppId           int                      `json:"appId,omitempty"`           // 用户申请推送业务时生成的appId，用于与获取authToken时传递的appId校验，一致才可以推送
	RegId           string                   `json:"regId,omitempty"`           // 应用订阅PUSH服务器得到的id   长度23个字符（regId，alias 两者需一个不为空，当两个不为空时，取regId）
	Alias           string                   `json:"alias,omitempty"`           // 别名 长度不超过40字符（regId，alias两者需一个不为空，当两个不为空时，取regId）
	NotifyType      int                      `json:"notifyType,omitempty"`      // 通知类型 1:无，2:响铃，3:振动，4:响铃和振动
	Title           string                   `json:"title,omitempty"`           // 通知标题（用于通知栏消息） 最大20个汉字（一个汉字等于两个英文字符，即最大不超过40个英文字符）
	Content         string                   `json:"content,omitempty"`         // 通知内容（用于通知栏消息） 最大50个汉字（一个汉字等于两个英文字符，即最大不超过100个英文字符）
	TimeToLive      int64                    `json:"timeToLive,omitempty"`      // 消息保留时长 单位：秒，取值至少60秒，最长7天。当值为空时，默认一天
	SkipType        int                      `json:"skipType,omitempty"`        // 点击跳转类型 1：打开APP首页 2：打开链接 3：自定义 4:打开app内指定页面
	SkipContent     string                   `json:"skipContent,omitempty"`     // 跳转内容 跳转类型为2时，跳转内容最大1000个字符，跳转类型为3或4时，跳转内容最大1024个字符，skipType传3需要在onNotificationMessageClicked回调函数中自己写处理逻辑。关于skipContent的内容可以参考【vivo推送常见问题汇总】 pushSDK版本号：480以上，不在支持skipType=3，自定义跳转统一使用skipType=4，详见【vivo推送常见问题汇总】中API接入问题的Q11中的intent uri示例。
	NetworkType     int                      `json:"networkType,omitempty"`     // 网络方式 -1：不限，1：wifi下发送，不填默认为-1
	Classification  int                      `json:"classification,omitempty"`  // 消息类型 0：运营类消息，1：系统类消息。不填默认为0
	Category        string                   `json:"category,omitempty"`        // 二级分类消息类型 https://dev.vivo.com.cn/documentCenter/doc/359
	ClientCustomMap map[string]string        `json:"clientCustomMap,omitempty"` // 客户端自定义键值对 自定义key和Value键值对个数不能超过10个，且长度不能超过1024字符, key和Value键值对总长度不能超过1024字符。app可以按照客户端SDK接入文档获取该键值对
	Extra           map[string]string        `json:"extra,omitempty"`           // 高级特性（详见目录：一.公共——5.高级特性 extra）
	RequestId       string                   `json:"requestId,omitempty"`       // 用户请求唯一标识 最大64字符
	PushMode        int                      `json:"pushMode,omitempty"`        // 推送模式 0：正式推送；1：测试推送，不填默认为0 备注： 1.测试推送，只能给web界面录入的测试用户推送；审核中应用，只能用测试推送 2.若未设置pushMode=1进行测试，文案相同时，将被当做重复推送的运营消息被去重
	AuditReview     []map[string]interface{} `json:"auditReview,omitempty"`     // 第三方审核结果，参见：基于第三方审核结果的消息推送
	NotifyId        int                      `json:"notifyId,omitempty"`        // 每条消息在通知显示时的唯一标识。不携带时，vpush自动为给每条消息生成一个唯一标识；
}

type SendRes struct {
	Result      int    `json:"result"` // 接口调用是否成功的状态码 0成功，非0失败
	Desc        string `json:"desc"`   // 文字描述接口调用情况
	TaskId      string `json:"taskId"` // 任务编号
	InvalidUser *struct {
		UserId int `json:"userid"` // userid为接入方传的regid或者alias
		Status int `json:"status"` // status有三种情况：1.userId不存在；2.卸载或者关闭了通知；3.90天不在线；4.非测试用户
	} `json:"invalidUser"` // 非法用户信息，包括status和userid
}
