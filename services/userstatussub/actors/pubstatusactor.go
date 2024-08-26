package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/userstatussub/services"

	"google.golang.org/protobuf/proto"
)

type PublishStatusActor struct {
	bases.BaseActor
}

func (actor *PublishStatusActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserInfo); ok {
		services.PublishStatus(ctx, req)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *PublishStatusActor) CreateInputObj() proto.Message {
	return &pbobjs.UserInfo{}
}
