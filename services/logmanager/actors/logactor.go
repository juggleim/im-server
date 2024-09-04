package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/logmanager/services"

	"google.golang.org/protobuf/proto"
)

type LogServiceActor struct {
	bases.BaseActor
}

func (actor *LogServiceActor) OnReceive(ctx context.Context, input proto.Message) {
	var err error
	if req, ok := input.(*pbobjs.LogEntity); ok {
		switch req.GetLogOf().(type) {
		case *pbobjs.LogEntity_UserConnectLog:
			err = services.WriteUserConnectLog(req.GetUserConnectLog())
		case *pbobjs.LogEntity_ConnectionLog:
			err = services.WriteConnectLog(req.GetConnectionLog())
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
