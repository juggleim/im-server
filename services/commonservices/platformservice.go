package commonservices

import "im-server/commons/pbdefines/pbobjs"

type Platform string
type PushChannel string

const (
	Platform_Android Platform = "Android"
	Platform_IOS     Platform = "iOS"
	Platform_Web     Platform = "Web"
	Platform_PC      Platform = "PC"
	Platform_Harmony Platform = "Harmony"
	Platform_Server  Platform = "Server"
	Platform_Bot     Platform = "Bot"

	PushChannel_Apple  PushChannel = "Apple"
	PushChannel_Huawei PushChannel = "Huawei"
	PushChannel_Xiaomi PushChannel = "Xiaomi"
	PushChannel_OPPO   PushChannel = "Oppo"
	PushChannel_VIVO   PushChannel = "Vivo"
	PushChannel_Jpush  PushChannel = "Jpush"
	PushChannel_FCM    PushChannel = "FCM"
	PushChannel_Meizu  PushChannel = "Meizu"
	PushChannel_HONOR  PushChannel = "Honor"
	PushChannel_Getui  PushChannel = "Getui"
)

type AndroidPushConf struct {
	AppKey      string         `json:"app_key"`
	PushChannel string         `json:"push_channel,omitempty"`
	Package     string         `json:"package"`
	Extra       map[string]any `json:"extra"`
}
type HuaweiPushConf struct {
	AppId      string `json:"app_id"`
	AppSecret  string `json:"app_secret"`
	BadgeClass string `json:"badge_class,omitempty"`
}

func (conf *HuaweiPushConf) Valid() bool {
	return conf.AppId != "" && conf.AppSecret != ""
}

type XiaomiPushConf struct {
	AppSecret string `json:"app_secret"`
	ChannelId string `json:"channel_id"`
}

func (conf *XiaomiPushConf) Valid() bool {
	return conf.AppSecret != ""
}

type OppoPushConf struct {
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
	ChannelId    string `json:"channel_id"`
}

func (conf *OppoPushConf) Valid() bool {
	return conf.AppKey != "" && conf.MasterSecret != ""
}

