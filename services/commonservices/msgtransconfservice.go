package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/transengines"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func GetTransEngine(appkey string) transengines.ITransEngine {
	appInfo, exist := GetAppInfo(appkey)
	if exist && appInfo != nil {
		if appInfo.TransEngine == nil {
			lock := appLocks.GetLocks(appInfo.AppKey)
			lock.Lock()
			defer lock.Unlock()
			loadTransEngine(appInfo)
		}
		if appInfo.TransEngine != nil {
			return appInfo.TransEngine
		}
	}
	return transengines.DefaultTransEngine
}

func loadTransEngine(appInfo *AppInfo) {
	if appInfo.TransEngineConf == "" {
		appInfo.TransEngine = transengines.DefaultTransEngine
		return
	}
	transConf := &TransEngineConf{}
	err := tools.JsonUnMarshal([]byte(appInfo.TransEngineConf), transConf)
	if err != nil {
		appInfo.TransEngine = transengines.DefaultTransEngine
		return
	}
	if transConf.Channel == "baidu" && transConf.BdTransEngine != nil && transConf.BdTransEngine.ApiKey != "" && transConf.BdTransEngine.SecretKey != "" {
		appInfo.TransEngine = &transengines.BdTransEngine{
			AppKey:    appInfo.AppKey,
			ApiKey:    transConf.BdTransEngine.ApiKey,
			SecretKey: transConf.BdTransEngine.SecretKey,
		}
	} else if transConf.Channel == "deepl" && transConf.DeeplTransEngine != nil && transConf.DeeplTransEngine.AuthKey != "" {
		appInfo.TransEngine = &transengines.DeeplTransEngine{
			AppKey:  appInfo.AppKey,
			AuthKey: transConf.DeeplTransEngine.AuthKey,
		}
	} else {
		if transConf.BdTransEngine != nil && transConf.BdTransEngine.ApiKey != "" && transConf.BdTransEngine.SecretKey != "" {
			appInfo.TransEngine = &transengines.BdTransEngine{
				AppKey:    appInfo.AppKey,
				ApiKey:    transConf.BdTransEngine.ApiKey,
				SecretKey: transConf.BdTransEngine.SecretKey,
			}
		} else if transConf.DeeplTransEngine != nil && transConf.DeeplTransEngine.AuthKey != "" {
			appInfo.TransEngine = &transengines.DeeplTransEngine{
				AppKey:  appInfo.AppKey,
				AuthKey: transConf.DeeplTransEngine.AuthKey,
			}
		} else {
			appInfo.TransEngine = transengines.DefaultTransEngine
		}
	}
}

type TransEngineConf struct {
	Channel          string                         `json:"channel,omitempty"`
	BdTransEngine    *transengines.BdTransEngine    `json:"baidu,omitempty"`
	DeeplTransEngine *transengines.DeeplTransEngine `json:"deepl,omitempty"`
}

type MsgTransConfs struct {
	Confs map[string][]string
}

func loadMsgTransConfs(appInfo *AppInfo) {
	lock := appLocks.GetLocks(appInfo.AppKey)
	lock.Lock()
	defer lock.Unlock()
	if appInfo.MsgTransConfs == nil {
		dao := dbs.MsgTransConfDao{}
		var startId int64 = 0
		msgTransConfs := &MsgTransConfs{
			Confs: make(map[string][]string),
		}
		for i := 0; i < 3; i++ {
			confs, err := dao.QueryConfs(appInfo.AppKey, startId, 1000)
			if err != nil {
				break
			}
			for _, conf := range confs {
				if conf.ID > startId {
					startId = conf.ID
				}
				var jsonPaths []string
				if arr, exist := msgTransConfs.Confs[conf.MsgType]; exist {
					jsonPaths = arr
				} else {
					jsonPaths = []string{}
				}
				jsonPaths = append(jsonPaths, conf.JsonPath)
				msgTransConfs.Confs[conf.MsgType] = jsonPaths
			}
			if len(confs) < 1000 {
				break
			}
		}
		appInfo.MsgTransConfs = msgTransConfs
	}
}

func GetMsgTransConfs(appkey, msgType string) []string {
	appInfo, exist := GetAppInfo(appkey)
	if exist && appInfo != nil {
		if appInfo.MsgTransConfs == nil {
			loadMsgTransConfs(appInfo)
		}
		if appInfo.MsgTransConfs != nil {
			return appInfo.MsgTransConfs.Confs[msgType]
		}
	}
	return []string{}
}

func TranslateMsg(ctx context.Context, lans []string, downMsg *pbobjs.DownMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	transEngine := GetTransEngine(appkey)
	if transEngine == nil || transEngine == transengines.DefaultTransEngine {
		return
	}
	jsonPaths := GetMsgTransConfs(appkey, downMsg.MsgType)
	if len(lans) <= 0 || len(jsonPaths) <= 0 {
		return
	}
	downMsg.TransMsgMap = map[string]*pbobjs.TransMsgContent{}
	for _, lan := range lans {
		var pushData *pbobjs.PushData
		if downMsg.PushData != nil {
			pushData = &pbobjs.PushData{}
		}
		downMsg.TransMsgMap[lan] = &pbobjs.TransMsgContent{
			MsgContent: downMsg.MsgContent,
			PushData:   pushData,
		}
	}
	wg := &sync.WaitGroup{}
	var lock *sync.RWMutex
	if len(jsonPaths) > 0 {
		lock = &sync.RWMutex{}
	}
	for _, jsonPath := range jsonPaths {
		wg.Add(1)
		jsonP := jsonPath
		go func() {
			defer wg.Done()
			text := getStrByPath(downMsg.MsgContent, jsonP)
			afterTransMap := transEngine.Translate(text, lans)
			if lock != nil {
				lock.Lock()
				defer lock.Unlock()
			}
			for lan, msgContent := range downMsg.TransMsgMap {
				if afterTrans, exist := afterTransMap[lan]; exist {
					newContent, err := setStrByPath(msgContent.MsgContent, jsonP, afterTrans)
					if err == nil {
						msgContent.MsgContent = newContent
					}
				}
			}
		}()
	}
	if downMsg.PushData != nil {
		if downMsg.PushData.Title != "" {
			wg.Add(1)
			go func() {
				wg.Done()
				afterTransMap := transEngine.Translate(downMsg.PushData.Title, lans)
				for lan, msgContent := range downMsg.TransMsgMap {
					if afterTrans, exist := afterTransMap[lan]; exist {
						msgContent.PushData.Title = afterTrans
					}
				}
			}()
		}
		if downMsg.PushData.PushText != "" {
			wg.Add(1)
			go func() {
				wg.Done()
				afterTransMap := transEngine.Translate(downMsg.PushData.PushText, lans)
				for lan, msgContent := range downMsg.TransMsgMap {
					if afterTrans, exist := afterTransMap[lan]; exist {
						msgContent.PushData.PushText = afterTrans
					}
				}
			}()
		}
	}
	wg.Wait()
}

func getStrByPath(jsonData []byte, path string) string {
	result := gjson.GetBytes(jsonData, path)
	return result.String()
}

func setStrByPath(jsonData []byte, path, value string) ([]byte, error) {
	nb, err := sjson.SetBytes(jsonData, path, value)
	return nb, err
}
