package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/statussubscriptions/services"

	"google.golang.org/protobuf/proto"
)

type SyncSubChangeActor struct {
	bases.BaseActor
}

func (actor *SyncSubChangeActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.SubRelChangeReq); ok {
		services.SyncSubChange(ctx, req)
	}
}

func (actor *SyncSubChangeActor) CreateInputObj() proto.Message {
	return &pbobjs.SubRelChangeReq{}
}
