package actors

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type QryGlobalConversActor struct {
	bases.BaseActor
}

func (actor *QryGlobalConversActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.QryGlobalConversReq); ok {
		logs.WithContext(ctx).Infof("start:%d\torder:%d\tcount:%d\ttarget_id:%s\tchannel_type:%d", req.Start, req.Order, req.Count, req.TargetId, req.ChannelType)
		resp := services.QryGlobalConvers(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, 0, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		fmt.Println("p_msg, upMsg is illigal.")
	}
}

func (actor *QryGlobalConversActor) CreateInputObj() proto.Message {
	return &pbobjs.QryGlobalConversReq{}
}
