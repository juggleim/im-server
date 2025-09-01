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

type RtcConverBindActor struct {
	bases.BaseActor
}

func (actor *RtcConverBindActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.RtcConverBindReq); ok {
		logs.WithContext(ctx).Infof("action:%v\tconver_id:%s\tchannel_type:%d\tsub_channel:%s\trtc_room_id:%s", req.Action, req.ConverId, req.ChannelType, req.SubChannel, req.RtcRoomId)
		code := services.RtcConverBind(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *RtcConverBindActor) CreateInputObj() proto.Message {
	return &pbobjs.RtcConverBindReq{}
}
