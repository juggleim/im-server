package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"

	"google.golang.org/protobuf/proto"
)

type AddBotActor struct {
	bases.BaseActor
}

func (actor *AddBotActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserInfo); ok {
		logs.WithContext(ctx).Infof("bot_id:%d\tnickname:%s", req.UserId, req.Nickname)
		code := services.AddUser(ctx, req.UserId, req.Nickname, req.UserPortrait, req.ExtFields, req.Settings, pbobjs.UserType_Bot)
		ack := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illigal")
	}
}

func (actor *AddBotActor) CreateInputObj() proto.Message {
	return &pbobjs.UserInfo{}
}
