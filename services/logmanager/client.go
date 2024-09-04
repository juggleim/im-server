package logmanager

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"time"
)

func WriteUserConnectLog(ctx context.Context, msg *pbobjs.UserConnectLog) {
	msg.Timestamp = time.Now().UnixMilli()
	data, _ := tools.PbMarshal(&pbobjs.LogEntity{
		LogOf: &pbobjs.LogEntity_UserConnectLog{
			UserConnectLog: msg,
		},
	})
	bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		AppKey:       msg.AppKey,
		Session:      msg.Session,
		Method:       "vlog",
		TargetId:     msg.UserId,
		AppDataBytes: data,
	})
}

func WriteConnectionLog(ctx context.Context, msg *pbobjs.ConnectionLog) {
	msg.Timestamp = time.Now().UnixMilli()
	data, _ := tools.PbMarshal(&pbobjs.LogEntity{
		LogOf: &pbobjs.LogEntity_ConnectionLog{
			ConnectionLog: msg,
		},
	})
	bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		AppKey:       msg.AppKey,
		Session:      msg.Session,
		Method:       "vlog",
		TargetId:     msg.Session,
		AppDataBytes: data,
	})
}

func WriteBusinessLog(ctx context.Context, msg *pbobjs.BusinessLog) {
	// clusters.AsyncRpcCall(ctx, "vlog", "", &pbobjs.LogEntity{
	// 	LogOf: &pbobjs.LogEntity_BusinessLog{BusinessLog: msg},
	// })
}
