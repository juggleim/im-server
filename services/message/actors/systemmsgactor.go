package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type SystemMsgActor struct {
	bases.BaseActor
}

func (actor *SystemMsgActor) OnReceive(ctx context.Context, input proto.Message) {
	isFromApi := bases.GetIsFromApiFromCtx(ctx)
	isFromApp := bases.GetIsFromAppFromCtx(ctx)
	if isFromApi || isFromApp {
		if upMsg, ok := input.(*pbobjs.UpMsg); ok {
			userId := bases.GetRequesterIdFromCtx(ctx)
			exts := bases.GetExtsFromCtx(ctx)
			if receiverId, exist := exts[commonservices.RpcExtKey_RealTargetId]; exist {
				logs.WithContext(ctx).WithField("method", "s_msg").Infof("sender:%s\treceiver:%s", userId, receiverId)
				code, msgId, sendTime, msgSeq := services.SendSystemMsg(ctx, userId, receiverId, upMsg)
				userPubAck := bases.CreateUserPubAckWraper(ctx, code, msgId, sendTime, msgSeq, "")
				actor.Sender.Tell(userPubAck, actorsystem.NoSender)
				logs.WithContext(ctx).Infof("code:%d", code)
			} else {
				logs.WithContext(ctx).Errorf("have no receiver")
			}
		} else {
			logs.WithContext(ctx).Errorf("upMsg is illigal. upMsg:%v", upMsg)
		}
	} else {
		userPubAck := bases.CreateUserPubAckWraper(ctx, errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC, "", 0, 0, "")
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
	}
}

func (actor *SystemMsgActor) CreateInputObj() proto.Message {
	return &pbobjs.UpMsg{}
}
