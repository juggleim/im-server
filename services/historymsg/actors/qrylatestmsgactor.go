package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type QryLatestMsgActor struct {
	bases.BaseActor
}

func (actor *QryLatestMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if qryMsgIndexReq, ok := input.(*pbobjs.QryLatestMsgReq); ok {
		if qryMsgIndexReq != nil {
			latest := services.QryLatestHisMsg(ctx, bases.GetAppKeyFromCtx(ctx), qryMsgIndexReq.ConverId, qryMsgIndexReq.ChannelType)
			ack := bases.CreateQueryAckWraper(ctx, 0, &pbobjs.QryLatestMsgResp{
				ConverId:    qryMsgIndexReq.ConverId,
				ChannelType: qryMsgIndexReq.ChannelType,
				MsgSeqNo:    latest.LatestMsgSeq,
				MsgTime:     latest.LatestMsgTime,
				MsgId:       latest.LatestMsgId,
			})
			actor.Sender.Tell(ack, actorsystem.NoSender)
		}
	}
}

func (actor *QryLatestMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.QryLatestMsgReq{}
}
