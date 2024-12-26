package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/models"
	"im-server/services/commonservices"
)

func SendGrpNotify(ctx context.Context, grpId string, notify *models.GroupNotify) {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	bs, _ := tools.JsonMarshal(notify)
	flag := commonservices.SetStoreMsg(0)
	commonservices.AsyncGroupMsgOverUpstream(ctx, requestId, grpId, &pbobjs.UpMsg{
		MsgType:    models.GroupNotifyMsgType,
		MsgContent: bs,
		Flags:      flag,
	}, &bases.MarkFromApiOption{})
}

func SendFriendNotify(ctx context.Context, targetId string, notify *models.FriendNotify) {
	bs, _ := tools.JsonMarshal(notify)
	flag := commonservices.SetStoreMsg(0)
	commonservices.AsyncPrivateMsgOverUpstream(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, &pbobjs.UpMsg{
		MsgType:    models.FriendNotifyMsgType,
		MsgContent: bs,
		Flags:      flag,
	})
}

func SendFriendApplyNotify(ctx context.Context, targetId string, notify *models.FriendApplyNotify) {
	bs, _ := tools.JsonMarshal(notify)
	flag := commonservices.SetStoreMsg(0)
	flag = commonservices.SetCountMsg(flag)
	commonservices.AsyncSystemMsg(ctx, models.SystemFriendApplyConverId, targetId, &pbobjs.UpMsg{
		MsgType:    models.FriendApplicationMsgType,
		MsgContent: bs,
		Flags:      flag,
	})
}
