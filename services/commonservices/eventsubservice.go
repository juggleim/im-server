package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/msgdefines"
)

func SubPrivateMsg(ctx context.Context, targetId string, msg *pbobjs.DownMsg) {
	if msgdefines.IsCmdMsg(msg.Flags) {
		return
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	appInfo, ok := GetAppInfo(appkey)
	if ok && appInfo.EventSubSwitchObj != nil && appInfo.EventSubSwitchObj.PrivateMsgSubSwitch > 0 {
		if bases.GetIsFromApiFromCtx(ctx) && !appInfo.IsSubApiMsg {
			return
		}
		bases.AsyncRpcCall(ctx, "msg_sub", targetId, &pbobjs.SubMsgs{
			SubMsgs: []*pbobjs.SubMsg{
				{
					Platform: bases.GetPlatformFromCtx(ctx),
					Msg:      msg,
				},
			},
		})
	}
}

func SubGroupMsg(ctx context.Context, targetId string, msg *pbobjs.DownMsg) {
	if msgdefines.IsCmdMsg(msg.Flags) {
		return
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	appInfo, ok := GetAppInfo(appkey)
	if ok && appInfo.EventSubConfigObj != nil && appInfo.EventSubSwitchObj.GroupMsgSubSwitch > 0 {
		if bases.GetIsFromApiFromCtx(ctx) && !appInfo.IsSubApiMsg {
			return
		}
		bases.AsyncRpcCall(ctx, "msg_sub", targetId, &pbobjs.SubMsgs{
			SubMsgs: []*pbobjs.SubMsg{
				{
					Platform: bases.GetPlatformFromCtx(ctx),
					Msg:      msg,
				},
			},
		})
	}
}

func SubOnlineEvent(ctx context.Context, targetId string, msg *pbobjs.OnlineOfflineMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	appInfo, ok := GetAppInfo(appkey)
	if ok && appInfo.EventSubSwitchObj != nil && appInfo.EventSubSwitchObj.OnlineSubSwitch > 0 {
		bases.AsyncRpcCall(ctx, "online_offline_sub", targetId, msg)
	}
}

func SubOfflineEvent(ctx context.Context, targetId string, msg *pbobjs.OnlineOfflineMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	appInfo, ok := GetAppInfo(appkey)
	if ok && appInfo.EventSubSwitchObj != nil && appInfo.EventSubSwitchObj.OfflineSubSwitch > 0 {
		bases.AsyncRpcCall(ctx, "online_offline_sub", targetId, msg)
	}
}
