package services

import (
	"context"
	"strings"
	"sync"

	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"

	"google.golang.org/protobuf/proto"
)

func QryUserStatus(ctx context.Context, userIds []string) *pbobjs.UserStatusList {
	ret := &pbobjs.UserStatusList{
		Items: []*pbobjs.UserStatusItem{},
	}
	if len(userIds) == 0 {
		return ret
	}
	trimmed := make([]string, 0, len(userIds))
	for _, id := range userIds {
		id = strings.TrimSpace(id)
		if id != "" {
			trimmed = append(trimmed, id)
		}
	}
	if len(trimmed) == 0 {
		return ret
	}

	tmpMap := sync.Map{}
	groups := bases.GroupTargets("qry_online_status", trimmed)
	var wg sync.WaitGroup
	for _, ids := range groups {
		wg.Add(1)
		uIds := ids
		go func() {
			defer wg.Done()
			_, resp, err := bases.SyncRpcCall(ctx, "qry_online_status", uIds[0], &pbobjs.UserOnlineStatusReq{
				UserIds: uIds,
			}, func() proto.Message {
				return &pbobjs.UserOnlineStatusResp{}
			})
			if err == nil {
				onlineResp, ok := resp.(*pbobjs.UserOnlineStatusResp)
				if ok && onlineResp != nil && len(onlineResp.Items) > 0 {
					for _, item := range onlineResp.Items {
						if item != nil && item.UserId != "" {
							tmpMap.Store(item.UserId, item)
						}
					}
				}
			}
		}()
	}
	wg.Wait()

	for _, uid := range trimmed {
		it := &pbobjs.UserStatusItem{UserId: uid}
		if v, ok := tmpMap.Load(uid); ok {
			on := v.(*pbobjs.UserOnlineItem)
			it.OnlineStatus = &pbobjs.UserOnlineItem{
				UserId:   on.UserId,
				IsOnline: on.IsOnline,
			}
		} else {
			it.OnlineStatus = &pbobjs.UserOnlineItem{
				UserId:   uid,
				IsOnline: false,
			}
		}
		ret.Items = append(ret.Items, it)
	}
	return ret
}
