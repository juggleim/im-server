package commonservices

import (
	"im-server/services/commonservices/dbs"

	"github.com/kataras/i18n"
)

var nilI18n *i18n.I18n

func init() {
	langMap := i18n.LangMap(map[string]i18n.Map{})
	langMap["nil"] = i18n.Map(map[string]interface{}{})
	loader := i18n.KV(langMap, i18n.DefaultLoaderConfig)
	nilI18n, _ = i18n.New(loader, []string{}...)
}

func loadI18nKeys(appInfo *AppInfo) {
	lock := appLocks.GetLocks(appInfo.AppKey)
	lock.Lock()
	defer lock.Unlock()
	if appInfo.I18nKeys == nil {
		langMap := i18n.LangMap(map[string]i18n.Map{})
		dao := dbs.I18nKeyDao{}
		var startId int64 = 0
		for i := 0; i < 3; i++ {
			kvs, err := dao.Query(appInfo.AppKey, startId, 10000)
			if err != nil {
				break
			}
			for _, kv := range kvs {
				if startId < kv.ID {
					startId = kv.ID
				}
				kvMap := i18n.Map(map[string]interface{}{})
				if old, exist := langMap[kv.Lang]; exist {
					kvMap = old
				} else {
					langMap[kv.Lang] = kvMap
				}
				kvMap[kv.Key] = kv.Value
			}
			if len(kvs) < 10000 {
				break
			}
		}
		if len(langMap) > 0 {
			loader := i18n.KV(langMap, i18n.DefaultLoaderConfig)
			appInfo.I18nKeys, _ = i18n.New(loader, []string{}...)
		} else {
			appInfo.I18nKeys = nilI18n
		}
	}
}

func GetI18nValue(appkey string, lang, key, defaultVal string) string {
	appInfo, exist := GetAppInfo(appkey)
	if exist && appInfo != nil {
		if appInfo.I18nKeys == nil {
			loadI18nKeys(appInfo)
		}
		if appInfo.I18nKeys != nil && appInfo.I18nKeys != nilI18n {
			val := appInfo.I18nKeys.Tr(lang, key)
			if val == "" {
				return defaultVal
			}
			return val
		}
	}
	return defaultVal
}
