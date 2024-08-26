package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
)

func Dispatch2Conversation(ctx context.Context, groupId string, req *pbobjs.UpdLatestMsgReq) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	groupContainer, exist := GetGroupMembersFromCache(ctx, appkey, groupId)
	memberIds := []string{}
	if exist {
		memberMap := groupContainer.GetMemberMap()
		for memberId := range memberMap {
			memberIds = append(memberIds, memberId)
		}
	}
	groups := bases.GroupTargets("upd_latest_msg", memberIds)
	data, _ := tools.PbMarshal(req)
	for _, ids := range groups {
		bases.UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
			RpcMsgType:   pbobjs.RpcMsgType_UserPub,
			AppKey:       bases.GetAppKeyFromCtx(ctx),
			Session:      bases.GetSessionFromCtx(ctx),
			Method:       "upd_latest_msg",
			RequesterId:  bases.GetRequesterIdFromCtx(ctx),
			ReqIndex:     bases.GetSeqIndexFromCtx(ctx),
			Qos:          bases.GetQosFromCtx(ctx),
			AppDataBytes: data,
			TargetId:     ids[0],
			GroupId:      groupId,
			TargetIds:    ids,
		})
	}

}
