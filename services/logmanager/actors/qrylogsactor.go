package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/logmanager/services"

	"google.golang.org/protobuf/proto"
)

type QryLogsActor struct {
	bases.BaseActor
}

func (actor *QryLogsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryServerLogsReq); ok {
		appkey := bases.GetAppKeyFromCtx(ctx)
		ret := &pbobjs.QryServerLogsResp{
			Logs: []string{},
		}
		var logs []services.LogEntity
		var err error
		if req.LogType == string(services.ServerLogType_UserConnect) {
			logs, err = services.QryUserConnectLogs(appkey, req.UserId, req.Start, req.Count)
		} else if req.LogType == string(services.ServerLogType_Connect) {
			logs, err = services.QryConnectLogs(appkey, req.Session, req.Start, req.Count)
		} else if req.LogType == string(services.ServerLogType_Business) {
			logs, err = services.QryBusinessLogs(appkey, req.Session, req.Index, req.Start, req.Count)
		}
		if err == nil {
			for _, log := range logs {
				ret.Logs = append(ret.Logs, log.Value)
			}
		}
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, ret)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Info("input is illegal")
	}
}

func (actor *QryLogsActor) CreateInputObj() proto.Message {
	return &pbobjs.QryServerLogsReq{}
}
