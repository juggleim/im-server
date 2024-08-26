package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"
	"sync"

	"google.golang.org/protobuf/proto"
)

type QryUserStatusActor struct {
	bases.BaseActor
}

func (actor *QryUserStatusActor) OnReceive(ctx context.Context, input proto.Message) {
	if uIds, ok := input.(*pbobjs.UserIdsReq); ok {
		resp := &pbobjs.UserInfos{
			UserInfos: []*pbobjs.UserInfo{},
		}
		if len(uIds.UserIds) > 0 {
			wg := sync.WaitGroup{}
			disMap := bases.GroupTargets("inner_qry_user_status", uIds.UserIds)
			for nodeName, ids := range disMap {
				if nodeName == bases.GetCluster().GetCurrentNode().Name {
					rpcResp := services.QryUserStatus(ctx, ids)
					resp.UserInfos = append(resp.UserInfos, rpcResp.UserInfos...)
				} else {
					wg.Add(1)
					tarIds := ids
					go func() {
						defer wg.Done()
						code, rpcRespObj, err := bases.SyncRpcCall(ctx, "inner_qry_user_status", tarIds[0], &pbobjs.UserIdsReq{}, func() proto.Message {
							return &pbobjs.UserInfos{}
						}, &bases.TargetIdsOption{
							TargetIds: tarIds,
						})
						if err == nil && code == errs.IMErrorCode_SUCCESS && rpcRespObj != nil {
							rpcResp := rpcRespObj.(*pbobjs.UserInfos)
							resp.UserInfos = append(resp.UserInfos, rpcResp.UserInfos...)
						}
					}()
				}
			}
			wg.Wait()
		}
		queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, resp)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)

	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryUserStatusActor) CreateInputObj() proto.Message {
	return &pbobjs.UserIdsReq{}
}

type InnerQryUserStatusActor struct {
	bases.BaseActor
}

func (actor *InnerQryUserStatusActor) OnReceive(ctx context.Context, input proto.Message) {
	if _, ok := input.(*pbobjs.UserIdsReq); ok {
		userId := bases.GetRequesterIdFromCtx(ctx)
		targetIds := bases.GetTargetIdsFromCtx(ctx)
		logs.WithContext(ctx).Infof("user_id:%s\ttarget_ids:%v", userId, targetIds)
		resp := services.QryUserStatus(ctx, targetIds)
		queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, resp)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}
func (actor *InnerQryUserStatusActor) CreateInputObj() proto.Message {
	return &pbobjs.UserIdsReq{}
}