type VivoPushConf struct {
	AppId     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

func (conf *VivoPushConf) Valid() bool {
	return conf.AppId != "" && conf.AppKey != "" && conf.AppSecret != ""
}

type JPushConf struct {
	AppKey       string        `json:"app_key"`
	MasterSecret string        `json:"master_secret"`
	BadgeClass   string        `json:"badge_class,omitempty"`
	Options      *JPushOptions `json:"options,omitempty"`
}

type JPushOptions struct {
	Classification    int                     `json:"classification,omitempty"`
	ThirdPartyChannel *JPushThirdPartyChannel `json:"third_party_channel,omitempty"`
}

type JPushThirdPartyChannel struct {
	Huawei *JPushHuaweiChannel `json:"huawei,omitempty"`
	Xiaomi *JPushXiaomiChannel `json:"xiaomi,omitempty"`
	Honor  *JPushHonorChannel  `json:"honor,omitempty"`
	Oppo   *JPushOppoChannel   `json:"oppo,omitempty"`
	Vivo   *JPushVivoChannel   `json:"vivo,omitempty"`
	Meizu  *JPushMeizuChannel  `json:"meizu,omitempty"`
	Fcm    *JPushFcmChannel    `json:"fcm,omitempty"`
	Nio    *JPushNioChannel    `json:"nio,omitempty"`
	Asus   *JPushAsusChannel   `json:"asus,omitempty"`
	Hmos   *JPushHmosChannel   `json:"hmos,omitempty"`
}

type JPushHuaweiChannel struct {
	Distribution          string                 `json:"distribution,omitempty"`
	DistributionFcm       string                 `json:"distribution_fcm,omitempty"`
	DistributionCustomize string                 `json:"distribution_customize,omitempty"`
	ChannelId             string                 `json:"channel_id,omitempty"`
	Importance            string                 `json:"importance,omitempty"`
	Category              string                 `json:"category,omitempty"`
	Sound                 string                 `json:"sound,omitempty"`
	DefaultSound          *bool                  `json:"default_sound,omitempty"`
	Urgency               string                 `json:"urgency,omitempty"`
	ReceiptId             string                 `json:"receipt_id,omitempty"`
	TargetUserType        *int                   `json:"target_user_type,omitempty"`
	LargeIcon             string                 `json:"large_icon,omitempty"`
	SmallIconUri          string                 `json:"small_icon_uri,omitempty"`
	Style                 int                    `json:"style,omitempty"`
	BigText               string                 `json:"big_text,omitempty"`
	Inbox                 map[string]interface{} `json:"inbox,omitempty"`
	OnlyUseVendorStyle    *bool                  `json:"only_use_vendor_style,omitempty"`
	AuditResponse         map[string]interface{} `json:"auditResponse,omitempty"`
	HwPushType            int                    `json:"hw_push_type,omitempty"`
	HwLivePayload         map[string]interface{} `json:"hw_live_payload,omitempty"`
}

type JPushXiaomiChannel struct {
	Distribution       string `json:"distribution,omitempty"`
	DistributionFcm    string `json:"distribution_fcm,omitempty"`
	ChannelId          string `json:"channel_id,omitempty"`
	SkipQuota          *bool  `json:"skip_quota,omitempty"`
	SmallIconColor     string `json:"small_icon_color,omitempty"`
	Style              int    `json:"style,omitempty"`
	BigText            string `json:"big_text,omitempty"`
	OnlyUseVendorStyle *bool  `json:"only_use_vendor_style,omitempty"`
	MiTemplateId       string `json:"mi_template_id,omitempty"`
	MiTemplateParam    string `json:"mi_template_param,omitempty"`
	MiPushType         int    `json:"mi_push_type,omitempty"`
	VoipExtraData      string `json:"voip_extraData,omitempty"`
}

type JPushHonorChannel struct {
	Distribution          string `json:"distribution,omitempty"`
	DistributionFcm       string `json:"distribution_fcm,omitempty"`
	DistributionCustomize string `json:"distribution_customize,omitempty"`
	Importance            string `json:"importance,omitempty"`
	TargetUserType        *int   `json:"target_user_type,omitempty"`
	LargeIcon             string `json:"large_icon,omitempty"`
	SmallIconUri          string `json:"small_icon_uri,omitempty"`
	Style                 int    `json:"style,omitempty"`
	BigText               string `json:"big_text,omitempty"`
	OnlyUseVendorStyle    *bool  `json:"only_use_vendor_style,omitempty"`
	HonorPushType         int    `json:"honor_push_type,omitempty"`
	VoipExtraData         string `json:"voip_extraData,omitempty"`
}

type JPushOppoChannel struct {
	Distribution             string                 `json:"distribution,omitempty"`
	DistributionFcm          string                 `json:"distribution_fcm,omitempty"`
	BadgeOperationType       *int                   `json:"badge_operation_type,omitempty"`
	ChannelId                string                 `json:"channel_id,omitempty"`
	Category                 string                 `json:"category,omitempty"`
	NotifyLevel              int                    `json:"notify_level,omitempty"`
	PrivateContentParameters map[string]string      `json:"private_content_parameters,omitempty"`
	PrivateMsgTemplateId     string                 `json:"private_msg_template_id,omitempty"`
	PrivateTitleParameters   map[string]string      `json:"private_title_parameters,omitempty"`
	SkipQuota                *bool                  `json:"skip_quota,omitempty"`
	LargeIcon                string                 `json:"large_icon,omitempty"`
	Style                    int                    `json:"style,omitempty"`
	BigText                  string                 `json:"big_text,omitempty"`
	BigPicPath               string                 `json:"big_pic_path,omitempty"`
	OnlyUseVendorStyle       *bool                  `json:"only_use_vendor_style,omitempty"`
	AuditResponse            map[string]interface{} `json:"auditResponse,omitempty"`
	BadgeMessageCount        *int                   `json:"badge_message_count,omitempty"`
	OpPushType               int                    `json:"op_push_type,omitempty"`
	OpIntelligentIntent      map[string]interface{} `json:"op_intelligent_intent,omitempty"`
	OpDeleteIntentData       map[string]interface{} `json:"op_delete_intent_data,omitempty"`
	VoipExtraData            string                 `json:"voip_extraData,omitempty"`
}

type JPushVivoChannel struct {
	Distribution        string                 `json:"distribution,omitempty"`
	DistributionFcm     string                 `json:"distribution_fcm,omitempty"`
	Classification      *int                   `json:"classification,omitempty"`
	PushMode            *int                   `json:"push_mode,omitempty"`
	Category            string                 `json:"category,omitempty"`
	CallbackId          string                 `json:"callback_id,omitempty"`
	AuditResponse       map[string]interface{} `json:"auditResponse,omitempty"`
	AddBadge            bool                   `json:"add_badge,omitempty"`
	VivoPushType        int                    `json:"vivo_push_type,omitempty"`
	VivoInappMsg        map[string]interface{} `json:"vivo_inapp_msg,omitempty"`
	VoipExtraData       string                 `json:"voip_extraData,omitempty"`
	ExtensionExpireShow *bool                  `json:"extensionExpireShow,omitempty"`
}

type JPushMeizuChannel struct {
	Distribution    string `json:"distribution,omitempty"`
	DistributionFcm string `json:"distribution_fcm,omitempty"`
}

type JPushFcmChannel struct {
	Distribution    string `json:"distribution,omitempty"`
	DistributionFcm string `json:"distribution_fcm,omitempty"`
}

type JPushNioChannel struct {
	Distribution    string `json:"distribution,omitempty"`
	DistributionFcm string `json:"distribution_fcm,omitempty"`
	ChannelId       string `json:"channel_id,omitempty"`
}

type JPushAsusChannel struct {
	Distribution    string `json:"distribution,omitempty"`
	DistributionFcm string `json:"distribution_fcm,omitempty"`
}

type JPushHmosChannel struct {
	Distribution string `json:"distribution,omitempty"`
}

func (conf *JPushConf) Valid() bool {
	return conf.AppKey != "" && conf.MasterSecret != ""
}

type HonorPushConf struct {
	AppId      string `json:"app_id"`
	AppKey     string `json:"app_key"`
	AppSecret  string `json:"app_secret"`
	BadgeClass string `json:"badge_class,omitempty"`
}

func (conf *HonorPushConf) Valid() bool {
	return conf.AppId != "" && conf.AppKey != "" && conf.AppSecret != ""
}

type GetuiPushConf struct {
	AppId        string `json:"app_id"`
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
}

func (conf *GetuiPushConf) Valid() bool {
	return conf.AppId != "" && conf.AppKey != "" && conf.MasterSecret != ""
}

func Str2PushChannel(str string) pbobjs.PushChannel {
	switch str {
	case string(PushChannel_Apple):
		return pbobjs.PushChannel_Apple
	case string(PushChannel_Huawei):
		return pbobjs.PushChannel_Huawei
	case string(PushChannel_Xiaomi):
		return pbobjs.PushChannel_Xiaomi
	case string(PushChannel_OPPO):
		return pbobjs.PushChannel_Oppo
	case string(PushChannel_VIVO):
		return pbobjs.PushChannel_Vivo
	case string(PushChannel_Jpush):
		return pbobjs.PushChannel_JPush
	case string(PushChannel_FCM):
		return pbobjs.PushChannel_FCM
	case string(PushChannel_Meizu):
		return pbobjs.PushChannel_Meizhu
	case string(PushChannel_HONOR):
		return pbobjs.PushChannel_Honor
	case string(PushChannel_Getui):
		return pbobjs.PushChannel_Getui

	default:
		return pbobjs.PushChannel_DefaultChannel
	}
}

func PushChannel2Str(pushChannel pbobjs.PushChannel) string {
	switch pushChannel {
	case pbobjs.PushChannel_Apple:
		return string(PushChannel_Apple)
	case pbobjs.PushChannel_Huawei:
		return string(PushChannel_Huawei)
	case pbobjs.PushChannel_Xiaomi:
		return string(PushChannel_Xiaomi)
	case pbobjs.PushChannel_Oppo:
		return string(PushChannel_OPPO)
	case pbobjs.PushChannel_Vivo:
		return string(PushChannel_VIVO)
	case pbobjs.PushChannel_JPush:
		return string(PushChannel_Jpush)
	case pbobjs.PushChannel_FCM:
		return string(PushChannel_FCM)
	case pbobjs.PushChannel_Meizhu:
		return string(PushChannel_Meizu)
	case pbobjs.PushChannel_Honor:
		return string(PushChannel_HONOR)
	case pbobjs.PushChannel_Getui:
		return string(PushChannel_Getui)
	default:
		return ""
	}
}

func Str2Platform(str string) pbobjs.Platform {
	switch str {
	case string(Platform_Android):
		return pbobjs.Platform_Android
	case string(Platform_IOS):
		return pbobjs.Platform_iOS
	case string(Platform_Web):
		return pbobjs.Platform_Web
	case string(Platform_PC):
		return pbobjs.Platform_PC
	default:
		return pbobjs.Platform_DefaultPlatform
	}
}
func Platform2Str(platform pbobjs.Platform) string {
	switch platform {
	case pbobjs.Platform_Android:
		return string(Platform_Android)
	case pbobjs.Platform_iOS:
		return string(Platform_IOS)
	case pbobjs.Platform_Web:
		return string(Platform_Web)
	case pbobjs.Platform_PC:
		return string(Platform_PC)
	default:
		return ""
	}
}
