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
		queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, ret)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryGroupInfoByIdsActor) CreateInputObj() proto.Message {
	return &pbobjs.GroupIdsReq{}
}
