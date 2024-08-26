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

type DelUserTagsActor struct {
	bases.BaseActor
}

func (actor *DelUserTagsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserTagList); ok {
		err := services.DelUserTags(ctx, req)
		if err != nil {
			logs.WithContext(ctx).Errorf("GetUserTags failed:%v", err)
			ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_DEFAULT, nil)
			actor.Sender.Tell(ack, actorsystem.NoSender)
		} else {
			ack := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, nil)
			actor.Sender.Tell(ack, actorsystem.NoSender)
		}
		logs.WithContext(ctx).Infof("user_id:%s\t%v", bases.GetTargetIdFromCtx(ctx), req)
	} else {
		logs.WithContext(ctx).Infof("input is illegal")
	}
}

func (actor *DelUserTagsActor) CreateInputObj() proto.Message {
	return &pbobjs.UserTagList{}
}
