package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/services"

	"google.golang.org/protobuf/proto"
)

type UpdPushStatusActor struct {
	bases.BaseActor
}

func (actor *UpdPushStatusActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserPushStatus); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		appkey := bases.GetAppKeyFromCtx(ctx)
		if services.UserStatusCacheContains(appkey, userId) {
			userStatus := services.GetUserStatus(appkey, userId)
			if req.CanPush {
				userStatus.SetPushStatus(1)
			} else {
				userStatus.SetPushStatus(0)
			}
		}
	}
}

func (actor *UpdPushStatusActor) CreateInputObj() proto.Message {
	return &pbobjs.UserPushStatus{}
}
