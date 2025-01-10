package commonservices

import (
	"im-server/commons/tools"
	"im-server/services/commonservices/sms"
)

func GetSmsEngine(appkey string) sms.ISmsEngine {
	appInfo, exist := GetAppInfo(appkey)
	if exist && appInfo != nil {
		if appInfo.SmsEngine == nil {
			lock := appLocks.GetLocks(appkey)
			lock.Lock()
			defer lock.Unlock()
			loadSmsEngine(appInfo)
		}
		if appInfo.SmsEngine != nil {
			return appInfo.SmsEngine
		}
	}
	return sms.DefaultSmsEngine
}

func loadSmsEngine(appInfo *AppInfo) {
	if appInfo.SmsEngineConf == "" {
		appInfo.SmsEngine = sms.DefaultSmsEngine
		return
	}
	smsConf := &SmsEngineConf{}
	err := tools.JsonUnMarshal([]byte(appInfo.SmsEngineConf), smsConf)
	if err != nil {
		appInfo.SmsEngine = sms.DefaultSmsEngine
		return
	}
	if smsConf.Channel == "baidu" && smsConf.BdSmsEngine != nil && smsConf.BdSmsEngine.ApiKey != "" && smsConf.BdSmsEngine.SecretKey != "" {
		appInfo.SmsEngine = smsConf.BdSmsEngine
	} else {
		appInfo.SmsEngine = sms.DefaultSmsEngine
		return
	}
}

type SmsEngineConf struct {
	Channel     string           `json:"channel,omitempty"`
	BdSmsEngine *sms.BdSmsEngine `json:"baidu,omitempty"`
}
