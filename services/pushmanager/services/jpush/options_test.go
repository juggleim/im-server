package jpush

import (
	"encoding/json"
	"im-server/services/commonservices"
	"reflect"
	"testing"
)

func TestOptionsMarshalAdditionalFields(t *testing.T) {
	trueValue := true
	falseValue := false
	options := Options{
		TargetEvent:              []string{"jg_app_show"},
		TestMessage:              &trueValue,
		ReceiptId:                "receipt-id",
		ActivePush:               &trueValue,
		NeedBackup:               &trueValue,
		BusinessOperationCode:    "operation-code",
		TestModel:                &trueValue,
		Notification3rdVer:       "v2",
		AutoTruncation:           &falseValue,
		MktEnable:                &trueValue,
		NotificationSwitchFilter: &trueValue,
	}

	data, err := json.Marshal(options)
	if err != nil {
		t.Fatalf("marshal options: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal options: %v", err)
	}

	want := map[string]interface{}{
		"target_event":               []interface{}{"jg_app_show"},
		"test_message":               true,
		"receipt_id":                 "receipt-id",
		"active_push":                true,
		"need_backup":                true,
		"business_operation_code":    "operation-code",
		"test_model":                 true,
		"notification_3rd_ver":       "v2",
		"auto_truncation":            false,
		"mkt_enable":                 true,
		"notification_switch_filter": true,
	}
	for key, wantValue := range want {
		gotValue, ok := payload[key]
		if !ok {
			t.Errorf("missing JSON field %q", key)
			continue
		}
		if !reflect.DeepEqual(gotValue, wantValue) {
			t.Errorf("JSON field %q = %#v, want %#v", key, gotValue, wantValue)
		}
	}
}

