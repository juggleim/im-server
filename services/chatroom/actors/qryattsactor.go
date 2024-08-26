package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/chatroom/services"
	"im-server/services/commonservices/logs"
	"time"

	"google.golang.org/protobuf/proto"
)

type QryAttsActor struct {
	bases.BaseActor
}

func (actor *QryAttsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.ChatroomInfo); ok {
		logs.WithContext(ctx).Infof("user_id:%s\tchat_id:%s", bases.GetRequesterIdFromCtx(ctx), req.ChatId)
		code, atts := services.QryChatAtts(ctx, req)
		qrAck := bases.CreateQueryAckWraper(ctx, code, atts)
		actor.Sender.Tell(qrAck, actorsystem.NoSender)
	} else {
		userPubAck := bases.CreateUserPubAckWraper(ctx, errs.IMErrorCode_PBILLEGAL, "", time.Now().UnixMilli(), 0)
		actor.Sender.Tell(userPubAck, actorsystem.NoSender)
		logs.WithContext(ctx).Error("input is illigal.")
	}
}

func (actor *QryAttsActor) CreateInputObj() proto.Message {
	return &pbobjs.ChatroomInfo{}
}
