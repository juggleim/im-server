package logmanager

import (
	"context"
	"im-server/commons/configures"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/logmanager/services"
	"time"
)

func checkSwitch(appkey string) bool {
	if !configures.Config.Log.Visual {
		return false
	}
	appInfo, exist := commonservices.GetAppInfo(appkey)
	if exist && appInfo != nil {
		return appInfo.OpenVisualLog
	}
	return false
}

func WriteUserConnectLog(ctx context.Context, msg *pbobjs.UserConnectLog) {
	if checkSwitch(msg.AppKey) {
		msg.Timestamp = time.Now().UnixMilli()
		services.WriteUserConnectLog(msg)
	}
}

func WriteConnectionLog(ctx context.Context, msg *pbobjs.ConnectionLog) {
	if checkSwitch(msg.AppKey) {
		msg.Timestamp = time.Now().UnixMilli()
		services.WriteConnectLog(msg)
	}
}

func WriteBusinessLog(ctx context.Context, msg *pbobjs.BusinessLog) {
	// clusters.AsyncRpcCall(ctx, "vlog", "", &pbobjs.LogEntity{
	// 	LogOf: &pbobjs.LogEntity_BusinessLog{BusinessLog: msg},
	// })
}
