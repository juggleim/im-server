package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/apimodels"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
)

func SendGrpNotify(ctx context.Context, grpId string, notify *apimodels.GroupNotify) {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	bs, _ := tools.JsonMarshal(notify)
	flag := msgdefines.SetStoreMsg(0)
	commonservices.AsyncGroupMsgOverUpstream(ctx, requestId, grpId, &pbobjs.UpMsg{
		MsgType:    apimodels.GroupNotifyMsgType,
		MsgContent: bs,
		Flags:      flag,
	}, &bases.MarkFromApiOption{})
}

func SendFriendNotify(ctx context.Context, targetId string, notify *apimodels.FriendNotify) {
	bs, _ := tools.JsonMarshal(notify)
	flag := msgdefines.SetStoreMsg(0)
	commonservices.AsyncPrivateMsgOverUpstream(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, &pbobjs.UpMsg{
		MsgType:    apimodels.FriendNotifyMsgType,
		MsgContent: bs,
		Flags:      flag,
	})
}

func SendFriendApplyNotify(ctx context.Context, targetId string, notify *apimodels.FriendApplyNotify) {
	bs, _ := tools.JsonMarshal(notify)
	flag := msgdefines.SetStoreMsg(0)
	flag = msgdefines.SetCountMsg(flag)
	commonservices.AsyncSystemMsg(ctx, apimodels.SystemFriendApplyConverId, targetId, &pbobjs.UpMsg{
		MsgType:    apimodels.FriendApplicationMsgType,
		MsgContent: bs,
		Flags:      flag,
	})
}
