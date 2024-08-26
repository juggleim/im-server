package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"sort"
	"strings"
)

func SaveConversation(ctx context.Context, userId string, msg *pbobjs.DownMsg) {
	if IsStoreMsg(msg.Flags) {
		bases.AsyncRpcCall(ctx, "add_conver", userId, &pbobjs.Conversation{
			TargetId:    msg.TargetId,
			ChannelType: msg.ChannelType,
			Msg:         msg,
		})
	}
}

func GetConversationId(fromId, targetId string, channelType pbobjs.ChannelType) string {
	if channelType == pbobjs.ChannelType_Private || channelType == pbobjs.ChannelType_System {
		array := []string{fromId, targetId}
		sort.Strings(array)
		return strings.Join(array, ":")
	} else if channelType == pbobjs.ChannelType_BroadCast {
		return fromId
	} else {
		return targetId
	}
}
