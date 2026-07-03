package getuipush

// AuthReq 鉴权请求
// sign = sha256(appkey + timestamp + mastersecret)
type AuthReq struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	AppKey    string `json:"appkey"`
}

type AuthResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ExpireTime string `json:"expire_time"`
		Token      string `json:"token"`
	} `json:"data"`
}

// ToSingleCIDReq 执行 cid 单推请求体
type ToSingleCIDReq struct {
	RequestID   string       `json:"request_id"`
	GroupName   string       `json:"group_name,omitempty"`
	Settings    *Settings    `json:"settings,omitempty"`
	Audience    *AudienceCID `json:"audience"`
	PushMessage *PushMessage `json:"push_message"`
	PushChannel any          `json:"push_channel,omitempty"`
}

type Settings struct {
	TTL      int64 `json:"ttl,omitempty"` // 毫秒
	Strategy any   `json:"strategy,omitempty"`
}

type AudienceCID struct {
	CID []string `json:"cid"`
}

type PushMessage struct {
	Notification *Notification `json:"notification,omitempty"`
	Transmission string        `json:"transmission,omitempty"`
}

type Notification struct {
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
	ClickType string `json:"click_type,omitempty"`
	URL       string `json:"url,omitempty"`
}

type ToSingleCIDResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	// 典型返回：{"taskid":{"cid":"status"}}
	Data map[string]map[string]string `json:"data,omitempty"`
}
