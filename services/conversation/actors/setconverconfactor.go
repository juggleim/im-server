package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type SetGlobalConverConfActor struct {
	bases.BaseActor
}

func (actor *SetGlobalConverConfActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.SetConverConfReq); ok {
		logs.WithContext(ctx).Infof("conver_id:%s\tchannel_type:%d\tsub_channel:%s\titem_type:%d\titem_key:%s", req.ConverId, req.ChannelType, req.SubChannel, req.ItemType, req.ItemKey)
		code := services.SetGlobalConverConf(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *SetGlobalConverConfActor) CreateInputObj() proto.Message {
	return &pbobjs.SetConverConfReq{}
}
