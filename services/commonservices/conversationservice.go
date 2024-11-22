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

func BatchSaveConversation(ctx context.Context, convers []*pbobjs.Conversation) {
	method := "batch_add_conver"
	grps := map[string][]*pbobjs.Conversation{}
	for _, conver := range convers {
		if conver.Msg == nil || !IsStoreMsg(conver.Msg.Flags) {
			continue
		}
		node := bases.GetCluster().GetTargetNode(method, conver.UserId)
		if node != nil {
			var arr []*pbobjs.Conversation
			if existArr, ok := grps[node.Name]; ok {
				arr = existArr
			} else {
				arr = []*pbobjs.Conversation{}
			}
			arr = append(arr, &pbobjs.Conversation{
				UserId: conver.UserId,
				Msg:    conver.Msg,
			})
			grps[node.Name] = arr
		}
	}
	for _, items := range grps {
		bases.AsyncRpcCall(ctx, method, items[0].UserId, &pbobjs.BatchAddConvers{
			Convers: items,
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
