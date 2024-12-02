package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/models"
	"im-server/services/commonservices"

	"google.golang.org/protobuf/proto"
)

func AppSyncRpcCall(ctx context.Context, method string, requesterId, targetId string, req proto.Message, respFactor func() proto.Message) (errs.IMErrorCode, proto.Message, error) {
	ctx = context.WithValue(ctx, bases.CtxKey_RequesterId, requesterId)
	ctx = context.WithValue(ctx, bases.CtxKey_IsFromApp, true)
	code, resp, err := bases.SyncRpcCall(ctx, method, targetId, req, respFactor)
	return code, resp, err
}

func AppAsyncRpcCall(ctx context.Context, method string, requesterId, targetId string, req proto.Message) {
	ctx = context.WithValue(ctx, bases.CtxKey_RequesterId, requesterId)
	ctx = context.WithValue(ctx, bases.CtxKey_IsFromApp, true)
	bases.AsyncRpcCall(ctx, method, targetId, req)
}

func AppAsyncRpcCallWithSender(ctx context.Context, method string, requesterId, targetId string, req proto.Message, sender actorsystem.ActorRef) {
	ctx = context.WithValue(ctx, bases.CtxKey_RequesterId, requesterId)
	ctx = context.WithValue(ctx, bases.CtxKey_IsFromApp, true)
	bases.AsyncRpcCallWithSender(ctx, method, targetId, req, sender)
}

func SendGrpNotify(ctx context.Context, grpId string, notify *models.GroupNotify) {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	bs, _ := tools.JsonMarshal(notify)
	flag := commonservices.SetStoreMsg(0)
	commonservices.GroupMsgFromApi(ctx, requestId, grpId, &pbobjs.UpMsg{
		MsgType:    models.GroupNotifyMsgType,
		MsgContent: bs,
		Flags:      flag,
	}, false)
}

func SendFriendNotify(ctx context.Context, targetId string, notify *models.FriendNotify) {
	bs, _ := tools.JsonMarshal(notify)
	flag := commonservices.SetStoreMsg(0)
	commonservices.AsyncPrivateMsg(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, &pbobjs.UpMsg{
		MsgType:    models.FriendNotifyMsgType,
		MsgContent: bs,
		Flags:      flag,
	})
}
