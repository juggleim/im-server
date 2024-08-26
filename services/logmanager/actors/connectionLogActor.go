package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/logmanager/services"
	"time"

	"google.golang.org/protobuf/proto"
)

type LogServiceActor struct {
	bases.BaseActor
}

func (actor *LogServiceActor) OnReceive(ctx context.Context, input proto.Message) {
	var err error
	if req, ok := input.(*pbobjs.LogEntity); ok {
		switch req.GetLogOf().(type) {
		case *pbobjs.LogEntity_ConnectionLog:
			data := req.GetConnectionLog()
			err = services.WriteConnectionLog(data.AppKey, data.UserId, time.Now().UnixMilli(), data.String())
		case *pbobjs.LogEntity_DisconnectionLog:
			data := req.GetDisconnectionLog()
			err = services.WriteConnectionLog(data.AppKey, data.UserId, time.Now().UnixMilli(), data.String())
		case *pbobjs.LogEntity_SdkRequestLog:
			data := req.GetSdkRequestLog()
			err = services.WriteSdkLog(data.AppKey, data.Session, data.Index, time.Now().UnixMilli(), data.String())
		case *pbobjs.LogEntity_SdkResponseLog:
			data := req.GetSdkResponseLog()
			err = services.WriteSdkLog(data.AppKey, data.Session, data.Index, time.Now().UnixMilli(), data.String())
		case *pbobjs.LogEntity_BusinessLog:
			data := req.GetBusinessLog()
			err = services.WriteBusinessLog(data.AppKey, data.Session, data.Index, time.Now().UnixMilli(), data.String())
		default:
		}

		if err != nil {
			logs.WithContext(ctx).Errorf("write log error: %+v", err)
		}

	} else {
		logs.WithContext(ctx).Infof("input is illegal")
	}
}

func (actor *LogServiceActor) CreateInputObj() proto.Message {
	return &pbobjs.LogEntity{}
}
