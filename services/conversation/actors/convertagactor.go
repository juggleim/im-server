package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/conversation/services"

	"google.golang.org/protobuf/proto"
)

type TagAddConversActor struct {
	bases.BaseActor
}

func (actor *TagAddConversActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.TagConvers); ok {
		logs.WithContext(ctx).Infof("user_id:%s\ttag:%s\ttag_name:%s\tconvers:%v", bases.GetRequesterIdFromCtx(ctx), req.Tag, req.TagName, req.Convers)
		code := services.TagAddConvers(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *TagAddConversActor) CreateInputObj() proto.Message {
	return &pbobjs.TagConvers{}
}

type TagDelConversActor struct {
	bases.BaseActor
}

func (actor *TagDelConversActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.TagConvers); ok {
		logs.WithContext(ctx).Infof("user_id:%s\ttag:%s\ttag_name:%s\tconvers:%v", bases.GetRequesterIdFromCtx(ctx), req.Tag, req.TagName, req.Convers)
		code := services.TagDelConvers(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *TagDelConversActor) CreateInputObj() proto.Message {
	return &pbobjs.TagConvers{}
}

type QryUserConverTagsActor struct {
	bases.BaseActor
}

func (actor *QryUserConverTagsActor) OnReceive(ctx context.Context, input proto.Message) {
	if _, ok := input.(*pbobjs.Nil); ok {
		logs.WithContext(ctx).Infof("user_id:%s", bases.GetRequesterIdFromCtx(ctx))
		resp, code := services.QryUserConverTags(ctx)
		qryAck := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *QryUserConverTagsActor) CreateInputObj() proto.Message {
	return &pbobjs.Nil{}
}

type DelUserConverTagsActor struct {
	bases.BaseActor
}

func (actor *DelUserConverTagsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserConverTags); ok {
		logs.WithContext(ctx).Infof("user_id:%s\treq:%v", bases.GetRequesterIdFromCtx(ctx), req)
		code := services.DelUserConverTags(ctx, req)
		qryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(qryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Error("input is illegal")
	}
}

func (actor *DelUserConverTagsActor) CreateInputObj() proto.Message {
	return &pbobjs.UserConverTags{}
}
