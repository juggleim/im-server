package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"

	"google.golang.org/protobuf/proto"
)

type QryConnectCount struct {
	bases.BaseActor
}

func (actor *QryConnectCount) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryConnectCountReq); ok {
		logs.WithContext(ctx).Infof("start:%d\tend:%d", req.Start, req.End)
		appkey := bases.GetAppKeyFromCtx(ctx)
		items := commonservices.QryConncurrentConnect(appkey, req.Start, req.End)
		resp := &pbobjs.QryConnectCountResp{
			Items: []*pbobjs.ConnectCountItem{},
		}
		for _, item := range items {
			resp.Items = append(resp.Items, &pbobjs.ConnectCountItem{
				TimeMark: item.TimeMark,
				Count:    item.Count,
			})
		}
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, resp)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Info("input is illegal")
	}
}

func (actor *QryConnectCount) CreateInputObj() proto.Message {
	return &pbobjs.QryConnectCountReq{}
}
