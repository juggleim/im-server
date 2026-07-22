package jpush

import "im-server/services/commonservices"

type Options struct {
	SendNo                   int                `json:"sendno,omitempty"`
	TimeLive                 int                `json:"time_to_live,omitempty"`
	OverrideMsgId            int64              `json:"override_msg_id,omitempty"`
	ApnsProduction           bool               `json:"apns_production"`
	ApnsCollapseId           string             `json:"apns_collapse_id,omitempty"`
	BigPushDuration          int                `json:"big_push_duration,omitempty"`
	ThirdPartyChannel        *ThirdPartyChannel `json:"third_party_channel,omitempty"`
	Classification           int                `json:"classification,omitempty"`
	TargetEvent              []string           `json:"target_event,omitempty"`
	TestMessage              *bool              `json:"test_message,omitempty"`
	ReceiptId                string             `json:"receipt_id,omitempty"`
	ActivePush               *bool              `json:"active_push,omitempty"`
	NeedBackup               *bool              `json:"need_backup,omitempty"`
	BusinessOperationCode    string             `json:"business_operation_code,omitempty"`
	TestModel                *bool              `json:"test_model,omitempty"`
	Notification3rdVer       string             `json:"notification_3rd_ver,omitempty"`
	AutoTruncation           *bool              `json:"auto_truncation,omitempty"`
	MktEnable                *bool              `json:"mkt_enable,omitempty"`
	NotificationSwitchFilter *bool              `json:"notification_switch_filter,omitempty"`
}

type ThirdPartyChannel struct {
	Huawei *commonservices.JPushHuaweiChannel `json:"huawei,omitempty"`
	Xiaomi *commonservices.JPushXiaomiChannel `json:"xiaomi,omitempty"`
	Honor  *commonservices.JPushHonorChannel  `json:"honor,omitempty"`
	Oppo   *commonservices.JPushOppoChannel   `json:"oppo,omitempty"`
	Vivo   *commonservices.JPushVivoChannel   `json:"vivo,omitempty"`
	Meizu  *commonservices.JPushMeizuChannel  `json:"meizu,omitempty"`
	Fcm    *commonservices.JPushFcmChannel    `json:"fcm,omitempty"`
	Nio    *commonservices.JPushNioChannel    `json:"nio,omitempty"`
	Asus   *commonservices.JPushAsusChannel   `json:"asus,omitempty"`
	Hmos   *commonservices.JPushHmosChannel   `json:"hmos,omitempty"`
}
