package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"

	"google.golang.org/protobuf/proto"
)

type SetUserSettingActor struct {
	bases.BaseActor
}

func (actor *SetUserSettingActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.UserInfo); ok {
		userId := bases.GetTargetIdFromCtx(ctx)
		code := services.SetUserSettings(ctx, userId, req)
		queryAck := bases.CreateQueryAckWraper(ctx, code, nil)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *SetUserSettingActor) CreateInputObj() proto.Message {
	return &pbobjs.UserInfo{}
}

type GetUserSettingActor struct {
	bases.BaseActor
}

func (actor *GetUserSettingActor) OnReceive(ctx context.Context, input proto.Message) {
	userId := bases.GetTargetIdFromCtx(ctx)
	appkey := bases.GetAppKeyFromCtx(ctx)
	user, exist := services.GetUserInfo(appkey, userId)
	ret := &pbobjs.UserInfo{}
	if exist && user != nil {
		ret.UserId = user.UserId
		ret.Settings = commonservices.Map2KvItems(user.SettingFields)
	}
	queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, ret)
	actor.Sender.Tell(queryAck, actorsystem.NoSender)
}

func (actor *GetUserSettingActor) CreateInputObj() proto.Message {
	return &pbobjs.UserIdReq{}
}
