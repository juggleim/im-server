package actors

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/usermanager/services"
	"sync"

	"google.golang.org/protobuf/proto"
)

type QryUserInfoByIdsActor struct {
	bases.BaseActor
}

func (actor *QryUserInfoByIdsActor) OnReceive(ctx context.Context, input proto.Message) {
	if userIdsReq, ok := input.(*pbobjs.UserIdsReq); ok {
		if len(userIdsReq.AttTypes) <= 0 {
			userIdsReq.AttTypes = append(userIdsReq.AttTypes, int32(commonservices.AttItemType_Att))
		}
		appkey := bases.GetAppKeyFromCtx(ctx)
		ret := &pbobjs.UserInfosResp{
			UserInfoMap: make(map[string]*pbobjs.UserInfo),
		}
		if userIdsReq.NoDispatch {
			for _, userId := range userIdsReq.UserIds {
				userInfo, exist := services.GetUserInfo(appkey, userId)
				if exist && userInfo != nil {
					pbUserInfo := &pbobjs.UserInfo{
						UserId:   userId,
						UserType: pbobjs.UserType(userInfo.UserType),
					}
					for _, attType := range userIdsReq.AttTypes {
						if attType == int32(commonservices.AttItemType_Att) {
							pbUserInfo.Nickname = userInfo.Nickname
							pbUserInfo.UserPortrait = userInfo.UserPortrait
							pbUserInfo.UpdatedTime = userInfo.UpdatedTime
							pbUserInfo.ExtFields = commonservices.Map2KvItems(userInfo.ExtFields)
						} else if attType == int32(commonservices.AttItemType_Setting) {
							pbUserInfo.Settings = commonservices.Map2KvItems(userInfo.SettingFields)
						}
					}
					ret.UserInfoMap[userId] = pbUserInfo
				}
			}
		} else {
			groups := bases.GroupTargets("qry_user_info_by_ids", userIdsReq.UserIds)
			currentNodeName := ""
			if bases.GetCluster() != nil && bases.GetCluster().GetCurrentNode() != nil {
				currentNodeName = bases.GetCluster().GetCurrentNode().Name
			}
			waitGroup := sync.WaitGroup{}
			tmpMap := sync.Map{}
			for nodeName, userIds := range groups {
				if nodeName == currentNodeName {
					for _, userId := range userIds {
						userInfo, exist := services.GetUserInfo(appkey, userId)
						if exist && userInfo != nil {
							pbUserInfo := &pbobjs.UserInfo{
								UserId: userId,
							}
							for _, attType := range userIdsReq.AttTypes {
								if attType == int32(commonservices.AttItemType_Att) {
									pbUserInfo.Nickname = userInfo.Nickname
									pbUserInfo.UserPortrait = userInfo.UserPortrait
									pbUserInfo.UpdatedTime = userInfo.UpdatedTime
									pbUserInfo.ExtFields = commonservices.Map2KvItems(userInfo.ExtFields)
								} else if attType == int32(commonservices.AttItemType_Setting) {
									pbUserInfo.Settings = commonservices.Map2KvItems(userInfo.SettingFields)
								}
							}
							ret.UserInfoMap[userId] = pbUserInfo
						}
					}
				} else { //todo for cluster
					waitGroup.Add(1)
					curUserIds := userIds
					go func() {
						waitGroup.Done()
						_, respObj, err := bases.SyncRpcCall(ctx, "qry_user_info_by_ids", curUserIds[0], &pbobjs.UserIdsReq{
							UserIds:    curUserIds,
							AttTypes:   userIdsReq.AttTypes,
							NoDispatch: true,
						}, func() proto.Message {
							return &pbobjs.UserInfosResp{}
						})
						if err == nil {
							resp, ok := respObj.(*pbobjs.UserInfosResp)
							if ok {
								for uid, uinfo := range resp.UserInfoMap {
									tmpMap.Store(uid, uinfo)
								}
							}
						}
					}()
				}
			}
			waitGroup.Wait()
			tmpMap.Range(func(key, value any) bool {
				uid := key.(string)
				uinfo := value.(*pbobjs.UserInfo)
				ret.UserInfoMap[uid] = uinfo
				return true
			})
		}
		queryAck := bases.CreateQueryAckWraper(ctx, errs.IMErrorCode_SUCCESS, ret)
		actor.Sender.Tell(queryAck, actorsystem.NoSender)
	} else {
		logs.WithContext(ctx).Errorf("input is illigal")
	}
}

func (actor *QryUserInfoByIdsActor) CreateInputObj() proto.Message {
	return &pbobjs.UserIdsReq{}
}
