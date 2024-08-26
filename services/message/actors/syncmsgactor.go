package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type SyncMsgActor struct {
	bases.BaseActor
}

func (actor *SyncMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	if syncMsg, ok := input.(*pbobjs.SyncMsgReq); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tsync_time:%d\tsendbox_sync_time:%d\tcontain_sendbox:%v", bases.GetRequesterIdFromCtx(ctx), syncMsg.SyncTime, syncMsg.SendBoxSyncTime, syncMsg.ContainsSendBox)
		code, messages := services.SyncMessages(ctx, syncMsg)
		rpcMsg := bases.CreateQueryAckWraper(ctx, code, messages)
		actor.Sender.Tell(rpcMsg, actorsystem.NoSender)
		logs.WithContext(ctx).Infof("result_len:%d", len(messages.Msgs))
	} else {
		logs.WithContext(ctx).Error("Failed to decode.")
		rpcMsg := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, nil)
		actor.Sender.Tell(rpcMsg, actorsystem.NoSender)
	}
}

func (actor *SyncMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.SyncMsgReq{}
}
