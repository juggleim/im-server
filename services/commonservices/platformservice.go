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

	PushChannel_Apple  PushChannel = "Apple"
	PushChannel_Huawei PushChannel = "Huawei"
	PushChannel_Xiaomi PushChannel = "Xiaomi"
	PushChannel_OPPO   PushChannel = "Oppo"
	PushChannel_VIVO   PushChannel = "Vivo"
	PushChannel_Jpush  PushChannel = "Jpush"
	PushChannel_FCM    PushChannel = "FCM"
)

type AndroidPushConf struct {
	AppKey      string         `json:"app_key"`
	PushChannel string         `json:"push_channel,omitempty"`
	Package     string         `json:"package"`
	Extra       map[string]any `json:"extra"`
}
type HuaweiPushConf struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func (conf *HuaweiPushConf) Valid() bool {
	return conf.AppId != "" && conf.AppSecret != ""
}

type XiaomiPushConf struct {
	AppSecret string `json:"app_secret"`
}

func (conf *XiaomiPushConf) Valid() bool {
	return conf.AppSecret != ""
}

type OppoPushConf struct {
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
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
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
}

func (conf *JPushConf) Valid() bool {
	return conf.AppKey != "" && conf.MasterSecret != ""
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
