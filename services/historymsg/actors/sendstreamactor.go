package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/historymsg/services"

	"google.golang.org/protobuf/proto"
)

type SendStreamActor struct {
	bases.BaseActor
}

func (actor *SendStreamActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.StreamMsg); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		targetId := req.TargetId
		code := services.SendStreamMsg(ctx, req)
		logs.WithContext(ctx).Infof("sender:%s\treceiver:%s\tseq:%d\tis_finished:%v\tstream_id:%s\tcontent:%s\tcode:%d", userId, targetId, req.Seq, req.IsFinished, req.StreamMsgId, string(req.PartialContent), code)
	}
}

func (actor *SendStreamActor) CreateInputObj() proto.Message {
	return &pbobjs.StreamMsg{}
}
