package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/connectmanager/services"

	"google.golang.org/protobuf/proto"
)

type BanUsersActor struct {
	bases.BaseActor
}

func (actor *BanUsersActor) OnReceive(ctx context.Context, input proto.Message) {
	if banReq, ok := input.(*pbobjs.BanUsersReq); ok {
		logs.WithContext(ctx).WithField("method", "ban_users").Infof("is_add:%v\tusers:%v", banReq.IsAdd, banReq.BanUsers)
		if len(banReq.BanUsers) > 0 {
			if banReq.IsAdd {
				services.BanUsers(ctx, banReq.BanUsers)
			} else {
				services.UnBanUsers(ctx, banReq.BanUsers)
			}
		}
	} else {
		logs.WithContext(ctx).Errorf("ban_users, input is illigal.")
	}
	ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, nil)
	actor.Sender.Tell(ack, actorsystem.NoSender)
}

func (actor *BanUsersActor) CreateInputObj() proto.Message {
	return &pbobjs.BanUsersReq{}
}