func TestThirdPartyChannelMarshalVendorFields(t *testing.T) {
	falseValue := false
	zero := 0
	one := 1
	channels := ThirdPartyChannel{
		Huawei: &commonservices.JPushHuaweiChannel{
			Distribution:          "secondary_push",
			DistributionFcm:       "jpush",
			DistributionCustomize: "first_ospush",
			ChannelId:             "channel-id",
			Importance:            "NORMAL",
			Category:              "IM",
			Sound:                 "/raw/shake",
			DefaultSound:          &falseValue,
			Urgency:               "NORMAL",
			ReceiptId:             "receipt-id",
			TargetUserType:        &zero,
			LargeIcon:             "large-icon",
			SmallIconUri:          "small-icon",
			Style:                 2,
			BigText:               "big-text",
			Inbox:                 map[string]interface{}{"line": "content"},
			OnlyUseVendorStyle:    &falseValue,
			AuditResponse:         map[string]interface{}{"code": 200},
			HwPushType:            7,
			HwLivePayload:         map[string]interface{}{"activityId": 15},
		},
		Xiaomi: &commonservices.JPushXiaomiChannel{
			Distribution:       "jpush",
			DistributionFcm:    "fcm",
			ChannelId:          "channel-id",
			SkipQuota:          &falseValue,
			SmallIconColor:     "#ffffff",
			Style:              1,
			BigText:            "big-text",
			OnlyUseVendorStyle: &falseValue,
			MiTemplateId:       "template-id",
			MiTemplateParam:    `{"name":"value"}`,
			MiPushType:         3,
			VoipExtraData:      "voip-data",
		},
		Honor: &commonservices.JPushHonorChannel{
			Distribution:          "secondary_push",
			DistributionFcm:       "jpush",
			DistributionCustomize: "first_ospush",
			Importance:            "NORMAL",
			TargetUserType:        &zero,
			LargeIcon:             "large-icon",
			SmallIconUri:          "small-icon",
			Style:                 1,
			BigText:               "big-text",
			OnlyUseVendorStyle:    &falseValue,
			HonorPushType:         3,
			VoipExtraData:         "voip-data",
		},
		Oppo: &commonservices.JPushOppoChannel{
			Distribution:             "ospush",
			DistributionFcm:          "secondary_fcm_push",
			BadgeOperationType:       &one,
			ChannelId:                "channel-id",
			Category:                 "IM",
			NotifyLevel:              16,
			PrivateContentParameters: map[string]string{"name": "value"},
			PrivateMsgTemplateId:     "template-id",
			PrivateTitleParameters:   map[string]string{"title": "value"},
			SkipQuota:                &falseValue,
			LargeIcon:                "large-icon",
			Style:                    3,
			BigText:                  "big-text",
			BigPicPath:               "big-picture",
			OnlyUseVendorStyle:       &falseValue,
			AuditResponse:            map[string]interface{}{"code": 200},
			BadgeMessageCount:        &zero,
			OpPushType:               7,
			OpIntelligentIntent:      map[string]interface{}{"id": "intent-id"},
			OpDeleteIntentData:       map[string]interface{}{"id": "intent-id"},
			VoipExtraData:            "voip-data",
		},
		Vivo: &commonservices.JPushVivoChannel{
			Distribution:        "jpush",
			DistributionFcm:     "secondary_pns_push",
			Classification:      &zero,
			PushMode:            &zero,
			Category:            "IM",
			CallbackId:          "callback-id",
			AuditResponse:       map[string]interface{}{"code": 200},
			AddBadge:            true,
			VivoPushType:        3,
			VivoInappMsg:        map[string]interface{}{"content": "message"},
			VoipExtraData:       "voip-data",
			ExtensionExpireShow: &falseValue,
		},
		Meizu: &commonservices.JPushMeizuChannel{
			Distribution:    "jpush",
			DistributionFcm: "pns",
		},
		Fcm: &commonservices.JPushFcmChannel{
			Distribution:    "jpush",
			DistributionFcm: "fcm",
		},
		Nio: &commonservices.JPushNioChannel{
			Distribution:    "jpush",
			DistributionFcm: "pns",
			ChannelId:       "channel-id",
		},
		Asus: &commonservices.JPushAsusChannel{
			Distribution:    "jpush",
			DistributionFcm: "fcm",
		},
		Hmos: &commonservices.JPushHmosChannel{
			Distribution: "jpush",
		},
	}

	data, err := json.Marshal(channels)
	if err != nil {
		t.Fatalf("marshal third-party channels: %v", err)
	}

	var payload map[string]map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal third-party channels: %v", err)
	}

	wantFields := map[string][]string{
		"huawei": {
			"distribution", "distribution_fcm", "distribution_customize", "channel_id",
			"importance", "category", "sound", "default_sound", "urgency", "receipt_id",
			"target_user_type", "large_icon", "small_icon_uri", "style", "big_text", "inbox",
			"only_use_vendor_style", "auditResponse", "hw_push_type", "hw_live_payload",
		},
		"xiaomi": {
			"distribution", "distribution_fcm", "channel_id", "skip_quota", "small_icon_color",
			"style", "big_text", "only_use_vendor_style", "mi_template_id", "mi_template_param",
			"mi_push_type", "voip_extraData",
		},
		"honor": {
			"distribution", "distribution_fcm", "distribution_customize", "importance",
			"target_user_type", "large_icon", "small_icon_uri", "style", "big_text",
			"only_use_vendor_style", "honor_push_type", "voip_extraData",
		},
		"oppo": {
			"distribution", "distribution_fcm", "badge_operation_type", "channel_id", "category",
			"notify_level", "private_content_parameters", "private_msg_template_id",
			"private_title_parameters", "skip_quota", "large_icon", "style", "big_text",
			"big_pic_path", "only_use_vendor_style", "auditResponse", "badge_message_count",
			"op_push_type", "op_intelligent_intent", "op_delete_intent_data", "voip_extraData",
		},
		"vivo": {
			"distribution", "distribution_fcm", "classification", "push_mode", "category",
			"callback_id", "auditResponse", "add_badge", "vivo_push_type", "vivo_inapp_msg",
			"voip_extraData", "extensionExpireShow",
		},
		"meizu": {"distribution", "distribution_fcm"},
		"fcm":   {"distribution", "distribution_fcm"},
		"nio":   {"distribution", "distribution_fcm", "channel_id"},
		"asus":  {"distribution", "distribution_fcm"},
		"hmos":  {"distribution"},
	}
	for vendor, fields := range wantFields {
		vendorPayload, ok := payload[vendor]
		if !ok {
			t.Errorf("missing vendor %q", vendor)
			continue
		}
		for _, field := range fields {
			if _, ok := vendorPayload[field]; !ok {
				t.Errorf("vendor %q missing JSON field %q", vendor, field)
			}
		}
	}
}
