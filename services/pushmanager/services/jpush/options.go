package jpush

import "im-server/services/commonservices"

type Options struct {
	SendNo            int                `json:"sendno,omitempty"`
	TimeLive          int                `json:"time_to_live,omitempty"`
	OverrideMsgId     int64              `json:"override_msg_id,omitempty"`
	ApnsProduction    bool               `json:"apns_production"`
	ApnsCollapseId    string             `json:"apns_collapse_id,omitempty"`
	BigPushDuration   int                `json:"big_push_duration,omitempty"`
	ThirdPartyChannel *ThirdPartyChannel `json:"third_party_channel,omitempty"`
	Classification    int                `json:"classification,omitempty"`
}

type ThirdPartyChannel struct {
	Huawei *commonservices.JPushHuaweiChannel `json:"huawei,omitempty"`
	Xiaomi *commonservices.JPushXiaomiChannel `json:"xiaomi,omitempty"`
	Honor  *commonservices.JPushHonorChannel  `json:"honor,omitempty"`
	Oppo   *commonservices.JPushOppoChannel   `json:"oppo,omitempty"`
	Vivo   *commonservices.JPushVivoChannel   `json:"vivo,omitempty"`
	Meizu  *commonservices.JPushMeizuChannel  `json:"meizu,omitempty"`
}
