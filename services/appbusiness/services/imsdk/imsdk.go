package imsdk

import (
	"fmt"
	"im-server/commons/caches"
	"im-server/commons/configures"
	"im-server/services/commonservices"
	"sync"
	"time"

	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

var sdkCache *caches.LruCache
var lock *sync.RWMutex

func init() {
	sdkCache = caches.NewLruCacheWithAddReadTimeout("jimsdk_cache", 1000, nil, 5*time.Minute, 5*time.Minute)
	lock = &sync.RWMutex{}
}

func GetImSdk(appkey string) *juggleimsdk.JuggleIMSdk {
	if val, exist := sdkCache.Get(appkey); exist {
		return val.(*juggleimsdk.JuggleIMSdk)
	} else {
		lock.Lock()
		defer lock.Unlock()
		if val, exist := sdkCache.Get(appkey); exist {
			return val.(*juggleimsdk.JuggleIMSdk)
		} else {
			if appinfo, exist := commonservices.GetAppInfo(appkey); exist {
				sdk := juggleimsdk.NewJuggleIMSdk(appkey, appinfo.AppSecret, fmt.Sprintf("http://127.0.0.1:%d", configures.Config.ApiGateway.HttpPort))
				sdkCache.Add(appkey, sdk)
				return sdk
			}
			return nil
		}
	}
}
