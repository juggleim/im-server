package logmanager

import (
	"context"
	"im-server/commons/pbdefines/pbobjs"
	"time"
)

func WriteConnectionLog(ctx context.Context, msg *pbobjs.ConnectionLog) {
	// clusters.AsyncRpcCall(ctx, "vlog", "", &pbobjs.LogEntity{
	// 	LogOf: &pbobjs.LogEntity_ConnectionLog{ConnectionLog: msg},
	// })
}
func WriteDisconnectionLog(ctx context.Context, msg *pbobjs.DisconnectionLog) {
	// clusters.AsyncRpcCall(ctx, "vlog", "", &pbobjs.LogEntity{
	// 	LogOf: &pbobjs.LogEntity_DisconnectionLog{DisconnectionLog: msg},
	// })
}
func WriteSdkRequestLog(ctx context.Context, msg *pbobjs.SdkRequestLog) {
	// clusters.AsyncRpcCall(ctx, "vlog", "", &pbobjs.LogEntity{
	// 	LogOf: &pbobjs.LogEntity_SdkRequestLog{SdkRequestLog: msg},
	// })
}

func WriteSdkResponseLog(ctx context.Context, msg *pbobjs.SdkResponseLog) {
	// clusters.AsyncRpcCall(ctx, "vlog", "", &pbobjs.LogEntity{
	// 	LogOf: &pbobjs.LogEntity_SdkResponseLog{SdkResponseLog: msg},
	// })
}

func WriteBusinessLog(ctx context.Context, msg *pbobjs.BusinessLog) {
	// clusters.AsyncRpcCall(ctx, "vlog", "", &pbobjs.LogEntity{
	// 	LogOf: &pbobjs.LogEntity_BusinessLog{BusinessLog: msg},
	// })
}

func LogTimestamp() string {
	return time.Now().Format("060102150405.000")
}
