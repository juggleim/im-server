package convercallers

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
)

var method = "upd_latest_msg"

func UpdLatestMsgBody(ctx context.Context, userId, targetId string, channelType pbobjs.ChannelType, latestMsgId string, msg *pbobjs.DownMsg) {
	bases.AsyncRpcCall(ctx, method, userId, &pbobjs.UpdLatestMsgReq{
		TargetId:    targetId,
		ChannelType: channelType,
		LatestMsgId: latestMsgId,
		Action:      pbobjs.UpdLatestMsgAction_UpdMsg,
		Msg:         msg,
	})
}
