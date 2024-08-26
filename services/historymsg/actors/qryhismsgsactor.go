package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"
	"time"

	"google.golang.org/protobuf/proto"
)

type QryHistoryMsgsActor struct {
	bases.BaseActor
}

func (actor *QryHistoryMsgsActor) OnReceive(ctx context.Context, input proto.Message) {
	if qryHisMsgsReq, ok := input.(*pbobjs.QryHisMsgsReq); ok {
		appkey := bases.GetAppKeyFromCtx(ctx)
		startTime := qryHisMsgsReq.StartTime
		isPositiveOrder := false
		logs.WithContext(ctx).Infof("conver_id:%s\tchannel_type:%v\tstart_time:%d\tcount:%d\torder:%d", qryHisMsgsReq.TargetId, qryHisMsgsReq.ChannelType, qryHisMsgsReq.StartTime, qryHisMsgsReq.Count, qryHisMsgsReq.Order)
		if qryHisMsgsReq.Order == 0 { //0:倒序;1:正序;
			if startTime <= 0 {
				startTime = time.Now().UnixMilli()
			}
		} else {
			isPositiveOrder = true
		}
		code, resp := services.QryHisMsgs(ctx, appkey, qryHisMsgsReq.TargetId, qryHisMsgsReq.ChannelType, startTime, qryHisMsgsReq.Count, isPositiveOrder, qryHisMsgsReq.MsgTypes)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		msgCount := 0
		if resp != nil {
			msgCount = len(resp.Msgs)
		}
		logs.WithContext(ctx).Infof("result_code:%d\tlen:%d", code, msgCount)
	} else {
		qryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
		logs.WithContext(ctx).Error("failed to decode input")
	}
}

func (actor *QryHistoryMsgsActor) CreateInputObj() proto.Message {
	return &pbobjs.QryHisMsgsReq{}
}
