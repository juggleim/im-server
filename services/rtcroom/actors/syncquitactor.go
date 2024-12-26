package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/rtcroom/services"

	"google.golang.org/protobuf/proto"
)

type SyncQuitActor struct {
	bases.BaseActor
}

func (actor *SyncQuitActor) OnReceive(ctx context.Context, input proto.Message) {
	userId := bases.GetRequesterIdFromCtx(ctx)
	logs.WithContext(ctx).Infof("user_id:%s", userId)
	services.SyncQuitWhenConnectKicked(ctx, userId)
}

func (actor *SyncQuitActor) CreateInputObj() proto.Message {
	return &pbobjs.Nil{}
}
