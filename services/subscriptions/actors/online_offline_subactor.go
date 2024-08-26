package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/subscriptions/services"

	"google.golang.org/protobuf/proto"
)

type OnlineOfflineSubActor struct {
	bases.BaseActor
}

func (actor *OnlineOfflineSubActor) OnReceive(ctx context.Context, input proto.Message) {
	services.OnlineOfflineHandle(ctx, input.(*pbobjs.OnlineOfflineMsg))
}

func (actor *OnlineOfflineSubActor) CreateInputObj() proto.Message {
	return &pbobjs.OnlineOfflineMsg{}
}
