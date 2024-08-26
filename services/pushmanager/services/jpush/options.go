package jpush

type Options struct {
	SendNo          int    `json:"sendno,omitempty"`
	TimeLive        int    `json:"time_to_live,omitempty"`
	OverrideMsgId   int64  `json:"override_msg_id,omitempty"`
	ApnsProduction  bool   `json:"apns_production"`
	ApnsCollapseId  string `json:"apns_collapse_id,omitempty"`
	BigPushDuration int    `json:"big_push_duration,omitempty"`
}
