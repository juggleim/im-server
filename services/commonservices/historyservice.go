package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
)

func SaveHistoryMsg(ctx context.Context, senderId, targetId string, channelType pbobjs.ChannelType, downMsg *pbobjs.DownMsg, memberCount int) {
	if IsStoreMsg(downMsg.Flags) {
		addHisMsgReq := pbobjs.AddHisMsgReq{
			SenderId:         senderId,
			TargetId:         targetId,
			ChannelType:      channelType,
			SendTime:         downMsg.MsgTime,
			Msg:              downMsg,
			GroupMemberCount: int32(memberCount),
		}
		data, _ := tools.PbMarshal(&addHisMsgReq)

		bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
			RpcMsgType:   pbobjs.RpcMsgType_UserPub,
			AppKey:       bases.GetAppKeyFromCtx(ctx),
			Session:      bases.GetSessionFromCtx(ctx),
			Method:       "add_hismsg",
			RequesterId:  bases.GetRequesterIdFromCtx(ctx),
			ReqIndex:     bases.GetSeqIndexFromCtx(ctx),
			Qos:          bases.GetQosFromCtx(ctx),
			AppDataBytes: data,
			TargetId:     GetConversationId(senderId, targetId, channelType),
		})
	}
}
