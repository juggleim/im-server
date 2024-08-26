package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/group/services"
	"sync"

	"google.golang.org/protobuf/proto"
)

type QryGroupInfoByIdsActor struct {
	bases.BaseActor
}

func (actor *QryGroupInfoByIdsActor) OnReceive(ctx context.Context, input proto.Message) {
	if grpIdsReq, ok := input.(*pbobjs.GroupIdsReq); ok {
		appkey := bases.GetAppKeyFromCtx(ctx)
		ret := &pbobjs.GroupInfosResp{
			GroupInfoMap: make(map[string]*pbobjs.GroupInfo),
		}
		if grpIdsReq.NoDispatch {
			for _, groupId := range grpIdsReq.GroupIds {
				groupInfo, exist := services.GetGroupInfoFromCache(ctx, appkey, groupId)
				if exist && groupInfo != nil {
					ret.GroupInfoMap[groupId] = &pbobjs.GroupInfo{
						GroupId:       groupId,
						GroupName:     groupInfo.GroupName,
						GroupPortrait: groupInfo.GroupPortrait,
						UpdatedTime:   groupInfo.UpdatedTime,
						ExtFields:     commonservices.Map2KvItems(groupInfo.ExtFields),
					}
				}
			}
		} else {
			groups := bases.GroupTargets("qry_group_info_by_ids", grpIdsReq.GroupIds)
			currentNodeName := ""
			if bases.GetCluster() != nil && bases.GetCluster().GetCurrentNode() != nil {
				currentNodeName = bases.GetCluster().GetCurrentNode().Name
			}
			waitGroup := sync.WaitGroup{}
			tmpMap := sync.Map{}
			for nodeName, groupIds := range groups {
				if nodeName == currentNodeName {
					for _, groupId := range groupIds {
						groupInfo, exist := services.GetGroupInfoFromCache(ctx, appkey, groupId)
						if exist && groupInfo != nil {
							ret.GroupInfoMap[groupId] = &pbobjs.GroupInfo{
								GroupId:       groupId,
								GroupName:     groupInfo.GroupName,
								GroupPortrait: groupInfo.GroupPortrait,
								UpdatedTime:   groupInfo.UpdatedTime,
								ExtFields:     commonservices.Map2KvItems(groupInfo.ExtFields),
							}
						}
					}
				} else { //todo for cluster
					waitGroup.Add(1)
					curGroupIds := groupIds
					go func() {
						waitGroup.Done()
						_, respObj, err := bases.SyncRpcCall(ctx, "qry_group_info_by_ids", curGroupIds[0], &pbobjs.GroupIdsReq{
							GroupIds:   curGroupIds,
							NoDispatch: true,
						}, func() proto.Message {
							return &pbobjs.GroupInfosResp{}
						})
						if err == nil {
							resp, ok := respObj.(*pbobjs.GroupInfosResp)
							if ok {
								for gid, ginfo := range resp.GroupInfoMap {
									tmpMap.Store(gid, ginfo)
								}
							}
						}
					}()
				}
			}
			waitGroup.Wait()
			tmpMap.Range(func(key, value any) bool {
				gid := key.(string)
				ginfo := value.(*pbobjs.GroupInfo)
				ret.GroupInfoMap[gid] = ginfo
				return true
			})
		}
		queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, ret)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryGroupInfoByIdsActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupIdsReq{}
}
