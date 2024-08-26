package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/pushmanager/services"

	"google.golang.org/protobuf/proto"
)

type QryUserTagsActor struct {
	bases.BaseActor
}

func (actor *QryUserTagsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserIds); ok {
		res, err := services.GetUserTags(ctx, req.UserIds)
		if err != nil {
			logs.WithContext(ctx).Errorf("GetUserTags failed:%v", err)
			ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_DEFAULT, nil)
			actor.Sender.Tell(ack, actorsystem.NoSender)
		} else {
			ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, res)
			actor.Sender.Tell(ack, actorsystem.NoSender)
		}
		logs.WithContext(ctx).Infof("user_id:%s\t%v", bases.GetTargetIdFromCtx(ctx), req)
	} else {
		logs.WithContext(ctx).Infof("input is illegal")
	}
}

func (actor *QryUserTagsActor) CreateInputObj() proto.Message {
	return &pbobjs.UserIds{}
}
