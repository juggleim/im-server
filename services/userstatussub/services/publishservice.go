package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"time"
)

func PublishStatus(ctx context.Context, req *pbobjs.UserInfo) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	resp := &pbobjs.UserInfos{
		UserInfos: []*pbobjs.UserInfo{},
	}
	resp.UserInfos = append(resp.UserInfos, req)
	rels := GetSubRelationsFromCache(appkey, req.UserId)
	targetIds := rels.GetSubscriptions()

	for _, targetId := range targetIds {
		rpcMsg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "ustatus", resp)
		rpcMsg.Qos = 0
		bases.UnicastRouteWithNoSender(rpcMsg)
		time.Sleep(5 * time.Millisecond)
	}
	return errs.IMErrorCode_SUCCESS
}
