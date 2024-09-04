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
	if req, ok := input.(*pbobjs.LogEntities); ok {
		for _, entity := range req.Entities {
			switch entity.GetLogOf().(type) {
			case *pbobjs.LogEntity_UserConnectLog:
				err = services.WriteUserConnectLog(entity.GetUserConnectLog())
			case *pbobjs.LogEntity_ConnectionLog:
				err = services.WriteConnectLog(entity.GetConnectionLog())
			default:
			}
			if err != nil {
				logs.WithContext(ctx).Errorf("write log error: %+v", err)
			}
		}

	} else {
		logs.WithContext(ctx).Infof("input is illegal")
	}
}

func (actor *LogServiceActor) CreateInputObj() proto.Message {
	return &pbobjs.LogEntities{}
}
